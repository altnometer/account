package model_test

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/altnometer/account/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FormOKer", func() {
	var (
		f                  model.RegForm
		name, pwd, pwdConf string
		NameIsInSetB4      func(name string) bool
	)
	BeforeEach(func() {
		name = "unameЯйцоЖЭ"
		pwd = "ka88dk#яфэюж"
		pwdConf = "ka88dk#яфэюж"
		NameIsInSetB4 = model.NameIsInSet
	})
	AfterEach(func() {
		model.NameIsInSet = NameIsInSetB4
	})

	JustBeforeEach(func() {
		f = model.RegForm{name, pwd, pwdConf}
	})
	Context("Pwd and PwdConf do not match", func() {
		BeforeEach(func() {
			pwdConf = "abc"
		})
		It("returns 400 code and error", func() {
			code, err := f.OK()
			Expect(code).To(Equal(400))
			Expect(err.Error()).To(ContainSubstring("PWD_PWDCONF_NO_MATCH"))
		})
	})
	Context("username contains a reserved name", func() {
		rand.Seed(time.Now().Unix())
		rns := model.ReservedUsernames
		uname := fmt.Sprintf("z%sж", rns[rand.Intn(len(rns)-1)])
		// uname := rns[rand.Intn(len(rns)-1)]
		BeforeEach(func() {
			name = uname
		})
		It("returns 400 code and error", func() {
			code, err := f.OK()
			Expect(code).To(Equal(400))
			Expect(err.Error()).To(ContainSubstring("NO_RESERVED_NAMES_ALLOWED"))
		})

	})
	Context("username exceeds max length", func() {
		uname := strings.Repeat("й", model.MaxUserNameLength+1)
		BeforeEach(func() {
			name = uname
		})
		It("returns 400 code and error", func() {
			code, err := f.OK()
			Expect(code).To(Equal(400))
			Expect(err.Error()).To(ContainSubstring("NAME_TOO_LONG"))
		})
	})
	Context("password exceeds max length", func() {
		pwdLong := strings.Repeat("й", model.MaxPasswordLength+1)
		BeforeEach(func() {
			pwd = pwdLong
			pwdConf = pwdLong
		})
		It("returns 400 code and error", func() {
			code, err := f.OK()
			Expect(code).To(Equal(400))
			Expect(err.Error()).To(ContainSubstring("PWD_TOO_LONG"))
		})
	})
	Context("password is less than min length", func() {
		pwdShort := strings.Repeat("й", model.MinPasswordLength-1)
		BeforeEach(func() {
			pwd = pwdShort
			pwdConf = pwdShort
		})
		It("returns 400 code and error", func() {
			code, err := f.OK()
			Expect(code).To(Equal(400))
			Expect(err.Error()).To(ContainSubstring("PWD_TOO_SHORT"))
		})
	})
	Context("username is an invalid utf8 string", func() {
		uname := "zйфж\xbd"
		BeforeEach(func() {
			name = uname
		})
		It("returns 400 code and error", func() {
			code, err := f.OK()
			Expect(code).To(Equal(400))
			Expect(err.Error()).To(ContainSubstring("NAME_INVALID_UTF8_STRING"))
		})
	})
	Context("username contains new line char", func() {
		uname := "zйфж\n"
		BeforeEach(func() {
			name = uname
		})
		It("returns 400 code and error", func() {
			code, err := f.OK()
			Expect(code).To(Equal(400))
			Expect(err.Error()).To(ContainSubstring("NAME_NEWLINE_NOT_ALLOWED"))
		})
	})
	Context("when username already exists", func() {
		BeforeEach(func() {
			model.NameIsInSet = func(string) bool {
				return true
			}
		})
		It("returns 400 code and error", func() {
			code, err := f.OK()
			Expect(code).To(Equal(400))
			Expect(err.Error()).To(ContainSubstring("NAME_ALREADY_EXISTS"))
		})
	})
})
