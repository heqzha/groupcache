package dcache

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/heqzha/dcache/testpb"
)

var (
	once                    sync.Once
	stringGroup, protoGroup Handler
	stringc                 = make(chan string)
	setstringc              = make(chan string)
	delstringc              = make(chan string)
	dummyCtx                Context
	// cacheFills is the number of times stringGroup or
	// protoGroup's Getter have been called. Read using the
	// cacheFills function.
	cacheFills AtomicInt
	cacheSets  AtomicInt
	cacheDels  AtomicInt
)

const (
	stringGroupName = "string-group"
	protoGroupName  = "proto-group"
	testMessageType = "google3/net/groupcache/go/test_proto.TestMessage"
	fromChan        = "from-chan"
	toChan          = "to-chan"
	cacheSize       = 1 << 20
)

func testSetup() {
	strHandler := HandlerFuncs{
		Getter: func(ctx Context, key string, dest Sink) error {
			if key == fromChan {
				key = <-stringc
			}
			cacheFills.Add(1)
			return dest.SetString("ECHO:" + key)
		},
		Setter: func(ctx Context, key string, dest Sink) error {
			if key == toChan {
				s, err := dest.GetString()
				if err == nil {
					setstringc <- s
				}
			}
			cacheSets.Add(1)
			return nil
		},
		Deller: func(ctx Context, key string, dest Sink) error {
			if key == toChan {
				s, err := dest.GetString()
				if err == nil {
					delstringc <- s
				}
			}
			cacheDels.Add(1)
			return nil
		},
	}

	stringGroup = NewGroup(stringGroupName, cacheSize, strHandler)

	protoHandler := HandlerFuncs{
		Getter: func(ctx Context, key string, dest Sink) error {
			if key == fromChan {
				key = <-stringc
			}
			cacheFills.Add(1)
			return dest.SetProto(&testpb.TestMessage{
				Name: proto.String("ECHO:" + key),
				City: proto.String("SOME-CITY"),
			})
		},
		Setter: func(ctx Context, key string, dest Sink) error {
			if key == fromChan {
				key = <-stringc
			}
			cacheSets.Add(1)
			return dest.SetProto(&testpb.TestMessage{
				Name: proto.String("ECHO:" + key),
				City: proto.String("SOME-CITY"),
			})
		},
		Deller: func(ctx Context, key string, dest Sink) error {
			if key == fromChan {
				key = <-stringc
			}
			cacheDels.Add(1)
			return dest.SetProto(&testpb.TestMessage{
				Name: proto.String("ECHO:" + key),
				City: proto.String("SOME-CITY"),
			})
		},
	}

	protoGroup = NewGroup(protoGroupName, cacheSize, protoHandler)
}

// tests that a Getter's Get method is only called once with two
// outstanding callers.  This is the string variant.
func TestGetDupSuppressString(t *testing.T) {
	once.Do(testSetup)
	// Start two getters. The first should block (waiting reading
	// from stringc) and the second should latch on to the first
	// one.
	resc := make(chan string, 2)
	for i := 0; i < 2; i++ {
		go func() {
			var s string
			if err := stringGroup.Get(dummyCtx, fromChan, StringSink(&s)); err != nil {
				resc <- "ERROR:" + err.Error()
				return
			}
			resc <- s
		}()
	}

	// Wait a bit so both goroutines get merged together via
	// singleflight.
	// TODO(bradfitz): decide whether there are any non-offensive
	// debug/test hooks that could be added to singleflight to
	// make a sleep here unnecessary.
	time.Sleep(250 * time.Millisecond)

	// Unblock the first getter, which should unblock the second
	// as well.
	stringc <- "foo"

	for i := 0; i < 2; i++ {
		select {
		case v := <-resc:
			if v != "ECHO:foo" {
				t.Errorf("got %q; want %q", v, "ECHO:foo")
			}
			t.Logf("got %q!!", v)
		case <-time.After(5 * time.Second):
			t.Errorf("timeout waiting on getter #%d of 2", i+1)
		}
	}
}

func TestSetString(t *testing.T) {
	once.Do(testSetup)
	resc := make(chan string, 2)
	for i := 0; i < 2; i++ {
		go func() {
			var s string
			if err := stringGroup.Get(dummyCtx, fromChan, StringSink(&s)); err != nil {
				resc <- "ERROR:" + err.Error()
				return
			}
			resc <- s
		}()
	}

	for i := 1; i < 3; i++ {
		go func() {
			var s = strconv.Itoa(i)
			if err := stringGroup.Set(dummyCtx, toChan, StringSink(&s)); err != nil {
				resc <- "ERROR:" + err.Error()
				return
			}
			resc <- s
		}()
	}
	//TODO handle resc
}
