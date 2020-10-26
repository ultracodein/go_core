module go_core/go_lesson_04/main

go 1.15

replace gosearch/pkg/crawler/spider => ../../../crawler/pkg/spider

replace gosearch/pkg/crawler => ../../../crawler/cmd/crawler

replace gosearch/pkg/index => ../../../index/pkg/index

require (
	gosearch/pkg/crawler v0.0.0-00010101000000-000000000000
	gosearch/pkg/crawler/spider v0.0.0-00010101000000-000000000000
	gosearch/pkg/index v0.0.0-00010101000000-000000000000
)
