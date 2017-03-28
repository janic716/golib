package overmind

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
)

const CONTROLMESSAGE_MAGIC uint32 = 0xFEA64952

const (
	PingMessageType uint32 = iota + 1
	PongMessageType
	DataMessageType
	BroadcastMessageType
	SuccessMessageType
	FailureMessageType
	MultiMessageType
)

type ControlMessage interface {
	ReadOne(br net.Conn) error
	WriteOne(bw net.Conn) error
	SetSign([]byte) ControlMessage
	SetCondition([]byte) ControlMessage
	SetData([]byte) ControlMessage
	SetProduct([]byte) ControlMessage
	SetMultiSigns([]string) ControlMessage
	SetSource(source *string) ControlMessage
	GetType() uint32
	GetSign() []byte
	GetCondition() []byte
	GetData() []byte
	GetProduct() []byte
	GetMultiSigns() []string
	GetSource() string
}

var (
	PingControlMessage ControlMessage = &controlMessage{
		magic: CONTROLMESSAGE_MAGIC,
		body:  ControlMsgBody{ControlType: proto.Uint32(PingMessageType)},
	}
	PongControlMessage ControlMessage = &controlMessage{
		magic: CONTROLMESSAGE_MAGIC,
		body:  ControlMsgBody{ControlType: proto.Uint32(PongMessageType)},
	}
	SuccessControlMessage ControlMessage = &controlMessage{
		magic: CONTROLMESSAGE_MAGIC,
		body:  ControlMsgBody{ControlType: proto.Uint32(SuccessMessageType)},
	}
	FailureControlMessage ControlMessage = &controlMessage{
		magic: CONTROLMESSAGE_MAGIC,
		body:  ControlMsgBody{ControlType: proto.Uint32(FailureMessageType)},
	}
)

type controlMessage struct {
	magic      uint32
	bodyLength uint32
	body       ControlMsgBody
}

func NewControlMessage() ControlMessage {
	cMessage := controlMessage{}
	cMessage.magic = CONTROLMESSAGE_MAGIC
	cMessage.body = ControlMsgBody{}
	cMessage.body.ControlType = proto.Uint32(DataMessageType)
	return &cMessage
}

func NewBroadcastControlMessage() ControlMessage {
	cMessage := controlMessage{}
	cMessage.magic = CONTROLMESSAGE_MAGIC
	cMessage.body = ControlMsgBody{}
	cMessage.body.ControlType = proto.Uint32(BroadcastMessageType)
	return &cMessage
}

func NewMultiControlMessage() ControlMessage {
	cMessage := controlMessage{}
	cMessage.magic = CONTROLMESSAGE_MAGIC
	cMessage.body = ControlMsgBody{}
	cMessage.body.ControlType = proto.Uint32(MultiMessageType)
	return &cMessage
}

func (control *controlMessage) ReadOne(br net.Conn) (err error) {
	var header []byte = make([]byte, 4)
	if _, err = br.Read(header); err != nil {
		return fmt.Errorf("read header err: %s", err)
	}
	control.magic = binary.BigEndian.Uint32(header[0:4])
	if control.magic != CONTROLMESSAGE_MAGIC {
		return fmt.Errorf("illegal magic_num: %d", control.magic)
	}

	var (
		length     []byte = make([]byte, 4)
		bufBytes   []byte
		body       []byte
		readNBytes int
	)
	if _, err = br.Read(length); err != nil {
		return fmt.Errorf("read bodylength err: %s", err)
	}
	control.bodyLength = binary.BigEndian.Uint32(length[0:4])
	body = make([]byte, control.bodyLength)
	for bufBytes = body; len(bufBytes) > 0; bufBytes = bufBytes[readNBytes:] {
		if readNBytes, err = br.Read(bufBytes); err != nil {
			return fmt.Errorf("read controlMgs body err: %s", err)
		}
	}
	if err = proto.Unmarshal(body, &control.body); err != nil {
		return fmt.Errorf("protp Unmarshal body err: %s", err)
	}
	return
}

func (control *controlMessage) WriteOne(bw net.Conn) (err error) {
	var buf *bytes.Buffer = new(bytes.Buffer)
	if err = binary.Write(buf, binary.BigEndian, control.magic); err != nil {
		return fmt.Errorf("write magic err: %s", err)
	}
	body, err := proto.Marshal(&control.body)
	if err != nil {
		return fmt.Errorf("proto Marshal body err: %s", err)
	}
	control.bodyLength = uint32(len(body))
	if err = binary.Write(buf, binary.BigEndian, control.bodyLength); err != nil {
		return fmt.Errorf("write bodylength err: %s", err)
	}
	if err = binary.Write(buf, binary.BigEndian, body); err != nil {
		return fmt.Errorf("write body err: %s", err)
	}
	if _, err = bw.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("write buffer err: %s", err)
	}
	return
}

func (control *controlMessage) SetSign(sign []byte) ControlMessage {
	var length int = len(sign)
	if length > 0 {
		control.body.Sign = make([]byte, length)
		copy(control.body.Sign, sign)
	}
	return control
}

func (control *controlMessage) SetCondition(condition []byte) ControlMessage {
	var length int = len(condition)
	if length > 0 {
		control.body.Condition = make([]byte, length)
		copy(control.body.Condition, condition)
	}
	return control
}

func (control *controlMessage) SetData(data []byte) ControlMessage {
	var length int = len(data)
	if length > 0 {
		control.body.Data = make([]byte, length)
		copy(control.body.Data, data)
	}
	return control
}

func (control *controlMessage) SetProduct(product []byte) ControlMessage {
	var length int = len(product)
	if length > 0 {
		control.body.Product = make([]byte, length)
		copy(control.body.Product, product)
	}
	return control
}

func (control *controlMessage) SetMultiSigns(multiSigns []string) ControlMessage {
	var length int = len(multiSigns)
	if length > 0 {
		control.body.MultiSigns = make([]string, length)
		copy(control.body.MultiSigns, multiSigns)
	}
	return control
}

func (control *controlMessage) SetSource(source *string) ControlMessage {
	control.body.Source = source
	return control
}

func (control *controlMessage) GetType() uint32 {
	return control.body.GetControlType()
}

func (control *controlMessage) GetSign() []byte {
	return control.body.GetSign()
}

func (control *controlMessage) GetCondition() []byte {
	return control.body.GetCondition()
}

func (control *controlMessage) GetData() []byte {
	return control.body.GetData()
}

func (control *controlMessage) GetProduct() []byte {
	return control.body.GetProduct()
}

func (control *controlMessage) GetMultiSigns() []string {
	return control.body.GetMultiSigns()
}

func (control *controlMessage) GetSource() string {
	return control.body.GetSource()
}
