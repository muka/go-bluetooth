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

// BlueZ D-Bus LE Advertising API Description
// Advertising packets are structured data which is broadcast on the LE Advertising
// channels and available for all devices in range.  Because of the limited space
// available in LE Advertising packets (31 bytes), each packet's contents must be
// carefully controlled.
// BlueZ acts as a store for the Advertisement Data which is meant to be sent.
// It constructs the correct Advertisement Data from the structured
// data and configured the kernel to send the correct advertisement.
// Advertisement Data objects are registered freely and then referenced by BlueZ
// when constructing the data sent to the kernel.
package advertising
