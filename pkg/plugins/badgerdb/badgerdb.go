package badgerdb

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/egevorkyan/flufik/pkg/logging"
)

type FlufikBadger struct {
	badgerDB *badger.DB
}

func NewFlufikBadgerDB(dbName string) *FlufikBadger {
	var b FlufikBadger
	opts := badger.DefaultOptions(dbName)
	opts.Logger = nil
	db, err := badger.Open(opts)
	if err != nil {
		logging.ErrorHandler("fatal: ", err)
	}
	b.badgerDB = db
	return &b
}

func (b *FlufikBadger) UpdateDb(kv map[string]string) error {
	txn := b.badgerDB.NewTransaction(true)
	defer txn.Discard()
	for k, v := range kv {
		if err := txn.Set([]byte(k), []byte(v)); err != nil {
			return err
		}
	}
	if err := txn.Commit(); err != nil {
		return err
	}
	return nil
}

func (b *FlufikBadger) Get(k string) ([]byte, error) {
	txn := b.badgerDB.NewTransaction(true)
	defer txn.Discard()
	item, err := txn.Get([]byte(k))
	if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (b *FlufikBadger) Remove(k string) error {
	txn := b.badgerDB.NewTransaction(true)
	defer txn.Discard()
	if err := txn.Delete([]byte(k)); err != nil {
		return err
	}
	if err := txn.Commit(); err != nil {
		return err
	}
	return nil
}

func (b *FlufikBadger) Close() {
	b.badgerDB.Close()
}
