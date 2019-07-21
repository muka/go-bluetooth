package profile

// Media hierarchy
// ===============
//
// Service		org.bluez
// Interface	org.bluez.Media1
// Object path	[variable prefix]/{hci0,hci1,...}

// void RegisterEndpoint(object endpoint, dict properties)
//
//   Register a local end point to sender, the sender can
//   register as many end points as it likes.
//
//   Note: If the sender disconnects the end points are
//   automatically unregistered.
//
//   possible properties:
//
//     string UUID:
//
//       UUID of the profile which the endpoint
//       is for.
//
//     byte Codec:
//
//       Assigned number of codec that the
//       endpoint implements. The values should
//       match the profile specification which
//       is indicated by the UUID.
//
//     array{byte} Capabilities:
//
//       Capabilities blob, it is used as it is
//       so the size and byte order must match.
//
//   Possible Errors: org.bluez.Error.InvalidArguments
//        org.bluez.Error.NotSupported - emitted
//        when interface for the end-point is
//        disabled.
//
// void UnregisterEndpoint(object endpoint)
//
//   Unregister sender end point.
//
// void RegisterPlayer(object player, dict properties)
//
//   Register a media player object to sender, the sender
//   can register as many objects as it likes.
//
//   Object must implement at least
//   org.mpris.MediaPlayer2.Player as defined in MPRIS 2.2
//   spec:
//
//   http://specifications.freedesktop.org/mpris-spec/latest/
//
//   Note: If the sender disconnects its objects are
//   automatically unregistered.
//
//   Possible Errors: org.bluez.Error.InvalidArguments
//        org.bluez.Error.NotSupported
//
// void UnregisterPlayer(object player)
//
//   Unregister sender media player.
