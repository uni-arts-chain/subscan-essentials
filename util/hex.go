package util

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
)

// Add `0x` prefix
func AddHex(s string) string {
	if strings.TrimSpace(s) == "" {
		return ""
	}
	if strings.HasPrefix(s, "0x") {
		return s
	}
	return strings.ToLower("0x" + s)
}

func TrimHex(s string) string {
	return strings.TrimPrefix(s, "0x")
}

func IntToHex(i interface{}) string {
	return fmt.Sprintf("%x", i)
}

func HexToNumStr(v string) string {
	return U256(v).String()
}

func HexToBytes(s string) []byte {
	s = strings.TrimPrefix(s, "0x")
	c := make([]byte, hex.DecodedLen(len(s)))
	_, _ = hex.Decode(c, []byte(s))
	return c
}

func BytesToHex(b []byte) string {
	c := make([]byte, hex.EncodedLen(len(b)))
	hex.Encode(c, b)
	return string(c)
}

func IntToEncode64Hex(value int) string {
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, uint64(value))
	return BytesToHex(bs)
}