package store

import (
	"encoding/binary"
	"errors"

	"github.com/dgraph-io/badger/v4"
)

const nextUnusedKey = byte(0x00)

type Store struct {
	db *badger.DB
}

func New(path string) (*Store, error) {
	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		return nil, err
	}
	err = db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte{nextUnusedKey})
		if errors.Is(err, badger.ErrKeyNotFound) {
			// do something here
			firstKey := make([]byte, 8)
			binary.BigEndian.PutUint64(firstKey, 1)
			return txn.Set([]byte{nextUnusedKey}, firstKey)
		}
		return err

	})
	if err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Set(longUrl []byte) (tinyUrl []byte, err error) {
	err = s.db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte{nextUnusedKey})
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			tinyUrl = val
			return nil
		})
		if err != nil {
			return err
		}

		if err := txn.Set(tinyUrl, longUrl); err != nil {
			return err
		}

		newValue := make([]byte, 8)
		binary.BigEndian.PutUint64(newValue, binary.BigEndian.Uint64(tinyUrl)+1)
		return txn.Set([]byte{nextUnusedKey}, newValue)
	})
	return tinyUrl, err
}

func (s *Store) Get(tinyUrl []byte) (longUrl []byte, err error) {
	err = s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(tinyUrl)
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			longUrl = val
			return nil
		})
	})
	return longUrl, err
}
