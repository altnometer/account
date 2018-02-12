package system_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/altnometer/account/model"

	"github.com/altnometer/account/kafka"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func randStr(length int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	abc := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ~!@#$%^&*(),./[]{}|=-<>"
	s := make([]byte, length)
	for i := 0; i < length; i++ {
		s[i] = abc[rand.Intn(len(abc))]
	}
	return string(s)

}

var _ = Describe("System: the webserver", func() {
	var (
		f                  *url.Values
		name, pwd, pwdConf string
	)
	urlPath := fmt.Sprintf("http://%s/register", serverAddress)
	ts := &http.Transport{
		IdleConnTimeout:   2 * time.Second,
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport: ts,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	// if you do not bring kafka up for every test run,
	// same name will cause "name already exist" error.
	BeforeEach(func() {
		name = randStr(14)
		pwd = "ka88dk;ad"
		pwdConf = "ka88dk;ad"
		f = &url.Values{}
		f.Add("Name", name)
		f.Add("Pwd", pwd)
		f.Add("PwdConf", pwdConf)

	})
	Context("with request to register", func() {
		It("sends Account msg to kafka stream", func() {
			res, err := client.PostForm(urlPath, *f)
			Expect(err).NotTo(HaveOccurred())
			defer res.Body.Close()

			c := make(chan []byte, 256)

			go func() {
				err = kafka.ConsumeAccMsgs(c)
				Expect(err).NotTo(HaveOccurred())
			}()
			var acc model.Account
			msgSent := false
			for msg := range c {
				err = json.Unmarshal(msg, &acc)
				Expect(err).NotTo(HaveOccurred())
				if acc.Name == name {
					msgSent = true
					break
				}
			}
			Expect(msgSent).To(Equal(true))
		})
		It("redirects to a correct url", func() {
			// res, err := client.Post(url, "application/x-www-form-urlencoded", strings.NewReader(f.Encode()))
			res, err := client.PostForm(urlPath, *f)
			Expect(err).NotTo(HaveOccurred())
			defer res.Body.Close()
			bodyBytes, err := ioutil.ReadAll(res.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(bodyBytes)).To(Equal(""))

			Expect(res.StatusCode).To(Equal(302))

			newUrl, err := res.Location()
			Expect(err).NotTo(HaveOccurred())
			Expect(newUrl.Path).To(Equal("/"))
		})
	})
	Context("submitting duplicate user name", func() {
		It("returns NAME_ALREADY_EXISTS error", func() {
			res, err := client.PostForm(urlPath, *f)
			Expect(err).NotTo(HaveOccurred())
			defer res.Body.Close()
			// nameSet := model.GetNameSet()
			time.Sleep(10 * time.Millisecond)
			res, err = client.PostForm(urlPath, *f)
			Expect(err).NotTo(HaveOccurred())
			bodyBytes, err := ioutil.ReadAll(res.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(res.StatusCode).To(Equal(400))
			Expect(string(bodyBytes)).To(ContainSubstring("NAME_ALREADY_EXISTS"))
		})
	})
})
