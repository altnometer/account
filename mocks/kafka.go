package mocks

import (
	"github.com/altnometer/account/model"
)

// KafkaSyncProducer is a mock of ISyncProducer interface.
type KafkaSyncProducer struct {
	SendAccMsgCall struct {
		Receives struct {
			Acc *model.Account
		}
		Returns struct {
			Error error
		}
	}
	InitMySyncProducerCall struct {
		Returns struct {
			Error error
		}
	}
	InitMySyncProducerCalled bool
}

// InitMySyncProducer is a mock method for KafkaSyncProducer.
func (p *KafkaSyncProducer) InitMySyncProducer() error {
	p.InitMySyncProducerCalled = true
	return p.InitMySyncProducerCall.Returns.Error
}

// SendAccMsg is a mock method for KafkaSyncProducer.
func (p *KafkaSyncProducer) SendAccMsg(acc *model.Account) error {
	p.SendAccMsgCall.Receives.Acc = acc
	return p.SendAccMsgCall.Returns.Error
}
