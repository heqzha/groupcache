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

type SrvGroup map[string]*SrvTable

func (s *SrvGroup) NewGroup(group string) *SrvTable {
	if group == "" {
		panic("missing group name")
	}
	(*s)[group] = &SrvTable{}
	return (*s)[group]
}

func (s *SrvGroup) GetTable(group string) *SrvTable {
	if group == "" {
		panic("missing group name")
	}
	table, ok := (*s)[group]
	if ok && table != nil {
		return table
	}
	return s.NewGroup(group)
}

func (s *SrvGroup) SetTable(group string, table *SrvTable) {
	if group == "" {
		panic("missing group name")
	}
	(*s)[group] = table
}

type SGM struct {
	myAddr string
	group  SrvGroup
	clock  VClock
	mutex  *sync.RWMutex
}

func (s *SGM) Init(myAddr string) {
	if s == nil {
		s = new(SGM)
	}
	s.myAddr = myAddr
	s.group = SrvGroup{}
	s.clock = VClock{}
	s.mutex = &sync.RWMutex{}
}

func (s *SGM) Register(group string, addrs ...string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	table := s.group.GetTable(group)
	for _, addr := range addrs {
		(*table)[addr] = true
	}

	s.clock.Tick(s.myAddr)
}

func (s *SGM) Unregister(group string, addr string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	table := s.group.GetTable(group)
	_, eixst := (*table)[addr]
	if eixst {
		delete((*table), addr)
	}
	s.clock.Tick(s.myAddr)
}

func (s *SGM) GetGroups() []string {
	groups := []string{}
	for g := range s.group {
		groups = append(groups, g)
	}
	return groups
}

func (s *SGM) Merge(t SGM) Condition {
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
		for tg, tt := range t.group {
			st := s.group.GetTable(tg)
			st.Add((*tt))
		}
		return Concurrent
	} else if s.clock.Compare(t.clock, Descendant) {
		//s is older than t
		for tg, tt := range t.group {
			ttc := tt.Clone()
			s.group.SetTable(tg, &ttc)
		}
		return Descendant
	} else if s.clock.Compare(t.clock, Ancestor) {
		//s is newer than t
		return Ancestor
	}
	return Condition(-1)
}

func (s *SGM) CompareReadable(t SGM) string {
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

func (s *SGM) GetTable(group string) *SrvTable {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.group.GetTable(group)
}

func (s *SGM) GetClock() VClock {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.clock
}

func (s *SGM) Dump() ([]byte, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	b := new(bytes.Buffer)
	enc := gob.NewEncoder(b)
	if err := enc.Encode(s.group); err != nil {
		return nil, err
	}
	if err := enc.Encode(s.clock); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (s *SGM) Load(buf []byte) error {
	r := bytes.NewBuffer(buf)
	dec := gob.NewDecoder(r)
	if err := dec.Decode(&s.group); err != nil {
		return err
	}
	if err := dec.Decode(&s.clock); err != nil {
		return err
	}
	return nil
}

type HashFn func(data []byte) uint32

type RepTable map[int]string

type SGHash struct {
	group    map[string]*RepTable
	replicas int
	hashes   map[string][]int
	hashfn   HashFn
	mutex    *sync.RWMutex
}

func (s *SGHash) Init(replicas int, fn HashFn) {
	if s == nil {
		s = new(SGHash)
	}
	s.group = map[string]*RepTable{}
	s.replicas = replicas
	s.hashes = map[string][]int{}
	s.hashfn = fn
	s.mutex = &sync.RWMutex{}
	if s.hashfn == nil {
		s.hashfn = crc32.ChecksumIEEE
	}
}

func (s *SGHash) Load(t SrvGroup) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for tg, tt := range t {
		for addr := range *tt {
			for i := 0; i < s.replicas; i++ {
				hash := int(s.hashfn([]byte(strconv.Itoa(i) + addr)))
				if s.hashes[tg] == nil {
					s.hashes[tg] = make([]int, len(*tt)*s.replicas)
				}
				s.hashes[tg] = append(s.hashes[tg], hash)
				if s.group[tg] == nil {
					s.group[tg] = new(RepTable)
				}
				(*s.group[tg])[hash] = addr
			}
		}
		sort.Ints(s.hashes[tg])
	}
}

func (s *SGHash) Pick(group, key string) string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if len(s.hashes) == 0 || len(s.hashes[group]) == 0 {
		return ""
	}
	hash := int(s.hashfn([]byte(key)))

	idx := sort.Search(len(s.hashes[group]), func(i int) bool { return s.hashes[group][i] >= hash })

	if idx == len(s.hashes[group]) {
		idx = 0
	}

	return (*s.group[group])[s.hashes[group][idx]]
}
