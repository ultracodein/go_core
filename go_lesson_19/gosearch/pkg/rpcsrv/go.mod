module gosearch/pkg/rpcsrv

replace (
	gosearch/pkg/crawler => ../crawler
	gosearch/pkg/engine => ../engine
	gosearch/pkg/index => ../index
	gosearch/pkg/index/hash => ../index/hash
	gosearch/pkg/storage => ../storage
	gosearch/pkg/storage/memstore => ../storage/memstore
)

go 1.15

require (
	github.com/gorilla/mux v1.8.0
	gosearch/pkg/crawler v0.0.0-00010101000000-000000000000
	gosearch/pkg/engine v0.0.0-00010101000000-000000000000
	gosearch/pkg/index v0.0.0-00010101000000-000000000000
	gosearch/pkg/storage v0.0.0-00010101000000-000000000000
)
