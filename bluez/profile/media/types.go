package media

import "github.com/godbus/dbus/v5"

// Item map to array{objects, properties}
type Item struct {
	Object   dbus.ObjectPath
	Property map[string]interface{}
}

// Track map to a media track
type Track struct {
	// Track title name
	Title string
	// Track artist name
	Artist string
	// Track album name
	Album string
	// Track genre name
	Genre string
	// Number of tracks in total
	NumberOfTracks uint32
	// Track number
	TrackNumber uint32
	// Track duration in milliseconds
	Duration uint32
}
