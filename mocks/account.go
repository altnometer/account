package mocks

import "net/http"

// Account is a mock of modle.Account
type Account struct {
	SendToKafkaCalled bool
	SendToKafkaCall   struct {
		Receives struct {
			R *http.Request
		}
		Returns struct {
			Error error
		}
	}
}

// SendToKafka is a mock method to implement AccSender interface.
func (a *Account) SendToKafka(r *http.Request) error {
	a.SendToKafkaCalled = true
	a.SendToKafkaCall.Receives.R = r
	return a.SendToKafkaCall.Returns.Error
}
