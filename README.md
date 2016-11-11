# bluez-client

Golang bluetooth client based on bluez DBus interfaces

See here for reference https://git.kernel.org/cgit/bluetooth/bluez.git/tree/doc

About
---

The code is still unstable but mostly working

Features
---

- [x] Adapter on/off via `rfkill`
- [x] Handle systemd `bluetooth.service` unit
- [x] Discovery
- [x] Adapter support
- [x] Device support
- [ ] GATT Service and characteristics interface

Usage
---

Run `go run main.go` to get an overview of API

Setup
---

Tested with

- golang >= v1.6.2
- bluez bluetooth >= v5.41
- ubuntu 16.10

Use `glide install` to install dependencies

If `bluetoothd -v` is < 5.42 then

1. add `--experimental` to the bluetooth unit

```
#/etc/systemd/system/bluetooth.target.wants/bluetooth.service
ExecStart=/usr/lib/bluetooth/bluetoothd --experimental
```

2.  reload the unit `systemctl daemon-reload`
3. restart the service `service bluetooth restart`


Todo
---

 - Add Device read / write
 - Refactor DBus interface implementation
 - Add unit tests
 - Add Travis integration
 - Add and generate docs

License
---

MIT
