// Code generated by github.com/vektah/dataloaden, DO NOT EDIT.

package dataloader

import (
	"sync"
	"time"

<<<<<<< HEAD
<<<<<<< HEAD
	"github.com/mikeydub/go-gallery/db/gen/coredb"
=======
	"github.com/mikeydub/go-gallery/db/sqlc/coregen"
>>>>>>> 93a3a41 (Add indexer models)
=======
	"github.com/mikeydub/go-gallery/db/sqlc/coregen"
>>>>>>> a4e9c3f (Add indexer models)
	"github.com/mikeydub/go-gallery/service/persist"
)

// TokensLoaderByIDConfig captures the config to create a new TokensLoaderByID
type TokensLoaderByIDConfig struct {
	// Fetch is a method that provides the data for the loader
<<<<<<< HEAD
<<<<<<< HEAD
	Fetch func(keys []persist.DBID) ([][]coredb.Token, []error)
=======
	Fetch func(keys []persist.DBID) ([][]coregen.Token, []error)
>>>>>>> 93a3a41 (Add indexer models)
=======
	Fetch func(keys []persist.DBID) ([][]coregen.Token, []error)
>>>>>>> a4e9c3f (Add indexer models)

	// Wait is how long wait before sending a batch
	Wait time.Duration

	// MaxBatch will limit the maximum number of keys to send in one batch, 0 = not limit
	MaxBatch int
}

// NewTokensLoaderByID creates a new TokensLoaderByID given a fetch, wait, and maxBatch
func NewTokensLoaderByID(config TokensLoaderByIDConfig) *TokensLoaderByID {
	return &TokensLoaderByID{
		fetch:    config.Fetch,
		wait:     config.Wait,
		maxBatch: config.MaxBatch,
	}
}

// TokensLoaderByID batches and caches requests
type TokensLoaderByID struct {
	// this method provides the data for the loader
<<<<<<< HEAD
<<<<<<< HEAD
	fetch func(keys []persist.DBID) ([][]coredb.Token, []error)
=======
	fetch func(keys []persist.DBID) ([][]coregen.Token, []error)
>>>>>>> 93a3a41 (Add indexer models)
=======
	fetch func(keys []persist.DBID) ([][]coregen.Token, []error)
>>>>>>> a4e9c3f (Add indexer models)

	// how long to done before sending a batch
	wait time.Duration

	// this will limit the maximum number of keys to send in one batch, 0 = no limit
	maxBatch int

	// INTERNAL

	// lazily created cache
<<<<<<< HEAD
<<<<<<< HEAD
	cache map[persist.DBID][]coredb.Token
=======
	cache map[persist.DBID][]coregen.Token
>>>>>>> 93a3a41 (Add indexer models)
=======
	cache map[persist.DBID][]coregen.Token
>>>>>>> a4e9c3f (Add indexer models)

	// the current batch. keys will continue to be collected until timeout is hit,
	// then everything will be sent to the fetch method and out to the listeners
	batch *tokensLoaderByIDBatch

	// mutex to prevent races
	mu sync.Mutex
}

type tokensLoaderByIDBatch struct {
	keys    []persist.DBID
<<<<<<< HEAD
<<<<<<< HEAD
	data    [][]coredb.Token
=======
	data    [][]coregen.Token
>>>>>>> 93a3a41 (Add indexer models)
=======
	data    [][]coregen.Token
>>>>>>> a4e9c3f (Add indexer models)
	error   []error
	closing bool
	done    chan struct{}
}

// Load a Token by key, batching and caching will be applied automatically
<<<<<<< HEAD
<<<<<<< HEAD
func (l *TokensLoaderByID) Load(key persist.DBID) ([]coredb.Token, error) {
=======
func (l *TokensLoaderByID) Load(key persist.DBID) ([]coregen.Token, error) {
>>>>>>> 93a3a41 (Add indexer models)
=======
func (l *TokensLoaderByID) Load(key persist.DBID) ([]coregen.Token, error) {
>>>>>>> a4e9c3f (Add indexer models)
	return l.LoadThunk(key)()
}

// LoadThunk returns a function that when called will block waiting for a Token.
// This method should be used if you want one goroutine to make requests to many
// different data loaders without blocking until the thunk is called.
<<<<<<< HEAD
<<<<<<< HEAD
func (l *TokensLoaderByID) LoadThunk(key persist.DBID) func() ([]coredb.Token, error) {
	l.mu.Lock()
	if it, ok := l.cache[key]; ok {
		l.mu.Unlock()
		return func() ([]coredb.Token, error) {
=======
func (l *TokensLoaderByID) LoadThunk(key persist.DBID) func() ([]coregen.Token, error) {
	l.mu.Lock()
	if it, ok := l.cache[key]; ok {
		l.mu.Unlock()
		return func() ([]coregen.Token, error) {
>>>>>>> 93a3a41 (Add indexer models)
=======
func (l *TokensLoaderByID) LoadThunk(key persist.DBID) func() ([]coregen.Token, error) {
	l.mu.Lock()
	if it, ok := l.cache[key]; ok {
		l.mu.Unlock()
		return func() ([]coregen.Token, error) {
>>>>>>> a4e9c3f (Add indexer models)
			return it, nil
		}
	}
	if l.batch == nil {
		l.batch = &tokensLoaderByIDBatch{done: make(chan struct{})}
	}
	batch := l.batch
	pos := batch.keyIndex(l, key)
	l.mu.Unlock()

<<<<<<< HEAD
<<<<<<< HEAD
	return func() ([]coredb.Token, error) {
		<-batch.done

		var data []coredb.Token
=======
=======
>>>>>>> a4e9c3f (Add indexer models)
	return func() ([]coregen.Token, error) {
		<-batch.done

		var data []coregen.Token
<<<<<<< HEAD
>>>>>>> 93a3a41 (Add indexer models)
=======
>>>>>>> a4e9c3f (Add indexer models)
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
			l.mu.Lock()
			l.unsafeSet(key, data)
			l.mu.Unlock()
		}

		return data, err
	}
}

// LoadAll fetches many keys at once. It will be broken into appropriate sized
// sub batches depending on how the loader is configured
<<<<<<< HEAD
<<<<<<< HEAD
func (l *TokensLoaderByID) LoadAll(keys []persist.DBID) ([][]coredb.Token, []error) {
	results := make([]func() ([]coredb.Token, error), len(keys))
=======
func (l *TokensLoaderByID) LoadAll(keys []persist.DBID) ([][]coregen.Token, []error) {
	results := make([]func() ([]coregen.Token, error), len(keys))
>>>>>>> 93a3a41 (Add indexer models)
=======
func (l *TokensLoaderByID) LoadAll(keys []persist.DBID) ([][]coregen.Token, []error) {
	results := make([]func() ([]coregen.Token, error), len(keys))
>>>>>>> a4e9c3f (Add indexer models)

	for i, key := range keys {
		results[i] = l.LoadThunk(key)
	}

<<<<<<< HEAD
<<<<<<< HEAD
	tokens := make([][]coredb.Token, len(keys))
=======
	tokens := make([][]coregen.Token, len(keys))
>>>>>>> 93a3a41 (Add indexer models)
=======
	tokens := make([][]coregen.Token, len(keys))
>>>>>>> a4e9c3f (Add indexer models)
	errors := make([]error, len(keys))
	for i, thunk := range results {
		tokens[i], errors[i] = thunk()
	}
	return tokens, errors
}

// LoadAllThunk returns a function that when called will block waiting for a Tokens.
// This method should be used if you want one goroutine to make requests to many
// different data loaders without blocking until the thunk is called.
<<<<<<< HEAD
<<<<<<< HEAD
func (l *TokensLoaderByID) LoadAllThunk(keys []persist.DBID) func() ([][]coredb.Token, []error) {
	results := make([]func() ([]coredb.Token, error), len(keys))
	for i, key := range keys {
		results[i] = l.LoadThunk(key)
	}
	return func() ([][]coredb.Token, []error) {
		tokens := make([][]coredb.Token, len(keys))
=======
func (l *TokensLoaderByID) LoadAllThunk(keys []persist.DBID) func() ([][]coregen.Token, []error) {
	results := make([]func() ([]coregen.Token, error), len(keys))
	for i, key := range keys {
		results[i] = l.LoadThunk(key)
	}
	return func() ([][]coregen.Token, []error) {
		tokens := make([][]coregen.Token, len(keys))
>>>>>>> 93a3a41 (Add indexer models)
=======
func (l *TokensLoaderByID) LoadAllThunk(keys []persist.DBID) func() ([][]coregen.Token, []error) {
	results := make([]func() ([]coregen.Token, error), len(keys))
	for i, key := range keys {
		results[i] = l.LoadThunk(key)
	}
	return func() ([][]coregen.Token, []error) {
		tokens := make([][]coregen.Token, len(keys))
>>>>>>> a4e9c3f (Add indexer models)
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
<<<<<<< HEAD
<<<<<<< HEAD
func (l *TokensLoaderByID) Prime(key persist.DBID, value []coredb.Token) bool {
=======
func (l *TokensLoaderByID) Prime(key persist.DBID, value []coregen.Token) bool {
>>>>>>> 93a3a41 (Add indexer models)
=======
func (l *TokensLoaderByID) Prime(key persist.DBID, value []coregen.Token) bool {
>>>>>>> a4e9c3f (Add indexer models)
	l.mu.Lock()
	var found bool
	if _, found = l.cache[key]; !found {
		// make a copy when writing to the cache, its easy to pass a pointer in from a loop var
		// and end up with the whole cache pointing to the same value.
<<<<<<< HEAD
<<<<<<< HEAD
		cpy := make([]coredb.Token, len(value))
=======
		cpy := make([]coregen.Token, len(value))
>>>>>>> 93a3a41 (Add indexer models)
=======
		cpy := make([]coregen.Token, len(value))
>>>>>>> a4e9c3f (Add indexer models)
		copy(cpy, value)
		l.unsafeSet(key, cpy)
	}
	l.mu.Unlock()
	return !found
}

// Clear the value at key from the cache, if it exists
func (l *TokensLoaderByID) Clear(key persist.DBID) {
	l.mu.Lock()
	delete(l.cache, key)
	l.mu.Unlock()
}

<<<<<<< HEAD
<<<<<<< HEAD
func (l *TokensLoaderByID) unsafeSet(key persist.DBID, value []coredb.Token) {
	if l.cache == nil {
		l.cache = map[persist.DBID][]coredb.Token{}
=======
func (l *TokensLoaderByID) unsafeSet(key persist.DBID, value []coregen.Token) {
	if l.cache == nil {
		l.cache = map[persist.DBID][]coregen.Token{}
>>>>>>> 93a3a41 (Add indexer models)
=======
func (l *TokensLoaderByID) unsafeSet(key persist.DBID, value []coregen.Token) {
	if l.cache == nil {
		l.cache = map[persist.DBID][]coregen.Token{}
>>>>>>> a4e9c3f (Add indexer models)
	}
	l.cache[key] = value
}

// keyIndex will return the location of the key in the batch, if its not found
// it will add the key to the batch
func (b *tokensLoaderByIDBatch) keyIndex(l *TokensLoaderByID, key persist.DBID) int {
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

func (b *tokensLoaderByIDBatch) startTimer(l *TokensLoaderByID) {
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

func (b *tokensLoaderByIDBatch) end(l *TokensLoaderByID) {
	b.data, b.error = l.fetch(b.keys)
	close(b.done)
}
