package part_v

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"log"
)

// IntToHex int64 转字节数组
func IntToHex(num int64) []byte {
	buff := bytes.NewBuffer(make([]byte, 0))
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
func DeserializeBlock(d []byte) *Block {
	var target Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	_ = decoder.Decode(&target)

	return &target
}

// ReverseBytes reverses a byte array
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}
