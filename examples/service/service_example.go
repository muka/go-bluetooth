package service_example

import (
	"time"

	"github.com/muka/go-bluetooth/api/service"
	log "github.com/sirupsen/logrus"
)

func Run(adapterID string, mode string, hwaddr string) error {

	log.SetLevel(log.TraceLevel)

	if mode == "client" {
		return client(adapterID, hwaddr)
	} else {
		return serve(adapterID)
	}
}

func serve(adapterID string) error {

	a, err := service.NewApp(adapterID)
	if err != nil {
		return err
	}
	defer a.Close()

	a.SetName("go_bluetooth")

	log.Infof("HW address %s", a.Adapter().Properties.Address)

	service1, err := a.NewService()
	if err != nil {
		return err
	}

	char1, err := service1.NewChar()
	if err != nil {
		return err
	}

	char1.OnRead(service.CharReadCallback(func(c *service.Char, options map[string]interface{}) ([]byte, error) {
		log.Warnf("GOT READ REQUEST")
		return []byte{42}, nil
	}))
	char1.OnWrite(service.CharWriteCallback(func(c *service.Char, value []byte) ([]byte, error) {
		log.Warnf("GOT WRITE REQUEST")
		return value, nil
	}))

	err = service1.AddChar(char1)
	if err != nil {
		return err
	}

	err = a.AddService(service1)
	if err != nil {
		return err
	}

	err = a.Run()
	if err != nil {
		return err
	}

	log.Infof("Exposed service %s", service1.Properties.UUID)

	timeout := uint32(6 * 3600) // 6h
	log.Infof("Advertising for %ds...", timeout)
	cancel, err := a.Advertise(timeout)
	if err != nil {
		return err
	}

	defer cancel()

	wait := make(chan bool)
	go func() {
		time.Sleep(time.Duration(timeout) * time.Second)
		wait <- true
	}()

	<-wait

	return nil
}
