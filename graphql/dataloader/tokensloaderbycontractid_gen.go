// Code generated by github.com/gallery-so/dataloaden, DO NOT EDIT.

package dataloader

import (
	"context"
	"sync"
	"time"

	"github.com/mikeydub/go-gallery/db/gen/coredb"
)

type TokensLoaderByContractIDSettings interface {
	getContext() context.Context
	getWait() time.Duration
	getMaxBatchOne() int
	getMaxBatchMany() int
	getDisableCaching() bool
	getPublishResults() bool
	getSubscriptionRegistry() *[]interface{}
	getMutexRegistry() *[]*sync.Mutex
}

func (l *TokensLoaderByContractID) setContext(ctx context.Context) {
	l.ctx = ctx
}

func (l *TokensLoaderByContractID) setWait(wait time.Duration) {
	l.wait = wait
}

func (l *TokensLoaderByContractID) setMaxBatch(maxBatch int) {
	l.maxBatch = maxBatch
}

func (l *TokensLoaderByContractID) setDisableCaching(disableCaching bool) {
	l.disableCaching = disableCaching
}

func (l *TokensLoaderByContractID) setPublishResults(publishResults bool) {
	l.publishResults = publishResults
}

// NewTokensLoaderByContractID creates a new TokensLoaderByContractID with the given settings, functions, and options
func NewTokensLoaderByContractID(
	settings TokensLoaderByContractIDSettings, fetch func(ctx context.Context, keys []coredb.GetTokensByContractIdBatchPaginateParams) ([][]coredb.Token, []error),
	opts ...func(interface {
		setContext(context.Context)
		setWait(time.Duration)
		setMaxBatch(int)
		setDisableCaching(bool)
		setPublishResults(bool)
	}),
) *TokensLoaderByContractID {
	loader := &TokensLoaderByContractID{
		ctx:                  settings.getContext(),
		wait:                 settings.getWait(),
		disableCaching:       settings.getDisableCaching(),
		publishResults:       settings.getPublishResults(),
		subscriptionRegistry: settings.getSubscriptionRegistry(),
		mutexRegistry:        settings.getMutexRegistry(),
		maxBatch:             settings.getMaxBatchMany(),
	}

	for _, opt := range opts {
		opt(loader)
	}

	// Set this after applying options, in case a different context was set via options
	loader.fetch = func(keys []coredb.GetTokensByContractIdBatchPaginateParams) ([][]coredb.Token, []error) {
		return fetch(loader.ctx, keys)
	}

	if loader.subscriptionRegistry == nil {
		panic("subscriptionRegistry may not be nil")
	}

	if loader.mutexRegistry == nil {
		panic("mutexRegistry may not be nil")
	}

	// No cache functions here; caching isn't very useful for dataloaders that return slices. This dataloader can
	// still send its results to other cache-priming receivers, but it won't register its own cache-priming function.

	return loader
}

// TokensLoaderByContractID batches and caches requests
type TokensLoaderByContractID struct {
	// context passed to fetch functions
	ctx context.Context

	// this method provides the data for the loader
	fetch func(keys []coredb.GetTokensByContractIdBatchPaginateParams) ([][]coredb.Token, []error)

	// how long to wait before sending a batch
	wait time.Duration

	// this will limit the maximum number of keys to send in one batch, 0 = no limit
	maxBatch int

	// whether this dataloader will cache results
	disableCaching bool

	// whether this dataloader will publish its results for others to cache
	publishResults bool

	// a shared slice where dataloaders will register and invoke caching functions.
	// the same slice should be passed to every dataloader.
	subscriptionRegistry *[]interface{}

	// a shared slice, parallel to the subscription registry, that holds a reference to the
	// cache mutex for the subscription's dataloader
	mutexRegistry *[]*sync.Mutex

	// INTERNAL

	// lazily created cache
	cache map[coredb.GetTokensByContractIdBatchPaginateParams][]coredb.Token

	// typed cache functions
	//subscribers []func([]coredb.Token)
	subscribers []tokensLoaderByContractIDSubscriber

	// functions used to cache published results from other dataloaders
	cacheFuncs []interface{}

	// the current batch. keys will continue to be collected until timeout is hit,
	// then everything will be sent to the fetch method and out to the listeners
	batch *tokensLoaderByContractIDBatch

	// mutex to prevent races
	mu sync.Mutex

	// only initialize our typed subscription cache once
	once sync.Once
}

type tokensLoaderByContractIDBatch struct {
	keys    []coredb.GetTokensByContractIdBatchPaginateParams
	data    [][]coredb.Token
	error   []error
	closing bool
	done    chan struct{}
}

// Load a Token by key, batching and caching will be applied automatically
func (l *TokensLoaderByContractID) Load(key coredb.GetTokensByContractIdBatchPaginateParams) ([]coredb.Token, error) {
	return l.LoadThunk(key)()
}

// LoadThunk returns a function that when called will block waiting for a Token.
// This method should be used if you want one goroutine to make requests to many
// different data loaders without blocking until the thunk is called.
func (l *TokensLoaderByContractID) LoadThunk(key coredb.GetTokensByContractIdBatchPaginateParams) func() ([]coredb.Token, error) {
	l.mu.Lock()
	if !l.disableCaching {
		if it, ok := l.cache[key]; ok {
			l.mu.Unlock()
			return func() ([]coredb.Token, error) {
				return it, nil
			}
		}
	}
	if l.batch == nil {
		l.batch = &tokensLoaderByContractIDBatch{done: make(chan struct{})}
	}
	batch := l.batch
	pos := batch.keyIndex(l, key)
	l.mu.Unlock()

	return func() ([]coredb.Token, error) {
		<-batch.done

		var data []coredb.Token
		if pos < len(batch.data) {
			data = batch.data[pos]
		}

		var err error
		// its convenient to be able to return a single error for everything
		if len(batch.error) == 1 {
			err = batch.error[0]
		} else if batch.error != nil {
			err = batch.error[pos]
		}

		if err == nil {
			if !l.disableCaching {
				l.mu.Lock()
				l.unsafeSet(key, data)
				l.mu.Unlock()
			}

			if l.publishResults {
				l.publishToSubscribers(data)
			}
		}

		return data, err
	}
}

// LoadAll fetches many keys at once. It will be broken into appropriate sized
// sub batches depending on how the loader is configured
func (l *TokensLoaderByContractID) LoadAll(keys []coredb.GetTokensByContractIdBatchPaginateParams) ([][]coredb.Token, []error) {
	results := make([]func() ([]coredb.Token, error), len(keys))

	for i, key := range keys {
		results[i] = l.LoadThunk(key)
	}

	tokens := make([][]coredb.Token, len(keys))
	errors := make([]error, len(keys))
	for i, thunk := range results {
		tokens[i], errors[i] = thunk()
	}
	return tokens, errors
}

// LoadAllThunk returns a function that when called will block waiting for a Tokens.
// This method should be used if you want one goroutine to make requests to many
// different data loaders without blocking until the thunk is called.
func (l *TokensLoaderByContractID) LoadAllThunk(keys []coredb.GetTokensByContractIdBatchPaginateParams) func() ([][]coredb.Token, []error) {
	results := make([]func() ([]coredb.Token, error), len(keys))
	for i, key := range keys {
		results[i] = l.LoadThunk(key)
	}
	return func() ([][]coredb.Token, []error) {
		tokens := make([][]coredb.Token, len(keys))
		errors := make([]error, len(keys))
		for i, thunk := range results {
			tokens[i], errors[i] = thunk()
		}
		return tokens, errors
	}
}

// Prime the cache with the provided key and value. If the key already exists, no change is made
// and false is returned.
// (To forcefully prime the cache, clear the key first with loader.clear(key).prime(key, value).)
func (l *TokensLoaderByContractID) Prime(key coredb.GetTokensByContractIdBatchPaginateParams, value []coredb.Token) bool {
	if l.disableCaching {
		return false
	}
	l.mu.Lock()
	var found bool
	if _, found = l.cache[key]; !found {
		// make a copy when writing to the cache, its easy to pass a pointer in from a loop var
		// and end up with the whole cache pointing to the same value.
		cpy := make([]coredb.Token, len(value))
		copy(cpy, value)
		l.unsafeSet(key, cpy)
	}
	l.mu.Unlock()
	return !found
}

// Clear the value at key from the cache, if it exists
func (l *TokensLoaderByContractID) Clear(key coredb.GetTokensByContractIdBatchPaginateParams) {
	if l.disableCaching {
		return
	}
	l.mu.Lock()
	delete(l.cache, key)
	l.mu.Unlock()
}

func (l *TokensLoaderByContractID) unsafeSet(key coredb.GetTokensByContractIdBatchPaginateParams, value []coredb.Token) {
	if l.cache == nil {
		l.cache = map[coredb.GetTokensByContractIdBatchPaginateParams][]coredb.Token{}
	}
	l.cache[key] = value
}

// keyIndex will return the location of the key in the batch, if its not found
// it will add the key to the batch
func (b *tokensLoaderByContractIDBatch) keyIndex(l *TokensLoaderByContractID, key coredb.GetTokensByContractIdBatchPaginateParams) int {
	for i, existingKey := range b.keys {
		if key == existingKey {
			return i
		}
	}

	pos := len(b.keys)
	b.keys = append(b.keys, key)
	if pos == 0 {
		go b.startTimer(l)
	}

	if l.maxBatch != 0 && pos >= l.maxBatch-1 {
		if !b.closing {
			b.closing = true
			l.batch = nil
			go b.end(l)
		}
	}

	return pos
}

func (b *tokensLoaderByContractIDBatch) startTimer(l *TokensLoaderByContractID) {
	time.Sleep(l.wait)
	l.mu.Lock()

	// we must have hit a batch limit and are already finalizing this batch
	if b.closing {
		l.mu.Unlock()
		return
	}

	l.batch = nil
	l.mu.Unlock()

	b.end(l)
}

func (b *tokensLoaderByContractIDBatch) end(l *TokensLoaderByContractID) {
	b.data, b.error = l.fetch(b.keys)
	close(b.done)
}

type tokensLoaderByContractIDSubscriber struct {
	cacheFunc func(coredb.Token)
	mutex     *sync.Mutex
}

func (l *TokensLoaderByContractID) publishToSubscribers(value []coredb.Token) {
	// Lazy build our list of typed cache functions once
	l.once.Do(func() {
		for i, subscription := range *l.subscriptionRegistry {
			if typedFunc, ok := subscription.(*func(coredb.Token)); ok {
				// Don't invoke our own cache function
				if !l.ownsCacheFunc(typedFunc) {
					l.subscribers = append(l.subscribers, tokensLoaderByContractIDSubscriber{cacheFunc: *typedFunc, mutex: (*l.mutexRegistry)[i]})
				}
			}
		}
	})

	// Handling locking here (instead of in the subscribed functions themselves) isn't the
	// ideal pattern, but it's an optimization that allows the publisher to iterate over slices
	// without having to acquire the lock many times.
	for _, s := range l.subscribers {
		s.mutex.Lock()
		for _, v := range value {
			s.cacheFunc(v)
		}
		s.mutex.Unlock()
	}
}

func (l *TokensLoaderByContractID) registerCacheFunc(cacheFunc interface{}, mutex *sync.Mutex) {
	l.cacheFuncs = append(l.cacheFuncs, cacheFunc)
	*l.subscriptionRegistry = append(*l.subscriptionRegistry, cacheFunc)
	*l.mutexRegistry = append(*l.mutexRegistry, mutex)
}

func (l *TokensLoaderByContractID) ownsCacheFunc(f *func(coredb.Token)) bool {
	for _, cacheFunc := range l.cacheFuncs {
		if cacheFunc == f {
			return true
		}
	}

	return false
}
