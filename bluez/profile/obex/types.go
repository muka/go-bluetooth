package obex

import "github.com/godbus/dbus/v5"

//VCardItem vcard-listing data where every entry consists of a pair of strings containing the vcard handle and the contact name. For example:"1.vcf" : "John"
type VCardItem struct {
	Vcard string
	Name  string
}

//Message map to array{object, dict}
type Message struct {
	Path dbus.ObjectPath
	Dict map[string]interface{}
}
