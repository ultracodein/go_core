module gosearch/pkg/cache/btcache

go 1.15

replace (
	gosearch/pkg/crawler => ../../crawler
	gosearch/pkg/index => ../../index
	gosearch/pkg/index/hash => ../../index/hash
	gosearch/pkg/storage => ../../storage
	gosearch/pkg/storage/btstore => ../../storage/btstore
)

require (
	gosearch/pkg/crawler v0.0.0-00010101000000-000000000000
	gosearch/pkg/index v0.0.0-00010101000000-000000000000
	gosearch/pkg/index/hash v0.0.0-00010101000000-000000000000
	gosearch/pkg/storage v0.0.0-00010101000000-000000000000
	gosearch/pkg/storage/btstore v0.0.0-00010101000000-000000000000
)
