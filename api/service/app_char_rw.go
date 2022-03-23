package service

import (
	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
)

// Confirm This method doesn't expect a reply so it is just a
// confirmation that value was received.
//
// Possible Errors: org.bluez.Error.Failed
func (s *Char) Confirm() *dbus.Error {
	log.Debug("Char.Confirm")
	return nil
}

// StartNotify Starts a notification session from this characteristic
// if it supports value notifications or indications.
// Calls a notification callback to send a start signal
//
// Possible Errors: org.bluez.Error.Failed
// 		 org.bluez.Error.NotPermitted
// 		 org.bluez.Error.InProgress
// 		 org.bluez.Error.NotSupported
func (s *Char) StartNotify() *dbus.Error {
	log.Debug("Char.StartNotify")
	if s.notifyCallback != nil {
		log.Trace("Used notify callback")

		err := s.notifyCallback(s, true)
		if err != nil {
			return dbus.MakeFailedError(err)
		}
	} else {
		log.Trace("Nothing to notify")
	}
	return nil
}

// StopNotify This method will cancel any previous StartNotify
// transaction. Note that notifications from a
// characteristic are shared between sessions thus
// calling StopNotify will release a single session.
// Calls a notification callback to send a stop signal
//
// Possible Errors: org.bluez.Error.Failed
func (s *Char) StopNotify() *dbus.Error {
	log.Debug("Char.StopNotify")
	if s.notifyCallback != nil {
		log.Trace("Used notify callback")
		err := s.notifyCallback(s, false)
		if err != nil {
			return dbus.MakeFailedError(err)
		}
	} else {
		log.Trace("Nothing to notify")
	}
	return nil
}

// ReadValue Issues a request to read the value of the
// characteristic and returns the value if the
// operation was successful.
//
// Possible options: "offset": uint16 offset
// 			"device": Object Device (Server only)
//
// Possible Errors: org.bluez.Error.Failed
// 		 org.bluez.Error.InProgress
// 		 org.bluez.Error.NotPermitted
// 		 org.bluez.Error.NotAuthorized
// 		 org.bluez.Error.InvalidOffset
// 		 org.bluez.Error.NotSupported
func (s *Char) ReadValue(options map[string]interface{}) ([]byte, *dbus.Error) {

	log.Debug("Characteristic.ReadValue")
	if s.readCallback != nil {
		b, err := s.readCallback(s, options)
		if err != nil {
			return nil, dbus.MakeFailedError(err)
		}
		return b, nil
	}

	return s.Properties.Value, nil
}

//WriteValue Issues a request to write the value of the
// characteristic.
//
// Possible options: "offset": Start offset
// 			"device": Device path (Server only)
// 			"link": Link type (Server only)
// 			"prepare-authorize": boolean Is prepare
// 							 authorization
// 							 request
//
// Possible Errors: org.bluez.Error.Failed
// 		 org.bluez.Error.InProgress
// 		 org.bluez.Error.NotPermitted
// 		 org.bluez.Error.InvalidValueLength
// 		 org.bluez.Error.NotAuthorized
// 		 org.bluez.Error.NotSupported
func (s *Char) WriteValue(value []byte, options map[string]interface{}) *dbus.Error {

	log.Trace("Characteristic.WriteValue")

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
