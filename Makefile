
all: gen/clean gen/run

bluetoothd:
	sudo killall bluetoothd && \
	sudo bluetoothd -Edn P hostname

run/example/service:
	go run examples/service/*.go

run/example/client:
	go run examples/service/*.go client

gen/clean:
	rm -rf src/gen
	mkdir -p src/gen

gen/run:
	git submodule update
	go run gen/srcgen/main.go

gen: gen/clean gen/run

test/switch:
	sudo go test api/switch*
