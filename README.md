# go-bluetooth

Go bluetooth API for Linux-based Bluez DBus interface.

[![GoDoc](https://godoc.org/github.com/muka/go-bluetooth?status.svg)](https://godoc.org/github.com/muka/go-bluetooth)

<img align="center" width="90" src="./gopher.png">

## Features

The library is a wrapper to the Bluez DBus API and some high level API to ease the interaction.

High level features supported:

- [x] Client code generation from bluez documentation
- [x] Shell wrappers for `rfkill`, `btmgmt`, `hciconfig`, `hcitool`
- [x] An `hci` socket basic API (inspired by [go-ble/ble](https://github.com/go-ble/ble))
- [x] Expose bluetooth service from go code [*unstable*]
- [x] Pairing and authentication support (via agent)
- [x] Beaconing send & receive (iBeacon and Eddystone)

## Code generation

The code structure follow this pattern:

 - `./api` contains wrappers for the DBus Api
 - `./bluez` contains the actual implementation, generated from the `bluez` documentation

Use `make gen` to re-generate go sources. There is also a commodity bluez JSON file available in the root folder for reference.

Generated code has `gen_` prefix. If an API file exists with the same filename but _without_ the prefix, generation will be skipped for that API.

## Development setup

1. Clone the repository

  `git clone https://github.com/muka/go-bluetooth.git`

1. Retrieve the bluetooth API and generate GO code

  `make all`

2. Run an example

  `go-bluetooth discovery`

The `examples/` folder offer an overview of the API.

## DBus configuration setup

In order to interact with DBus, propert configurations must be installed in the system. For a development setup, the repository provides example configurations.

```sh
  cd $GOPATH/src/github.com/muka/go-bluetooth
  make dev/dbus/install
```

## Requirements

The library is tested with

- golang `1.14.1`
- bluez bluetooth `v5.54`

### Development notes

-   Give access to `hciconfig` to any user and avoid `sudo` (may have [security implications](https://www.insecure.ws/linux/getcap_setcap.html))

    ```
    sudo setcap 'cap_net_raw,cap_net_admin+eip' `which hciconfig`
    ```
- Monitor Bluetooth activity

  `sudo btmon`

- Monitor DBus activity

    `sudo dbus-monitor --system "type=error"`

- Start `bluetoothd` with experimental features and verbose debug messages `make bluetoothd`

- Enable LE advertisement (on a single pc ensure to use at least 2x bluetooth adapter)

  ```bash

    sudo btmgmt -i 0 power off
    sudo btmgmt -i 0 name "my go app"
    sudo btmgmt -i 0 le on    
    sudo btmgmt -i 0 connectable on
    sudo btmgmt -i 0 advertising on
    sudo btmgmt -i 0 power on

  ```

## References

- [Standard GATT characteristics descriptions](https://www.bluetooth.com/specifications/gatt/)
- [bluez git](https://git.kernel.org/cgit/bluetooth/bluez.git/tree/doc)
- [GATT services specs](https://www.bluetooth.com/specifications/gatt/services)
- [A C++ implementation](https://github.com/nettlep/gobbledegook)
- [DBus specifications](https://dbus.freedesktop.org/doc/dbus-specification.html#type-system)
- [SensorTag specs](http://processors.wiki.ti.com/images/a/a8/BLE_SensorTag_GATT_Server.pdf)

## License

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
