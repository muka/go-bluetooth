
.PHONY: gen

BLUEZ_VERSION ?= 5.50
FILTER ?=
DEV_HOST ?= minion

DOCKER_PARAMS :=  --privileged -it --rm \
									--net=host \
								  -v /dev:/dev \
									-v /var/run/dbus:/var/run/dbus \
									-v /sys/class/bluetooth:/sys/class/bluetooth \
									-v /var/lib/bluetooth:/var/lib/bluetooth \
									opny/bluez-${BLUEZ_VERSION}

all: bluez/checkout gen/clean gen/run

bluez/init:
	git submodule init
	git submodule update

bluez/checkout:
	cd src/bluez && git fetch && git checkout ${BLUEZ_VERSION}

service/bluetoothd/logs:
	journalctl -u bluetooth -f

service/bluetoothd/start: service/bluetoothd/stop
	sudo bluetoothd -E -d -n -P hostname

service/bluetoothd/stop:
	sudo killall bluetoothd || true

gen/clean:
	rm -f `ls bluez/profile/*/gen_* -1` || true

gen/run: bluez/checkout
	BLUEZ_VERSION=${BLUEZ_VERSION} FILTER=${FILTER} go run gen/srcgen/main.go full

gen: gen/run

test/api:
	sudo go test github.com/muka/go-bluetooth/api

build: gen
	CGO_ENABLED=0 go build -o go-bluetooth ./main.go

dev/dbus/install: dev/dbus/link dev/dbus/reload

dev/dbus/link:
	sudo ln -s `pwd`/env/dbus/dbus-go-bluetooth-dev.conf /etc/dbus-1/system.d/
	sudo ln -s `pwd`/env/dbus/dbus-go-bluetooth-service.conf /etc/dbus-1/system.d/

dev/dbus/reload:
	dbus-send --system --type=method_call \
		--dest=org.freedesktop.DBus / org.freedesktop.DBus.ReloadConfig

dev/kill:
	ssh ${DEV_HOST} "killall go-bluetooth" || true

dev/exec: dev/kill
	ssh ${DEV_HOST} "~/go-bluetooth service server --adapterID hci0"

dev/run: dev/cp dev/exec

dev/cp: build dev/kill
	scp go-bluetooth ${DEV_HOST}:~/

dev/logs:
	ssh ${DEV_HOST} "journalctl -u bluetooth.service -f"

docker/bluetoothd/init:
	sudo addgroup bluetooth || true
	sudo adduser `id -nu` bluetooth || true
	sudo ln -s `pwd`/src/bluetooth.conf /etc/dbus-1/system.d/

docker/service/setup:
	./bin/btmgmt power off
	./bin/btmgmt le on
	./bin/btmgmt bredr off
	./bin/btmgmt power on

docker/btmgmt:
	./bin/btmgmt

docker/bluetoothd/build:
	docker build ./env/bluez --build-arg BLUEZ_VERSION=${BLUEZ_VERSION} -t opny/bluez-${BLUEZ_VERSION}

docker/bluetoothd/push:
	docker push opny/bluez-${BLUEZ_VERSION}

docker/bluetoothd/run: service/bluetoothd/stop
	docker run --name bluez_bluetoothd \
		${DOCKER_PARAMS}

bluez-5.50/gen:
	BLUEZ_VERSION=5.50 make gen/clean gen

bluez-5.53/gen:
	BLUEZ_VERSION=5.53 make gen/clean gen

bluez-5.54/gen:
	BLUEZ_VERSION=5.54 make gen/clean gen
