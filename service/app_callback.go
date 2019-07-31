package service

import "github.com/godbus/dbus"

//CallbackNotRegistered callback not registered
const CallbackNotRegistered = -1

//CallbackFunctionError callback reported an error
const CallbackFunctionError = -2

//GattWriteCallback A callback we can register to handle write requests
type GattWriteCallback func(app *Application, serviceObjPath dbus.ObjectPath, charObjPath dbus.ObjectPath, value []byte) error

//GattDescriptorWriteCallback A callback we can register to handle descriptor write requests
type GattDescriptorWriteCallback func(app *Application, serviceObjPath dbus.ObjectPath, charObjPath dbus.ObjectPath, descObjPath dbus.ObjectPath, value []byte) error

//GattReadCallback A callback we can register to handle read requests
type GattReadCallback func(app *Application, serviceObjPath dbus.ObjectPath, charObjPath dbus.ObjectPath) ([]byte, error)

//GattDescriptorReadCallback A callback we can register to handle descriptor ead requests
type GattDescriptorReadCallback func(app *Application, serviceObjPath dbus.ObjectPath, charObjPath dbus.ObjectPath, descObjPath dbus.ObjectPath) ([]byte, error)

// CallbackError error from a callback
type CallbackError struct {
	msg  string
	code int
}

func (e *CallbackError) Error() string {
	return e.msg
}

//NewCallbackError create a new callback error
func NewCallbackError(code int, msg string) *CallbackError {
	result := &CallbackError{msg: msg, code: code}
	return result
}

//HandleRead Handle application read
func (app *Application) HandleRead(serviceObjPath dbus.ObjectPath, charObjPath dbus.ObjectPath) ([]byte, *CallbackError) {
	if app.config.ReadFunc == nil {
		b := make([]byte, 0)
		return b, NewCallbackError(CallbackNotRegistered, "No callback registered.")
	}

	var cberr *CallbackError
	b, err := app.config.ReadFunc(app, serviceObjPath, charObjPath)
	if err != nil {
		if callbackErr, ok := err.(*CallbackError); ok {
			// If a CallbackError is returned, simply pass it through
			cberr = callbackErr
		} else {
			cberr = NewCallbackError(CallbackFunctionError, err.Error())
		}
	}

	return b, cberr
}

// HandleWrite handle application write
func (app *Application) HandleWrite(serviceObjPath dbus.ObjectPath, charObjPath dbus.ObjectPath, value []byte) *CallbackError {
	if app.config.WriteFunc == nil {
		return NewCallbackError(CallbackNotRegistered, "No callback registered.")
	}

	err := app.config.WriteFunc(app, serviceObjPath, charObjPath, value)
	if err != nil {
		if callbackErr, ok := err.(*CallbackError); ok {
			// If a CallbackError is returned, simply pass it through
			return callbackErr
		}

		return NewCallbackError(CallbackFunctionError, err.Error())
	}

	return nil
}

//HandleDescriptorRead handle descriptor read
func (app *Application) HandleDescriptorRead(serviceObjPath dbus.ObjectPath, charObjPath dbus.ObjectPath, descObjPath dbus.ObjectPath) ([]byte, *CallbackError) {
	if app.config.DescReadFunc == nil {
		b := make([]byte, 0)
		return b, NewCallbackError(CallbackNotRegistered, "No callback registered.")
	}

	var cberr *CallbackError
	b, err := app.config.DescReadFunc(app, serviceObjPath, charObjPath, descObjPath)
	if err != nil {
		if callbackErr, ok := err.(*CallbackError); ok {
			// If a CallbackError is returned, simply pass it through
			cberr = callbackErr
		} else {
			cberr = NewCallbackError(CallbackFunctionError, err.Error())
		}
	}

	return b, cberr
}

//HandleDescriptorWrite handle descriptor write
func (app *Application) HandleDescriptorWrite(serviceObjPath dbus.ObjectPath, charObjPath dbus.ObjectPath, descObjPath dbus.ObjectPath, value []byte) *CallbackError {
	if app.config.DescWriteFunc == nil {
		return NewCallbackError(CallbackNotRegistered, "No callback registered.")
	}

	err := app.config.DescWriteFunc(app, serviceObjPath, charObjPath, descObjPath, value)
	if err != nil {
		if callbackErr, ok := err.(*CallbackError); ok {
			// If a CallbackError is returned, simply pass it through
			return callbackErr
		}

		return NewCallbackError(CallbackFunctionError, err.Error())
	}

	return nil
}
