/*
 * Copyright (c) 2020.
 * Project:qitmeer
 * File:urlv4_test.go
 * Date:7/7/20 3:30 PM
 * Author:Jin
 * Email:lochjin@gmail.com
 */

package qnode

import (
	"crypto/ecdsa"
	"errors"
	"net"
	"reflect"
	"strings"
	"testing"

	"github.com/Qitmeer/qitmeer/p2p/crypto"
	"github.com/Qitmeer/qitmeer/p2p/qnr"
)

func init() {
	lookupIPFunc = func(name string) ([]net.IP, error) {
		if name == "node.example.org" {
			return []net.IP{{33, 44, 55, 66}}, nil
		}
		return nil, errors.New("no such host")
	}
}

var parseNodeTests = []struct {
	input      string
	wantError  string
	wantResult *Node
}{
	// Records
	{
		input: "qnr:-IS4QGrdq0ugARp5T2BZ41TrZOqLc_oKvZoPuZP5--anqWE_J-Tucc1xgkOL7qXl0puJgT7qc2KSvcupc4NCb0nr4tdjgmlkgnY0gmlwhH8AAAGJc2VjcDI1NmsxoQM6UUF2Rm-oFe1IH_rQkRCi00T2ybeMHRSvw1HDpRvjPYN1ZHCCdl8",
		wantResult: func() *Node {
			testKey, _ := crypto.HexToECDSA("45a915e4d060149eb4365960e6a7a45f334393093061116b197e3240065ff2d8")
			var r qnr.Record
			r.Set(qnr.IP{127, 0, 0, 1})
			r.Set(qnr.UDP(30303))
			r.SetSeq(99)
			SignV4(&r, testKey)
			n, _ := New(ValidSchemes, &r)
			return n
		}(),
	},
	// Invalid Records
	{
		input:     "qnr:",
		wantError: "EOF", // could be nicer
	},
	{
		input:     "qnr:x",
		wantError: "illegal base64 data at input byte 0",
	},
	{
		input:     "qnr:-EmGZm9vYmFyY4JpZIJ2NIJpcIR_AAABiXNlY3AyNTZrMaEDOlFBdkZvqBXtSB_60JEQotNE9sm3jB0Ur8NRw6Ub4z2DdWRwgnZf",
		wantError: qnr.ErrInvalidSig.Error(),
	},
	// Complete node URLs with IP address and ports
	{
		input:     "qnode://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@invalid.:3",
		wantError: `no such host`,
	},
	{
		input:     "qnode://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@127.0.0.1:foo",
		wantError: `invalid port`,
	},
	{
		input:     "qnode://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@127.0.0.1:3?discport=foo",
		wantError: `invalid discport in query`,
	},
	{
		input: "qnode://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@127.0.0.1:52150",
		wantResult: NewV4(
			hexPubkey("1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439"),
			net.IP{127, 0, 0, 1},
			52150,
			52150,
		),
	},
	{
		input: "qnode://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@[::]:52150",
		wantResult: NewV4(
			hexPubkey("1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439"),
			net.ParseIP("::"),
			52150,
			52150,
		),
	},
	{
		input: "qnode://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@[2001:db8:3c4d:15::abcd:ef12]:52150",
		wantResult: NewV4(
			hexPubkey("1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439"),
			net.ParseIP("2001:db8:3c4d:15::abcd:ef12"),
			52150,
			52150,
		),
	},
	{
		input: "qnode://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@127.0.0.1:52150?discport=22334",
		wantResult: NewV4(
			hexPubkey("1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439"),
			net.IP{0x7f, 0x0, 0x0, 0x1},
			52150,
			22334,
		),
	},
	// Incomplete node URLs with no address
	{
		input: "qnode://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439",
		wantResult: NewV4(
			hexPubkey("1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439"),
			nil, 0, 0,
		),
	},
	// Invalid URLs
	{
		input:     "",
		wantError: errMissingPrefix.Error(),
	},
	{
		input:     "1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439",
		wantError: errMissingPrefix.Error(),
	},
	{
		input:     "01010101",
		wantError: errMissingPrefix.Error(),
	},
	{
		input:     "qnode://01010101@123.124.125.126:3",
		wantError: `invalid public key (wrong length, want 128 hex chars)`,
	},
	{
		input:     "qnode://01010101",
		wantError: `invalid public key (wrong length, want 128 hex chars)`,
	},
	{
		input:     "http://foobar",
		wantError: errMissingPrefix.Error(),
	},
	{
		input:     "://foo",
		wantError: errMissingPrefix.Error(),
	},
}

func hexPubkey(h string) *ecdsa.PublicKey {
	k, err := parsePubkey(h)
	if err != nil {
		panic(err)
	}
	return k
}

func TestParseNode(t *testing.T) {
	for _, test := range parseNodeTests {
		n, err := Parse(ValidSchemes, test.input)
		if test.wantError != "" {
			if err == nil {
				t.Errorf("test %q:\n  got nil error, expected %#q", test.input, test.wantError)
				continue
			} else if !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("test %q:\n  got error %#q, expected %#q", test.input, err.Error(), test.wantError)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("test %q:\n  unexpected error: %v", test.input, err)
				continue
			}
			if !reflect.DeepEqual(n, test.wantResult) {
				t.Errorf("test %q:\n  result mismatch:\ngot:  %#v\nwant: %#v", test.input, n, test.wantResult)
			}
		}
	}
}

func TestNodeString(t *testing.T) {
	for i, test := range parseNodeTests {
		if test.wantError == "" && strings.HasPrefix(test.input, "qnode://") {
			str := test.wantResult.String()
			if str != test.input {
				t.Errorf("test %d: Node.String() mismatch:\ngot:  %s\nwant: %s", i, str, test.input)
			}
		}
	}
}