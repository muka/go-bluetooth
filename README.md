# go-bluetooth

Golang bluetooth client based on bluez DBus interfaces

See here for reference https://git.kernel.org/cgit/bluetooth/bluez.git/tree/doc

## Status

The current API is unstable and may change in the future.

The features implemented are

- [x] Discovery
- [x] Adapter support
- [x] Device support (SensorTag example)
- [x] GATT Service and characteristics interface
- [x] Adapter on/off via `rfkill`
- [x] Handle systemd `bluetooth.service` unit
- [x] Expose `hciconfig` basic API
- [ ] Expose bluetooth services via bluez DBus API [See #7](https://github.com/muka/go-bluetooth/issues/7)

## Examples

Check `main.go` or in `examples/` for an initial overview of the API

## Setup

The library has been tested with

- golang `1.8` (minimum `v1.6`)
- bluez bluetooth `v5.46` (**Note** `v5.43` is the minimum supported one)
- ubuntu 16.10 kernel `4.8.0-27-generic`
- raspbian and hypirot (debian 8) armv7 `4.4.x`  

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
    ln -s `pwd`/scripts/dbus-dev.conf /etc/dbus1/system.d/go-bluetooth.config
    ```
- Monitor activity

    `sudo dbus-monitor --system "type=error"`

## TODO List / Help wanted

-   Add docs with examples
-   Add Device read / write and custom data converters
-   Unit tests coverage

## License

MIT License
