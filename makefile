assembly:
	go build -gcflags="-S" ./test/component-tests/large > dev.S 2>&1

.PHONY: test
test:
	go test ./... --cover -coverprofile=reports/coverage.out --covermode atomic --coverpkg=./...

show-coverage-report: test
	go tool cover -html=reports/coverage.out

benchmark:
	go test -benchmem -run=^$$ -bench . ./test/benchmarks/large --cpuprofile ./reports/benchmark-cpu.pprof --memprofile ./reports/benchmark-mem.pprof -blockprofile ./reports/benchmark-block.pprof

pprof:
	go tool pprof --http=:8080 ./reports/benchmark-cpu.pprof

pgo:
	go tool pprof -proto ./reports/benchmark-cpu.pprof > default.pgo
