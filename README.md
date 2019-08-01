# go-bluetooth

Golang bluetooth client based on bluez DBus interfaces based on Bluez reference documentation https://git.kernel.org/cgit/bluetooth/bluez.git/tree/doc

## Setup

The go Bluez API is generated from the documentation, run `make gen` to generate go sources.

`make gen`

**Note** Ensure to install proper dbus rules on the system. For a dev setup, you can use the library configuration as follow

```sh
  cd $GOPATH/src/github.com/muka/go-bluetooth
  sudo ln -s `pwd`/scripts/dbus-go-bluetooth-service.conf /etc/dbus-1/system.d/
  sudo ln -s `pwd`/scripts/dbus-go-bluetooth-dev.conf /etc/dbus-1/system.d/
  # Reload dbus to load new policies
  dbus-send --system --type=method_call --dest=org.freedesktop.DBus / org.freedesktop.DBus.ReloadConfig
```

## Features

The features implemented are

- [x] Discovery
- [x] Adapter support
- [x] Client device support
- [x] GATT Service and characteristics interface
- [x] Shell wrapper to `rfkill`, `btmgmt`, `hciconfig`, `hcitool`
- [x] Handle systemd `bluetooth.service` unit
- [x] An `hci` basic API (based on a fork of [go-ble/ble](https://github.com/muka/ble))
- [x] Bluetooth services via bluez GATT API (requires 2 bluetooth adapters)
- [x] Pairing support
- [x] Authentication support
- [x] API generation from bluez documentation

## Examples

The `examples/` folder offer an API overview. Use `go run main.go` to list the available examples.

## Requirements

The library is tested with

- golang `1.11`
- bluez bluetooth `v5.50`

### Development notes

-  Standard GATT characteristics descriptions can be found on https://www.bluetooth.com/specifications/gatt/

-   Give access to `hciconfig` to any user and avoid `sudo` (may have [security implications](https://www.insecure.ws/linux/getcap_setcap.html))

    ```
    sudo setcap 'cap_net_raw,cap_net_admin+eip' `which hciconfig`
    ```
- Monitor Bluetooth activity

  `sudo btmon`

- Monitor DBus activity

    `sudo dbus-monitor --system "type=error"`

- Start `bluetoothd` with experimental features and verbose debug messages

    `sudo service bluetooth stop && sudo bluetoothd -Edn P hostname`

- Enable LE advertisement (on a single pc ensure to use at least 2x bluetooth adapter)

  ```bash

    sudo btmgmt -i 0 power off
    sudo btmgmt -i 0 name "my go app"
    sudo btmgmt -i 0 le on    
    sudo btmgmt -i 0 connectable on
    sudo btmgmt -i 0 advertising on
    sudo btmgmt -i 0 power on

  ```

## Contributing

Feel free to open an issue and/or a PR to contribute. If you would like to help improve the library without coding directly, you can also consider to contribute by providing some hardware to test on.

## References

- https://git.kernel.org/cgit/bluetooth/bluez.git/tree/doc
- https://www.bluetooth.com/specifications/gatt/services
- http://events.linuxfoundation.org/sites/events/files/slides/Bluetooth%20on%20Modern%20Linux_0.pdf
- https://github.com/nettlep/gobbledegook

## License

Apache2 License
