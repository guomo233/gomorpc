package gomorpc

import (
	"net"
	"sync"
	"time"

	pb "gomorpc/gomorpcpb"

	"github.com/golang/protobuf/proto"
)

type reply struct {
	ret proto.Message
	done chan uint32
}

type client struct {
	codec codec
	id uint32
	mu sync.Mutex // for replys
	pending map[uint32]*reply
}

type Client interface {
	Call(string, string, proto.Message, proto.Message, time.Duration) (chan uint32, error)
}

func NewClient(addr string, codeType CodeType, timeout time.Duration, loadBalancer LoadBalancer) (Client, error) {
	// resolve domain
	if addr[:7] == "domain:" {
		hosts, err := _resolverRegister.Resolve(addr)
		if err != nil {
			return nil, err
		}

		if loadBalancer == nil {
			addr = hosts[0]
		} else {
			addr = loadBalancer.Get(hosts)
		}
	}

	// connect
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return nil, err
	}

	// send Header
	header := &header{magic_number, codeType}
	headerHandler := &headerHandler{conn}
	if err := headerHandler.Write(header); err != nil {
		conn.Close()
		return nil, err
	}

	client := &client{
		codec: NewCodec(conn, codeType),
		pending: make(map[uint32]*reply),
	}

	go client.handleReply()

	return client, nil
}

// timeout == 0 means no limit
func (c *client) Call(service, method string, arg, ret proto.Message, timeout time.Duration) (chan uint32, error) {
	mh := &pb.MHeader{Service: service, Method: method, Id: c.id, Timeout: uint64(timeout)}

	if err := c.codec.Write(mh, arg); err != nil {
		c.codec.Close()
		return nil, err
	}

	c.mu.Lock()
	reply := &reply{ret, make(chan uint32, 1)}
	c.pending[c.id] = reply
	c.mu.Unlock()

	done := make(chan uint32, 1)
	go c.handleTimeout(timeout, reply.done, done)
	
	c.id += 1
	return done, nil
}

func (c *client) handleTimeout(timeout time.Duration, replyDone, done chan uint32) {
	if timeout == 0 {
		done <- <- replyDone	
	} else {
		select {
		case done <- <- replyDone:
		case <- time.After(timeout):
			done <- ERROR_TIMEOUT
		}
	}
}

func (c *client) handleReply() {
	defer c.codec.Close()

	for {
		mh := &pb.MHeader{}
		if err := c.codec.ReadHeader(mh); err != nil {
			panic(err)
		}

		c.mu.Lock()
		reply := c.pending[mh.Id]
		delete(c.pending, mh.Id)
		c.mu.Unlock()

		if err := c.codec.ReadMessage(reply.ret); err != nil {
			panic(err)
		}

		reply.done <- mh.Error
	}
}

var _ Client = (*client)(nil)
