package main

import (
	"fmt"

	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/gob"
)

// This default DB is for doing stuff like saving "globals", if you need your

// ErrBucketNotFound indicates that a database file doesn't have it's string bucket.
var ErrBucketNotFound = fmt.Errorf("could not find strings bucket in database")

type counter int

func (m *machine) loadBoltWords() error {

	var err error
	m.db, err = storm.Open(".piscdb", storm.Codec(gob.Codec))
	if err != nil {
		return err
	}

	m.addGoWord("incr-counter", "( key -- newval )", GoWord(func(m *machine) error {
		key := m.popValue().String()
		tx, err := m.db.Begin(true)
		defer tx.Rollback()
		var c counter
		err = tx.Get("counter", key, &c)
		if err == storm.ErrNotFound {
			tx.Set("counter", key, 0)
			c = 0
			err = nil
		}
		if err != nil {
			return err
		}
		c++
		err = tx.Set("counter", key, c)
		if err != nil {
			return err
		}
		err = tx.Commit()
		if err != nil {
			return err
		}
		m.pushValue(Integer(c))
		return nil
	}))

	m.addGoWord("save-str", " ( key val -- ) ",
		GoWord(func(m *machine) error {
			val := m.popValue().String()
			key := m.popValue().String()
			return m.db.Set("strings", key, val)
		}))

	m.addGoWord("load-str", " ( key -- val ) ",
		GoWord(func(m *machine) error {
			key := m.popValue().String()
			var val string
			err := m.db.Get("strings", key, &val)
			if err != nil {
				return err
			}
			m.pushValue(String(val))
			return nil
		}))

	m.addGoWord("with-bucket", " ( str quot -- .. ) ", GoWord(func(m *machine) error {
		quot := m.popValue().(*quotation)
		bucketName := m.popValue().String()
		tx, err := m.db.Begin(true)
		if err != nil {
			return err
		}
		defer tx.Rollback()
		dict := make(Dict)

		// ( key val -- )
		dict["put-int"] = GoFunc(func(m *machine) error {
			val := m.popValue().(Integer)
			fmt.Println("value?", val)
			key := m.popValue().String()
			tx.Set(bucketName, key, int(val))
			return nil
		})

		// ( key -- val )
		dict["get-int"] = GoFunc(func(m *machine) error {
			var val int
			key := m.popValue().String()
			err := tx.Get(bucketName, key, &val)
			if err == storm.ErrNotFound {
				err = nil
				val = 0
			}
			if err != nil {
				return err
			}
			m.pushValue(Integer(val))
			return nil
		})

		// ( key val -- )
		dict["put-str"] = GoFunc(func(m *machine) error {
			val := m.popValue().String()
			key := m.popValue().String()
			return tx.Set(bucketName, key, val)
		})

		// ( key -- val )
		dict["get-str"] = GoFunc(func(m *machine) error {
			var val string
			key := m.popValue().String()
			err := tx.Get(bucketName, key, &val)
			if err == storm.ErrNotFound {
				err = nil
				val = ""
			}
			if err != nil {
				return err
			}
			m.pushValue(String(val))
			return nil
		})

		m.pushValue(dict)
		m.pushValue(quot)
		err = m.executeString("with", quot.toCode().codePositions[0])
		if err != nil {
			return err
		} else {
			tx.Commit()
			return nil
		}
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
