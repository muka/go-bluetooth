# go-bluetooth

Golang bluetooth client based on bluez DBus interfaces

See here for reference https://git.kernel.org/cgit/bluetooth/bluez.git/tree/doc

## Features

The features implemented are

- [x] Discovery
- [x] Adapter support
- [x] Device support (SensorTag example)
- [x] GATT Service and characteristics interface
- [x] Adapter on/off via `rfkill`
- [x] Handle systemd `bluetooth.service` unit
- [x] Expose `hciconfig` basic API
- [x] Expose bluetooth services via bluez GATT API
- [x] Basic pairing support

## Examples

The `examples/` folder offer an overview of library

- `agent` a simple agent to support pairing
- `btmgmt` interface to CLI btmgmt
- `discovery` find devices around
- `hci_updown` HCI based communication example
- `obex_push` send file to a device
- `sensortag_info` Obtain data from a TI SensorTag
- `sensortag_temperature` Obtain temperature from a TI SensorTag
- `service` expose a bluetooth device with corresponding services
- `show_miband_info` show informations for MiBand2
- `watch_changes` register for notifications from a TI SensorTag

**Note** Ensure to install proper dbus rules on the system. For a dev setup use

```
sudo ln -s `pwd`/scripts/dbus-go-bluetooth-service.conf /etc/dbus-1/system.d/
sudo ln -s `pwd`/scripts/dbus-go-bluetooth-dev.conf /etc/dbus-1/system.d/
```


## Setup

The library has been tested with

- golang `1.11.4` (starting from `v1.6`)
- bluez bluetooth `v5.50` (starting from `v5.43`)

### bluez upgrade

Bluez, the linux bluetooth implementation, has introduced GATT support from `v5.43`

Ensure you are using an up to date version with `bluetoothd -v`

See in `scripts/` how to upgrade bluez

### Development notes

-   Give access to `hciconfig` to any user (may have [security implications](https://www.insecure.ws/linux/getcap_setcap.html))

    ```
    sudo setcap 'cap_net_raw,cap_net_admin+eip' `which hciconfig`
    ```
- Create a dbus profile

    ```sh
    ln -s `pwd`/scripts/dbus-dev.conf /etc/dbus-1/system.d/go-bluetooth.config
    ```
- Monitor activity

    `sudo dbus-monitor --system "type=error"`

- View `bluetoothd` debug messages

    `sudo service bluetooth stop && sudo bluetoothd -Edn P hostname`

- Enable LE advertisement (to use a single pc, you will need 2 bluetooth adapter)

  ```bash

    sudo btmgmt -i 0 power off
    sudo btmgmt -i 0 name "my go app"
    sudo btmgmt -i 0 le on    
    sudo btmgmt -i 0 connectable on
    sudo btmgmt -i 0 advertising on
    sudo btmgmt -i 0 power on

  ```

## References

- https://git.kernel.org/cgit/bluetooth/bluez.git/tree/doc
- https://www.bluetooth.com/specifications/gatt/services
- http://events.linuxfoundation.org/sites/events/files/slides/Bluetooth%20on%20Modern%20Linux_0.pdf
- https://github.com/nettlep/gobbledegook

## License

MIT License
