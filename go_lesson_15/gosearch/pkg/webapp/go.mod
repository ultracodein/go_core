module gosearch/pkg/webapp

go 1.14

replace (
	gosearch/pkg/crawler => ../crawler
	gosearch/pkg/index => ../index
	gosearch/pkg/index/hash => ../index/hash
	gosearch/pkg/storage => ../storage
	gosearch/pkg/storage/btstore => ../storage/btstore
	gosearch/pkg/storage/btstore/btree => ../storage/btstore/btree
)

require (
	gosearch/pkg/crawler v0.0.0-00010101000000-000000000000
	gosearch/pkg/index v0.0.0-00010101000000-000000000000
	gosearch/pkg/index/hash v0.0.0-00010101000000-000000000000
	gosearch/pkg/storage v0.0.0-00010101000000-000000000000
	gosearch/pkg/storage/btstore v0.0.0-00010101000000-000000000000
	gosearch/pkg/storage/btstore/btree v0.0.0-00010101000000-000000000000
)
