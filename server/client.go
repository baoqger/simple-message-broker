package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/baoqger/simple-message-broker/hashmap"
)

const defaultBufSize = 65536

type client struct {
	mu   sync.Mutex
	cid  uint64
	opts clientOpts
	conn net.Conn
	bw   *bufio.Writer
	br   *bufio.Reader
	srv  *Server
	subs *hashmap.HashMap
	cstats
	parseState
}

type clientOpts struct {
	Verbose     bool `json:"verbose"`
	Pedantic    bool `json:"pedantic"`
	SslRequired bool `json:"ssl_required"`
}

type cstats struct {
	nr int
	nb int
	nm int
}

type subscription struct {
	client  *client
	subject []byte
	queue   []byte
	sid     []byte
	nm      int64
	max     int64
}

func (c *client) readLoop() {
	b := make([]byte, defaultBufSize)
	for {
		n, err := c.conn.Read(b)
		if err != nil {
			c.closeConnection()
			return
		}
		if err := c.parse(b[:n]); err != nil {
			log.Printf("Parse Error: %v\n", err)
			c.closeConnection()
			return
		}
	}
}

func (c *client) parse(buf []byte) error {
	var i int
	var b byte

	c.nr++
	c.nb += len(buf)
	for i, b = range buf {
		switch c.state {
		case OP_START:
			switch b {
			case 'C', 'c':
				c.state = OP_C
			case 'P', 'p':
				c.state = OP_P
			case 'S', 's':
				c.state = OP_S
			case 'U', 'u':
				c.state = OP_U
			default:
				goto parseErr
			}
		}
	}
	return nil
parseErr:
	return fmt.Errorf("parse Error [%d]: '%s'", c.state, buf[i:])
}

func (c *client) closeConnection() {
	if c.conn == nil {
		return
	}

	c.conn.Close() // close connection
	c.conn = nil

	if c.srv != nil {
		subs := c.subs.All() // all the subscriptions
		for _, s := range subs {
			sub := s.(*subscription)
			c.srv.sl.Remove(sub.subject, sub)
		}
	}
}
