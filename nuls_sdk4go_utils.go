package nuls_sdk4go

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math"
	"math/big"
)

// big endian 32 bytes => hex
func GetHexPrivKey(value *big.Int) (hexpriv string, err error) {
	hexpriv = ""
	bytesbuf := [32]byte{0}
	oBytes := value.Bytes()
	count := len(oBytes)

	if count == 32 {
		hexpriv = hex.EncodeToString(oBytes)
		return hexpriv, nil
	} else if count > 32 {
		return hexpriv, errors.New("The number is too large!")
	} else {
		var j int = 0
		for i := 32 - count; i < 32; i++ {
			bytesbuf[i] = oBytes[j]
			j++
		}
	}
	return hex.EncodeToString(bytesbuf[:]), nil
}

//return 32 bytes    little endian
func BigInteger2Bytes(value *big.Int) (r [32]byte, err error) {
	bytesbuf := [32]byte{0}
	oBytes := value.Bytes()
	count := len(oBytes)
	if count > 32 {
		return bytesbuf, errors.New("The number is too large!")
	}
	var j int
	j = 0
	for i := count - 1; i >= 0; i-- {
		bytesbuf[j] = oBytes[i]
		j++
	}
	return bytesbuf, nil
}

//const HASH_LENGTH = 32

func calcHashTwice(data []byte) (hashbytes [32]byte, err error) {
	bytesbuf := [32]byte{0}
	h := sha256.New()
	_, err = h.Write(data)
	if err != nil {
		return bytesbuf, err
	}
	r := h.Sum(nil)

	h2 := sha256.New()

	_, err = h2.Write(r)
	if err != nil {
		return bytesbuf, err
	}
	r = h2.Sum(nil)

	if len(r) != 32 {
		return bytesbuf, errors.New("calc hash error")
	}

	for i := 0; i < 32; i++ {
		bytesbuf[i] = r[i]
	}
	return bytesbuf, nil

}

func bufferWriteBytes(buffer *bytes.Buffer, data []byte, size int) (err error) {
	for i := int(0); i < size; i++ {
		err = buffer.WriteByte(data[i])
	}
	return err
}

//for encode int
func sizeOfVarInt(value int64) int {
	if value < 0 {
		// 1 marker + 8 entity bytes
		return 9
	}
	if value < 253 {
		// 1 entity byte
		return 1
	}
	if value <= 0xFFFF {
		// 1 marker + 2 entity bytes
		return 3
	}
	if value <= 0xFFFFFFFF {
		// 1 marker + 4 entity bytes
		return 5
	}
	// 1 marker + 8 entity bytes
	return 9
}

func bufferWriteEncodedVarInt(buffer *bytes.Buffer, value int64) (err error) {
	size := sizeOfVarInt(value)
	if size == 1 {
		return buffer.WriteByte(byte(value))
	}
	if size == 3 {
		var data = [3]byte{253, byte(value), byte(value >> 8)}
		return bufferWriteBytes(buffer, data[:], 3)
	}
	if size == 5 {
		err = buffer.WriteByte(byte(254))
		if err != nil {
			return err
		}
		//小端
		return bufferWriteUInt32(buffer, uint32(value))
	}

	//default size == 9
	err = buffer.WriteByte(byte(254))

	if err != nil {
		return err
	}
	return bufferWriteUInt64(buffer, uint64(value))
}

func bufferWriteBytesWithLength(buffer *bytes.Buffer, data []byte, size int) (err error) {
	err = bufferWriteEncodedVarInt(buffer, int64(size))
	if err != nil {
		return err
	}
	return bufferWriteBytes(buffer, data, size)
}

func bufferWriteCharArray(buffer *bytes.Buffer, data string) error {
	dataBytes := []byte(data)
	return bufferWriteBytes(buffer, dataBytes, len(dataBytes))
}

func bufferWriteString(buffer *bytes.Buffer, data string) error {
	dataBytes := []byte(data)
	strLen := len(data)
	if strLen < 0xFE {
		res := bufferWriteInt8(buffer, uint8(strLen+1))
		if res != nil {
			return res
		}
	} else {
		res := bufferWriteInt8(buffer, 0xFF)
		if res != nil {
			return res
		}
		res = bufferWriteUInt64(buffer, uint64(strLen+1))
		if res != nil {
			return res
		}
	}
	return bufferWriteBytes(buffer, dataBytes, len(dataBytes))
}

func bufferWriteInt8(buffer *bytes.Buffer, data uint8) error {
	return buffer.WriteByte(byte(data))
}

func bufferWriteUInt32(buffer *bytes.Buffer, data uint32) error {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, data)
	_, err := buffer.Write(bs)
	return err
}

func bufferWriteUInt16(buffer *bytes.Buffer, data uint16) error {
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, data)
	_, err := buffer.Write(bs)
	return err
}

func bufferWriteUInt64(buffer *bytes.Buffer, data uint64) error {
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, data)
	_, err := buffer.Write(bs)
	return err
}

func bufferWriteInt64(buffer *bytes.Buffer, data int64) error {
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, uint64(data))
	_, err := buffer.Write(bs)
	return err
}

func bufferWriteFloat32(buffer *bytes.Buffer, data float32) error {
	return bufferWriteUInt32(buffer, math.Float32bits(data))
}

func bufferWriteFloat64(buffer *bytes.Buffer, data float64) error {
	return bufferWriteUInt64(buffer, math.Float64bits(data))
}

func bufferWriteBool(buffer *bytes.Buffer, data bool) error {
	var n uint8
	if data {
		n = 1
	} else {
		n = 0
	}
	return bufferWriteInt8(buffer, n)
}

func bufferWriteBigInt(buffer *bytes.Buffer, data *big.Int) (err error) {
	bytesbuf := [32]byte{0}
	oBytes := data.Bytes()
	count := len(oBytes)
	if count > 32 {
		return errors.New("The number is too large!")
	}
	var j int
	j = 0
	for i := count - 1; i >= 0; i-- {
		bytesbuf[j] = oBytes[i]
		j++
	}

	for i := int(0); i < 32; i++ {
		err = buffer.WriteByte(bytesbuf[i])
	}
	return err
}
