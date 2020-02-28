package set

import (
	"sync"

	"encoding/json"
)

type void struct{}
var empty void

func getKeys(mp map[interface{}]void) []interface{} {
	keys := make([]interface{}, 0, len(mp))
	for k := range mp {
		keys = append(keys, k)
	}
	return keys
}

func mapFromKeys(ks []interface{}) map[interface{}]void {
	mp := map[interface{}]void{}
	for _, k := range ks {
		addToMap(mp, k)
	}
	return mp
}

func addToMap(mp map[interface{}]void, v interface{}) {
	switch v.(type) {
	case int8:
		mp[int64(v.(int8))] = empty
	case int:
		mp[int64(v.(int))] = empty
	case int32:
		mp[int64(v.(int32))] = empty
	case float64:
		mp[int64(v.(float64))] = empty
	default:
		mp[v] = empty
	}
}

type Set interface {
	Merge(v Set)

	Add(v ...interface{})

	Remove(v ...interface{})

	Members()[]interface{}

	Exists(v interface{}) bool
}

type set struct {
	mu sync.Mutex
	vals map[interface{}]void
}

func New() Set {
	return &set{
		vals: map[interface{}]void{},
	}
}

func (s *set) Merge(v Set) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var vals map[interface{}]void
	d, ok := v.(*set)
	if ok {
		vals = d.vals
	}
	for k := range vals {
		addToMap(s.vals, k)
	}
}

func (s *set) Members()[]interface{} {
	return getKeys(s.vals)
}

func (s *set) Add(v ...interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, val := range v {
		_, ok := s.vals[val]
		if !ok {
			addToMap(s.vals, val)
		}
	}
}

func (s *set) Remove(v ...interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, val := range v {
		delete(s.vals, val)
	}
}

func (s set) Exists(v interface{}) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.vals[v]
	return exists
}

func (s *set) UnmarshalJSON(b []byte) error {
	var sarr []interface{}

	if err := json.Unmarshal(b, &sarr); err != nil {
		return err
	}

	mp := mapFromKeys(sarr)
	*s = set{vals: mp}

	return nil
}

func (s set) MarshalJSON() ([]byte, error) {
	return json.Marshal(getKeys(s.vals))
}
