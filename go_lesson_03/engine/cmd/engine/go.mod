module go_core/go_lesson_02/main

go 1.15

replace crawler/pkg/spider => ../../../crawler/pkg/spider

replace engine/pkg/scanners/base => ../../pkg/scanners/base

replace engine/pkg/scanners/mock => ../../pkg/scanners/mock

require (
	crawler/pkg/spider v0.0.0-00010101000000-000000000000
	engine/pkg/scanners/base v0.0.0-00010101000000-000000000000
	engine/pkg/scanners/mock v0.0.0-00010101000000-000000000000
)
