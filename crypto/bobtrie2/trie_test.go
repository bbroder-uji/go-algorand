// Copyright (C) 2019-2023 Algorand, Inc.
// This file is part of go-algorand
//
// go-algorand is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// go-algorand is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with go-algorand.  If not, see <https://www.gnu.org/licenses/>.

package bobtrie2

import (
	"bytes"
	"fmt"
	"github.com/algorand/go-algorand/crypto"
	"github.com/stretchr/testify/require"
	"testing"
)

func verifyNewTrie(t *testing.T, mt *Trie) {
	require.NotNil(t, mt)
	require.NotNil(t, mt.db)
	require.Nil(t, mt.rootHash)
}

func TestMakeTrie(t *testing.T) {
	mt, err := MakeTrie()
	require.NoError(t, err)
	verifyNewTrie(t, mt)

}

func TestNodeSerialization(t *testing.T) {
	rn := &RootNode{}
	hash := crypto.Hash([]byte("rootnode"))
	rn.child = hash
	data, err := serializeRootNode(rn)
	require.NoError(t, err)
	expected := []byte{0x0, 0x73, 0x88, 0xa, 0x4c, 0x76, 0x1d, 0x89, 0x6f, 0x2c, 0x2, 0xbf, 0xd9, 0x44, 0x5, 0x1, 0xa9, 0xd1, 0x17, 0x47, 0x8c, 0x68, 0xab, 0x94, 0xb4, 0x93, 0xcb,
		0xab, 0xa9, 0x7c, 0x8, 0xcb, 0xfc}
	require.Equal(t, expected, data)
	rn2, err := deserializeRootNode(data)
	require.NoError(t, err)
	require.Equal(t, rn, rn2)

	ln := &LeafNode{}
	ln.keyEnd = []byte("leafendkey")
	for i := range ln.keyEnd {
		ln.keyEnd[i] &= 0x0f
	}
	ln.valueHash = crypto.Hash([]byte("leafvalue"))
	data, err = serializeLeafNode(ln)
	require.NoError(t, err)
	expected = []byte{0x4, 0x9a, 0xf2, 0xee, 0x24, 0xf9, 0xd3, 0xde, 0x8d, 0xdb, 0x45, 0x71, 0x82, 0x90, 0xca, 0x38, 0x42, 0xad, 0x8e, 0xcf, 0x81, 0x56, 0x17, 0x16, 0x55, 0x42, 0x73, 0x6, 0xaa, 0xd0, 0x16, 0x87, 0x45, 0xc5, 0x16, 0x5e, 0x4b, 0x59}
	require.Equal(t, expected, data)
	ln2, err := deserializeLeafNode(data)
	require.NoError(t, err)
	require.Equal(t, ln, ln2)
	ln.keyEnd = []byte("leafendke")
	for i := range ln.keyEnd {
		ln.keyEnd[i] &= 0x0f
	}
	data, err = serializeLeafNode(ln)
	require.NoError(t, err)
	expected = []byte{0x3, 0x9a, 0xf2, 0xee, 0x24, 0xf9, 0xd3, 0xde, 0x8d, 0xdb, 0x45, 0x71, 0x82, 0x90, 0xca, 0x38, 0x42, 0xad, 0x8e, 0xcf, 0x81, 0x56, 0x17, 0x16, 0x55, 0x42, 0x73, 0x6, 0xaa, 0xd0, 0x16, 0x87, 0x45, 0xc5, 0x16, 0x5e, 0x4b, 0x50}
	require.Equal(t, expected, data)
	ln3, err := deserializeLeafNode(data)
	require.NoError(t, err)
	require.Equal(t, ln, ln3)

	bn := &BranchNode{}
	bn.children[0] = crypto.Hash([]byte("branchchild0"))
	bn.children[1] = crypto.Hash([]byte("branchchild1"))
	bn.children[2] = crypto.Hash([]byte("branchchild2"))
	bn.children[3] = crypto.Hash([]byte("branchchild3"))
	bn.children[4] = crypto.Hash([]byte("branchchild4"))
	bn.children[5] = crypto.Hash([]byte("branchchild5"))
	bn.children[6] = crypto.Hash([]byte("branchchild6"))
	bn.children[7] = crypto.Hash([]byte("branchchild7"))
	bn.children[8] = crypto.Hash([]byte("branchchild8"))
	bn.children[9] = crypto.Hash([]byte("branchchild9"))
	bn.children[10] = crypto.Hash([]byte("branchchild10"))
	bn.children[11] = crypto.Hash([]byte("branchchild11"))
	bn.children[12] = crypto.Hash([]byte("branchchild12"))
	bn.children[13] = crypto.Hash([]byte("branchchild13"))
	bn.children[14] = crypto.Hash([]byte("branchchild14"))
	bn.children[15] = crypto.Hash([]byte("branchchild15"))
	bn.valueHash = crypto.Hash([]byte("branchvalue"))
	data, err = serializeBranchNode(bn)
	require.NoError(t, err)
	expected = []byte{0x5, 0xe8, 0x31, 0x2c, 0x27, 0xec, 0x3d, 0x32, 0x7, 0x48, 0xab, 0x13, 0xed, 0x2f, 0x67, 0x94, 0xb3, 0x34, 0x8f, 0x1e, 0x14, 0xe5, 0xac, 0x87, 0x6e, 0x7, 0x68, 0xd6, 0xf6, 0x92, 0x99, 0x4b, 0xc8, 0x2e, 0x93, 0xde, 0xf1, 0x72, 0xc8, 0x55, 0xbb, 0x7e, 0xd1, 0x1d, 0x38, 0x6, 0xd2, 0x97, 0xd7, 0x2, 0x2, 0x86, 0x93, 0x37, 0x57, 0xce, 0xa4, 0xc5, 0x7e, 0x4c, 0xd4, 0x50, 0x94, 0x2e, 0x75, 0xeb, 0xcd, 0x9b, 0x80, 0xa2, 0xf5, 0xf3, 0x15, 0x4a, 0xf2, 0x62, 0x6, 0x7d, 0x6d, 0xdd, 0xe9, 0x20, 0xe1, 0x1a, 0x95, 0x3b, 0x2b, 0xb9, 0xc1, 0xaf, 0x3e, 0xcb, 0x72, 0x1d, 0x3f, 0xad, 0xe9, 0xa6, 0x30, 0xc6, 0xc5, 0x65, 0xf, 0x86, 0xb2, 0x3a, 0x5b, 0x47, 0xcb, 0x29, 0x31, 0xf7, 0x8a, 0xdf, 0xe0, 0x41, 0x6b, 0x11, 0xc0, 0xd, 0xbc, 0x80, 0xa7, 0x48, 0x97, 0x21, 0xbd, 0xee, 0x6f, 0x36, 0xf4, 0x7b, 0x6d, 0x68, 0xa1, 0x43, 0x31, 0x90, 0xf8, 0x56, 0x69, 0x4c, 0xee, 0x88, 0x76, 0x9c, 0xd1, 0xde, 0xe4, 0xbd, 0x64, 0x7d, 0x18, 0xce, 0xd6, 0xdb, 0xf8, 0x85, 0x84, 0x88, 0x5d, 0x7e, 0xda, 0xe0, 0xf2, 0xa0, 0x6d, 0x24, 0x4f, 0xcf, 0xb, 0x8c, 0x34, 0x57, 0x2a, 0x13, 0x22, 0xd9, 0x8d, 0x79, 0x8, 0xa4, 0x22, 0x91, 0x45, 0x64, 0x7b, 0xf3, 0xad, 0xe8, 0x9b, 0x5f, 0x7c, 0x5c, 0xbd, 0x9, 0xd3, 0xc7, 0x3, 0xe2, 0xef, 0x6b, 0x8, 0x8, 0x98, 0x52, 0xb, 0xd1, 0x6a, 0x5a, 0x18, 0x89, 0x44, 0x4f, 0xf1, 0xb0, 0x37, 0xd9, 0x7f, 0x99, 0x3f, 0x6a, 0x84, 0x46, 0x83, 0x2c, 0x91, 0x58, 0xa8, 0xb3, 0xda, 0xd8, 0x26, 0x2e, 0x8a, 0x4, 0x8f, 0x81, 0xa5, 0xf3, 0xef, 0x46, 0x34, 0x4a, 0x8f, 0x6a, 0x61, 0x2f, 0x3, 0x26, 0x9d, 0xe6, 0x77, 0xee, 0xec, 0xe2, 0xa4, 0x84, 0x38, 0x6b, 0x6e, 0x7e, 0xf0, 0xef, 0xaa, 0x29, 0xa5, 0x13, 0x0, 0xef, 0xff, 0xdf, 0xb5, 0xd7, 0x4e, 0x41, 0x75, 0x4d, 0x2, 0x84, 0x20, 0xe2, 0x18, 0x50, 0x52, 0xae, 0xf4, 0xea, 0xeb, 0x84, 0xb3, 0x91, 0x85, 0xa8, 0xa, 0xba, 0xc9, 0x31, 0x9f, 0x5e, 0x3e, 0xf8, 0xb5, 0xf4, 0x4b, 0xf8, 0xf2, 0xf0, 0x76, 0xa1, 0x6d, 0xec, 0x57, 0x65, 0xbd, 0x2e, 0x78, 0xbe, 0xf4, 0x7c, 0xe4, 0xf2, 0x45, 0xc0, 0xaf, 0x94, 0xb, 0x45, 0x1b, 0xd3, 0xcf, 0x9f, 0x17, 0x7e, 0x1a, 0x52, 0x6d, 0x18, 0xe5, 0x1a, 0x7c, 0xd9, 0x9d, 0xef, 0x8a, 0xe3, 0xe9, 0xe6, 0xf6, 0x76, 0x5e, 0x12, 0xbf, 0xd2, 0xe8, 0xaa, 0x8, 0x88, 0x15, 0x81, 0x99, 0x4e, 0xa3, 0x12, 0x98, 0xc1, 0xb3, 0xde, 0x42, 0x53, 0x2, 0x29, 0x82, 0x87, 0xfe, 0x3d, 0x8, 0xe0, 0xc2, 0x3, 0x70, 0x56, 0xd, 0x9, 0xad, 0xe4, 0x1a, 0xa5, 0xf6, 0x4, 0xdb, 0x63, 0xd0, 0x49, 0x6b, 0x5b, 0xa2, 0x56, 0xb1, 0xd1, 0x4b, 0x56, 0xc3, 0x7e, 0x4b, 0xec, 0xb5, 0xdb, 0xd4, 0xd9, 0xe1, 0x20, 0x99, 0x80, 0x71, 0x9, 0x72, 0x3b, 0xc, 0x8b, 0x56, 0x4, 0x94, 0xe6, 0x4e, 0x35, 0xd, 0x3e, 0x7, 0x8b, 0x86, 0x73, 0x62, 0x5f, 0x61, 0x8d, 0x70, 0x68, 0x86, 0xe8, 0x65, 0xbe, 0x18, 0xa8, 0x4a, 0xac, 0x6d, 0x81, 0x15, 0xde, 0x1b, 0xe1, 0xb3, 0xe8, 0x6a, 0x46, 0xdf, 0xdc, 0xf1, 0x6, 0x3c, 0xa6, 0x1c, 0xc9, 0xcd, 0x12, 0x5e, 0x5f, 0x28, 0xd1, 0x71, 0x6e, 0x9f, 0xc7, 0xdc, 0x77, 0x98, 0x47, 0x7, 0x94, 0x38, 0x4, 0xc4, 0xc4, 0xfe, 0x17, 0x12, 0x1b, 0xcf, 0x96, 0xd8, 0xb1, 0xf2, 0x1e, 0x81, 0xab, 0x15, 0x86, 0x75, 0x5a, 0x39, 0x13, 0xdb, 0xe, 0x1a, 0xd9, 0xa9, 0x70, 0x7d, 0xdd, 0xaf, 0x64, 0x12, 0x27, 0xe5, 0x97, 0xa1, 0x34, 0xb8, 0x1a, 0x61, 0x48, 0x29, 0x61, 0x62, 0xe4, 0x40, 0xba, 0x5, 0x44, 0x24, 0x51, 0xc1, 0x9b, 0x8e, 0x62, 0xf2, 0x1c, 0x6f, 0xd6, 0x8, 0x3, 0xbe, 0x88, 0xf}
	require.Equal(t, expected, data)
	bn2, err := deserializeBranchNode(data)
	require.NoError(t, err)
	require.Equal(t, bn, bn2)

	en := &ExtensionNode{}
	en.sharedKey = []byte("extensionkey")
	for i := range en.sharedKey {
		en.sharedKey[i] &= 0x0f
	}
	en.child = crypto.Hash([]byte("extensionchild"))
	data, err = serializeExtensionNode(en)
	require.NoError(t, err)
	expected = []byte{0x2, 0xa7, 0xa7, 0xc, 0x66, 0xad, 0xa, 0xc3, 0xef, 0xd6, 0x24, 0x4b, 0x78, 0x46, 0xbb, 0x4, 0x39, 0x28, 0xb9, 0xe2, 0xcf, 0xe0, 0x3e, 0x35, 0xa3, 0x91, 0x8e,
		0x83, 0xad, 0x36, 0x8, 0xb7, 0x5b, 0x58, 0x45, 0xe3, 0x9f, 0xeb, 0x59}
	require.Equal(t, expected, data)
	en2, err := deserializeExtensionNode(data)
	require.NoError(t, err)
	require.Equal(t, en, en2)
	en.sharedKey = []byte("extensionke")
	for i := range en.sharedKey {
		en.sharedKey[i] &= 0x0f
	}
	data, err = serializeExtensionNode(en)
	require.NoError(t, err)
	expected = []byte{0x1, 0xa7, 0xa7, 0xc, 0x66, 0xad, 0xa, 0xc3, 0xef, 0xd6, 0x24, 0x4b, 0x78, 0x46, 0xbb, 0x4, 0x39, 0x28, 0xb9, 0xe2, 0xcf, 0xe0, 0x3e, 0x35, 0xa3, 0x91, 0x8e,
		0x83, 0xad, 0x36, 0x8, 0xb7, 0x5b, 0x58, 0x45, 0xe3, 0x9f, 0xeb, 0x50}
	require.Equal(t, expected, data)
	en3, err := deserializeExtensionNode(data)
	require.NoError(t, err)
	require.Equal(t, en, en3)

	broken := []byte{0x6, 0xa7, 0xa7, 0xc, 0x66, 0xad, 0xa, 0xc3, 0xef, 0xd6, 0x24, 0x4b, 0x78, 0x46, 0xbb, 0x4, 0x39, 0x28, 0xb9, 0xe2, 0xcf, 0xe0, 0x3e, 0x35, 0xa3, 0x91, 0x8e,
		0x83, 0xad, 0x36, 0x8, 0xb7, 0x5b, 0x58, 0x45, 0xe3, 0x9f, 0xeb, 0x50}
	_, err = deserializeExtensionNode(broken)
	require.Error(t, err)
	expected = []byte{0x1, 0xa7, 0xa7, 0xc, 0x66, 0xad, 0xa, 0xc3, 0xef, 0xd6, 0x24, 0x4b, 0x78, 0x46, 0xbb, 0x4, 0x39, 0x28, 0xb9, 0xe2, 0xcf, 0xe0, 0x3e, 0x35, 0xa3, 0x91, 0x8e,
		0x83, 0xad, 0x36, 0x8, 0xb7, 0x5b, 0x58, 0x45, 0xe3, 0x9f, 0xeb, 0x50}
	_, err = deserializeLeafNode(expected)
	require.Error(t, err)
	_, err = deserializeBranchNode(expected)
	require.Error(t, err)
	_, err = deserializeRootNode(expected)
	require.Error(t, err)

}

func TestTrieAdd(t *testing.T) {
	mt, err := MakeTrie()
	require.NoError(t, err)
	verifyNewTrie(t, mt)
	mt.Add([]byte{0x01, 0x02, 0x03}, []byte{0x04, 0x05, 0x06})
	fmt.Println(mt.rootHash)
}

func TestNibbleUtilities(t *testing.T) {
	if false {
		fmt.Println("TestNibbleUtilities")
	}
	sampleNibbles := []nibbles{
		{0x0, 0x1, 0x2, 0x3, 0x4},
		{0x4, 0x1, 0x2, 0x3, 0x4},
		{0x0, 0x0, 0x2, 0x3, 0x5},
		{0x0, 0x1, 0x2, 0x3, 0x4, 0x5},
		{},
	}

	sampleNibblesPacked := [][]byte{
		{0x01, 0x23, 0x40},
		{0x41, 0x23, 0x40},
		{0x00, 0x23, 0x50},
		{0x01, 0x23, 0x45},
		{},
	}

	sampleNibblesShifted1 := []nibbles{
		{0x1, 0x2, 0x3, 0x4},
		{0x1, 0x2, 0x3, 0x4},
		{0x0, 0x2, 0x3, 0x5},
		{0x1, 0x2, 0x3, 0x4, 0x5},
		{},
	}

	sampleNibblesShifted2 := []nibbles{
		{0x2, 0x3, 0x4},
		{0x2, 0x3, 0x4},
		{0x2, 0x3, 0x5},
		{0x2, 0x3, 0x4, 0x5},
		{},
	}

	for i, n := range sampleNibbles {
		b, half, err := n.pack()
		require.NoError(t, err)
		if half {
			require.True(t, b[len(b)-1]&0x0f == 0x00)
		}
		require.True(t, bytes.Equal(b, sampleNibblesPacked[i]))

		unp, err := unpack(b, half)
		require.NoError(t, err)
		require.True(t, bytes.Equal(unp, n))

	}
	badNibbles := []nibbles{
		{0x12, 0x02, 0x03, 0x04},
	}
	_, _, err := badNibbles[0].pack()
	require.Error(t, err)

	for i, n := range sampleNibbles {
		require.True(t, bytes.Equal(shiftNibbles(n, 1), sampleNibblesShifted1[i]))
		require.True(t, bytes.Equal(shiftNibbles(n, 2), sampleNibblesShifted2[i]))
	}

	sampleSharedNibbles := [][]nibbles{
		{{0x0, 0x1, 0x2, 0x9, 0x2}, {0x0, 0x1, 0x2}},
		{{0x4, 0x1}, {0x4, 0x1}},
		{{0x9, 0x2, 0x3}, {}},
		{{0x0}, {0x0}},
		{{}, {}},
	}
	for i, n := range sampleSharedNibbles {
		shared := sharedNibbles(n[0], sampleNibbles[i])
		require.True(t, bytes.Equal(shared, n[1]))
		shared = sharedNibbles(sampleNibbles[i], n[0])
		require.True(t, bytes.Equal(shared, n[1]))
	}
	require.True(t, bytes.Equal(shiftNibbles(sampleNibbles[0], -2), sampleNibbles[0]))
	require.True(t, bytes.Equal(shiftNibbles(sampleNibbles[0], -1), sampleNibbles[0]))
	require.True(t, bytes.Equal(shiftNibbles(sampleNibbles[0], 0), sampleNibbles[0]))
}
