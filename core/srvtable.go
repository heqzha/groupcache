package core

import (
	"bytes"
	"encoding/gob"
	"hash/crc32"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type SrvTable map[string]bool

func (st *SrvTable) Add(t SrvTable) {
	for addr, v := range t {
		(*st)[addr] = v
	}
}

func (st *SrvTable) Clone() SrvTable {
	nst := make(SrvTable, len(*st))
	for addr, v := range *st {
		nst[addr] = v
	}
	return nst
}

func (st *SrvTable) String() string {
	addrs := []string{}
	for addr, _ := range *st {
		addrs = append(addrs, addr)
	}
	return strings.Join(addrs, ";")
}

type STM struct {
	myAddr string
	table  SrvTable
	clock  VClock
	mutex  *sync.RWMutex
}

func (s *STM) Init(myAddr string) {
	if s == nil {
		s = new(STM)
	}
	s.myAddr = myAddr
	s.table = SrvTable{}
	s.clock = VClock{}
	s.mutex = &sync.RWMutex{}
}

func (s *STM) Register(addrs ...string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, addr := range addrs {
		s.table[addr] = true
	}
	s.clock.Tick(s.myAddr)
}

func (s *STM) Unregister(addr string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, eixst := s.table[addr]
	if eixst {
		delete(s.table, addr)
	}
	s.clock.Tick(s.myAddr)
}

func (s *STM) Merge(t STM) Condition {
	s.mutex.Lock()
	defer func() {
		s.clock.Merge(t.clock)
		s.mutex.Unlock()
	}()
	if s.clock.Compare(t.clock, Equal) {
		//s and t are equal
		return Equal
	} else if s.clock.Compare(t.clock, Concurrent) {
		//s and t are concurrent
		s.table.Add(t.table)
		return Concurrent
	} else if s.clock.Compare(t.clock, Descendant) {
		//s is older than t
		s.table = t.table.Clone()
		return Descendant
	} else if s.clock.Compare(t.clock, Ancestor) {
		//s is newer than t
		return Ancestor
	}
	return Condition(-1)
}

func (s *STM) CompareReadable(t STM) string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.clock.Compare(t.clock, Equal) {
		return "Equal"
	} else if s.clock.Compare(t.clock, Concurrent) {
		return "Concurrent"
	} else if s.clock.Compare(t.clock, Descendant) {
		return "Descendant"
	} else if s.clock.Compare(t.clock, Ancestor) {
		return "Ancestor"
	}
	return "Unknown"
}

func (s *STM) GetTable() SrvTable {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.table
}

func (s *STM) GetClock() VClock {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.clock
}

func (s *STM) Dump() ([]byte, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	b := new(bytes.Buffer)
	enc := gob.NewEncoder(b)
	if err := enc.Encode(s.myAddr); err != nil {
		return nil, err
	}
	if err := enc.Encode(s.table); err != nil {
		return nil, err
	}
	if err := enc.Encode(s.clock); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (s *STM) Load(buf []byte) error {
	r := bytes.NewBuffer(buf)
	dec := gob.NewDecoder(r)
	if err := dec.Decode(&s.myAddr); err != nil {
		return err
	}
	if err := dec.Decode(&s.table); err != nil {
		return err
	}
	if err := dec.Decode(&s.clock); err != nil {
		return err
	}
	return nil
}

type HashFn func(data []byte) uint32

type STHash struct {
	repTable map[int]string
	replicas int
	hashes   []int
	hashfn   HashFn
}

func (s *STHash) Init(replicas int, fn HashFn) {
	if s == nil {
		s = new(STHash)
	}
	s.repTable = map[int]string{}
	s.replicas = replicas
	s.hashes = []int{}
	s.hashfn = fn
	if s.hashfn == nil {
		s.hashfn = crc32.ChecksumIEEE
	}
}

func (s *STHash) Load(t SrvTable) {
	for addr, _ := range t {
		for i := 0; i < s.replicas; i++ {
			hash := int(s.hashfn([]byte(strconv.Itoa(i) + addr)))
			s.hashes = append(s.hashes, hash)
			s.repTable[hash] = addr
		}
	}
	sort.Ints(s.hashes)
}

func (s *STHash) Pick(key string) string {
	if len(s.hashes) == 0 {
		return ""
	}
	hash := int(s.hashfn([]byte(key)))

	idx := sort.Search(len(s.hashes), func(i int) bool { return s.hashes[i] >= hash })

	if idx == len(s.hashes) {
		idx = 0
	}

	return s.repTable[s.hashes[idx]]
}
