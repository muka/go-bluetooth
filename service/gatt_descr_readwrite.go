package service

import "github.com/godbus/dbus"

//ReadValue read a value
func (s *GattDescriptor1) ReadValue(options map[string]interface{}) ([]byte, *dbus.Error) {
	b, err := s.config.characteristic.config.service.config.app.HandleDescriptorRead(
		s.config.characteristic.config.service.Path(), s.config.characteristic.Path(),
		s.Path())

	var dberr *dbus.Error
	if err != nil {
		if err.code == -1 {
			// No registered callback, so we'll just use our stored value
			b = s.properties.Value
		} else {
			dberr = dbus.NewError(err.Error(), nil)
		}
	}

	return b, dberr
}

//WriteValue write a value
func (s *GattDescriptor1) WriteValue(value []byte, options map[string]interface{}) *dbus.Error {
	err := s.config.characteristic.config.service.config.app.HandleDescriptorWrite(
		s.config.characteristic.config.service.Path(), s.config.characteristic.Path(),
		s.Path(), value)

	if err != nil {
		if err.code == -1 {
			// No registered callback, so we'll just store this value
			s.UpdateValue(value)
			return nil
		}
		dberr := dbus.NewError(err.Error(), nil)
		return dberr
	}

	return nil
}

//UpdateValue update a descriptor value
func (s *GattDescriptor1) UpdateValue(value []byte) error {
	s.properties.Value = value
	err := s.PropertiesInterface.Instance().Set(s.Interface(), "Value", dbus.MakeVariant(value))
	if err != nil {
		return err
	}
	return nil
}
