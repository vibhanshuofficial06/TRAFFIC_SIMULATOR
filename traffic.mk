# build & test automation
build:
	go build TrafficSimulator.go

test: build
	@echo Run1 - TrafficSimulator.go
	./TrafficSimulator