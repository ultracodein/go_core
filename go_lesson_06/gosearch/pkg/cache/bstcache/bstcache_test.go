package bstcache

import (
	"gosearch/pkg/index"
	"gosearch/pkg/index/hash"
	"gosearch/pkg/storage"
	"gosearch/pkg/storage/bststore"
	"reflect"
	"testing"
)

func TestCacher_CreateAndLoad(t *testing.T) {
	cacheFiles := map[string]string{
		"index":   "index.bin",
		"storage": "storage.bin",
	}
	cache := New(cacheFiles)
	var index index.Interface = hash.New()
	var storage storage.Interface = bststore.New()

	err := cache.Create(&index, &storage)
	if err != nil {
		t.Fatalf(err.Error())
	}
	indexLoaded, _, err := cache.Load()
	if err != nil {
		t.Fatalf(err.Error())
	}
	if !reflect.DeepEqual(index, *indexLoaded) {
		t.Fatalf("загруженный индекс отличается от исходного")
	}
}
