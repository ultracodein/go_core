package btcache

import (
	"gosearch/pkg/index"
	"gosearch/pkg/index/hash"
	"gosearch/pkg/storage"
	"gosearch/pkg/storage/btstore"
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
	var storage storage.Interface = btstore.New()

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
