package ProtocolMultiplexer

type Message struct {

}
type ProtocolMux interface {
	Init(service Service)
	Call(message Message) Message
}

type Service interface {
	//ReadTag returns the muxing identifier in the request or reply message.
	//Multiple goroutines may call ReadTag concurrently
	ReadTag(message Message) int64

	//Send sends a request to the remote service
	//Send must not be called concurrently with is self
	Send(message Message)

	//Recv waits for and returns a reply message from the remote service
	//Recv must not be called concurrently with is self
	Recv() Message
}

