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

// UsersLoaderByIDConfig captures the config to create a new UsersLoaderByID
type UsersLoaderByIDConfig struct {
	// Fetch is a method that provides the data for the loader
<<<<<<< HEAD
<<<<<<< HEAD
	Fetch func(keys []persist.DBID) ([][]coredb.User, []error)
=======
	Fetch func(keys []persist.DBID) ([][]coregen.User, []error)
>>>>>>> 93a3a41 (Add indexer models)
=======
	Fetch func(keys []persist.DBID) ([][]coregen.User, []error)
>>>>>>> a4e9c3f (Add indexer models)

	// Wait is how long wait before sending a batch
	Wait time.Duration

	// MaxBatch will limit the maximum number of keys to send in one batch, 0 = not limit
	MaxBatch int
}

// NewUsersLoaderByID creates a new UsersLoaderByID given a fetch, wait, and maxBatch
func NewUsersLoaderByID(config UsersLoaderByIDConfig) *UsersLoaderByID {
	return &UsersLoaderByID{
		fetch:    config.Fetch,
		wait:     config.Wait,
		maxBatch: config.MaxBatch,
	}
}

// UsersLoaderByID batches and caches requests
type UsersLoaderByID struct {
	// this method provides the data for the loader
<<<<<<< HEAD
<<<<<<< HEAD
	fetch func(keys []persist.DBID) ([][]coredb.User, []error)
=======
	fetch func(keys []persist.DBID) ([][]coregen.User, []error)
>>>>>>> 93a3a41 (Add indexer models)
=======
	fetch func(keys []persist.DBID) ([][]coregen.User, []error)
>>>>>>> a4e9c3f (Add indexer models)

	// how long to done before sending a batch
	wait time.Duration

	// this will limit the maximum number of keys to send in one batch, 0 = no limit
	maxBatch int

	// INTERNAL

	// lazily created cache
<<<<<<< HEAD
<<<<<<< HEAD
	cache map[persist.DBID][]coredb.User
=======
	cache map[persist.DBID][]coregen.User
>>>>>>> 93a3a41 (Add indexer models)
=======
	cache map[persist.DBID][]coregen.User
>>>>>>> a4e9c3f (Add indexer models)

	// the current batch. keys will continue to be collected until timeout is hit,
	// then everything will be sent to the fetch method and out to the listeners
	batch *usersLoaderByIDBatch

	// mutex to prevent races
	mu sync.Mutex
}

type usersLoaderByIDBatch struct {
	keys    []persist.DBID
<<<<<<< HEAD
<<<<<<< HEAD
	data    [][]coredb.User
=======
	data    [][]coregen.User
>>>>>>> 93a3a41 (Add indexer models)
=======
	data    [][]coregen.User
>>>>>>> a4e9c3f (Add indexer models)
	error   []error
	closing bool
	done    chan struct{}
}

// Load a User by key, batching and caching will be applied automatically
<<<<<<< HEAD
<<<<<<< HEAD
func (l *UsersLoaderByID) Load(key persist.DBID) ([]coredb.User, error) {
=======
func (l *UsersLoaderByID) Load(key persist.DBID) ([]coregen.User, error) {
>>>>>>> 93a3a41 (Add indexer models)
=======
func (l *UsersLoaderByID) Load(key persist.DBID) ([]coregen.User, error) {
>>>>>>> a4e9c3f (Add indexer models)
	return l.LoadThunk(key)()
}

// LoadThunk returns a function that when called will block waiting for a User.
// This method should be used if you want one goroutine to make requests to many
// different data loaders without blocking until the thunk is called.
<<<<<<< HEAD
<<<<<<< HEAD
func (l *UsersLoaderByID) LoadThunk(key persist.DBID) func() ([]coredb.User, error) {
	l.mu.Lock()
	if it, ok := l.cache[key]; ok {
		l.mu.Unlock()
		return func() ([]coredb.User, error) {
=======
func (l *UsersLoaderByID) LoadThunk(key persist.DBID) func() ([]coregen.User, error) {
	l.mu.Lock()
	if it, ok := l.cache[key]; ok {
		l.mu.Unlock()
		return func() ([]coregen.User, error) {
>>>>>>> 93a3a41 (Add indexer models)
=======
func (l *UsersLoaderByID) LoadThunk(key persist.DBID) func() ([]coregen.User, error) {
	l.mu.Lock()
	if it, ok := l.cache[key]; ok {
		l.mu.Unlock()
		return func() ([]coregen.User, error) {
>>>>>>> a4e9c3f (Add indexer models)
			return it, nil
		}
	}
	if l.batch == nil {
		l.batch = &usersLoaderByIDBatch{done: make(chan struct{})}
	}
	batch := l.batch
	pos := batch.keyIndex(l, key)
	l.mu.Unlock()

<<<<<<< HEAD
<<<<<<< HEAD
	return func() ([]coredb.User, error) {
		<-batch.done

		var data []coredb.User
=======
=======
>>>>>>> a4e9c3f (Add indexer models)
	return func() ([]coregen.User, error) {
		<-batch.done

		var data []coregen.User
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
func (l *UsersLoaderByID) LoadAll(keys []persist.DBID) ([][]coredb.User, []error) {
	results := make([]func() ([]coredb.User, error), len(keys))
=======
func (l *UsersLoaderByID) LoadAll(keys []persist.DBID) ([][]coregen.User, []error) {
	results := make([]func() ([]coregen.User, error), len(keys))
>>>>>>> 93a3a41 (Add indexer models)
=======
func (l *UsersLoaderByID) LoadAll(keys []persist.DBID) ([][]coregen.User, []error) {
	results := make([]func() ([]coregen.User, error), len(keys))
>>>>>>> a4e9c3f (Add indexer models)

	for i, key := range keys {
		results[i] = l.LoadThunk(key)
	}

<<<<<<< HEAD
<<<<<<< HEAD
	users := make([][]coredb.User, len(keys))
=======
	users := make([][]coregen.User, len(keys))
>>>>>>> 93a3a41 (Add indexer models)
=======
	users := make([][]coregen.User, len(keys))
>>>>>>> a4e9c3f (Add indexer models)
	errors := make([]error, len(keys))
	for i, thunk := range results {
		users[i], errors[i] = thunk()
	}
	return users, errors
}

// LoadAllThunk returns a function that when called will block waiting for a Users.
// This method should be used if you want one goroutine to make requests to many
// different data loaders without blocking until the thunk is called.
<<<<<<< HEAD
<<<<<<< HEAD
func (l *UsersLoaderByID) LoadAllThunk(keys []persist.DBID) func() ([][]coredb.User, []error) {
	results := make([]func() ([]coredb.User, error), len(keys))
	for i, key := range keys {
		results[i] = l.LoadThunk(key)
	}
	return func() ([][]coredb.User, []error) {
		users := make([][]coredb.User, len(keys))
=======
func (l *UsersLoaderByID) LoadAllThunk(keys []persist.DBID) func() ([][]coregen.User, []error) {
	results := make([]func() ([]coregen.User, error), len(keys))
	for i, key := range keys {
		results[i] = l.LoadThunk(key)
	}
	return func() ([][]coregen.User, []error) {
		users := make([][]coregen.User, len(keys))
>>>>>>> 93a3a41 (Add indexer models)
=======
func (l *UsersLoaderByID) LoadAllThunk(keys []persist.DBID) func() ([][]coregen.User, []error) {
	results := make([]func() ([]coregen.User, error), len(keys))
	for i, key := range keys {
		results[i] = l.LoadThunk(key)
	}
	return func() ([][]coregen.User, []error) {
		users := make([][]coregen.User, len(keys))
>>>>>>> a4e9c3f (Add indexer models)
		errors := make([]error, len(keys))
		for i, thunk := range results {
			users[i], errors[i] = thunk()
		}
		return users, errors
	}
}

// Prime the cache with the provided key and value. If the key already exists, no change is made
// and false is returned.
// (To forcefully prime the cache, clear the key first with loader.clear(key).prime(key, value).)
<<<<<<< HEAD
<<<<<<< HEAD
func (l *UsersLoaderByID) Prime(key persist.DBID, value []coredb.User) bool {
=======
func (l *UsersLoaderByID) Prime(key persist.DBID, value []coregen.User) bool {
>>>>>>> 93a3a41 (Add indexer models)
=======
func (l *UsersLoaderByID) Prime(key persist.DBID, value []coregen.User) bool {
>>>>>>> a4e9c3f (Add indexer models)
	l.mu.Lock()
	var found bool
	if _, found = l.cache[key]; !found {
		// make a copy when writing to the cache, its easy to pass a pointer in from a loop var
		// and end up with the whole cache pointing to the same value.
<<<<<<< HEAD
<<<<<<< HEAD
		cpy := make([]coredb.User, len(value))
=======
		cpy := make([]coregen.User, len(value))
>>>>>>> 93a3a41 (Add indexer models)
=======
		cpy := make([]coregen.User, len(value))
>>>>>>> a4e9c3f (Add indexer models)
		copy(cpy, value)
		l.unsafeSet(key, cpy)
	}
	l.mu.Unlock()
	return !found
}

// Clear the value at key from the cache, if it exists
func (l *UsersLoaderByID) Clear(key persist.DBID) {
	l.mu.Lock()
	delete(l.cache, key)
	l.mu.Unlock()
}

<<<<<<< HEAD
<<<<<<< HEAD
func (l *UsersLoaderByID) unsafeSet(key persist.DBID, value []coredb.User) {
	if l.cache == nil {
		l.cache = map[persist.DBID][]coredb.User{}
=======
func (l *UsersLoaderByID) unsafeSet(key persist.DBID, value []coregen.User) {
	if l.cache == nil {
		l.cache = map[persist.DBID][]coregen.User{}
>>>>>>> 93a3a41 (Add indexer models)
=======
func (l *UsersLoaderByID) unsafeSet(key persist.DBID, value []coregen.User) {
	if l.cache == nil {
		l.cache = map[persist.DBID][]coregen.User{}
>>>>>>> a4e9c3f (Add indexer models)
	}
	l.cache[key] = value
}

// keyIndex will return the location of the key in the batch, if its not found
// it will add the key to the batch
func (b *usersLoaderByIDBatch) keyIndex(l *UsersLoaderByID, key persist.DBID) int {
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

func (b *usersLoaderByIDBatch) startTimer(l *UsersLoaderByID) {
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

func (b *usersLoaderByIDBatch) end(l *UsersLoaderByID) {
	b.data, b.error = l.fetch(b.keys)
	close(b.done)
}
