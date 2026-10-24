// Package stash provides caches for storing and retrieving values by key.
//
// [ExpiringCache] supports concurrent use and expires values after a configurable
// time-to-live. Its zero value is ready to use with a one-minute time-to-live.
package stash
