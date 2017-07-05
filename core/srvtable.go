package core

import (
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type SrvTable map[string]bool

func (st *SrvTable) Add(addr string) {

}

func (st *SrvTable) Merge(t *SrvTable) error {
	return nil
}

func (st *SrvTable) Encode() ([]byte, error) {
	return nil, nil
}

func (st *SrvTable) Decode(buf []byte) error {
	return nil
}

type STM struct {
	myAddr string
	table  SrvTable
	clock  *VClock
	mutex  *sync.RWMutex
}

func (s *STM) Init(myAddr string) {
	s = &STM{
		myAddr: myAddr,
		table:  SrvTable{},
		clock:  &VClock{},
		mutex:  &sync.RWMutex{},
	}
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

}

func (s *STM) Get(key string) string {
	return ""
}

func (s *STM) Merge(t *STM) error {
	return nil
}

func (s *STM) Dump() ([]byte, error) {
	return nil, nil
}

func (s *STM) Load(buf []byte) error {
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
	s = &STHash{
		repTable: map[int]string{},
		replicas: replicas,
		hashes:   []int{},
		hashfn:   fn,
	}
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
