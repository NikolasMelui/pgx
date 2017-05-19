package pgproto3

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
)

type Close struct {
	ObjectType byte // 'S' = prepared statement, 'P' = portal
	Name       string
}

func (*Close) Frontend() {}

func (dst *Close) Decode(src []byte) error {
	if len(src) < 2 {
		return &invalidMessageFormatErr{messageType: "Close"}
	}

	dst.ObjectType = src[0]
	rp := 1

	idx := bytes.IndexByte(src[rp:], 0)
	if idx != len(src[rp:])-1 {
		return &invalidMessageFormatErr{messageType: "Close"}
	}

	dst.Name = string(src[rp : len(src)-1])

	return nil
}

func (src *Close) MarshalBinary() ([]byte, error) {
	var bigEndian BigEndianBuf
	buf := &bytes.Buffer{}

	buf.WriteByte('C')
	buf.Write(bigEndian.Uint32(0))

	buf.WriteByte(src.ObjectType)
	buf.WriteString(src.Name)
	buf.WriteByte(0)

	binary.BigEndian.PutUint32(buf.Bytes()[1:5], uint32(buf.Len()-1))

	return buf.Bytes(), nil
}

func (src *Close) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type       string
		ObjectType string
		Name       string
	}{
		Type:       "Close",
		ObjectType: string(src.ObjectType),
		Name:       src.Name,
	})
}