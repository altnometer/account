package account

import "sync"

// Account holds core user details.
type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Pwd  string `json:"pwd"`
}

// UNameSet stores unique usernames.
type UNameSet struct {
	sync.RWMutex
	M map[string]struct{}
}

// IsInSet checks if a username is in UNameSet.
func (ns *UNameSet) IsInSet(name string) bool {
	ns.RLock()
	defer ns.RUnlock()
	if _, ok := ns.M[name]; !ok {
		return false
	}
	return true
}

// AddToSet  adds a username to UNameSet.
func (ns *UNameSet) AddToSet(name string) {
	ns.Lock()
	defer ns.Unlock()
	ns.M[name] = struct{}{}
}

// UNames caches unique usernames. Used to check duplicate names.
var UNames = UNameSet{M: make(map[string]struct{})}

// NameExists  incapsulates duplicate name checking api.
var NameExists = func(name string) bool {
	return UNames.IsInSet(name)
}
