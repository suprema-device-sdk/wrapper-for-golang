package main

import (
	"log"

	devicesdk "github.com/suprema-device-sdk/wrapper-for-golang/windows-wrapper"
)

const (
	FOUNDED      uint8 = 1
	ACCEPTED     uint8 = 2
	CONNECTED    uint8 = 3
	DISCONNECTED uint8 = 4
)

type ConnectionStatus struct {
	DeviceId uint32
	Status   uint8
}

type DeviceInterface struct {
	Context                 uintptr
	ChannelConnectionStatus chan<- ConnectionStatus
	ChannelDevicesDiscover  chan<- []uint32
	ChannelEventReceived    chan<- devicesdk.BS2Event
}

func NewDeviceInterface(channelConnectionStatus chan<- ConnectionStatus, channelDevicesDiscover chan<- []uint32, channelEentReceived chan<- devicesdk.BS2Event) *DeviceInterface {
	d := new(DeviceInterface)
	d.ChannelConnectionStatus = channelConnectionStatus
	d.ChannelDevicesDiscover = channelDevicesDiscover
	d.ChannelEventReceived = channelEentReceived

	//runtime.SetFinalizer(d, func(d *DeviceInterface) { devicesdk.ReleaseContext(d.Context); d.Context = uintptr(0) })
	return d
}

func (d *DeviceInterface) Init() {
	if d.Context != 0 {
		devicesdk.ReleaseContext(d.Context)
	}

	context := devicesdk.AllocateContext()
	if ret := devicesdk.Initialize(context); devicesdk.BS_SDK_SUCCESS != ret {
		log.Printf("Error occurred on initializing(Error:%d)", ret)
		return
	}
	d.Context = context
	log.Printf("Successfully initialized  : Context = %v", d.Context)

	deviceFounded := func(devId devicesdk.BS2_DEVICE_ID) uintptr {
		d.ChannelConnectionStatus <- ConnectionStatus{devId, FOUNDED}
		return 0
	}
	deviceAccepted := func(devId devicesdk.BS2_DEVICE_ID) uintptr {
		d.ChannelConnectionStatus <- ConnectionStatus{devId, ACCEPTED}
		return 0
	}
	deviceConnected := func(devId devicesdk.BS2_DEVICE_ID) uintptr {
		d.ChannelConnectionStatus <- ConnectionStatus{devId, CONNECTED}
		return 0
	}
	deviceDisconnected := func(devId devicesdk.BS2_DEVICE_ID) uintptr {
		d.ChannelConnectionStatus <- ConnectionStatus{devId, DISCONNECTED}
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

	// typedef void (*OnAlarmFired)(BS2_DEVICE_ID deviceId, const BS2Event* event);
	// typedef void (*OnInputDetected)(BS2_DEVICE_ID deviceId, const BS2Event* event);
	// typedef void (*OnConfigChanged)(BS2_DEVICE_ID deviceId, uint32_t configMask);

	if ret := devicesdk.SetDeviceEventListener(context, deviceFounded, deviceAccepted, deviceConnected, deviceDisconnected); devicesdk.BS_SDK_SUCCESS != ret {
		log.Printf("could not set device-event-listener : %v", ret)
	}

	if ret := devicesdk.SetNotificationListener(context, alarmFired, inputDetected, configChanged); devicesdk.BS_SDK_SUCCESS != ret {
		log.Printf("could not set notification-listener :%v ", ret)
	}
}

func (d *DeviceInterface) Version() string {
	return devicesdk.Version()
}

func (d *DeviceInterface) Discover() {
	if ret := devicesdk.SearchDevices(d.Context); devicesdk.BS_SDK_SUCCESS != ret {
		log.Printf("Error occurred on searching devices: %v", ret)
	}

	var deviceList []uint32
	if ret := devicesdk.GetDevices(d.Context, &deviceList); devicesdk.BS_SDK_SUCCESS != ret {
		log.Printf("Error occurred on getting devices : %v", ret)
	}
	d.ChannelDevicesDiscover <- deviceList
}

func (d *DeviceInterface) Connect(devId uint32) {

	if err := devicesdk.ConnectDevice(d.Context, devId); devicesdk.BS_SDK_SUCCESS != err {
		log.Printf("Error occurred on connecting the device(%d): %v", devId, err)
	}

}

func (d *DeviceInterface) ConnectViaIP(ipAddr string, port uint16) {

	var deviceId uint32
	if err := devicesdk.ConnectDeviceViaIP(d.Context, ipAddr, port, &deviceId); devicesdk.BS_SDK_SUCCESS != err {
		log.Printf("Error occurred on connecting the device(%s:%d)", ipAddr, port)
	}

}

func (d *DeviceInterface) Disconnect(devId uint32) {

	if err := devicesdk.DisconnectDevice(d.Context, devId); devicesdk.BS_SDK_SUCCESS != err {
		log.Printf("Error occurred on disconnecting the device(%d): %v", devId, err)
	}
}

func (d *DeviceInterface) StartMonitoring(devId uint32) {

	fnOnReceived := func(deviceId uint32, event *devicesdk.BS2Event) uintptr {
		d.ChannelEventReceived <- *event
		return uintptr(0)
	}

	if ret := devicesdk.StartMonitoringLog(d.Context, devId, fnOnReceived); devicesdk.BS_SDK_SUCCESS != ret {
		log.Printf("could not start monigoring log : %v", ret)
	}
}

func (d *DeviceInterface) StopMonitoring(devId uint32) {

	if ret := devicesdk.StopMonitoringLog(d.Context, devId); devicesdk.BS_SDK_SUCCESS != ret {
		log.Printf("could not stop monitoring log : %v", ret)
	}
}

func (d *DeviceInterface) GetDeviceInfo(devId uint32) {
	var deviceInfo devicesdk.BS2SimpleDeviceInfo

	if ret := devicesdk.GetDeviceInfo(d.Context, devId, &deviceInfo); devicesdk.BS_SDK_SUCCESS != ret {
		log.Printf("could not stop monitoring log : %v", ret)
	} else {
		log.Printf("%d device info : %v", devId, deviceInfo)

	}

}

func (d *DeviceInterface) ClearDatabase(devId uint32) {
	if ret := devicesdk.ClearDatabase(d.Context, devId); devicesdk.BS_SDK_SUCCESS != ret {
		log.Printf("could not stop monitoring log : %v", ret)
	}
}

func (d *DeviceInterface) FactoryReset(devId uint32) {
	if ret := devicesdk.FactoryReset(d.Context, devId); devicesdk.BS_SDK_SUCCESS != ret {
		log.Printf("could not stop monitoring log : %v", ret)
	}
}

func (d *DeviceInterface) RebootDevice(devId uint32) {
	if ret := devicesdk.RebootDevice(d.Context, devId); devicesdk.BS_SDK_SUCCESS != ret {
		log.Printf("could not stop monitoring log : %v", ret)
	}
}

func (d *DeviceInterface) LockDevice(devId uint32) {
	if ret := devicesdk.LockDevice(d.Context, devId); devicesdk.BS_SDK_SUCCESS != ret {
		log.Printf("could not stop monitoring log : %v", ret)
	}
}

func (d *DeviceInterface) UnlockDevice(devId uint32) {
	if ret := devicesdk.UnlockDevice(d.Context, devId); devicesdk.BS_SDK_SUCCESS != ret {
		log.Printf("could not stop monitoring log : %v", ret)
	}
}
