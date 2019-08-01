
.PHONY: gen

DEBUG ?= 0

all: gen/clean gen/run

bluetoothd:
	sudo killall bluetoothd && \
	sudo bluetoothd -Edn P hostname

run/example/service:
	go run examples/service/*.go

run/example/client:
	go run examples/service/*.go client

gen/dev/clean:
	rm -rf src/gen
	mkdir -p src/gen

gen/dev/run:
	DEBUG=1 make gen

gen/dev: gen/dev/clean gen/dev/run

gen:
	git submodule update
	DEBUG=${DEBUG} go run gen/srcgen/main.go

test/switch:
	sudo go test api/switch*
