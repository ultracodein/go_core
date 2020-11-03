package cache

// Кэш для временного хранения индекса и хранилища.

import (
	"gosearch/pkg/index"
	"gosearch/pkg/storage"
)

// Interface определяет контракт кэширующего модуля
type Interface interface {
	Create(*index.Interface, *storage.Interface) error
	Load() (*index.Interface, *storage.Interface, error)
}
