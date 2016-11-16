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
- [ ] Handle `hciconfig` command via CLI
- [ ] Expose bluetooth services via bluez DBus API

Usage
---

Check in `main.go` or in `examples/` for an initial overview of the API

Setup
---

Tested with

- golang >= `v1.6.2`
- bluez bluetooth `v5.43` (**Note** this version is the minimum supported one!)
- ubuntu 16.10 kernel `4.8.0-27-generic`

Use `glide install` to install dependencies

Todo
---

 - Refactor DBus interface implementation
 - Add Device read / write
 - Add unit tests
 - Add Travis integration
 - Add and generate docs

License
---

MIT
