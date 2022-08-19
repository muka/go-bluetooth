
# go-bluetooth

<img style="float:right" width="70" src="./gopher.png"> Go bluetooth API for Linux-based Bluez DBus interface.

[![GoDoc](https://godoc.org/github.com/muka/go-bluetooth?status.svg)](https://godoc.org/github.com/muka/go-bluetooth)


## Bluez API versioning

The version on master tracks Ubuntu bluetoothd version (v5.64). There are dedicated branches/tags with various version

- [bluez 5.60](https://github.com/muka/go-bluetooth)
- [bluez 5.62](https://github.com/muka/go-bluetooth/tree/bluez-5.62)
- [bluez 5.63](https://github.com/muka/go-bluetooth/tree/bluez-5.63)
- [bluez 5.64](https://github.com/muka/go-bluetooth/tree/bluez-5.64)
- [bluez 5.65](https://github.com/muka/go-bluetooth/tree/bluez-5.65)

See branches for available versions.

## Features

The library is a wrapper to the Bluez DBus API and some high level API to ease the interaction.

High level features supported:

- [x] Client code generation from bluez documentation
- [x] Shell wrappers for `rfkill`, `btmgmt`, `hciconfig`, `hcitool`
- [x] An `hci` socket basic API (inspired by [go-ble/ble](https://github.com/go-ble/ble))
- [x] Expose bluetooth service from go code [*unstable*]
- [x] Pairing and authentication support (via agent)
- [x] Beaconing send & receive (iBeacon and Eddystone)
- [x] Mesh API support (since v5.53)

## Todo

- [ ] Generate dbus-compatible XML from documentation
- [ ] Generate mock service to test the library against

## Using this library

- [TinyGo bluetooth](https://github.com/tinygo-org/bluetooth)
- [pi-wifi](https://github.com/muka/pi-wifi)

Want to add your project? Please, add it to the [README.md](https://github.com/muka/go-bluetooth/edit/master/README.md) and open a PR!

## Running examples

Examples are available in `_examples` folder.

```sh
cd _examples
go run main.go
# print available example commands
# Example discovery
go run main.go discovery
```

## Development setup

1. Clone the repository

  `git clone https://github.com/muka/go-bluetooth.git`

1. Retrieve the bluetooth API and generate GO code

  ```sh
  make bluez/init bluez/checkout gen/clean gen/run
  ```

## Code generation

To generate code for a new version of bluez use the command

`BLUEZ_VERSION=5.65 make gen/clean gen`

Change to a version available in bluez. Note that generated code may be broken or incomplete and the generation script may need fine tuning.

### How generation works

The code structure follow this pattern:

 - `./api` contains wrappers for the DBus Api
 - `./bluez` contains the actual implementation, generated from the `bluez` documentation

Use `make gen` to re-generate go sources. There is also a commodity bluez JSON file available in the root folder for reference.

Generated code has `gen_` prefix. If an API file exists with the same filename but _without_ the prefix, generation will be skipped for that API.

## Requirements

The library is tested with

- golang `1.17.5`
- bluez bluetooth `v5.50`, `v5.54`, `v5.60`

### Development notes

- Inspect a service ObjectManager

  ```shell
  dbus-send --system --print-reply --dest=go.bluetooth /hci0/apps/0 org.freedesktop.DBus.ObjectManager.GetManagedObjects
  ```

- Retrieve char properties

  ```
  dbus-send --system --print-reply --dest=go.bluetooth /hci0/apps/0/service000003e8/char0  org.freedesktop.DBus.Properties.GetAll string:org.bluez.GattCharacteristic1
  ```

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
