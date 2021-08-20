package PubSub

/*
Options for slow goroutines

1. slow down event generation
2. drop events
3. queue an arbitary number of events
*/

type Event struct{}

type PubSub interface {
	Publish(e Event)

	Subscribe(c chan<- Event)

	Cancel(c chan<- Event)
}

