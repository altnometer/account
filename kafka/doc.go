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
//
package kafka
