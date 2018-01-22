package mocks

// KafkaSyncProducer is a mock of ISyncProducer interface.
type KafkaSyncProducer struct {
	SendAccMsgCall struct {
		Receives struct {
			Key, Val string
		}
		Returns struct {
			Error error
		}
	}
}

// SendAccMsg is a mock method for KafkaSyncProducer.
func (p *KafkaSyncProducer) SendAccMsg(key, val string) error {
	p.SendAccMsgCall.Receives.Key = key
	p.SendAccMsgCall.Receives.Val = val
	return p.SendAccMsgCall.Returns.Error
}
