package overmind

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
)

type Message interface {
	ReadOne(br *bufio.Reader) error
	WriteOne(bw *bufio.Writer) error
	SetSeqnum(uint32) Message
	SetLogid(uint64) Message
	SetType(uint8) Message
	SetBody([]byte) Message
	GetLength() uint32
	GetSeqnum() uint32
	GetLogid() uint64
	GetType() uint8
	GetBody() []byte
	ToBytes() ([]byte, error)
}

type message struct {
	magic   uint32
	version uint8
	length  uint32
	seqnum  uint32
	logid   uint64
	mtype   uint8
	body    []byte //序列化
}

func NewMessage() Message {
	return &message{
		magic:   NICE_MAGIC_NUM,
		version: 1,
	}
}

func (msg *message) ReadOne(br *bufio.Reader) (err error) {
	if err = msg.readHeader(br); err != nil {
		return
	}
	if msg.magic != NICE_MAGIC_NUM {
		err = fmt.Errorf("illegal magic_num: %d", msg.magic)
		return
	}
	if msg.length > MAX_BODY_SIZE {
		err = fmt.Errorf("illegal message size: %d", msg.length)
		return
	}
	return msg.readBody(br)
}

func (msg *message) readHeader(br *bufio.Reader) (err error) {
	// 解析包头
	defer func() {
		if r := recover(); r != nil {
			if err == nil {
				err = r.(error)
			}
		}
	}()

	var header []byte = make([]byte, 22)
	if _, err = br.Read(header); err != nil {
		return
	}
	msg.magic = binary.BigEndian.Uint32(header[0:4])
	msg.version = uint8(header[4:5][0])
	msg.length = binary.BigEndian.Uint32(header[5:9])
	msg.seqnum = binary.BigEndian.Uint32(header[9:13])
	msg.logid = binary.BigEndian.Uint64(header[13:21])
	msg.mtype = uint8(header[21:22][0])
	return
}

func (msg *message) readBody(br *bufio.Reader) (err error) {
	msg.body = make([]byte, msg.length)
	var (
		bufBytes   []byte
		readNBytes int
	)
	for bufBytes = msg.body; len(bufBytes) > 0; bufBytes = bufBytes[readNBytes:] {
		if readNBytes, err = br.Read(bufBytes); err != nil {
			return
		}
	}
	return
}

func (msg *message) WriteOne(bw *bufio.Writer) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if err == nil {
				err = r.(error)
			}
		}
	}()

	var buf *bytes.Buffer = new(bytes.Buffer)
	// 写buf
	if err = binary.Write(buf, binary.BigEndian, msg.magic); err != nil {
		return
	}
	if err = binary.Write(buf, binary.BigEndian, msg.version); err != nil {
		return
	}
	if err = binary.Write(buf, binary.BigEndian, msg.length); err != nil {
		return
	}
	if err = binary.Write(buf, binary.BigEndian, msg.seqnum); err != nil {
		return
	}
	if err = binary.Write(buf, binary.BigEndian, msg.logid); err != nil {
		return
	}
	if err = binary.Write(buf, binary.BigEndian, msg.mtype); err != nil {
		return
	}
	if err = binary.Write(buf, binary.BigEndian, msg.body); err != nil {
		return
	}
	// 往io写数据
	if _, err = bw.Write(buf.Bytes()); err != nil {
		return
	}
	return bw.Flush()
}

func (msg *message) SetSeqnum(seq uint32) Message {
	msg.seqnum = seq
	return msg
}

func (msg *message) SetLogid(logid uint64) Message {
	msg.logid = logid
	return msg
}

func (msg *message) SetType(mtype uint8) Message {
	msg.mtype = mtype
	return msg
}

func (msg *message) SetBody(body []byte) Message {
	var length int = len(body)
	if length > 0 {
		msg.length = uint32(length)
		msg.body = make([]byte, length)
		copy(msg.body, body)
	}
	return msg
}

func (msg *message) GetLength() uint32 {
	return msg.length
}

func (msg *message) GetSeqnum() uint32 {
	return msg.seqnum
}

func (msg *message) GetLogid() uint64 {
	return msg.logid
}

func (msg *message) GetType() uint8 {
	return msg.mtype
}

func (msg *message) GetBody() []byte {
	return msg.body
}

func (msg *message) ToBytes() (data []byte, err error) {
	buf := new(bytes.Buffer)
	bw := bufio.NewWriter(buf)
	if err = msg.WriteOne(bw); err != nil {
		return
	}
	return buf.Bytes(), nil
}
