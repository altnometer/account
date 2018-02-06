// Package kafka provides a producer and consumer API
// for message streaming.
//
//
// A kafka consumer usage.
//
// Call
// err := ConsumeAccMsgs(handler msgHandler)
// where msgHandler type is func(key, val []byte) error.
// ConsumeAccMsgs() blocks execution. Execute from a goroutine if required.
// An error returned by msgHadler would shutdown the consumer that can be
// initialize only once. You need to restart the process after the error.
//
package kafka
