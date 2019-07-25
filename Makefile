
bluetoothd:
	sudo killall bluetoothd; true
	sudo bluetoothd -Edn P hostname

run/example/service:
	go run examples/service/*.go

run/example/client:
	go run examples/service/*.go client

gen/clean:
	rm -rf test/out
	mkdir -p test/out

gen/run: 
	go run examples/srcgen/main.go
