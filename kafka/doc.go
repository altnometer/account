// Package kafka provides a producer and consumer API
// for message streaming.
//
//
// A kafka consumer usage.
//
// Get the consumer instance with kafka.GetConsumer().
// Call kafka.GetConsumer().ConsumeMsgs(handler msgHandler)
// where msgHandler type is func(key, val []byte) error.
package kafka
