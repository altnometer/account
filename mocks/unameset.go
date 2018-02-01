package mocks

// UNameSet mocks model.uNameSet stuct of model.NameSetHandler interface.
type UNameSet struct {
	IsInSetCall struct {
		Receives struct {
			Name string
		}
		Returns struct {
			Bool bool
		}
	}
	AddToSetCall struct {
		Receives struct {
			Name string
		}
	}
}

// IsInSet mocks a uNameSet method.
func (ns *UNameSet) IsInSet(name string) bool {
	ns.IsInSetCall.Receives.Name = name
	return ns.IsInSetCall.Returns.Bool
}

// AddToSet mocks a uNameSet method.
func (ns *UNameSet) AddToSet(name string) {
	ns.AddToSetCall.Receives.Name = name
}
