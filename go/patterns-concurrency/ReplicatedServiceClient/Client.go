package ReplicatedServiceClient

import (
	"sync"
	"time"
)

type Client struct {
	servers []string
	callOne func(string, Args) Reply

	mu     sync.Mutex
	prefer int
}

func (c *Client) Init(servers []string, callOne func(string, Args) Reply) {
	c.servers = servers
	c.callOne = callOne
}

func (c *Client) Call(args Args) Reply {
	type result struct {
		serverID int
		reply    Reply
	}
	const timeout = 1 * time.Second
	t := time.NewTimer(timeout)
	defer t.Stop()

	done := make(chan result, len(c.servers))
	for id := 0; id < len(c.servers); id++ {
		go func(id1 int) {
			done <- result{id1, c.callOne(c.servers[id1], args)}
		}(id)
		select {
		case r := <-done:
			return r.reply
		case <-t.C:
			//timeout
			t.Reset(timeout)
		}
	}
	r := <-done
	return r.reply
}
