package model

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("User name set", func() {
	var (
		name string
		ns   *uNameSet
	)
	BeforeEach(func() {
		ns = &uNameSet{m: make(map[string]struct{})}

	})
	Context("when IsInSet called", func() {

		It("returns true when name is in set", func() {
			ns.m[name] = struct{}{}
			Expect(ns.IsInSet(name)).To(Equal(true))
		})
		It("returns false when name is not in set", func() {
			Expect(ns.IsInSet(name)).To(Equal(false))
		})
	})
	Context("when AddToSet is called", func() {
		It("adds a name to name set", func() {
			ns.AddToSet(name)
			_, ok := ns.m[name]
			Expect(ok).To(Equal(true))
		})
	})

})
