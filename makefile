assembly:
	go build -gcflags="-S" ./test/component-tests/example > dev.S 2>&1

.PHONY: test
test:
	go test ./... --cover -coverprofile=reports/coverage.out --covermode atomic --coverpkg=./...

show-coverage-report: test
	go tool cover -html=reports/coverage.out

benchmark:
	go test -benchmem -run=^$$ -bench . ./test/benchmarks/maps --cpuprofile ./reports/benchmark-cpu-maps.pprof --memprofile ./reports/benchmark-mem-maps.pprof -blockprofile ./reports/benchmark-block-maps.pprof
	go test -benchmem -run=^$$ -bench . ./test/benchmarks/sets --cpuprofile ./reports/benchmark-cpu-sets.pprof --memprofile ./reports/benchmark-mem-sets.pprof -blockprofile ./reports/benchmark-block-sets.pprof

pprof:
	go tool pprof --http=:8080 ./reports/benchmark-cpu-sets.pprof ./reports/benchmark-cpu-maps.pprof

pgo:
	go tool pprof -proto ./reports/benchmark-cpu-sets.pprof ./reports/benchmark-cpu-maps.pprof > default.pgo
