package bluez

import (
	"reflect"

	"github.com/godbus/dbus/v5"
	"github.com/muka/go-bluetooth/util"
	log "github.com/sirupsen/logrus"
)

type WatchableClient interface {
	Client() *Client
	Path() dbus.ObjectPath
	ToProps() Properties
	GetWatchPropertiesChannel() chan *dbus.Signal
	SetWatchPropertiesChannel(chan *dbus.Signal)
}

// WatchProperties updates on property changes
func WatchProperties(wprop WatchableClient) (chan *PropertyChanged, error) {

	channel, err := wprop.Client().Register(wprop.Path(), PropertiesInterface)
	if err != nil {
		return nil, err
	}

	wprop.SetWatchPropertiesChannel(channel)
	ch := make(chan *PropertyChanged)

	go (func() {
		defer func() {
			if err := recover(); err != nil {
				log.Warnf("Recovering from panic in SetWatchPropertiesChannel: %s", err)
			}
		}()

		for {

			if channel == nil {
				break
			}

			sig := <-channel

			if sig == nil {
				return
			}

			if sig.Name != PropertiesChanged {
				continue
			}
			if sig.Path != wprop.Path() {
				continue
			}

			iface := sig.Body[0].(string)
			changes := sig.Body[1].(map[string]dbus.Variant)

			for field, val := range changes {

				// updates [*]Properties struct when a property change
				s := reflect.ValueOf(wprop.ToProps()).Elem()
				// exported field
				f := s.FieldByName(field)
				if f.IsValid() {
					// A Value can be changed only if it is
					// addressable and was not obtained by
					// the use of unexported struct fields.
					if f.CanSet() {
						x := reflect.ValueOf(val.Value())
						wprop.ToProps().Lock()
						// map[*]variant -> map[*]interface{}
						ok, err := util.AssignMapVariantToInterface(f, x)
						if err != nil {
							log.Errorf("Failed to set %s: %s", f.String(), err)
							continue
						}
						// direct assignment
						if !ok {
							f.Set(x)
						}
						wprop.ToProps().Unlock()
					}
				}

				propChanged := &PropertyChanged{
					Interface: iface,
					Name:      field,
					Value:     val.Value(),
				}
				ch <- propChanged
			}

		}
	})()

	return ch, nil
}

func UnwatchProperties(wprop WatchableClient, ch chan *PropertyChanged) error {
	defer func() {
		if err := recover(); err != nil {
			log.Warnf("Recovering from panic in UnwatchProperties: %s", err)
		}
	}()
	if wprop.GetWatchPropertiesChannel() != nil {
		wprop.GetWatchPropertiesChannel() <- nil
		err := wprop.Client().Unregister(wprop.Path(), PropertiesInterface, wprop.GetWatchPropertiesChannel())
		if err != nil {
			return err
		}
	}
	ch <- nil
	close(ch)
	return nil
}
