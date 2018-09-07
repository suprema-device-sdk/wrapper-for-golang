package main

import (
	"log"
	"sync"
	"time"

	devicesdk "github.com/suprema-device-sdk/wrapper-for-golang/windows-wrapper"
)

func main() {
	var wg sync.WaitGroup

	const maxTimeout int = 100

	wg.Add(1)
	go func(w *sync.WaitGroup, delay_ss int) {
		time.Sleep(time.Duration(delay_ss) * time.Second)
		wg.Done()
	}(&wg, maxTimeout)
	log.Printf("started. will be terminated after %d seconds.", maxTimeout)

	context := devicesdk.AllocateContext()
	log.Printf("allcated(context=%v).", context)
	defer devicesdk.ReleaseContext(context)

	if ret := devicesdk.Initialize(context); devicesdk.BS_SDK_SUCCESS != ret {
		log.Printf("Error occurred on initializing(Error:%d)", ret)
		return
	}
	log.Printf("initialized. ")

	fnOnReceived := func(deviceId uint32, event *devicesdk.BS2Event) uintptr {
		log.Printf("log received(%d) : %v ", deviceId, *event)
		return uintptr(0)
	}

	deviceFounded := func(devId devicesdk.BS2_DEVICE_ID) uintptr {
		log.Printf("deviceFounded: %v", devId)
		return 0
	}
	deviceAccepted := func(devId devicesdk.BS2_DEVICE_ID) uintptr {
		log.Printf("deviceAccepted: %v", devId)
		devicesdk.StartMonitoringLog(context, devId, fnOnReceived)
		log.Printf("start monitoring log from %d", devId)
		return 0
	}
	deviceConnected := func(devId devicesdk.BS2_DEVICE_ID) uintptr {
		log.Printf("deviceConnected: %v", devId)
		devicesdk.StartMonitoringLog(context, devId, fnOnReceived)
		log.Printf("start monitoring log from %d", devId)
		return 0
	}
	deviceDisconnected := func(devId devicesdk.BS2_DEVICE_ID) uintptr {
		log.Printf("deviceDisconnected: %v", devId)
		devicesdk.StopMonitoringLog(context, devId)
		log.Printf("stop monitoring log from %d", devId)
		return 0
	}
	alarmFired := func(devId devicesdk.BS2_DEVICE_ID, event *devicesdk.BS2Event) uintptr {
		log.Printf("alarmFired: %v", *event)
		return 0
	}
	inputDetected := func(devId devicesdk.BS2_DEVICE_ID, event *devicesdk.BS2Event) uintptr {
		log.Printf("inputDetected: %v", *event)
		return 0
	}
	configChanged := func(devId devicesdk.BS2_DEVICE_ID, configMask uint32) uintptr {
		log.Printf("configChanged: Device ID - %d , configmask = 0x%02x", devId, configMask)
		return 0
	}

	if ret := devicesdk.SetDeviceEventListener(context, deviceFounded, deviceAccepted, deviceConnected, deviceDisconnected); devicesdk.BS_SDK_SUCCESS != ret {
		log.Printf("could not set device-event-listener : %v", ret)
	}
	log.Printf("registered device-event-listener. ")

	if ret := devicesdk.SetNotificationListener(context, alarmFired, inputDetected, configChanged); devicesdk.BS_SDK_SUCCESS != ret {
		log.Printf("could not set notification-listener :%v ", ret)
	}

	log.Printf("registered notification-listener. ")

	wg.Wait()

}
