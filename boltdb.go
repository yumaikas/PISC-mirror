package main

import (
	"fmt"

	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/gob"
)

// This default DB is for doing stuff like saving "globals", if you need your
var defaultDB *storm.DB

// ErrBucketNotFound indicates that a database file doesn't have it's string bucket.
var ErrBucketNotFound = fmt.Errorf("could not find strings bucket in database")

func (m *machine) loadBoltWords() error {

	var err error
	var strKey = []byte("strings")
	defaultDB, err = storm.Open(".piscdb", storm.Codec(gob.Codec))
	if err != nil {
		return err
	}

	m.addGoWord("save-str", " ( key val -- ) ",
		GoWord(func(m *machine) error {
			val := m.popValue().String()
			key := m.popValue().String()
			return defaultDB.Set("strings", key, val)
		}))
	m.addGoWord("load-str", " ( key -- val ) ",
		GoWord(func(m *machine) error {
			key := m.popValue().String()
			var val *string
			err := defaultDB.Get("strings", key, val)
			if err != nil {
				return err
			}
			m.pushValue(String(*val))
			return nil
		}))

	// TODO: Implement individual DBs..
	/*
			m.addGoWord("<bolt-db>", `( filepath -- db )
		    A bolt db has the following words available to it:

		    - set-str (  )
		    `, GoWord(func(m *machine) error {
				path := m.popValue().(String).String()
				db, err := bolt.Open(path, 0600, nil)
				db.Update(func(tx *bolt.Tx) error {
					_, err := tx.CreateBucketIfNotExists([]byte("strings"))
					return err
				})
				if err != nil {
					return err
				}

				strKey := []byte("strings")
				stackDB := Dict{}
				stackDB["set-str"] = GoFunc(func(m *machine) error {
					val, ok := m.popValue().(String)
					key := m.popValue().String()
					if !ok {
						return err
					}
					db.Update(func(tx *bolt.Tx) error {
						b := tx.Bucket(strKey)
						if b == nil {
							return fmt.Errorf("unable to load strings bucket")
						}
						b.Put([]byte(key), []byte(val))
						return nil
					})
					return nil
				})
				stackDB["get-str"] = GoFunc(func(m *machine) error {
					key := m.popValue().String()
					db.View(func(tx *bolt.Tx) error {
						b := tx.Bucket(strKey)
						if b == nil {
							return fmt.Errorf("unable to load strings bucket")
						}
						val := b.Get([]byte(key))
						export := make([]byte, len(val))
						copy(val, export)
						m.pushValue(String(string(export)))
						return nil
					})
					return nil
				})

				stackDB["close"] = GoFunc(func(m *machine) error {
					return db.Close()
				})

				return nil
			}))
	*/

	return nil
}
