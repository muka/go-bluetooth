// WARNING: generated code, do not edit!
// Copyright © 2019 luca capra
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package profile

import (
	"github.com/godbus/dbus"
)

var (

	// NotReady map to org.bluez.Error.NotReady
	ErrNotReady = dbus.Error{
		Name: "org.bluez.Error.NotReady",
		Body: []interface{}{"NotReady"},
	}

	// InvalidArguments map to org.bluez.Error.InvalidArguments
	ErrInvalidArguments = dbus.Error{
		Name: "org.bluez.Error.InvalidArguments",
		Body: []interface{}{"InvalidArguments"},
	}

	// Failed map to org.bluez.Error.Failed
	ErrFailed = dbus.Error{
		Name: "org.bluez.Error.Failed",
		Body: []interface{}{"Failed"},
	}

	// NotPermitted map to org.bluez.Error.NotPermitted
	ErrNotPermitted = dbus.Error{
		Name: "org.bluez.Error.NotPermitted",
		Body: []interface{}{"NotPermitted"},
	}

	// DoesNotExist map to org.bluez.Error.DoesNotExist
	ErrDoesNotExist = dbus.Error{
		Name: "org.bluez.Error.DoesNotExist",
		Body: []interface{}{"DoesNotExist"},
	}

	// Rejected map to org.bluez.Error.Rejected
	ErrRejected = dbus.Error{
		Name: "org.bluez.Error.Rejected",
		Body: []interface{}{"Rejected"},
	}

	// NotConnected map to org.bluez.Error.NotConnected
	ErrNotConnected = dbus.Error{
		Name: "org.bluez.Error.NotConnected",
		Body: []interface{}{"NotConnected"},
	}

	// NotAcquired map to org.bluez.Error.NotAcquired
	ErrNotAcquired = dbus.Error{
		Name: "org.bluez.Error.NotAcquired",
		Body: []interface{}{"NotAcquired"},
	}

	// NotSupported map to org.bluez.Error.NotSupported
	ErrNotSupported = dbus.Error{
		Name: "org.bluez.Error.NotSupported",
		Body: []interface{}{"NotSupported"},
	}

	// NotAuthorized map to org.bluez.Error.NotAuthorized
	ErrNotAuthorized = dbus.Error{
		Name: "org.bluez.Error.NotAuthorized",
		Body: []interface{}{"NotAuthorized"},
	}

	// NotAvailable map to org.bluez.Error.NotAvailable
	ErrNotAvailable = dbus.Error{
		Name: "org.bluez.Error.NotAvailable",
		Body: []interface{}{"NotAvailable"},
	}

)
