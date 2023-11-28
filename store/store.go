package store

import (
	"encoding/binary"
	"errors"
	"fmt"

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

func (s *Store) Set(longUrl []byte) (tinyUrl uint64, err error) {
	var tinyUrlBytes []byte
	err = s.db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte{nextUnusedKey})
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			tinyUrlBytes = val
			tinyUrl = binary.BigEndian.Uint64(val)
			fmt.Printf("tinyUrl: %X\n", tinyUrl)
			return nil
		})
		if err != nil {
			return err
		}

		if err := txn.Set(tinyUrlBytes, longUrl); err != nil {
			return err
		}

		newValue := make([]byte, 8)
		binary.BigEndian.PutUint64(newValue, tinyUrl+1)
		return txn.Set([]byte{nextUnusedKey}, newValue)
	})
	fmt.Printf("tinyUrl: %X\n", tinyUrl)
	return tinyUrl, err
}

func (s *Store) Get(tinyUrl uint64) (longUrl []byte, err error) {
	err = s.db.View(func(txn *badger.Txn) error {
		tinyUrlBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(tinyUrlBytes, tinyUrl)

		item, err := txn.Get(tinyUrlBytes)
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
