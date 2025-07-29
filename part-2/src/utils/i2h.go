package utils

import (
	"bytes"
	"encoding/binary"
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
