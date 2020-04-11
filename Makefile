
.PHONY: gen

BLUEZ_VERSION ?= 5.54
FILTER ?=

all: bluez/checkout gen/clean gen/run

bluez/checkout:
	git submodule init
	git submodule update
	cd src/bluez && git checkout ${BLUEZ_VERSION}

service/bluetoothd/logs:
	journalctl -u bluetooth -f

service/bluetoothd/start: bluetoothd/stop
	sudo bluetoothd -E -d -n -P hostname

service/bluetoothd/stop:
	sudo killall bluetoothd || true

gen/clean:
	rm `ls bluez/profile/*/gen_* -1` || true

gen/run:
	git submodule update
	FILTER=${FILTER} go run gen/srcgen/main.go full

gen: gen/run

test/api:
	sudo go test github.com/muka/go-bluetooth/api

build: gen
	CGO_ENABLED=0 go build -o go-bluetooth ./main.go

dev/cp: build
	ssh minion "killall go-bluetooth" || true
	scp go-bluetooth minion:~/
	ssh minion "~/go-bluetooth service server --adapterID hci1"

dev/logs:
	ssh minion "journalctl -u bluetooth.service -f"

docker/bluetoothd/build:
	docker build ./env/bluez --build-arg BLUEZ_VERSION=${BLUEZ_VERSION} -t opny/bluez-${BLUEZ_VERSION}

docker/bluetoothd/run: service/bluetoothd/stop
	docker run -it --rm --name bluez_${BLUEZ_VERSION} \
		--privileged \
		--net=host \
	  -v /dev:/dev \
		-v /var/run/dbus:/var/run/dbus \
		-v /sys/class/bluetooth:/sys/class/bluetooth \
		-v /var/lib/bluetooth:/var/lib/bluetooth \
		opny/bluez-${BLUEZ_VERSION}
