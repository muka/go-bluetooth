
.PHONY: gen

FILTER ?=

all: gen/clean gen/run

bluetoothd/logs:
	journalctl -u bluetooth -f

bluetoothd/start:
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
	FILTER=${FILTER} go run gen/srcgen/main.go

test/api:
	sudo go test github.com/muka/go-bluetooth/api

test/linux:
	sudo go test -v github.com/muka/go-bluetooth/linux/btmgmt

build:
	CGO_ENABLED=0 go build -o go-bluetooth ./main.go
