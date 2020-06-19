package service

import (
	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
)

// Set the Read callback, called when a client attempt to read
func (s *Descr) OnRead(fx DescrReadCallback) *Descr {
	s.readCallback = fx
	return s
}

// Set the Write callback, called when a client attempt to write
func (s *Descr) OnWrite(fx DescrWriteCallback) *Descr {
	s.writeCallback = fx
	return s
}

//ReadValue read a value
func (s *Descr) ReadValue(options map[string]interface{}) ([]byte, *dbus.Error) {

	log.Trace("Descr.ReadValue")

	if s.readCallback != nil {
		b, err := s.readCallback(s, options)
		if err != nil {
			return nil, dbus.MakeFailedError(err)
		}
		return b, nil
	}

	return s.Properties.Value, nil
}

//WriteValue write a value
func (s *Descr) WriteValue(value []byte, options map[string]interface{}) *dbus.Error {

	log.Trace("Descr.WriteValue")

	val := value
	if s.writeCallback != nil {
		log.Trace("Used write callback")
		b, err := s.writeCallback(s, value)
		val = b
		if err != nil {
			return dbus.MakeFailedError(err)
		}
	} else {
		log.Trace("Store directly to value (no callback)")
	}

	// TODO update on Properties interface
	s.Properties.Value = val
	err := s.iprops.Instance().Set(s.Interface(), "Value", dbus.MakeVariant(value))

	return err
}
