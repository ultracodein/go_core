module gosearch/pkg/cache

go 1.15

replace (
	gosearch/pkg/crawler => ../crawler
	gosearch/pkg/index => ../index
	gosearch/pkg/index/hash => ../index/hash
	gosearch/pkg/storage => ../storage
)

require (
	gosearch/pkg/crawler v0.0.0-00010101000000-000000000000
	gosearch/pkg/index v0.0.0-00010101000000-000000000000
	gosearch/pkg/index/hash v0.0.0-00010101000000-000000000000
	gosearch/pkg/storage v0.0.0-00010101000000-000000000000
)
