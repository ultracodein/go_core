module gosearch/pkg/index

go 1.15

replace gosearch/pkg/crawler/spider => ../../../crawler/pkg/spider

replace gosearch/pkg/crawler => ../../../crawler/cmd/crawler

require (
	gosearch/pkg/crawler v0.0.0-00010101000000-000000000000
	gosearch/pkg/crawler/spider v0.0.0-00010101000000-000000000000
)
