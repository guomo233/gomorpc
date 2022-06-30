package gomorpc

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	pb "gomorpc/gomorpcpb"
	"io"
	"net"

	"github.com/golang/protobuf/proto"
)

type CodeType int

const (
	PROTO_BUF CodeType = iota
)

const (
	magic_number = 0x5599
	header_length = 34
)

const (
	ERROR_TIMEOUT = iota + 1
)

// if modify this, change HEADER_LENGTH !!!
type header struct {
	MagicNumber uint32 // must be exported, for json
	CodeType CodeType  // too
}

type codec interface {
	ReadHeader(*pb.MHeader) error
	ReadMessage(proto.Message) error
	Write(*pb.MHeader, proto.Message) error
	io.Closer
}

type protoBufCodec struct {
	conn io.Closer
	reader *bufio.Reader
	writer *bufio.Writer
}

type headerHandler struct {
	conn net.Conn
}

func (h *headerHandler) Write(header *header) error {
	headerBytes, err := json.Marshal(header)
	if err != nil {
		return fmt.Errorf("error to encode Header: %w", err)
	}

	if _, err := h.conn.Write(headerBytes); err != nil {
		return fmt.Errorf("error to send Header: %w", err)
	}

	return nil
}

func (h *headerHandler) ReadAndCheck(header *header) error {
	headerBytes := make([]byte, header_length)
	if _, err := io.ReadFull(h.conn, headerBytes); err != nil {
		return fmt.Errorf("error to read Header: %w", err)
	}

	if err := json.Unmarshal(headerBytes, header); err != nil {
		return fmt.Errorf("error to decode Header: %w", err)
	}

	if header.MagicNumber != magic_number {
		return fmt.Errorf("wrong magic number of Header")
	}

	return nil
}

func (p *protoBufCodec) ReadHeader(mh *pb.MHeader) (err error) {
	// read MHeader length
	lenBytes := make([]byte, 4)
	if _, err := io.ReadFull(p.reader, lenBytes); err != nil {
		return fmt.Errorf("error to read MHeader length: %w", err)
	}
	mhLen := binary.BigEndian.Uint32(lenBytes)

	// read MHeader
	mhBytes := make([]byte, mhLen)
	if _, err := io.ReadFull(p.reader, mhBytes); err != nil {
		return fmt.Errorf("error to read MHeader: %w", err)
	}
	proto.Unmarshal(mhBytes, mh)

	return nil
}

func (p *protoBufCodec) ReadMessage(msg proto.Message) (err error) {
	// read message length
	lenBytes := make([]byte, 4)
	if _, err := io.ReadFull(p.reader, lenBytes); err != nil {
		return fmt.Errorf("error to read message length: %w", err)
	}
	msgLen := binary.BigEndian.Uint32(lenBytes)

	// read message
	msgBytes := make([]byte, msgLen)
	if _, err := io.ReadFull(p.reader, msgBytes); err != nil {
		return fmt.Errorf("error to read message: %w", err)
	}
	proto.Unmarshal(msgBytes, msg)

	return nil
}

func (p *protoBufCodec) Write(mh *pb.MHeader, msg proto.Message) (err error) {
	// encode MHeader
	mhBytes, err := proto.Marshal(mh)
	if err != nil {
		return fmt.Errorf("error to encode MHeader: %w", err)
	}
	
	// send MHeader length
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(mhBytes)))
	if _, err := p.writer.Write(lenBytes); err != nil {
		return fmt.Errorf("error to send MHeader length: %w", err)
	}

	// send MHeader
	if _, err := p.writer.Write(mhBytes); err != nil {
		return fmt.Errorf("error to send MHeader: %w", err)
	}

	// encode message
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error to encode message: %w", err)
	}

	// send message length
	binary.BigEndian.PutUint32(lenBytes, uint32(len(msgBytes)))
	if _, err := p.writer.Write(lenBytes); err != nil {
		return fmt.Errorf("error to send message length: %w", err)
	}

	// send message
	if _, err := p.writer.Write(msgBytes); err != nil {
		return fmt.Errorf("error to send message: %w", err)
	}

	if err := p.writer.Flush(); err != nil {
		return fmt.Errorf("error to flush: %w", err)
	}

	return nil
}

func (p *protoBufCodec) Close() error {
	return p.conn.Close()
}

//var _ codec = (*ProtoBufCodec)(nil)

func NewCodec(conn net.Conn, codeType CodeType) codec {
	return &protoBufCodec{
		conn,
		bufio.NewReader(conn),
		bufio.NewWriter(conn),
	}
}
