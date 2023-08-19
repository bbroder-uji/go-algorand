// Copyright (C) 2018-2023 Algorand, Inc.
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

package statetrie

import (
	"errors"
	"fmt"
	"github.com/algorand/go-algorand/crypto"
	"github.com/cockroachdb/pebble"
	"runtime"
)

const (
	MaxKeyLength = 65535
)

var debugTrie = false

type Trie struct {
	db     *pebble.DB
	parent *Trie
	root   node
	dels   map[string]node
}

// MakeTrie creates a merkle trie
func MakeTrie(db *pebble.DB) (*Trie, error) {
	mt := &Trie{db: db}
	mt.dels = make(map[string]node)
	return mt, nil
}

func (mt *Trie) Close() error {
	return mt.db.Close()
}

// Trie Add adds the given key/value pair to the trie.
func (mt *Trie) Add(key nibbles, value []byte) (err error) {
	if len(key) == 0 {
		return errors.New("empty key not allowed")
	}

	if len(key) > MaxKeyLength {
		return errors.New("key too long")
	}
	if debugTrie {
		fmt.Printf("Add: %x : %v\n", key, crypto.Hash(value))
	}

	if mt.root == nil {
		stats.cryptohashes++
		stats.newrootnode++
		mt.root = makeLeafNode(key, crypto.Hash(value), nibbles{})
		mt.addNode(mt.root)
		return nil
	}

	stats.cryptohashes++
	replacement, err := mt.root.descendAdd(mt, nibbles{}, key, crypto.Hash(value))
	if err != nil {
		return err
	}
	stats.newrootnode++
	mt.root = replacement
	return nil
}

// Delete deletes the given key from the trie, if such element exists.
// if no such element exists, return false
func (mt *Trie) Delete(key nibbles) (bool, error) {
	var err error
	if len(key) == 0 {
		return false, errors.New("empty key not allowed")
	}
	if len(key) > MaxKeyLength {
		return false, errors.New("key too long")
	}
	if mt.root == nil {
		return false, nil
	}

	replacement, found, err := mt.root.descendDelete(mt, nibbles{}, key)
	if err == nil && found {
		mt.root = replacement
	}
	return found, err
}

func (mt *Trie) SetRoot(root crypto.Digest) {
	rootdigest := new(crypto.Digest)
	*rootdigest = root
	mt.root = makeDBNode(rootdigest, nibbles{})
}

// Provide the root hash for this trie
func (mt *Trie) RootHash() (crypto.Digest, error) {
	if mt.root == nil {
		return crypto.Digest{}, nil
	}
	if mt.root.getHash() != nil {
		return *(mt.root.getHash()), nil
	}

	err := mt.root.descendHash()
	if err != nil {
		return crypto.Digest{}, err
	}

	return *(mt.root.getHash()), nil
}

func (mt *Trie) child() *Trie {
	ch := &Trie{root: mt.root, parent: mt}
	ch.dels = make(map[string]node)
	return ch
}

func (mt *Trie) merge() error {
	if mt.parent == nil {
		return nil
	}
	mt.parent.root = mt.root
	for k, v := range mt.dels {
		mt.parent.dels[k] = v
	}
	return nil
}

func (mt *Trie) Commit() error {
	if mt.parent != nil {
		mt.merge()
		return mt.parent.Commit()
	}

	b := mt.db.NewBatch()
	if mt.root != nil {
		err := mt.root.descendHashWithCommit(b)
		if err != nil {
			return err
		}
	}
	for k := range mt.dels {
		options := &pebble.WriteOptions{}
		err := b.Delete([]byte(k), options)
		if err != nil {
			return err
		}
	}
	//err := b.SyncWait()
	options := &pebble.WriteOptions{}
	err := mt.db.Apply(b, options)
	if err != nil {
		return err
	}
	err = b.Close()
	if err != nil {
		return err
	}

	mt.dels = make(map[string]node)

	//	mt.root = makeDBNode(mt.root.getHash(), nibbles{})
	//    mt.root = nil
	return nil
}

func (mt *Trie) countNodes() string {
	if mt.root == nil {
		return "Empty trie"
	}
	var nc struct {
		branches int
		leaves   int
		exts     int
		dbnodes  int
	}

	count := func() func(n node) {
		innerCount := func(n node) {
			switch n.(type) {
			case *BranchNode:
				nc.branches++
			case *LeafNode:
				nc.leaves++
			case *ExtensionNode:
				nc.exts++
			case *DBNode:
				nc.dbnodes++
			}
		}
		return innerCount
	}()
	mt.root.lambda(count)

	var nmem struct {
		branches int
		leaves   int
		exts     int
		dbnodes  int
	}

	mem := func() func(n node) {
		innerCount := func(n node) {
			switch n.(type) {
			//estimates
			case *BranchNode:
				nmem.branches += 16*16 + 32 + 24 + len(n.(*BranchNode).key) + 8 + 32
			case *LeafNode:
				nmem.leaves += 24 + len(n.(*LeafNode).key) + 24 + len(n.(*LeafNode).keyEnd) + 32 + 8 + 32
			case *ExtensionNode:
				nmem.exts += 24 + len(n.(*ExtensionNode).key) + 24 + len(n.(*ExtensionNode).sharedKey) + 8 + 32
			case *DBNode:
				nmem.dbnodes += 24 + len(n.(*DBNode).key) + 8 + 32
			}
		}
		return innerCount
	}()
	mt.root.lambda(mem)

	return fmt.Sprintf("[nodes: total %d (branches: %d, leaves: %d, exts: %d, dbnodes: %d), mem: total %d (branches: %d, leaves: %d, exts: %d, dbnodes: %d)]",
		nc.branches+nc.leaves+nc.exts+nc.dbnodes,
		nc.branches, nc.leaves, nc.exts, nc.dbnodes,
		nmem.branches+nmem.leaves+nmem.exts+nmem.dbnodes,
		nmem.branches, nmem.leaves, nmem.exts, nmem.dbnodes)

}

func (mt *Trie) getNode(dbn dbnode) (node, error) {
	stats.getnode++
	dbKey := dbn.getDBKey()
	if dbKey != nil {
		if _, ok := mt.dels[string(dbKey)]; ok {
			return nil, nil // return nil if node is scheduled for deletion
		}

		stats.dbgets++

		dbbytes, closer, err := mt.db.Get(dbKey)
		if err != nil {
			fmt.Printf("\ndbKey panic: dbkey %x\n", dbKey)
			panic(err)
		}
		defer closer.Close()

		// this node may have dbnodes in it, but it is not a dbnode.
		trieNode, err := deserializeNode(dbbytes, dbn.getKey())
		if err != nil {
			panic(err)
		}
		if debugTrie {
			fmt.Printf("getNode %T (%x) : (%v)\n", trieNode, dbKey, trieNode)
		}
		return trieNode, nil
	}
	return nil, nil
}
func (mt *Trie) delNode(n node) {
	stats.delnode++
	dbKey := n.getDBKey()
	if dbKey != nil {
		// note this key for deletion on commit
		if debugTrie {
			fmt.Printf("delNode %T (%x) : (%v)\n", n, dbKey, n)
		}
		mt.dels[string(dbKey)] = n
	}
}
func (mt *Trie) addNode(n node) {
	stats.addnode++
	key := string(n.getKey())
	_, _, lineno, _ := runtime.Caller(1)
	if debugTrie {
		fmt.Printf("addNode %T (%x) : (%v) caller %d\n", n, key, n, lineno)
	}
	// this key is no longer to be deleted, as it will be overwritten.
	delete(mt.dels, key)
}

// Make a dot graph of the trie
func (mt *Trie) DotGraph(keysAdded [][]byte, valuesAdded [][]byte) string {
	var keys string
	for i := 0; i < len(keysAdded); i++ {
		keys += fmt.Sprintf("%x = %x\\n", keysAdded[i], valuesAdded[i])
	}
	fmt.Printf("root: %v\n", mt.root)
	return fmt.Sprintf("digraph trie { key [shape=box, label=\"key/value inserted:\\n%s\"];\n %s }\n", keys, mt.dotGraph(mt.root, nibbles{}))
}
func (mt *Trie) dotGraph(n node, path nibbles) string {

	switch tn := n.(type) {
	case *DBNode:
		n2, _ := mt.getNode(tn)
		return mt.dotGraph(n2, path)
	case *LeafNode:
		ln := tn
		return fmt.Sprintf("n%p [label=\"leaf\\nkeyEnd:%x\\nvalueHash:%s\" shape=box];\n", tn, ln.keyEnd, ln.valueHash)
	case *ExtensionNode:
		en := tn
		return fmt.Sprintf("n%p [label=\"extension\\nshKey:%x\" shape=box];\n", tn, en.sharedKey) +
			fmt.Sprintf("n%p -> n%p;\n", en, en.child) +
			mt.dotGraph(en.child, append(path, en.sharedKey...))
	case *BranchNode:
		bn := tn
		var indexesFilled string
		indexesFilled = "--"
		for i, child := range bn.children {
			if child != nil {
				indexesFilled += fmt.Sprintf("%x ", i)
			}
		}
		indexesFilled += "--"

		s := fmt.Sprintf("n%p [label=\"branch\\nindexesFilled:%s\\nvalueHash:%s\" shape=box];\n", tn, indexesFilled, bn.valueHash)
		for _, child := range bn.children {
			if child != nil {
				s += fmt.Sprintf("n%p -> n%p;\n", tn, child)
			}
		}
		for childrenIndex, child := range bn.children {
			if child != nil {
				s += mt.dotGraph(child, append(path, byte(childrenIndex)))
			}
		}
		return s
	default:
		return ""
	}
}
