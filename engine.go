package bloom

import (
    "github.com/syndtr/goleveldb/leveldb"
    "github.com/syndtr/goleveldb/leveldb/filter"
    "github.com/syndtr/goleveldb/leveldb/opt"
    "path"
    "sync"
)

type (
    Engine struct {
        mu      sync.RWMutex
        buckets map[string]*Bucket
        baseDir string
    }
)

func New(baseDir string) *Engine {
    e := &Engine{
        buckets: make(map[string]*Bucket),
        baseDir: baseDir,
    }
    return e
}
func (e *Engine) Bucket(name string) (*Bucket, error) {
    e.mu.Lock()
    defer e.mu.Unlock()
    // check exist
    b, ok := e.buckets[name]
    if ok {
        return b, nil
    }
    // new
    db, err := leveldb.OpenFile(path.Join(e.baseDir, name), &opt.Options{
        Filter: filter.NewBloomFilter(10),
    })
    if err != nil {
        return nil, err
    }
    b = &Bucket{
        db: db,
    }
    return b, nil
}
func (e *Engine) CloseBucket(name string) error {
    e.mu.Lock()
    defer e.mu.Unlock()
    if b, ok := e.buckets[name]; ok {
        if err := b.close(); err == nil {
            delete(e.buckets, name)
            return nil
        } else {
            return err
        }
    }
    return nil
}
func (e *Engine) CloseAll() {
    e.onStop()
}
func (e *Engine) onStop() {
    e.mu.Lock()
    defer e.mu.Unlock()
    for k, b := range e.buckets {
        if err := b.close(); err != nil {
        }
        delete(e.buckets, k)
    }
    e.buckets = map[string]*Bucket{}
}