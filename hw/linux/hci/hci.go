// credits: https://github.com/go-ble/ble/blob/master/linux/hci/socket/socket.go
// I attempted a PR but id not go through

package hci

import (
	"fmt"
	"io"
	"sync"
	"unsafe"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

func ioR(t, nr, size uintptr) uintptr {
	return (2 << 30) | (t << 8) | nr | (size << 16)
}

func ioW(t, nr, size uintptr) uintptr {
	return (1 << 30) | (t << 8) | nr | (size << 16)
}

func ioctl(fd, op, arg uintptr) error {
	if _, _, ep := unix.Syscall(unix.SYS_IOCTL, fd, op, arg); ep != 0 {
		return ep
	}
	return nil
}

const (
	ioctlSize     = 4
	hciMaxDevices = 16
	typHCI        = 72 // 'H'
)

var (
	hciUpDevice      = ioW(typHCI, 201, ioctlSize) // HCIDEVUP
	hciDownDevice    = ioW(typHCI, 202, ioctlSize) // HCIDEVDOWN
	hciResetDevice   = ioW(typHCI, 203, ioctlSize) // HCIDEVRESET
	hciGetDeviceList = ioR(typHCI, 210, ioctlSize) // HCIGETDEVLIST
	hciGetDeviceInfo = ioR(typHCI, 211, ioctlSize) // HCIGETDEVINFO
)

type devListRequest struct {
	devNum     uint16
	devRequest [hciMaxDevices]struct {
		id  uint16
		opt uint32
	}
}

// Socket implements a HCI User Channel as ReadWriteCloser.
type Socket struct {
	fd     int
	closed chan struct{}
	rmu    sync.Mutex
	wmu    sync.Mutex
}

// NewSocket returns a HCI User Channel of specified device id.
// If id is -1, the first available HCI device is returned.
func NewSocket(id int) (*Socket, error) {
	var err error
	// Create RAW HCI Socket.
	fd, err := unix.Socket(unix.AF_BLUETOOTH, unix.SOCK_RAW, unix.BTPROTO_HCI)
	if err != nil {
		return nil, errors.Wrap(err, "can't create socket")
	}

	if id != -1 {
		return open(fd, id)
	}

	req := devListRequest{devNum: hciMaxDevices}
	if err = ioctl(uintptr(fd), hciGetDeviceList, uintptr(unsafe.Pointer(&req))); err != nil {
		return nil, errors.Wrap(err, "can't get device list")
	}
	var msg string
	for id := 0; id < int(req.devNum); id++ {
		s, err := open(fd, id)
		if err == nil {
			return s, nil
		}
		msg = msg + fmt.Sprintf("(hci%d: %s)", id, err)
	}
	return nil, errors.Errorf("no devices available: %s", msg)
}

func open(fd, id int) (*Socket, error) {
	// Reset the device in case previous session didn't cleanup properly.
	if err := ioctl(uintptr(fd), hciDownDevice, uintptr(id)); err != nil {
		return nil, errors.Wrap(err, "can't down device")
	}
	if err := ioctl(uintptr(fd), hciUpDevice, uintptr(id)); err != nil {
		return nil, errors.Wrap(err, "can't up device")
	}

	// HCI User Channel requires exclusive access to the device.
	// The device has to be down at the time of binding.
	if err := ioctl(uintptr(fd), hciDownDevice, uintptr(id)); err != nil {
		return nil, errors.Wrap(err, "can't down device")
	}

	// Bind the RAW socket to HCI User Channel
	sa := unix.SockaddrHCI{Dev: uint16(id), Channel: unix.HCI_CHANNEL_USER}
	if err := unix.Bind(fd, &sa); err != nil {
		return nil, errors.Wrap(err, "can't bind socket to hci user channel")
	}

	// poll for 20ms to see if any data becomes available, then clear it
	pfds := []unix.PollFd{{Fd: int32(fd), Events: unix.POLLIN}}
	_, err := unix.Poll(pfds, 20)
	if err != nil {
		log.Error(err)
	}

	if pfds[0].Revents&unix.POLLIN > 0 {
		b := make([]byte, 100)
		_, err = unix.Read(fd, b)
		if err != nil {
			log.Error(err)
		}
	}

	return &Socket{fd: fd, closed: make(chan struct{})}, nil
}

func (s *Socket) Read(p []byte) (int, error) {
	s.rmu.Lock()
	n, err := unix.Read(s.fd, p)
	s.rmu.Unlock()
	// Close always sends a dummy command to wake up Read
	// bad things happen to the HCI state machines if they receive
	// a reply from that command, so make sure no data is returned
	// on a closed socket.
	//
	// note that if Write and Close are called concurrently it's
	// indeterminate which replies get through.
	select {
	case <-s.closed:
		return 0, io.EOF
	default:
	}
	return n, errors.Wrap(err, "can't read hci socket")
}

func (s *Socket) Write(p []byte) (int, error) {
	s.wmu.Lock()
	defer s.wmu.Unlock()
	n, err := unix.Write(s.fd, p)
	return n, errors.Wrap(err, "can't write hci socket")
}

func (s *Socket) Close() error {
	close(s.closed)
	_, err := s.Write([]byte{0x01, 0x09, 0x10, 0x00}) // no-op command to wake up the Read call if it's blocked
	if err != nil {
		log.Error(err)
	}

	s.rmu.Lock()
	defer s.rmu.Unlock()
	return errors.Wrap(unix.Close(s.fd), "can't close hci socket")
}

//Up turn up a HCI device by ID
func Up(id int) error {
	// Create RAW HCI Socket.
	fd, err := unix.Socket(unix.AF_BLUETOOTH, unix.SOCK_RAW, unix.BTPROTO_HCI)
	if err != nil {
		return errors.Wrap(err, "can't create socket")
	}
	if err := ioctl(uintptr(fd), hciUpDevice, uintptr(id)); err != nil {
		return errors.Wrap(err, "can't down device")
	}
	return unix.Close(fd)
}

//Down turn down a HCI device by ID
func Down(id int) error {
	// Create RAW HCI Socket.
	fd, err := unix.Socket(unix.AF_BLUETOOTH, unix.SOCK_RAW, unix.BTPROTO_HCI)
	if err != nil {
		return errors.Wrap(err, "can't create socket")
	}
	if err := ioctl(uintptr(fd), hciDownDevice, uintptr(id)); err != nil {
		return errors.Wrap(err, "can't down device")
	}
	return unix.Close(fd)
}

//List List HCI devices
func List() ([]int, error) {

	var err error

	// Create RAW HCI Socket.
	fd, err := unix.Socket(unix.AF_BLUETOOTH, unix.SOCK_RAW, unix.BTPROTO_HCI)
	if err != nil {
		return nil, errors.Wrap(err, "can't create socket")
	}

	req := devListRequest{devNum: hciMaxDevices}
	if err = ioctl(uintptr(fd), hciGetDeviceList, uintptr(unsafe.Pointer(&req))); err != nil {
		return nil, errors.Wrap(err, "can't get device list")
	}

	list := make([]int, 0)
	for id := 0; id < int(req.devNum); id++ {
		list = append(list, id)
	}

	return list, nil
}
