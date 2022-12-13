
.PHONY: gen

BLUEZ_VERSION ?= 5.60
FILTER ?=

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
	rm -f `ls bluez/profile/gen_* -1` || true

gen/run: bluez/checkout
	BLUEZ_VERSION=${BLUEZ_VERSION} FILTER=${FILTER} go run gen/srcgen/main.go full

gen: gen/run

build: gen
	CGO_ENABLED=0 go build -o go-bluetooth ./main.go

dev/kill:
	ssh ${DEV_HOST} "killall go-bluetooth" || true

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

bluez-5.55/gen:
	BLUEZ_VERSION=5.55 make gen/clean gen

bluez-5.62/gen:
	BLUEZ_VERSION=5.62 make gen/clean gen

bluez-5.60/gen:
	BLUEZ_VERSION=5.60 make gen/clean gen

bluez-5.64/gen:
	BLUEZ_VERSION=5.64 make gen/clean gen

bluez-5.65/gen:
	BLUEZ_VERSION=5.65 make gen/clean gen

bluez-5.66/gen:
	BLUEZ_VERSION=5.66 make gen/clean gen
