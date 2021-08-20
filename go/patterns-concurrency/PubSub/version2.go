package PubSub

type Server struct {
	publish   chan Event
	subscribe chan subReq
	cancel    chan subReq
}

type subReq struct {
	c  chan<- Event
	ok chan bool
}

func (s *Server) Init() {
	s.publish = make(chan Event)
	s.subscribe = make(chan subReq)
	s.cancel = make(chan subReq)
	go s.loop()
}
/* version 1 loop
func (s *Server) loop() {
	sub := make(map[chan<- Event]bool)
	for {
		select {
		case e := <-s.publish:
			for c := range sub {
				c <- e
			}

		case r := <-s.subscribe:
			if sub[r.c] {
				r.ok <- false
				break
			}
			sub[r.c] = true
			r.ok <- true

		case r := <-s.cancel:
			if !sub[r.c] {
				r.ok <- false
				break
			}
			close(r.c)
			delete(sub, r.c)
			r.ok <- true
		}
	}
}
*/
func (s *Server) Publish(e Event) {
	s.publish <- e
}

func (s *Server) Subscribe(c chan<- Event) {
	r := subReq{c: c, ok: make(chan bool)}
	s.subscribe <- r
	if !<-r.ok {
		panic("pubsub: already subscribed")
	}
}

func (s *Server) Cancel(c chan<- Event) {
	r := subReq{c: c, ok: make(chan bool)}
	s.cancel <- r
	if !<-r.ok {
		panic("pubsub: not subscribed")
	}
}

/* version 1
func helper(in <-chan Event, out chan<- Event) {
	var q []Event
	for {
		var sendOut chan<- Event
		var next Event
		if len(q) > 0 {
			sendOut = out
			next = q[0]
		}
		select {
		case e := <-in:
			q = append(q, e)
		case sendOut <- next:
			q = q[1:]
		}
	}
}
*/

//version 2
func helper(in <-chan Event, out chan<- Event) {
	var q []Event
	for in != nil && len(q) > 0 {
		var sendOut chan<- Event
		var next Event
		if len(q) > 0 {
			sendOut = out
			next = q[0]
		}
		select {
		case e, ok := <-in:
			if !ok {
				in = nil
				break
			}
			q = append(q, e)
		case sendOut <- next:
			q = q[1:]
		}
	}
	close(out)
}

// version 2
func (s *Server) loop() {
	sub := make(map[chan<- Event]chan<-Event)
	for {
		select {
		case e := <-s.publish:
			for _,c := range sub {
				c <- e
			}

		case r := <-s.subscribe:
			if sub[r.c]!=nil {
				r.ok <- false
				break
			}
			e := make(chan Event)
			go helper(e,r.c)
			sub[r.c] = e
			r.ok <- true

		case r := <-s.cancel:
			if sub[r.c]==nil {
				r.ok <- false
				break
			}
			close(sub[r.c])
			delete(sub, r.c)
			r.ok <- true
		}
	}
}
