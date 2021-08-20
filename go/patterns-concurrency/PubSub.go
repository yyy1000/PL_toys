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
/*
type Server struct {
	mu  sync.Mutex
	sub map[chan<- Event]bool
}

func (s *Server) Init() {
	s.sub = make(map[chan<- Event]bool)
}

func (s *Server) Publish(e Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for c := range s.sub {
		c <- e
	}
}

func (s *Server) Subscribe(c chan<- Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.sub[c] {
		panic("pubsub: already subscribed")
	}
	s.sub[c] = true
}

func (s *Server) Cancel(c chan<- Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.sub[c] {
		panic("pubsub: not subscibed")
	}
	close(c)
	delete(s.sub, c)
}
*/
