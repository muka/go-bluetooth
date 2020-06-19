package device

import (
	"errors"
	"fmt"
	"strings"

	"github.com/godbus/dbus/v5"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile/gatt"
)

func NewDevice(adapterID string, address string) (*Device1, error) {
	path := fmt.Sprintf("%s/%s/dev_%s", bluez.OrgBluezPath, adapterID, strings.Replace(address, ":", "_", -1))
	return NewDevice1(dbus.ObjectPath(path))
}

// GetCharacteristicsList return device characteristics object path list
func (d *Device1) GetCharacteristicsList() ([]dbus.ObjectPath, error) {

	var chars []dbus.ObjectPath

	om, err := bluez.GetObjectManager()
	if err != nil {
		return nil, err
	}

	list, err := om.GetManagedObjects()
	if err != nil {
		return nil, err
	}

	for path := range list {

		spath := string(path)

		// log.Debugf("%s=%s", string(d.Path()), spath)

		if !strings.HasPrefix(spath, string(d.Path())) {
			continue
		}

		charPos := strings.Index(spath, "char")
		if charPos == -1 {
			continue
		}

		if strings.Contains(spath[charPos:], "desc") {
			continue
		}

		chars = append(chars, path)
	}

	return chars, nil
}

// GetDescriptorList returns all descriptors
func (d *Device1) GetDescriptorList() ([]dbus.ObjectPath, error) {
	var descr []dbus.ObjectPath

	om, err := bluez.GetObjectManager()
	if err != nil {
		return nil, err
	}

	list, err := om.GetManagedObjects()
	if err != nil {
		return nil, err
	}
	for path := range list {

		spath := string(path)

		if !strings.HasPrefix(spath, string(d.Path())) {
			continue
		}

		charPos := strings.Index(spath, "char")
		if charPos == -1 {
			continue
		}

		if strings.Contains(spath[charPos:], "desc") {
			continue
		}

		descr = append(descr, path)
	}

	return descr, nil
}

//GetDescriptors returns all descriptors for a given characteristic
func (d *Device1) GetDescriptors(char *gatt.GattCharacteristic1) ([]*gatt.GattDescriptor1, error) {
	descrPaths, err := d.GetDescriptorList()
	if err != nil {
		return nil, err
	}

	descrFound := []*gatt.GattDescriptor1{}
	for _, path := range descrPaths {
		descr, err := gatt.NewGattDescriptor1(path)
		if err != nil {
			return nil, err
		}
		if char.Path() == descr.Properties.Characteristic {
			descrFound = append(descrFound, descr)
		}
	}

	if len(descrFound) == 0 {
		return nil, errors.New("descriptors not found")
	}

	return descrFound, nil
}

//GetCharacteristics return a list of characteristics
func (d *Device1) GetCharacteristics() ([]*gatt.GattCharacteristic1, error) {

	list, err := d.GetCharacteristicsList()
	if err != nil {
		return nil, err
	}

	chars := []*gatt.GattCharacteristic1{}
	for _, path := range list {

		char, err := gatt.NewGattCharacteristic1(path)
		if err != nil {
			return nil, err
		}
		chars = append(chars, char)
	}

	return chars, nil
}

//GetAllServicesAndUUID return a list of uuid's with their corresponding services
func (d *Device1) GetAllServicesAndUUID() ([]string, error) {

	list, err := d.GetCharacteristicsList()
	if err != nil {
		return nil, err
	}

	chars := map[dbus.ObjectPath]*gatt.GattCharacteristic1{}

	var deviceFound []string
	var uuidAndService string
	for _, path := range list {

		char, err := gatt.NewGattCharacteristic1(path)
		if err != nil {
			return nil, err
		}
		chars[path] = char

		props := chars[path].Properties
		cuuid := strings.ToUpper(props.UUID)
		service := string(props.Service)

		uuidAndService = fmt.Sprint(cuuid, ":", service)
		deviceFound = append(deviceFound, uuidAndService)
	}

	return deviceFound, nil
}

//GetCharByUUID return a GattService by its uuid, return nil if not found
func (d *Device1) GetCharByUUID(uuid string) (*gatt.GattCharacteristic1, error) {
	devices, err := d.GetCharsByUUID(uuid)
	if len(devices) > 0 {
		return devices[0], err
	}
	return nil, err
}

// GetCharsByUUID returns all characteristics that match the given UUID.
func (d *Device1) GetCharsByUUID(uuid string) ([]*gatt.GattCharacteristic1, error) {
	uuid = strings.ToUpper(uuid)

	list, err := d.GetCharacteristicsList()
	if err != nil {
		return nil, err
	}

	charsFound := []*gatt.GattCharacteristic1{}

	for _, path := range list {

		char, err := gatt.NewGattCharacteristic1(path)
		if err != nil {
			return nil, err
		}

		cuuid := strings.ToUpper(char.Properties.UUID)

		if cuuid == uuid {
			charsFound = append(charsFound, char)
		}
	}

	if len(charsFound) == 0 {
		return nil, errors.New("characteristic not found")
	}

	return charsFound, nil
}
