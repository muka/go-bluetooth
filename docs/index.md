

# Examples

## Running the service example

Requires two adapters. In two shell run in order

- `go run main.go service server`
- `go run main.go service client 00:1A:7D:DA:71:15 --adapterID hci1 `


## Development

Edit `/etc/systemd/system/bluetooth.target.wants/bluetooth.service` and update the `ExecStart` command adding those optios

`ExecStart=/usr/lib/bluetooth/bluetoothd -E -d -P hostname`

 - `-E` experimental mode
 - `-d` debug mode
 - `-P` no plugin

Afterwards run `systemctl daemon-reload` to reload the config and `service bluetooth restart` to restart the service.

To view logs `journalctl -u bluetooth -f`
