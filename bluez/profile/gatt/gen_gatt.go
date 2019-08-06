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

// BlueZ D-Bus GATT API description
// GATT local and remote services share the same high-level D-Bus API. Local
// refers to GATT based service exported by a BlueZ plugin or an external
// application. Remote refers to GATT services exported by the peer.
// BlueZ acts as a proxy, translating ATT operations to D-Bus method calls and
// Properties (or the opposite). Support for D-Bus Object Manager is mandatory for
// external services to allow seamless GATT declarations (Service, Characteristic
// and Descriptors) discovery. Each GATT service tree is required to export a D-Bus
// Object Manager at its root that is solely responsible for the objects that
// belong to that service.
// Releasing a registered GATT service is not defined yet. Any API extension
// should avoid breaking the defined API, and if possible keep an unified GATT
// remote and local services representation.
package gatt
