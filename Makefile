
.PHONY: gen

DEBUG ?= 0

all: gen/clean gen/run

bluetoothd:
	sudo killall bluetoothd || true && \
	sudo bluetoothd -E -d -n -P hostname

run/example/service:
	go run examples/service/*.go

run/example/client:
	go run examples/service/*.go client

gen/clean:
	rm `ls bluez/profile/*/gen_* -1`

gen:
	git submodule update
	DEBUG=${DEBUG} go run gen/srcgen/main.go

test/switch:
	sudo go test api/switch*
