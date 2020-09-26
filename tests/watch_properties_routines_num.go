/**
 * This script demonstrates the increase of go routines depending on the
 * UnwatchProperties change. The code is related to
 * https://github.com/muka/go-bluetooth/issues/113 and just documents
 * the issue to see if the fix brings any value.
 *
 * Functionally this script doesn't do anything beyond connecting to a
 * BLE device and repeatedly subscribing and unsubscribing to receive
 * notifications as fast as possible. Given the "as fast as possible"
 * part also caught a panic when the channel was close in inappropriate
 * moment of routine started in the WatchProperties which resulted in
 * an attempt to write to a closed channel
 */
package main

import (
	"errors"
	"flag"
	"fmt"
	"runtime"
	"time"

	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile/gatt"
)

func runWatchPropertiesTestIteration(char *gatt.GattCharacteristic1) (int, error) {
	ch, err := char.WatchProperties()
	if err != nil {
		return 0, err
	}

	err = char.StartNotify()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return 0, err
	}

	go func() {
		for e := range ch {
			if e == nil {
				return
			}
		}
	}()

	err = char.UnwatchProperties(ch)
	if err != nil {
		return 0, err
	}
	err = char.StopNotify()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return 0, err
	}

	// Optionally wait a little for the anonumous go routine to finish
	// If you don't wait here then the results in go routines count diff
	// will sometimes show +1 or -1 to average around 0 due to delay
	// in reception of nil on WatchProperties() created channel
	// time.Sleep(10 * time.Millisecond)
	return runtime.NumGoroutine(), nil
}

func main() {
	hciName := flag.String("hci", "hci1", "Name of your HCI device")
	devMac := flag.String("mac", "CA:AF:FE:00:BE:EF", "MAC of the device to connect for test")
	charUUID := flag.String("uuid", "76494b4a-305c-4368-aa02-923b1b709333", "UUID characteristic which notifies")
	flag.Parse()

	a, err := api.GetAdapter(*hciName)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting scan")
	err = a.StartDiscovery()
	if err != nil {
		panic(err)
	}
	time.Sleep(2 * time.Second)
	err = a.StopDiscovery()
	if err != nil {
		panic(err)
	}
	fmt.Println("Scan finished")

	dev, err := a.GetDeviceByAddress(*devMac)
	if err != nil {
		panic(err)
	}
	if dev == nil {
		msg := fmt.Sprintf("Device %s not found. Cannot connect\n", *devMac)
		panic(errors.New(msg))
	}

	fmt.Printf("Connecting to %s\n", dev.Properties.Address)
	err = dev.Connect()
	if err != nil {
		panic(err)
	}

	char, err := dev.GetCharByUUID(*charUUID)
	if err != nil {
		panic("search device failed: " + err.Error())
	}

	// Run test iterations
	iterCount := 100
	prev := runtime.NumGoroutine()
	startGoRoutineCount := prev
	for i := 1; i < iterCount; i++ {
		result, err := runWatchPropertiesTestIteration(char)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		diff := result - prev
		fmt.Printf("[Iteration %04d] Result: %d, diff: %d\n", i, result, diff)
		prev = result
	}

	endGoRoutineCount := runtime.NumGoroutine()
	diff := endGoRoutineCount - startGoRoutineCount
	fmt.Printf("Number of go routines changed by %d. Started with %d,"+
		"finished with %d\n", diff, startGoRoutineCount, endGoRoutineCount)
	fmt.Printf("Done\n")
}
