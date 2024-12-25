assembly:
	go tool compile -S x.go

.PHONY: test
test:
	go test ./... --cover -coverprofile=reports/coverage.out --covermode atomic --coverpkg=./...

show-coverage-report: test
	go tool cover -html=reports/coverage.out
