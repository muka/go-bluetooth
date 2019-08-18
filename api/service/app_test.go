package service

import (
	"testing"

	"github.com/muka/go-bluetooth/api"
	log "github.com/sirupsen/logrus"
)

func createTestApp(t *testing.T) *App {

	log.SetLevel(log.TraceLevel)

	a, err := NewApp(api.GetDefaultAdapterID())
	if err != nil {
		t.Fatal(err)
	}

	s1, err := a.NewService()
	if err != nil {
		t.Fatal(err)
	}

	c1, err := s1.NewChar()
	if err != nil {
		t.Fatal(err)
	}

	err = s1.AddChar(c1)
	if err != nil {
		t.Fatal(err)
	}

	c1.
		OnRead(CharReadCallback(func(c *Char, options map[string]interface{}) ([]byte, error) {
			return nil, nil
		})).
		OnWrite(CharWriteCallback(func(c *Char, value []byte) ([]byte, error) {
			return nil, nil
		}))

	d1, err := c1.NewDescr()
	if err != nil {
		t.Fatal(err)
	}

	err = c1.AddDescr(d1)
	if err != nil {
		t.Fatal(err)
	}

	err = a.AddService(s1)
	if err != nil {
		t.Fatal(err)
	}

	err = a.Run()
	if err != nil {
		t.Fatal(err)
	}

	return a
}

func TestApp(t *testing.T) {
	a := createTestApp(t)
	defer a.Close()
}
