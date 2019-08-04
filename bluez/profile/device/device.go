package device

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
)

func NewDevice(adapterID string, address string) (*Device1, error) {
	return NewDevice1(fmt.Sprintf("%s/%s/dev_%s", bluez.OrgBluezPath, adapterID, strings.Replace(address, ":", "_", -1)))
}

type PropertyChanged struct{}

func (d *Device1) OnChange() (chan PropertyChanged, error) {

	channel, err := d.client.Register(d.Path(), bluez.PropertiesInterface)
	if err != nil {
		return err
	}

	go (func() {
		for {

			if channel == nil {
				break
			}

			sig := <-channel

			if sig == nil {
				return
			}

			if sig.Name != bluez.PropertiesChanged {
				continue
			}
			if fmt.Sprint(sig.Path) != d.Path {
				continue
			}

			// for i := 0; i < len(sig.Body); i++ {
			// log.Printf("%s -> %s\n", reflect.TypeOf(sig.Body[i]), sig.Body[i])
			// }

			iface := sig.Body[0].(string)
			changes := sig.Body[1].(map[string]dbus.Variant)
			propertyChangedEvents := make([]PropertyChangedEvent, 0)
			for field, val := range changes {

				// updates [*]Properties struct
				d.lock.RLock()
				props := d.Properties
				d.lock.RUnlock()

				s := reflect.ValueOf(props).Elem()
				// exported field
				f := s.FieldByName(field)
				if f.IsValid() {
					// A Value can be changed only if it is
					// addressable and was not obtained by
					// the use of unexported struct fields.
					if f.CanSet() {
						x := reflect.ValueOf(val.Value())
						props.Lock()
						f.Set(x)
						props.Unlock()
					}
				}

				propChanged := PropertyChangedEvent{string(iface), field, val.Value(), props, d}
				propertyChangedEvents = append(propertyChangedEvents, propChanged)
			}

			for _, propChanged := range propertyChangedEvents {
				d.Emit("changed", propChanged)
			}

		}
	})()

	return nil
}

func (d *Device1) unwatchProperties() error {
	var err error
	d.Lock()
	defer d.Unlock()
	if d.watchPropertiesChannel != nil {
		err = d.client.Unregister(d.watchPropertiesChannel)
		close(d.watchPropertiesChannel)
		d.watchPropertiesChannel = nil
	}

	return err
}
