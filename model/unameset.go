package model

import (
	"encoding/json"
	"sync"
)

// NameSetHandler handls operations with user name set.
type NameSetHandler interface {
	IsInSet(name string) bool
	AddToSet(name string)
}

var nameSet *uNameSet
var onceNameSet sync.Once

// GetNameSet returns a uNameSet instance
var GetNameSet = func() NameSetHandler {
	onceNameSet.Do(func() {
		nameSet = &uNameSet{m: make(map[string]struct{})}
	})
	return nameSet
}

type uNameSet struct {
	sync.RWMutex
	m map[string]struct{}
}

func (ns *uNameSet) IsInSet(name string) bool {
	ns.RLock()
	defer ns.RUnlock()
	if _, ok := ns.m[name]; !ok {
		return false
	}
	return true
}

func (ns *uNameSet) AddToSet(name string) {
	ns.Lock()
	defer ns.Unlock()
	ns.m[name] = struct{}{}
}

// AddKafkaMsgToNameSet adds and entry to the user name set.
func AddKafkaMsgToNameSet(val []byte) error {
	acc := Account{}
	if err := json.Unmarshal(val, &acc); err != nil {
		return err

	}
	GetNameSet().AddToSet(acc.Name)
	return nil
}
