package model_test

import (
	"encoding/json"

	"github.com/altnometer/account/mocks"
	"github.com/altnometer/account/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AddKafkaMsgToNameSet", func() {
	var (
		id, name, pwd string
		acc           model.Account
		accBytes      []byte

		ms mocks.UNameSet

		GetNamesSetB4 func() model.NameSetHandler
	)
	BeforeEach(func() {
		id = "1234"
		name = "unameЯйцоЖЭ"
		pwd = "ka88dk#яфэюж"
		acc = model.Account{ID: id, Name: name, PwdHash: pwd}
		GetNamesSetB4 = model.GetNamesSet
		ms = mocks.UNameSet{}

	})
	AfterEach(func() {
		accBytes, _ = json.Marshal(acc)
		model.GetNamesSet = GetNamesSetB4
	})
	JustBeforeEach(func() {
		model.GetNamesSet = func() model.NameSetHandler {
			return &ms
		}
	})
	Context("when json decoding fails", func() {
		It("returns an error", func() {
			accBytes = []byte("")
			err := model.AddKafkaMsgToNameSet([]byte(name), accBytes)
			Expect(err).To(HaveOccurred())
		})
	})
	It("calls meth AddToSet", func() {
		err := model.AddKafkaMsgToNameSet([]byte(name), accBytes)
		Expect(err).NotTo(HaveOccurred())
		Expect(ms.AddToSetCall.Receives.Name).To(Equal(name))
	})
})
