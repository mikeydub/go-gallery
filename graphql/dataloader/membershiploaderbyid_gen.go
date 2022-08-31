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

// MembershipLoaderByIdConfig captures the config to create a new MembershipLoaderById
type MembershipLoaderByIdConfig struct {
	// Fetch is a method that provides the data for the loader
<<<<<<< HEAD
<<<<<<< HEAD
	Fetch func(keys []persist.DBID) ([]coredb.Membership, []error)
=======
	Fetch func(keys []persist.DBID) ([]coregen.Membership, []error)
>>>>>>> 93a3a41 (Add indexer models)
=======
	Fetch func(keys []persist.DBID) ([]coregen.Membership, []error)
>>>>>>> a4e9c3f (Add indexer models)

	// Wait is how long wait before sending a batch
	Wait time.Duration

	// MaxBatch will limit the maximum number of keys to send in one batch, 0 = not limit
	MaxBatch int
}

// NewMembershipLoaderById creates a new MembershipLoaderById given a fetch, wait, and maxBatch
func NewMembershipLoaderById(config MembershipLoaderByIdConfig) *MembershipLoaderById {
	return &MembershipLoaderById{
		fetch:    config.Fetch,
		wait:     config.Wait,
		maxBatch: config.MaxBatch,
	}
}

// MembershipLoaderById batches and caches requests
type MembershipLoaderById struct {
	// this method provides the data for the loader
<<<<<<< HEAD
<<<<<<< HEAD
	fetch func(keys []persist.DBID) ([]coredb.Membership, []error)
=======
	fetch func(keys []persist.DBID) ([]coregen.Membership, []error)
>>>>>>> 93a3a41 (Add indexer models)
=======
	fetch func(keys []persist.DBID) ([]coregen.Membership, []error)
>>>>>>> a4e9c3f (Add indexer models)

	// how long to done before sending a batch
	wait time.Duration

	// this will limit the maximum number of keys to send in one batch, 0 = no limit
	maxBatch int

	// INTERNAL

	// lazily created cache
<<<<<<< HEAD
<<<<<<< HEAD
	cache map[persist.DBID]coredb.Membership
=======
	cache map[persist.DBID]coregen.Membership
>>>>>>> 93a3a41 (Add indexer models)
=======
	cache map[persist.DBID]coregen.Membership
>>>>>>> a4e9c3f (Add indexer models)

	// the current batch. keys will continue to be collected until timeout is hit,
	// then everything will be sent to the fetch method and out to the listeners
	batch *membershipLoaderByIdBatch

	// mutex to prevent races
	mu sync.Mutex
}

type membershipLoaderByIdBatch struct {
	keys    []persist.DBID
<<<<<<< HEAD
<<<<<<< HEAD
	data    []coredb.Membership
=======
	data    []coregen.Membership
>>>>>>> 93a3a41 (Add indexer models)
=======
	data    []coregen.Membership
>>>>>>> a4e9c3f (Add indexer models)
	error   []error
	closing bool
	done    chan struct{}
}

// Load a Membership by key, batching and caching will be applied automatically
<<<<<<< HEAD
<<<<<<< HEAD
func (l *MembershipLoaderById) Load(key persist.DBID) (coredb.Membership, error) {
=======
func (l *MembershipLoaderById) Load(key persist.DBID) (coregen.Membership, error) {
>>>>>>> 93a3a41 (Add indexer models)
=======
func (l *MembershipLoaderById) Load(key persist.DBID) (coregen.Membership, error) {
>>>>>>> a4e9c3f (Add indexer models)
	return l.LoadThunk(key)()
}

// LoadThunk returns a function that when called will block waiting for a Membership.
// This method should be used if you want one goroutine to make requests to many
// different data loaders without blocking until the thunk is called.
<<<<<<< HEAD
<<<<<<< HEAD
func (l *MembershipLoaderById) LoadThunk(key persist.DBID) func() (coredb.Membership, error) {
	l.mu.Lock()
	if it, ok := l.cache[key]; ok {
		l.mu.Unlock()
		return func() (coredb.Membership, error) {
=======
func (l *MembershipLoaderById) LoadThunk(key persist.DBID) func() (coregen.Membership, error) {
	l.mu.Lock()
	if it, ok := l.cache[key]; ok {
		l.mu.Unlock()
		return func() (coregen.Membership, error) {
>>>>>>> 93a3a41 (Add indexer models)
=======
func (l *MembershipLoaderById) LoadThunk(key persist.DBID) func() (coregen.Membership, error) {
	l.mu.Lock()
	if it, ok := l.cache[key]; ok {
		l.mu.Unlock()
		return func() (coregen.Membership, error) {
>>>>>>> a4e9c3f (Add indexer models)
			return it, nil
		}
	}
	if l.batch == nil {
		l.batch = &membershipLoaderByIdBatch{done: make(chan struct{})}
	}
	batch := l.batch
	pos := batch.keyIndex(l, key)
	l.mu.Unlock()

<<<<<<< HEAD
<<<<<<< HEAD
	return func() (coredb.Membership, error) {
		<-batch.done

		var data coredb.Membership
=======
=======
>>>>>>> a4e9c3f (Add indexer models)
	return func() (coregen.Membership, error) {
		<-batch.done

		var data coregen.Membership
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
func (l *MembershipLoaderById) LoadAll(keys []persist.DBID) ([]coredb.Membership, []error) {
	results := make([]func() (coredb.Membership, error), len(keys))
=======
func (l *MembershipLoaderById) LoadAll(keys []persist.DBID) ([]coregen.Membership, []error) {
	results := make([]func() (coregen.Membership, error), len(keys))
>>>>>>> 93a3a41 (Add indexer models)
=======
func (l *MembershipLoaderById) LoadAll(keys []persist.DBID) ([]coregen.Membership, []error) {
	results := make([]func() (coregen.Membership, error), len(keys))
>>>>>>> a4e9c3f (Add indexer models)

	for i, key := range keys {
		results[i] = l.LoadThunk(key)
	}

<<<<<<< HEAD
<<<<<<< HEAD
	memberships := make([]coredb.Membership, len(keys))
=======
	memberships := make([]coregen.Membership, len(keys))
>>>>>>> 93a3a41 (Add indexer models)
=======
	memberships := make([]coregen.Membership, len(keys))
>>>>>>> a4e9c3f (Add indexer models)
	errors := make([]error, len(keys))
	for i, thunk := range results {
		memberships[i], errors[i] = thunk()
	}
	return memberships, errors
}

// LoadAllThunk returns a function that when called will block waiting for a Memberships.
// This method should be used if you want one goroutine to make requests to many
// different data loaders without blocking until the thunk is called.
<<<<<<< HEAD
<<<<<<< HEAD
func (l *MembershipLoaderById) LoadAllThunk(keys []persist.DBID) func() ([]coredb.Membership, []error) {
	results := make([]func() (coredb.Membership, error), len(keys))
	for i, key := range keys {
		results[i] = l.LoadThunk(key)
	}
	return func() ([]coredb.Membership, []error) {
		memberships := make([]coredb.Membership, len(keys))
=======
func (l *MembershipLoaderById) LoadAllThunk(keys []persist.DBID) func() ([]coregen.Membership, []error) {
	results := make([]func() (coregen.Membership, error), len(keys))
	for i, key := range keys {
		results[i] = l.LoadThunk(key)
	}
	return func() ([]coregen.Membership, []error) {
		memberships := make([]coregen.Membership, len(keys))
>>>>>>> 93a3a41 (Add indexer models)
=======
func (l *MembershipLoaderById) LoadAllThunk(keys []persist.DBID) func() ([]coregen.Membership, []error) {
	results := make([]func() (coregen.Membership, error), len(keys))
	for i, key := range keys {
		results[i] = l.LoadThunk(key)
	}
	return func() ([]coregen.Membership, []error) {
		memberships := make([]coregen.Membership, len(keys))
>>>>>>> a4e9c3f (Add indexer models)
		errors := make([]error, len(keys))
		for i, thunk := range results {
			memberships[i], errors[i] = thunk()
		}
		return memberships, errors
	}
}

// Prime the cache with the provided key and value. If the key already exists, no change is made
// and false is returned.
// (To forcefully prime the cache, clear the key first with loader.clear(key).prime(key, value).)
<<<<<<< HEAD
<<<<<<< HEAD
func (l *MembershipLoaderById) Prime(key persist.DBID, value coredb.Membership) bool {
=======
func (l *MembershipLoaderById) Prime(key persist.DBID, value coregen.Membership) bool {
>>>>>>> 93a3a41 (Add indexer models)
=======
func (l *MembershipLoaderById) Prime(key persist.DBID, value coregen.Membership) bool {
>>>>>>> a4e9c3f (Add indexer models)
	l.mu.Lock()
	var found bool
	if _, found = l.cache[key]; !found {
		l.unsafeSet(key, value)
	}
	l.mu.Unlock()
	return !found
}

// Clear the value at key from the cache, if it exists
func (l *MembershipLoaderById) Clear(key persist.DBID) {
	l.mu.Lock()
	delete(l.cache, key)
	l.mu.Unlock()
}

<<<<<<< HEAD
<<<<<<< HEAD
func (l *MembershipLoaderById) unsafeSet(key persist.DBID, value coredb.Membership) {
	if l.cache == nil {
		l.cache = map[persist.DBID]coredb.Membership{}
=======
func (l *MembershipLoaderById) unsafeSet(key persist.DBID, value coregen.Membership) {
	if l.cache == nil {
		l.cache = map[persist.DBID]coregen.Membership{}
>>>>>>> 93a3a41 (Add indexer models)
=======
func (l *MembershipLoaderById) unsafeSet(key persist.DBID, value coregen.Membership) {
	if l.cache == nil {
		l.cache = map[persist.DBID]coregen.Membership{}
>>>>>>> a4e9c3f (Add indexer models)
	}
	l.cache[key] = value
}

// keyIndex will return the location of the key in the batch, if its not found
// it will add the key to the batch
func (b *membershipLoaderByIdBatch) keyIndex(l *MembershipLoaderById, key persist.DBID) int {
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

func (b *membershipLoaderByIdBatch) startTimer(l *MembershipLoaderById) {
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

func (b *membershipLoaderByIdBatch) end(l *MembershipLoaderById) {
	b.data, b.error = l.fetch(b.keys)
	close(b.done)
}
