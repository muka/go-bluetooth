# go-bluetooth

Golang bluetooth client based on bluez DBus interfaces

See here for reference https://git.kernel.org/cgit/bluetooth/bluez.git/tree/doc

Status
---

The code is still unstable but mostly working. The API and code is of `alpha` quality and may change without notice

Features
---

- [x] Discovery
- [x] Adapter support
- [x] Device support (SensorTag example)
- [x] GATT Service and characteristics interface
- [x] Adapter on/off via `rfkill`
- [x] Handle systemd `bluetooth.service` unit
- [x] Expose `hciconfig` basic API
- [ ] Expose bluetooth services via bluez DBus API

Usage
---

Check in `main.go` or in `examples/` for an initial overview of the API

Setup
---

Tested with

- golang `1.8` (minimum `v1.6`)
- bluez bluetooth `v5.43` (**Note** this version is the minimum supported one!)
- ubuntu 16.10 kernel `4.8.0-27-generic`
- raspbian and hypirot (debian 8) armv7 `4.4.x`  

See in `scripts/` how to upgrade bluez to 5.43

Give access to `hciconfig` to any user (may have [security implications](https://www.insecure.ws/linux/getcap_setcap.html))

```
sudo setcap 'cap_net_raw,cap_net_admin+eip' `which hciconfig`
```

Todo
---

 - Refactor DBus interface implementation
 - Add Device read / write
 - Add unit tests
 - Add Travis integration
 - Add and generate docs with examples

License
---

MIT
