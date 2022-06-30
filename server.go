package gomorpc

import (
	"fmt"
	"log"
	"net"
	"reflect"
	"sync"
	"time"

	pb "gomorpc/gomorpcpb"

	"github.com/golang/protobuf/proto"
)

type service struct {
	instance reflect.Value
	methods map[string]reflect.Method
}

type Server interface {
	Server() error
	Register(interface{}) error
}

type server struct {
	addr string
	services map[string]*service
	mu sync.Mutex // for send reply
}

func NewServer(addr string) Server {
	return &server{
		addr: addr,
		services: make(map[string]*service),
	}
}

func (s *server) Server() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			return err
		}

		go s.handleConn(conn)
	}
}

// method should be like: Method(arg proto.Message, ret proto.Message)
func (s *server) Register(instance interface{}) error {
	// check is instance a ptr
	svcType := reflect.TypeOf(instance)
	if svcType.Kind() != reflect.Ptr {
		return fmt.Errorf("expect ptr but got struct")
	}

	service := &service{
		reflect.ValueOf(instance),
		make(map[string]reflect.Method),
	}

	for i := 0; i < svcType.NumMethod(); i++ {
		method := svcType.Method(i)

		// check is method valid
		if !method.IsExported() {
			return fmt.Errorf("method %s is not exported", method.Name)
		}

		if method.Type.NumIn() != 3 || method.Type.NumOut() != 0 {
			return fmt.Errorf("wrong arg or ret number")
		}

		messageType := reflect.TypeOf((*proto.Message)(nil)).Elem()
		if !method.Type.In(1).Implements(messageType) || !method.Type.In(2).Implements(messageType) {
			return fmt.Errorf("wrong type of arg")
		}

		service.methods[method.Name] = method
	}

	s.services[svcType.Elem().Name()] = service
	return nil
}

func (s *server) handleConn(conn net.Conn) {
	// read Header
	header := &header{}
	headerHandler := &headerHandler{conn}
	if err := headerHandler.ReadAndCheck(header); err != nil {
		conn.Close()
		log.Printf("[Server] [ERROR] : %s", err)
		return
	}

	codec := NewCodec(conn, header.CodeType)
	for {
		// read MHeader
		mh := &pb.MHeader{}
		if err := codec.ReadHeader(mh); err != nil {
			codec.Close()
			log.Printf("[Server] [ERROR] %s", err)
			return
		}

		// reflect arg and ret
		service := s.services[mh.Service]
		method := service.methods[mh.Method]
		argv := reflect.New(method.Type.In(1).Elem())
		retv := reflect.New(method.Type.In(2).Elem())

		// read arg
		argi := argv.Interface().(proto.Message)
		if err := codec.ReadMessage(argi); err != nil {
			codec.Close()
			log.Printf("[Server] [ERROR] %s", err)
			return
		}

		go s.handleCall(codec, method, mh, service.instance, argv, retv)
	}
}

func (s *server) handleCall(codec codec, method reflect.Method, mh *pb.MHeader, instance, argv, retv reflect.Value) {
	done := make(chan struct{}, 1)

	go func() {
		method.Func.Call([]reflect.Value{instance, argv, retv})
		done <- struct{}{}
	}()

	// check timeout
	if mh.Timeout == 0 {
		<- done
	} else {
		select {
		case <- time.After(time.Duration(mh.Timeout)):
			mh.Error = ERROR_TIMEOUT
		case <- done:
		}
	}

	// reply
	s.mu.Lock()
	defer s.mu.Unlock()
	reti := retv.Interface().(proto.Message)
	if err := codec.Write(mh, reti); err != nil {
		codec.Close()
		log.Printf("[Server] [ERROR] %s", err)
		return
	}
}

var _ Server = (*server)(nil)
