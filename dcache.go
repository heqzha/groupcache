package dcache

import (
	"errors"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/heqzha/dcache/caches"
	pb "github.com/heqzha/dcache/dcachepb"
	"github.com/heqzha/dcache/singleflight"
)

type Handler interface {
	Get(ctx Context, key string, dest Sink) error
	Set(ctx Context, key string, dest Sink) error
	Del(ctx Context, key string, dest Sink) error
}

type HandlerFuncs struct {
	Getter func(ctx Context, key string, dest Sink) error
	Setter func(ctx Context, key string, dest Sink) error
	Deller func(ctx Context, key string, dest Sink) error
}

func (f HandlerFuncs) Get(ctx Context, key string, dest Sink) error {
	return f.Getter(ctx, key, dest)
}

func (f HandlerFuncs) Set(ctx Context, key string, dest Sink) error {
	return f.Setter(ctx, key, dest)
}

func (f HandlerFuncs) Del(ctx Context, key string, dest Sink) error {
	return f.Deller(ctx, key, dest)
}

// CacheStats are returned by stats accessors on Group.
type CacheStats struct {
	Bytes     int64
	Items     int64
	Gets      int64
	Sets      int64
	Dels      int64
	Hits      int64
	Evictions int64
}

// cache is a wrapper around an *caches.Cache that adds synchronization,
// makes values always be ByteView, and counts the size of all keys and
// values.
type cache struct {
	mu                     sync.RWMutex
	nbytes                 int64 // of all keys and values
	ca                     *caches.LRUCache
	nhit, nget, nset, ndel int64
	nevict                 int64 // number of evictions
}

func (c *cache) stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return CacheStats{
		Bytes:     c.nbytes,
		Items:     c.itemsLocked(),
		Gets:      c.nget,
		Sets:      c.nset,
		Dels:      c.ndel,
		Hits:      c.nhit,
		Evictions: c.nevict,
	}
}

func (c *cache) set(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.nset++
	if c.ca == nil {
		c.ca = &caches.LRUCache{
			OnEvicted: func(key caches.Key, value interface{}) {
				val := value.(ByteView)
				c.nbytes -= int64(len(key.(string))) + int64(val.Len())
				c.nevict++
			},
		}
	}
	c.ca.Add(key, value)
	c.nbytes += int64(len(key)) + int64(value.Len())
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.nget++
	if c.ca == nil {
		return
	}
	vi, ok := c.ca.Get(key)
	if !ok {
		return
	}
	c.nhit++
	return vi.(ByteView), true
}

func (c *cache) del(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ndel++
	if c.ca != nil {
		c.ca.Remove(key)
	}
}

func (c *cache) removeOldest() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ndel++
	if c.ca != nil {
		c.ca.RemoveOldest()
	}
}

func (c *cache) bytes() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.nbytes
}

func (c *cache) items() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.itemsLocked()
}

func (c *cache) itemsLocked() int64 {
	if c.ca == nil {
		return 0
	}
	return int64(c.ca.Len())
}

// A Group is a cache namespace and associated data loaded spread over
// a group of 1 or more machines.
type Group struct {
	name       string
	handler    Handler
	peersOnce  sync.Once
	peers      PeerPicker
	cacheBytes int64 // limit for sum of mainCache and hotCache size

	// mainCache is a cache of the keys for which this process
	// (amongst its peers) is authoritative. That is, this cache
	// contains keys which consistent hash on to this process's
	// peer number.
	mainCache cache

	// hotCache contains keys/values for which this peer is not
	// authoritative (otherwise they would be in mainCache), but
	// are popular enough to warrant mirroring in this process to
	// avoid going over the network to fetch from a peer.  Having
	// a hotCache avoids network hotspotting, where a peer's
	// network card could become the bottleneck on a popular key.
	// This cache is used sparingly to maximize the total number
	// of key/value pairs that can be stored globally.
	hotCache cache

	// loadGroup ensures that each key is only fetched once
	// (either locally or remotely), regardless of the number of
	// concurrent callers.
	loadGroup flightGroup

	_ int32 // force Stats to be 8-byte aligned on 32-bit platforms

	// Stats are statistics on the group.
	Stats Stats
}

// Name returns the name of the group.
func (g *Group) Name() string {
	return g.name
}

func (g *Group) initPeers() {
	if g.peers == nil {
		g.peers = getPeers(g.name)
	}
}

//
// ─── GET SECTION ────────────────────────────────────────────────────────────────
//
func (g *Group) Get(ctx Context, key string, dest Sink) error {
	g.peersOnce.Do(g.initPeers)
	g.Stats.Gets.Add(1)
	if dest == nil {
		return errors.New("groupcache: nil dest Sink")
	}
	value, cacheHit := g.lookupCache(key)
	if cacheHit {
		g.Stats.CacheHits.Add(1)
		return setSinkView(dest, value)
	}
	// Optimization to avoid double unmarshalling or copying: keep
	// track of whether the dest was already populated. One caller
	// (if local) will set this; the losers will not. The common
	// case will likely be one caller.
	destPopulated := false
	value, destPopulated, err := g.load(ctx, key, dest)
	if err != nil {
		return err
	}
	if destPopulated {
		return nil
	}
	return setSinkView(dest, value)
}

// load loads key either by invoking the getter locally or by sending it to another machine.
func (g *Group) load(ctx Context, key string, dest Sink) (value ByteView, destPopulated bool, err error) {
	g.Stats.Loads.Add(1)
	viewi, err := g.loadGroup.Do(key, func() (interface{}, error) {
		// Check the cache again because singleflight can only dedup calls
		// that overlap concurrently.  It's possible for 2 concurrent
		// requests to miss the cache, resulting in 2 load() calls.  An
		// unfortunate goroutine scheduling would result in this callback
		// being run twice, serially.  If we don't check the cache again,
		// cache.nbytes would be incremented below even though there will
		// be only one entry for this key.
		//
		// Consider the following serialized event ordering for two
		// goroutines in which this callback gets called twice for hte
		// same key:
		// 1: Get("key")
		// 2: Get("key")
		// 1: lookupCache("key")
		// 2: lookupCache("key")
		// 1: load("key")
		// 2: load("key")
		// 1: loadGroup.Do("key", fn)
		// 1: fn()
		// 2: loadGroup.Do("key", fn)
		// 2: fn()
		if value, cacheHit := g.lookupCache(key); cacheHit {
			g.Stats.CacheHits.Add(1)
			return value, nil
		}
		g.Stats.LoadsDeduped.Add(1)
		var value ByteView
		var err error
		if peer, ok := g.peers.PickPeer(key); ok {
			value, err = g.getFromPeer(ctx, peer, key)
			if err == nil {
				g.Stats.PeerLoads.Add(1)
				return value, nil
			}
			g.Stats.PeerErrors.Add(1)
			// TODO(bradfitz): log the peer's error? keep
			// log of the past few for /groupcachez?  It's
			// probably boring (normal task movement), so not
			// worth logging I imagine.
		}
		value, err = g.getLocally(ctx, key, dest)
		if err != nil {
			g.Stats.LocalLoadErrs.Add(1)
			return nil, err
		}
		g.Stats.LocalLoads.Add(1)
		destPopulated = true // only one caller of load gets this return value
		g.populateCache(key, value, &g.mainCache)
		return value, nil
	})
	if err == nil {
		value = viewi.(ByteView)
	}
	return
}

func (g *Group) getLocally(ctx Context, key string, dest Sink) (ByteView, error) {
	err := g.handler.Get(ctx, key, dest)
	if err != nil {
		return ByteView{}, err
	}
	return dest.view()
}

func (g *Group) getFromPeer(ctx Context, peer ProtoHandler, key string) (ByteView, error) {
	req := &pb.GetRequest{
		Group: g.name,
		Key:   key,
	}
	res := &pb.GetResponse{}
	err := peer.Get(ctx, req, res)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: res.Value}
	// TODO(bradfitz): use res.MinuteQps or something smart to
	// conditionally populate hotCache.  For now just do it some
	// percentage of the time.
	if rand.Intn(10) == 0 {
		g.populateCache(key, value, &g.hotCache)
	}
	return value, nil
}

func (g *Group) lookupCache(key string) (value ByteView, ok bool) {
	if g.cacheBytes <= 0 {
		return
	}
	value, ok = g.mainCache.get(key)
	if ok {
		return
	}
	value, ok = g.hotCache.get(key)
	return
}

func (g *Group) populateCache(key string, value ByteView, cache *cache) {
	if g.cacheBytes <= 0 {
		return
	}
	cache.set(key, value)

	// Evict items from cache(s) if necessary.
	for {
		mainBytes := g.mainCache.bytes()
		hotBytes := g.hotCache.bytes()
		if mainBytes+hotBytes <= g.cacheBytes {
			return
		}

		// TODO(bradfitz): this is good-enough-for-now logic.
		// It should be something based on measurements and/or
		// respecting the costs of different resources.
		victim := &g.mainCache
		if hotBytes > mainBytes/8 {
			victim = &g.hotCache
		}
		victim.removeOldest()
	}
}

//
// ─── SET SECTION ────────────────────────────────────────────────────────────────
//
func (g *Group) Set(ctx Context, key string, dest Sink) error {
	g.peersOnce.Do(g.initPeers)
	g.Stats.Sets.Add(1)
	if dest == nil {
		return errors.New("groupcache: nil dest Sink")
	}
	value, destPopulated, err := g.save(ctx, key, dest)
	if err != nil {
		return err
	}
	if destPopulated {
		return nil
	}
	return setSinkView(dest, value)
}

func (g *Group) setLocally(ctx Context, key string, dest Sink) error {
	err := g.handler.Set(ctx, key, dest)
	if err != nil {
		return err
	}
	return nil
}

func (g *Group) save(ctx Context, key string, dest Sink) (value ByteView, destPopulated bool, err error) {
	g.Stats.Saves.Add(1)
	viewi, err := g.loadGroup.Do(key, func() (interface{}, error) {
		g.Stats.SavesDeduped.Add(1)
		var err error
		if peer, ok := g.peers.PickPeer(key); ok {
			value, err := dest.view()
			if err != nil {
				return nil, err
			}
			_, err = g.setToPeer(ctx, peer, key, value)
			if err == nil {
				g.Stats.PeerSets.Add(1)
				return value, nil
			}
			g.Stats.PeerErrors.Add(1)
			// TODO(bradfitz): log the peer's error? keep
			// log of the past few for /groupcachez?  It's
			// probably boring (normal task movement), so not
			// worth logging I imagine.
		}

		if err = g.setLocally(ctx, key, dest); err != nil {
			g.Stats.LocalSetErrs.Add(1)
			return nil, err
		}
		g.Stats.LocalSets.Add(1)
		destPopulated = true
		value, err := dest.view()
		if err != nil {
			return nil, err
		}
		g.populateCache(key, value, &g.mainCache)
		return value, nil
	})
	if err == nil {
		value = viewi.(ByteView)
	}
	return
}
func (g *Group) setToPeer(ctx Context, peer ProtoHandler, key string, val ByteView) (bool, error) {
	req := &pb.SetRequest{
		Group: g.name,
		Key:   key,
		Value: val.ByteSlice(),
	}
	res := &pb.SetResponse{}
	err := peer.Set(ctx, req, res)
	if err != nil {
		return false, err
	}
	if res.Status {
		// TODO(bradfitz): use res.MinuteQps or something smart to
		// conditionally populate hotCache.  For now just do it some
		// percentage of the time.
		if rand.Intn(10) == 0 {
			g.populateCache(key, val, &g.hotCache)
		}
	}
	return res.Status, nil
}

//
// ─── DEL SECTION ────────────────────────────────────────────────────────────────
//
func (g *Group) Del(ctx Context, key string, dest Sink) error {
	g.peersOnce.Do(g.initPeers)
	g.Stats.Sets.Add(1)
	if dest == nil {
		return errors.New("groupcache: nil dest Sink")
	}
	value, destPopulated, err := g.remove(ctx, key, dest)
	if err != nil {
		return err
	}
	if destPopulated {
		return nil
	}
	return setSinkView(dest, value)
}

func (g *Group) delLocally(ctx Context, key string, dest Sink) error {
	err := g.handler.Del(ctx, key, dest)
	if err != nil {
		return err
	}
	return nil
}

func (g *Group) remove(ctx Context, key string, dest Sink) (value ByteView, destPopulated bool, err error) {
	g.Stats.Removes.Add(1)
	viewi, err := g.loadGroup.Do(key, func() (interface{}, error) {
		g.Stats.RemovesDeduped.Add(1)
		var err error
		if peer, ok := g.peers.PickPeer(key); ok {
			value, err := dest.view()
			if err != nil {
				return nil, err
			}
			_, err = g.delToPeer(ctx, peer, key, value)
			if err == nil {
				g.Stats.PeerSets.Add(1)
				return value, nil
			}
			g.Stats.PeerErrors.Add(1)
			// TODO(bradfitz): log the peer's error? keep
			// log of the past few for /groupcachez?  It's
			// probably boring (normal task movement), so not
			// worth logging I imagine.
		}
		if err = g.delLocally(ctx, key, dest); err != nil {
			g.Stats.LocalRemoveErrs.Add(1)
			return nil, err
		}
		g.Stats.LocalRemoves.Add(1)
		destPopulated = true
		g.mainCache.del(key)
		g.hotCache.del(key)
		value, err := dest.view()
		if err != nil {
			return nil, err
		}
		return value, nil
	})
	if err == nil {
		value = viewi.(ByteView)
	}
	return
}

func (g *Group) delToPeer(ctx Context, peer ProtoHandler, key string, val ByteView) (bool, error) {
	req := &pb.DelRequest{
		Group: g.name,
		Key:   key,
	}
	res := &pb.DelResponse{}
	err := peer.Del(ctx, req, res)
	if err != nil {
		return false, err
	}
	return res.Status, nil
}

// flightGroup is defined as an interface which flightgroup.Group
// satisfies.  We define this so that we may test with an alternate
// implementation.
type flightGroup interface {
	// Done is called when Do is done.
	Do(key string, fn func() (interface{}, error)) (interface{}, error)
}

// An AtomicInt is an int64 to be accessed atomically.
type AtomicInt int64

// Add atomically adds n to i.
func (i *AtomicInt) Add(n int64) {
	atomic.AddInt64((*int64)(i), n)
}

// Get atomically gets the value of i.
func (i *AtomicInt) Get() int64 {
	return atomic.LoadInt64((*int64)(i))
}

func (i *AtomicInt) String() string {
	return strconv.FormatInt(i.Get(), 10)
}

// Stats are per-group statistics.
type Stats struct {
	Gets            AtomicInt // any Get request, including from peers
	Sets            AtomicInt
	Dels            AtomicInt
	CacheHits       AtomicInt // either cache was good
	PeerSets        AtomicInt
	PeerLoads       AtomicInt // either remote load or remote cache hit (not an error)
	PeerErrors      AtomicInt
	Removes         AtomicInt
	RemovesDeduped  AtomicInt // after singleflight
	Saves           AtomicInt
	SavesDeduped    AtomicInt // after singleflight
	Loads           AtomicInt // (gets - cacheHits)
	LoadsDeduped    AtomicInt // after singleflight
	LocalLoads      AtomicInt // total good local loads
	LocalLoadErrs   AtomicInt // total bad local loads
	LocalSets       AtomicInt // total good local sets
	LocalSetErrs    AtomicInt // total bad local sets
	LocalRemoves    AtomicInt // total good local removes
	LocalRemoveErrs AtomicInt // total bad local removes
	ServerRequests  AtomicInt // gets that came over the network from peers
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)

	initPeerServerOnce sync.Once
	initPeerServer     func()
)

// GetGroup returns the named group previously created with NewGroup, or
// nil if there's no such group.
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

func NewGroup(name string, cacheBytes int64, handler Handler) *Group {
	return newGroup(name, cacheBytes, handler, nil)
}

// If peers is nil, the peerPicker is called via a sync.Once to initialize it.
func newGroup(name string, cacheBytes int64, handler Handler, peers PeerPicker) *Group {
	if handler == nil {
		panic("nil Handler")
	}
	mu.Lock()
	defer mu.Unlock()
	initPeerServerOnce.Do(callInitPeerServer)
	if _, dup := groups[name]; dup {
		panic("duplicate registration of group " + name)
	}
	g := &Group{
		name:       name,
		handler:    handler,
		peers:      peers,
		cacheBytes: cacheBytes,
		loadGroup:  &singleflight.Group{},
	}
	if fn := newGroupHook; fn != nil {
		fn(g)
	}
	groups[name] = g
	return g
}

// newGroupHook, if non-nil, is called right after a new group is created.
var newGroupHook func(*Group)

// RegisterNewGroupHook registers a hook that is run each time
// a group is created.
func RegisterNewGroupHook(fn func(*Group)) {
	if newGroupHook != nil {
		panic("RegisterNewGroupHook called more than once")
	}
	newGroupHook = fn
}

// RegisterServerStart registers a hook that is run when the first
// group is created.
func RegisterServerStart(fn func()) {
	if initPeerServer != nil {
		panic("RegisterServerStart called more than once")
	}
	initPeerServer = fn
}

func callInitPeerServer() {
	if initPeerServer != nil {
		initPeerServer()
	}
}
