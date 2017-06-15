package caches

type Key interface{}

type Cache interface {
	Add(key Key, value interface{})
	Get(key Key) (interface{}, bool)
	Remove(key Key) interface{}
	RemoveOldest() interface{}
	Len() int
	Clear()
}
