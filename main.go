package main

import (
	"fmt"

	"github.com/ingonyama-zk/icicle-gnark/v3/wrappers/golang/runtime"

)

func main() {

	devices, e := runtime.GetRegisteredDevices()
	if e != runtime.Success {
		panic("Failed to load registered devices")
	}

	fmt.Println("len(devices): ", len(devices))
	fmt.Println("Devices: ", devices)

}