package utreexo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"sort"
)

// HasAccess reports whether we have access to the named file.
// Returns true if HasAccess, false if it doesn't.
// Does NOT tell us if the file exists or not.
// File might exist but may not be available to us
func HasAccess(fileName string) bool {
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func BtU32(b []byte) uint32 {
	if len(b) != 4 {
		log.Error(fmt.Sprintf("Got %x to BtU32 (%d bytes)\n", b, len(b)))
		return 0xffffffff
	}
	var i uint32
	buf := bytes.NewBuffer(b)
	binary.Read(buf, binary.BigEndian, &i)
	return i
}

// 4 byte Big Endian slice to uint32.  Returns ffffffff if something doesn't work.
func BtI32(b []byte) int32 {
	if len(b) != 4 {
		fmt.Printf("Got %x to ItU32 (%d bytes)\n", b, len(b))
		return -0x7fffffff
	}
	var i int32
	buf := bytes.NewBuffer(b)
	binary.Read(buf, binary.BigEndian, &i)
	return i
}

// BtI64 : 8 bytes to int64.  returns ffff. if it doesn't work.
func BtI64(b []byte) int64 {
	if len(b) != 8 {
		fmt.Printf("Got %x to BtI64 (%d bytes)\n", b, len(b))
		return -0x7fffffffffffffff
	}
	var i int64
	buf := bytes.NewBuffer(b)
	binary.Read(buf, binary.BigEndian, &i)
	return i
}

// U32tB converts uint32 to 4 bytes.  Always works.
func U32tB(i uint32) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, i)
	return buf.Bytes()
}

// I64tB : int64 to 8 bytes.  Always works.
func I64tB(i int64) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, i)
	return buf.Bytes()
}

// it'd be cool if you just had .sort() methods on slices of builtin types...
func sortUint32s(s []uint32) {
	sort.Slice(s, func(a, b int) bool { return s[a] < s[b] })
}

// PrefixLen16 puts a 2 byte length prefix in front of a byte slice
func PrefixLen16(b []byte) []byte {
	l := uint16(len(b))
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, l)
	return append(buf.Bytes(), b...)
}
