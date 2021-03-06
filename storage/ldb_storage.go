package storage

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type leveldbStorage struct {
	db *leveldb.DB
}

func OpenLeveldbStorage(path string) (Storage, error) {

	if db, err := leveldb.OpenFile(path, &opt.Options{Filter: filter.NewBloomFilter(1000)}); err != nil {
		return nil, err
	} else {
		return &leveldbStorage{db}, nil
	}
}

func (s *leveldbStorage) WALName() string {
	return "null"
}

func (s *leveldbStorage) Set(k, v []byte) error {
	return s.db.Put(k, v, nil)
}

func (s *leveldbStorage) Get(k []byte) ([]byte, error) {
	return s.db.Get(k, nil)
}

func (s *leveldbStorage) Delete(k []byte) error {
	return s.db.Delete(k, nil)
}

func (s *leveldbStorage) ForEach(fn func(k, v []byte) error) error {
	iter := s.db.NewIterator(nil, nil)
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		key := iter.Key()
		value := iter.Value()
		if err := fn(key, value); err != nil {
			return err
		}
	}
	iter.Release()
	return iter.Error()
}

func (s *leveldbStorage) Close() error {
	return s.db.Close()
}
