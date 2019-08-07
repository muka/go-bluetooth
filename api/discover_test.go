package api

//
// func TestDiscoverDevice(t *testing.T) {
//
// 	a, err := GetDefaultAdapter()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	discovery, cancel, err := Discover(a, nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	defer cancel()
//
// 	wait := make(chan error)
//
// 	go func() {
// 		for dev := range discovery {
// 			if dev == nil {
// 				return
// 			}
// 			wait <- nil
// 		}
// 	}()
//
// 	go func() {
// 		sleep := 5
// 		time.Sleep(time.Duration(sleep) * time.Second)
// 		log.Debugf("Discovery timeout exceeded (%ds)", sleep)
// 		wait <- nil
// 	}()
//
// 	err = <-wait
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// }
