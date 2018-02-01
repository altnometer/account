package model

import "sync"

type uNameSet struct {
	sync.RWMutex
	m map[string]struct{}
}

func (ns *uNameSet) isInSet(name string) bool {
	ns.RLock()
	defer ns.RUnlock()
	if _, ok := ns.m[name]; !ok {
		return false
	}
	return true
}

func (ns *uNameSet) addToSet(name string) {
	ns.Lock()
	defer ns.Unlock()
	ns.m[name] = struct{}{}
}

var uNames = uNameSet{m: make(map[string]struct{})}

// UNameExists  incapsulates duplicate name checking api.
var UNameExists = func(name string) bool {
	return uNames.isInSet(name)
}
