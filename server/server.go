package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/baoqger/simple-message-broker/hashmap"
	sublist "github.com/baoqger/simple-message-broker/sublist"
)

type info struct {
	Id           string `json:"server_id"`
	Version      string `json:"version"`
	Host         string `json:"host"`
	Port         uint   `json:"port"`
	AuthRequired bool   `json:"auth_required"`
	SslRequired  bool   `json:"ssl_required"`
	MaxPayload   int    `json:"max_payload"`
}

type Server struct {
	info     info
	infoJson []byte
	sl       *sublist.Sublist
	gcid     uint64
}

func New() *Server {
	s := &Server{
		info: info{
			Id:           genId(),
			Version:      VERSION,
			Host:         DEFAULT_HOST,
			Port:         DEFAULT_PORT,
			AuthRequired: false,
			SslRequired:  false,
			MaxPayload:   MAX_PAYLOAD_SIZE,
		},
		sl: sublist.New(),
	}
	b, err := json.Marshal(s.info)
	if err != nil {
		log.Fatalf("Err marshalling INFO JSON: %+v\n", err)
	}
	s.infoJson = []byte(fmt.Sprintf("INFO %s %s", b, CR_LF))
	return s
}

func (s *Server) AcceptLoop() {
	l, e := net.Listen("tcp", "0.0.0.0:4222")
	if e != nil {
		println(e)
		return
	}
	log.Println("Listening on ", l.Addr())
	for {
		conn, err := l.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				log.Printf("Accept error: %v", err)
			}
			continue
		}
		s.createClient(conn)
	}
}

func (s *Server) createClient(conn net.Conn) *client {
	c := &client{srv: s, conn: conn}
	c.cid = atomic.AddUint64(&s.gcid, 1)
	c.bw = bufio.NewWriterSize(c.conn, defaultBufSize)
	c.br = bufio.NewReaderSize(c.conn, defaultBufSize)
	c.subs = hashmap.New()

	s.sendInfo(c)
	go c.readLoop()
	return c
}

func (s *Server) sendInfo(c *client) {
	c.conn.Write(s.infoJson)
}
