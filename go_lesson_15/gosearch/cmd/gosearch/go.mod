module gosearch/cmd/gosearch

go 1.14

replace (
	gosearch/pkg/crawler => ../../pkg/crawler
	gosearch/pkg/crawler/membot => ../../pkg/crawler/membot
	gosearch/pkg/crawler/spider => ../../pkg/crawler/spider
	gosearch/pkg/engine => ../../pkg/engine
	gosearch/pkg/index => ../../pkg/index
	gosearch/pkg/index/hash => ../../pkg/index/hash
	gosearch/pkg/storage => ../../pkg/storage
	gosearch/pkg/storage/btstore => ../../pkg/storage/btstore
	gosearch/pkg/storage/btstore/btree => ../../pkg/storage/btstore/btree
	gosearch/pkg/webapp => ../../pkg/webapp
)

require (
	gosearch/pkg/crawler v0.0.0-00010101000000-000000000000
	gosearch/pkg/crawler/spider v0.0.0-00010101000000-000000000000
	gosearch/pkg/engine v0.0.0-00010101000000-000000000000
	gosearch/pkg/index v0.0.0-00010101000000-000000000000
	gosearch/pkg/index/hash v0.0.0-00010101000000-000000000000
	gosearch/pkg/storage v0.0.0-00010101000000-000000000000
	gosearch/pkg/storage/btstore v0.0.0-00010101000000-000000000000
	gosearch/pkg/storage/btstore/btree v0.0.0-00010101000000-000000000000 // indirect
	gosearch/pkg/webapp v0.0.0-00010101000000-000000000000
)
