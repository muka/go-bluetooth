
.PHONY: gen

BLUEZ_VERSION ?= 5.50
FILTER ?=

all: bluez/checkout gen/clean gen/run

bluez/checkout:
	git submodule init
	git submodule update
	cd src/bluez && git checkout ${BLUEZ_VERSION}

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
	rm `ls bluez/profile/*/gen_* -1` || true

gen/run:
	git submodule update
	FILTER=${FILTER} go run gen/srcgen/main.go

gen: gen/run

test/api:
	sudo go test github.com/muka/go-bluetooth/api

test/linux:
	sudo go test -v github.com/muka/go-bluetooth/linux/btmgmt

build: gen
	CGO_ENABLED=0 go build -o go-bluetooth ./main.go

dev/cp: build
	scp go-bluetooth minion:~/
	ssh minion "~/go-bluetooth service server --adapterID hci1"
