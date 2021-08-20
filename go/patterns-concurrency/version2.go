package PubSub

type Server struct {
	publish chan Event
	subscribe chan subReq
	cancel chan subReq
}

type subReq struct {
	c chan <- Event
	ok chan bool
}

func (s *Server) Init(){
	s.publish = make(chan Event)
	s.subscribe = make(chan subReq)
	s.cancel = make(chan subReq)
	go s.loop()
}

func (s *Server) loop(){
	sub := make(map[chan <- Event]bool)
	for {
		select {
		case e:= <-s.publish:
			for c:= range sub{
				c <- e
			}

		case r:= <-s.subscribe:
			if sub[r.c]{
				r.ok <- false
				break
			}
			sub[r.c]=true
			r.ok<-true

		case r:= <-s.cancel:
			if !sub[r.c]{
				r.ok <- false
				break
			}
			close(r.c)
			delete(sub,r.c)
			r.ok <- true
		}
	}
}

func (s *Server) Publish(e Event)  {
	s.publish <- e
}

func (s *Server) Subscribe(c chan <-Event){
	r := subReq{c: c,ok:make(chan bool)}
	s.subscribe <- r
	if !<-r.ok{
		panic("pubsub: already subscribed")
	}
}

func (s *Server) Cancel(c chan<-Event){
	r := subReq{c:c, ok:make(chan bool)}
	s.cancel <- r
	if !<-r.ok{
		panic("pubsub: not subscribed")
	}
}

