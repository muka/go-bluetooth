package service

import (
	"github.com/godbus/dbus"
	log "github.com/sirupsen/logrus"
)

//ReadValue read a value
func (s *GattCharacteristic1) ReadValue(options map[string]interface{}) ([]byte, *dbus.Error) {
	log.Debug("Characteristic.ReadValue")

	b, err := s.config.service.config.app.HandleRead(s.config.service.Path(), s.Path())

	var dberr *dbus.Error
	if err != nil {
		if err.code == CallbackNotRegistered {
			// No registered callback, so we'll just use our stored value
			b = s.properties.Value
		} else {
			dberr = dbus.NewError(err.Error(), nil)
		}
	}

	return b, dberr
}

//WriteValue write a value
func (s *GattCharacteristic1) WriteValue(value []byte, options map[string]interface{}) *dbus.Error {
	log.Debug("Characteristic.WriteValue")

	err := s.config.service.config.app.HandleWrite(s.config.service.Path(), s.Path(), value)

	if err != nil {
		if err.code == CallbackNotRegistered {
			// No registered callback, so we'll just store this value
			s.UpdateValue(value)
			return nil
		}
		dberr := dbus.NewError(err.Error(), nil)
		return dberr
	}

	return nil
}

//UpdateValue update a value
func (s *GattCharacteristic1) UpdateValue(value []byte) {
	s.properties.Value = value
	s.PropertiesInterface.Instance().Set(s.Interface(), "Value", dbus.MakeVariant(value))
}
