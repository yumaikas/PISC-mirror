package boltdb

import (
	"fmt"

	"pisc"

	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/gob"
)

var ModBoltDB = pisc.Module{
	Author:    "Andrew Owen",
	Name:      "BoltDBExt",
	License:   "MIT",
	DocString: "WIP: This module is for interfacing with BoltDB",
	Load:      loadBoltDB,
}

// This default DB is intended to act like a single file registry

// ErrBucketNotFound indicates that a database file doesn't have it's string bucket.
var ErrBucketNotFound = fmt.Errorf("could not find strings bucket in database")

type counter int

type kvDB (storm.DB)

type kvNode struct {
	node storm.Node
}

func (kv *kvDB) Close(m *pisc.Machine) error {
	return (*storm.DB)(kv).Close()
}

func (kv kvNode) WithTransaction(m *pisc.Machine) error {
	// Should be a quotation
	quot := m.PopValue()
	db := kv.node
	node, err := db.Begin(true)
	if err != nil {
		return err
	}
	transactionNode := kvNode{
		node: node,
	}
	kvWrapper := pisc.Dict{
		"save-int": pisc.GoFunc(transactionNode.SaveInt),
		"get-int":  pisc.GoFunc(transactionNode.GetInt),
		"save-str": pisc.GoFunc(transactionNode.SaveString),
		"get-str":  pisc.GoFunc(transactionNode.GetString),
	}
	m.PushValue(kvWrapper)
	m.PushValue(quot)
	err = m.ExecuteQuotation()
	if err != nil {
		node.Rollback()
		return err
	}
	return node.Commit()
}

// k v -- err?
func (kv kvNode) SaveInt(m *pisc.Machine) error {
	i := m.PopValue().(pisc.Integer)
	k := m.PopValue().String()
	err := kv.node.Set("GLOBAL", k, i)
	return err
}

func (kv kvNode) GetInt(m *pisc.Machine) error {
	var i int
	k := m.PopValue().String()
	err := kv.node.Get("GLOBAL", k, &i)
	m.PushValue(pisc.Integer(i))
	return err
}

// k v -- err?
func (kv kvNode) SaveString(m *pisc.Machine) error {
	v := m.PopValue().String()
	k := m.PopValue().String()
	err := kv.node.Set("GLOBAL", k, v)
	return err
}

func (kv kvNode) GetString(m *pisc.Machine) error {
	var s string
	k := m.PopValue().String()
	err := kv.node.Get("GLOBAL", k, &s)
	m.PushValue(pisc.String(s))
	return err
}

// This function builds a KV store based on Storm/BoltDB
// TODO: Look into making this not depend on Storm
func initKVStore(m *pisc.Machine) error {
	fileName := m.PopValue().String()
	db, err := storm.Open(fileName, storm.Codec(gob.Codec))
	if err != nil {
		return err
	}
	kv := kvNode{
		node: db,
	}

	kvWrapper := pisc.Dict{
		"close":    pisc.GoFunc((*kvDB)(db).Close),
		"save-int": pisc.GoFunc(kv.SaveInt),
		"get-int":  pisc.GoFunc(kv.GetInt),
		"save-str": pisc.GoFunc(kv.SaveString),
		"get-str":  pisc.GoFunc(kv.GetString),
	}
	m.PushValue(kvWrapper)
	return nil
}

func loadBoltDB(m *pisc.Machine) error {
	m.AddGoWord("<open-kv-at-path>", "( path -- kvstore )", initKVStore)
	return nil
}
