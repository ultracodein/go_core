package bstcache

import (
	"encoding/gob"
	"os"

	"gosearch/pkg/crawler"
	"gosearch/pkg/index"
	"gosearch/pkg/index/hash"
	"gosearch/pkg/storage"
	"gosearch/pkg/storage/bststore"
)

// Cacher - структура, реализующая контракт кэширующего модуля
type Cacher struct {
	indexFile, storageFile string
}

// New - конструктор.
func New(cacheFiles map[string]string) *Cacher {
	c := Cacher{
		indexFile:   cacheFiles["index"],
		storageFile: cacheFiles["storage"],
	}
	return &c
}

// Create сохраняет индекс и хранилище в файлы
func (c *Cacher) Create(index *index.Interface, storage *storage.Interface) error {
	indexFile, err := os.Create(c.indexFile)
	if err != nil {
		return err
	}
	defer indexFile.Close()
	gob.Register(&hash.Index{})
	gob.Register(&crawler.Document{})
	enc := gob.NewEncoder(indexFile)
	err = enc.Encode(index)
	if err != nil {
		return err
	}

	storageFile, err := os.Create(c.storageFile)
	if err != nil {
		return err
	}
	defer storageFile.Close()
	gob.Register(&bststore.DB{})
	enc = gob.NewEncoder(storageFile)
	err = enc.Encode(storage)
	if err != nil {
		return err
	}

	return nil
}

// Load загружает индекс и хранилище из файлов
func (c *Cacher) Load() (*index.Interface, *storage.Interface, error) {
	indexFile, err := os.Open(c.indexFile)
	if err != nil {
		return nil, nil, err
	}
	defer indexFile.Close()
	var indexData index.Interface
	gob.Register(&hash.Index{})
	gob.Register(&crawler.Document{})
	dec := gob.NewDecoder(indexFile)
	err = dec.Decode(&indexData)
	if err != nil {
		return nil, nil, err
	}

	storageFile, err := os.Open(c.storageFile)
	if err != nil {
		return nil, nil, err
	}
	defer storageFile.Close()
	var storageData storage.Interface
	gob.Register(&bststore.DB{})
	dec = gob.NewDecoder(storageFile)
	err = dec.Decode(&storageData)
	if err != nil {
		return nil, nil, err
	}

	return &indexData, &storageData, nil
}
