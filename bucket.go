package bloom

import (
    "github.com/syndtr/goleveldb/leveldb"
    "github.com/syndtr/goleveldb/leveldb/util"
    "log"
)

type (
    Bucket struct {
        db *leveldb.DB
    }
)

func (b *Bucket) close() error{
    return b.db.Close()
}
func (b *Bucket) AddString(key string) error {
    return b.db.Put([]byte(key), nil, nil)
}

func (b *Bucket) AddStrings(keys ...string) error {
    bt := &leveldb.Batch{}
    for _, key := range keys {
        bt.Put([]byte(key), nil)
    }
    return b.db.Write(bt, nil)
}
func (b *Bucket) Exist(key string) bool {
    _, err := b.db.Get([]byte(key), nil)
    if err == nil {
        return true
    }
    return false
}

func (b *Bucket) Count(from, to []byte) (n uint32) {
    it := b.db.NewIterator(&util.Range{Start: from, Limit: to}, nil)
    defer it.Release()

    for it.Next() {
        log.Println("form:", string(from), string(it.Key()), it.Value())
        n++
    }

    return
}

func (b *Bucket) RemoveString(key string) error {
    return b.db.Delete([]byte(key), nil)
}
