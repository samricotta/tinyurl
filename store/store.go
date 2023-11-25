package store

import (
	"encoding/binary"

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

		var newValue []byte
		binary.BigEndian.PutUint64(newValue, binary.BigEndian.Uint64(tinyUrl)+1)
		return txn.Set([]byte{nextUnusedKey}, newValue)
	})
	return tinyUrl, err
}

func (s *Store) Get(tinyUrl uint64) (longUrl []byte, err error) {
	err = s.db.View(func(txn *badger.Txn) error {
		var tinyUrlBytes []byte
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
