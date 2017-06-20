package main

import (
	"fmt"

	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/gob"
)

var ModBoltDB = PISCModule{
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

func (kv *kvDB) Close(m *machine) error {
	return (*storm.DB)(kv).Close()
}

func (kv kvNode) WithTransaction(m *machine) error {
	quot := m.popValue().(quotation)
	db := kv.node
	node, err := db.Begin(true)
	if err != nil {
		return err
	}
	transactionNode := kvNode{
		node: node,
	}
	kvWrapper := Dict{
		"save-int": GoFunc(transactionNode.SaveInt),
		"get-int":  GoFunc(transactionNode.GetInt),
		"save-str": GoFunc(transactionNode.SaveString),
		"get-str":  GoFunc(transactionNode.GetString),
	}
	m.pushValue(kvWrapper)
	m.pushValue(quot)
	err = m.executeQuotation()
	if err != nil {
		node.Rollback()
		return err
	}
	return node.Commit()
}

// k v -- err?
func (kv kvNode) SaveInt(m *machine) error {
	i := m.popValue().(Integer)
	k := m.popValue().String()
	err := kv.node.Set("GLOBAL", k, i)
	return err
}

func (kv kvNode) GetInt(m *machine) error {
	var i int
	k := m.popValue().String()
	err := kv.node.Get("GLOBAL", k, &i)
	m.pushValue(Integer(i))
	return err
}

// k v -- err?
func (kv kvNode) SaveString(m *machine) error {
	v := m.popValue().String()
	k := m.popValue().String()
	err := kv.node.Set("GLOBAL", k, v)
	return err
}

func (kv kvNode) GetString(m *machine) error {
	var s string
	k := m.popValue().String()
	err := kv.node.Get("GLOBAL", k, &s)
	m.pushValue(String(s))
	return err
}

// This function builds a KV store based on Storm/BoltDB
// TODO: Look into making this not depend on Storm
func initKVStore(m *machine) error {
	fileName := m.popValue().String()
	db, err := storm.Open(fileName, storm.Codec(gob.Codec))
	if err != nil {
		return err
	}
	kv := kvNode{
		node: db,
	}

	kvWrapper := Dict{
		"close":    GoFunc((*kvDB)(db).Close),
		"save-int": GoFunc(kv.SaveInt),
		"get-int":  GoFunc(kv.GetInt),
		"save-str": GoFunc(kv.SaveString),
		"get-str":  GoFunc(kv.GetString),
	}
	m.pushValue(kvWrapper)
	return nil
}

func loadBoltDB(m *machine) error {
	m.addGoWord("<open-kv-at-path>", "( path -- kvstore )", initKVStore)
	return nil
}
