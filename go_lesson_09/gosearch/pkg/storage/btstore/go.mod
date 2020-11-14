module gosearch/pkg/storage/btstore

go 1.15

replace gosearch/pkg/crawler => ../../crawler

replace gosearch/pkg/storage/btstore/btree => ./btree

require (
	gosearch/pkg/crawler v0.0.0-00010101000000-000000000000
	gosearch/pkg/storage/btstore/btree v0.0.0-00010101000000-000000000000
)
