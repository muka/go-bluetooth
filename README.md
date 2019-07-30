# go-bluetooth

Golang bluetooth client based on bluez DBus interfaces

See here for reference https://git.kernel.org/cgit/bluetooth/bluez.git/tree/doc

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

The `examples/` folder offer an API overview

- [agent](./examples/agent) a simple agent to support pairing
- [btmgmt](./examples/btmgmt) interface to CLI btmgmt
- [discovery](./examples/discovery) find devices around
- [hci_updown](./examples/hci_updown) HCI based communication example
- [obex_push](./examples/obex_push) send file to a device
- [sensortag_info](./examples/sensortag_info) Obtain data from a TI SensorTag
- [sensortag_temperature](./examples/sensortag_temperature) Obtain temperature from a TI SensorTag
- [service](./examples/service) expose a bluetooth device with corresponding services
- [show_miband_info](./examples/show_miband_info) show informations for MiBand2
- [watch_changes](./examples/watch_changes) register for notifications from a TI SensorTag

**Note** Ensure to install proper dbus rules on the system. For a dev setup use

```sh
  sudo ln -s `pwd`/scripts/dbus-go-bluetooth-service.conf /etc/dbus-1/system.d/
  sudo ln -s `pwd`/scripts/dbus-go-bluetooth-dev.conf /etc/dbus-1/system.d/
  # Reload dbus to load new policies
  dbus-send --system --type=method_call --dest=org.freedesktop.DBus / org.freedesktop.DBus.ReloadConfig
```

## Setup

The library has been tested with

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

MIT License
