package set

import (
	"sync"

	"encoding/json"
)

func getStringKeys(mp map[string]void) []string {
	keys := make([]string, 0, len(mp))
	for k := range mp {
		keys = append(keys, k)
	}
	return keys
}

func stringMapFromKeys(ks []string) map[string]void {
	mp := map[string]void{}
	for _, k := range ks {
		mp[k] = empty
	}
	return mp
}

//type StringSet interface {
//	Merge(v StringSet)
//
//	Add(v ...string)
//
//	Remove(v ...string)
//
//	Members() []string
//
//	Exists(v string) bool
//}

type StringSet struct {
	mu   sync.Mutex
	vals map[string]void
}

func NewStringSet() *StringSet {
	return &StringSet{
		vals: map[string]void{},
	}
}

func (s *StringSet) Merge(v StringSet) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k := range v.vals {
		s.vals[k] = empty
	}
}

func (s *StringSet) Members() []string {
	return getStringKeys(s.vals)
}

func (s *StringSet) Add(v ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, val := range v {
		_, ok := s.vals[val]
		if !ok {
			s.vals[val] = empty
		}
	}
}

func (s *StringSet) Remove(v ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, val := range v {
		delete(s.vals, val)
	}
}

func (s StringSet) Exists(v string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.vals[v]
	return exists
}

func (s *StringSet) UnmarshalJSON(b []byte) error {
	var sarr []string

	if err := json.Unmarshal(b, &sarr); err != nil {
		return err
	}

	mp := stringMapFromKeys(sarr)
	s.vals = mp
	//*s = StringSet{vals: mp}

	return nil
}

func (s StringSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(getStringKeys(s.vals))
}
