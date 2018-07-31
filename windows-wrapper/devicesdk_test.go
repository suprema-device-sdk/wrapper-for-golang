package devicesdk

import (
	"bytes"
	"os"
	"sync"
	"testing"
	"time"
)

type TestParameter struct {
	DeviceID   uint32
	IPAddr     string
	Port       uint16
	DeviceType uint32
	SDKVersion string
	ServerPort uint16
}

var parameter = &TestParameter{DeviceID: 541530986, IPAddr: "192.168.16.104", Port: 51211, DeviceType: 0x09 /*BS2_DEVICE_TYPE_BIOSTATION_A2*/, SDKVersion: "2.6.1.20", ServerPort: 51212}

func TestArguments(t *testing.T) {
	t.Logf("%v", os.Args)
}
func TestBytes(t *testing.T) {
	var bs []byte = []byte{50, 46, 54, 46, 49, 46, 50, 48}

	for i, b := range bs {
		t.Logf("%02d = %v", i, b)
	}
	var str string = string(bs)

	t.Logf("str = %v", str)
}

func TestVersion(t *testing.T) {
	version := Version()
	if testing.Verbose() {
		t.Logf("SDK Version = %s", version)
	}
	if version != parameter.SDKVersion {
		t.Error("version missmatched :" + version)
	}
}

func TestAllocteAndReleaseContext(t *testing.T) {
	sdkContext := AllocateContext()
	defer ReleaseContext(sdkContext)
	if testing.Verbose() {
		t.Logf("sdkContext = %v", sdkContext)
	}
}

func TestAllocteInitializeRelease(t *testing.T) {
	sdkContext := AllocateContext()
	defer ReleaseContext(sdkContext)

	if testing.Verbose() {
		t.Logf("sdkContext = %v", sdkContext)
	}
	if ret := Initialize(sdkContext); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

}

func TestMakePinCode(t *testing.T) {
	sdkContext := AllocateContext()
	defer ReleaseContext(sdkContext)

	plain := string([]byte("aBcDefgHI012345678Q!#@$!"))

	buf := new(bytes.Buffer)
	if ret := MakePinCode(sdkContext, plain, buf); BS_SDK_SUCCESS != ret {
		t.Errorf(" Error occurred : %v", ret)
	}
}

func TestSetMaxThreadCount(t *testing.T) {
	sdkContext := AllocateContext()
	ReleaseContext(sdkContext)

	for i := 0; i < 3; i++ {
		if BS_SDK_ERROR_INVALID_PARAM != (int16)(SetMaxThreadCount(sdkContext, 0)) {
			t.Errorf("expected invalid param, but excuted successfully.")
		}
	}

	for i := 3; i < 100; i++ {
		if BS_SDK_SUCCESS != SetMaxThreadCount(sdkContext, uint32(i)) {
			t.Errorf("max thread count can not set with %d ", i)
		}
	}

}

func TestComputeCRC16CCITT(t *testing.T) {
	data := []byte("1234ajsfdklwejiro!#@$")
	// buffer := new( bytes.Buffer)
	var checksum uint16
	if ret := ComputeCRC16CCITT(data, &checksum); BS_SDK_SUCCESS != ret {
		t.Errorf("can not compute CRC: %v", data)
	}
	if testing.Verbose() {
		t.Logf("[RESULT] : %x", checksum)
	}
}

func TestGetCardModel(t *testing.T) {
	type DeviceCardModel struct {
		Name        string
		CardModel   BS2_CARD_MODEL
		ReturnValue int16
	}

	devices := []DeviceCardModel{{"BSA2-OMPW", BS2_CARD_MODEL_OMPW, BS_SDK_SUCCESS},
		{"BSA2-OIPW", BS2_CARD_MODEL_OIPW, BS_SDK_SUCCESS},
		{"BSA2-OHPW", BS2_CARD_MODEL_OHPW, BS_SDK_SUCCESS},
		//{"BSA2-ODPW",BS2_CARD_MODEL_ODPW, BS_SDK_ERROR_NOT_SUPPORTED},
		//{"BSA2-OAPW",BS2_CARD_MODEL_OAPW, BS_SDK_ERROR_NOT_SUPPORTED},
		{"BS2-OMPW", BS2_CARD_MODEL_OMPW, BS_SDK_SUCCESS},
		{"BS2-OIPW", BS2_CARD_MODEL_OIPW, BS_SDK_SUCCESS},
		{"BS2-OHPW", BS2_CARD_MODEL_OHPW, BS_SDK_SUCCESS},
		//{"BS2-ODPW",BS2_CARD_MODEL_ODPW, BS_SDK_ERROR_NOT_SUPPORTED},
		//{"BS2-OAPW",BS2_CARD_MODEL_OAPW, BS_SDK_ERROR_NOT_SUPPORTED}
	}
	for _, d := range devices {
		var cardModel BS2_CARD_MODEL
		if ret := GetCardModel(d.Name, &cardModel); d.ReturnValue != ret {
			t.Fatalf("can not get card model :%v (errno =%v)", d.Name, ret)
		}

		if cardModel != d.CardModel {
			t.Errorf("got %v, but expected %v for %s", cardModel, d.CardModel, d.Name)
		}
	}

}

func TestGetDataEncryptKey(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred on initializing : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Errorf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}

	//keyInfo := new(BS2EncryptKey)
	var keyInfo BS2EncryptKey
	if ret := GetDataEncryptKey(context, deviceId, &keyInfo); ret != BS_SDK_SUCCESS {
		if ret == BS_WRAPPER_ERROR_NOT_SUPPORTED {
			t.Logf("this sdk is not supported for apis related to data-encryption.")
			return
		}
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("[Result] %v", keyInfo)
	}
}

func TestGetDataEncryptKeyBeforeConnected(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred on initializing : %v", ret)
	}

	if ret := SearchDevicesEx(context, parameter.IPAddr); BS_SDK_SUCCESS != ret {
		t.Fatalf(" Error occurred on searching Devices : %v", ret)
	}

	var deviceList []uint32
	if ret := GetDevices(context, &deviceList); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred on getting devices : %v", ret)
	}

	if 0 == len(deviceList) {
		t.Errorf("got empty result")
	}

	keyInfo := new(BS2EncryptKey)
	if ret := GetDataEncryptKey(context, deviceList[0], keyInfo); BS_SDK_ERROR_SOCKET_IS_NOT_CONNECTED != BS_SDK_SUCCESS {
		if ret == BS_WRAPPER_ERROR_NOT_SUPPORTED {
			t.Logf("this sdk is not supported for apis related to data-encryption.")
			return
		}
		t.Errorf("exptected BS_SDK_ERROR_SOCKET_IS_NOT_CONNECTED, but got %v", ret)
	}
}

func TestSetDataEncryptKey(t *testing.T) {
	t.SkipNow()
	context := AllocateContext()
	defer ReleaseContext(context)

	keyInfo := new(BS2EncryptKey)

	if ret := SetDataEncryptKey(context, parameter.DeviceID, keyInfo); ret != BS_SDK_SUCCESS {
		if BS_WRAPPER_ERROR_NOT_SUPPORTED == ret {
			ReleaseContext(context)
			return
		}
		t.Fatalf(" Error occurred : %v", ret)
	}
}

func TestRemoveDataEncryptKey(t *testing.T) {
	t.SkipNow()
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := RemoveDataEncryptKey(context, parameter.DeviceID); BS_SDK_SUCCESS != ret {
		if BS_WRAPPER_ERROR_NOT_SUPPORTED == ret {
			ReleaseContext(context)
			return
		}
		t.Fatalf("Error occurred : %v", ret)
	}
}

func TestSetDeviceEventListener(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	onDeviceFound := func(deviceId uint32) uintptr {
		if testing.Verbose() {
			t.Logf("found a device(id=%v)", deviceId)
		}
		return uintptr(0)
	}
	onDeviceAccepted := func(deviceId uint32) uintptr {
		if testing.Verbose() {
			t.Logf("accepted a device(id=%v)", deviceId)
		}
		return uintptr(0)
	}
	onDeviceConnected := func(deviceId uint32) uintptr {
		if testing.Verbose() {
			t.Logf("connected a device(id=%v)", deviceId)
		}
		return uintptr(0)
	}
	onDeviceDisconnected := func(deviceId uint32) uintptr {
		if testing.Verbose() {
			t.Logf("disconnected a device(id=%v)", deviceId)
		}
		return uintptr(0)
	}

	if ret := SetDeviceEventListener(context, onDeviceFound, onDeviceAccepted, onDeviceConnected, onDeviceDisconnected); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred on setting device-event-listener : %v", ret)
	}
}

func TestSearchDevices(t *testing.T) {

	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	onDeviceFound := func(deviceId uint32) uintptr {
		if testing.Verbose() {
			t.Logf("CALLBACK:found a device(id=%v)", deviceId)
		}
		return uintptr(0)
	}

	if ret := SetDeviceEventListener(context, onDeviceFound, nil, nil, nil); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred on setting device-event-listener : %v", ret)
	}

	if ret := SearchDevices(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred on searching devices: %v", ret)
	}

	var deviceList []uint32
	if ret := GetDevices(context, &deviceList); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred on getting devices : %v", ret)
	}

	if 0 == len(deviceList) {
		t.Errorf("got empty result")
	}
	if testing.Verbose() {
		for i, d := range deviceList {
			t.Logf("[RESULT]found device-%d : %d\n", i, d)
		}
	}

}

func TestSearchDevicesEx(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	onDeviceFound := func(deviceId uint32) uintptr {
		if testing.Verbose() {
			t.Logf("CALLBACK:found a device(id=%v)", deviceId)
		}
		return uintptr(0)
	}

	if ret := SetDeviceEventListener(context, onDeviceFound, nil, nil, nil); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred on setting device-event-listener : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("trying to SearchDevicesEx(%s)... ", parameter.IPAddr)
	}
	if ret := SearchDevicesEx(context, parameter.IPAddr); BS_SDK_SUCCESS != ret {
		t.Fatalf(" Error occurred : %v", ret)
	}

	var deviceList []uint32
	if ret := GetDevices(context, &deviceList); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == len(deviceList) {
		t.Errorf("got empty result")
	}
	if testing.Verbose() {
		for i, d := range deviceList {
			t.Logf("[RESULT]found device-%02d : %d\n", i, d)
		}
	}

}

func TestGetDevices(t *testing.T) {

	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	onDeviceFound := func(deviceId uint32) uintptr { t.Logf("CALLBACK:found a device(id=%v)", deviceId); return uintptr(0) }

	if ret := SetDeviceEventListener(context, onDeviceFound, nil, nil, nil); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred on setting device-event-listener : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("trying to SearchDevicesEx(%s)... ", parameter.IPAddr)
	}
	if ret := SearchDevicesEx(context, parameter.IPAddr); BS_SDK_SUCCESS != ret {
		t.Fatalf(" Error occurred : %v", ret)
	}

	var array []uint32
	if ret := GetDevices(context, &array); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if testing.Verbose() {

		for i, a := range array {
			t.Logf("[RESULT] devices-%02d : %d", i, a)
		}
	}

}

func TestConnectDevice(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred on initializing : %v", ret)
	}

	deviceFounded := func(devId BS2_DEVICE_ID) uintptr {
		if testing.Verbose() {
			t.Logf("CALLBACK:device found(ID=%d).", devId)
		}
		return 0
	}
	deviceAccepted := func(devId BS2_DEVICE_ID) uintptr {
		if testing.Verbose() {
			t.Logf("CALLBACK:device accepted(ID=%d).", devId)
		}
		return 0
	}
	deviceConnected := func(devId BS2_DEVICE_ID) uintptr {
		if testing.Verbose() {
			t.Logf("CALLBACK:device connted(ID=%d).", devId)
		}
		return 0
	}
	deviceDisconnected := func(devId BS2_DEVICE_ID) uintptr {
		if testing.Verbose() {
			t.Logf("CALLBACK:device disconected(ID=%d).", devId)
		}
		return 0
	}

	if ret := SetDeviceEventListener(context, deviceFounded, deviceAccepted, deviceConnected, deviceDisconnected); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred on Setting device-event-listener : %v", ret)
	}

	t.Logf("trying to SearchDevicesEx(%s)... ", parameter.IPAddr)
	if ret := SearchDevicesEx(context, parameter.IPAddr); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred on serching devices : %v", ret)
	}

	var array []uint32
	if ret := GetDevices(context, &array); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred on getting devices : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("[RESULT] connecting to the device(ID=%d)", array[0])
	}

	if ret := ConnectDevice(context, array[0]); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred on connecting a device : %v", ret)
	}

}

func TestConnectDeviceViaIP(t *testing.T) {

	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	deviceFounded := func(devId BS2_DEVICE_ID) uintptr {
		if testing.Verbose() {
			t.Logf("CALLBACK:device found(ID=%d).", devId)
		}
		return 0
	}
	deviceAccepted := func(devId BS2_DEVICE_ID) uintptr {
		if testing.Verbose() {
			t.Logf("CALLBACK:device accepted(ID=%d).", devId)
		}
		return 0
	}
	deviceConnected := func(devId BS2_DEVICE_ID) uintptr {
		if testing.Verbose() {
			t.Logf("CALLBACK:device connted(ID=%d).", devId)
		}
		return 0
	}
	deviceDisconnected := func(devId BS2_DEVICE_ID) uintptr {
		if testing.Verbose() {
			t.Logf("CALLBACK:device disconected(ID=%d).", devId)
		}
		return 0
	}

	if ret := SetDeviceEventListener(context, deviceFounded, deviceAccepted, deviceConnected, deviceDisconnected); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred on Setting device-event-listener : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Errorf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	if testing.Verbose() {
		t.Logf("[RESULT] connected with the device(ID=%d)", deviceId)
	}

}

func TestDisconnectDevice(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	deviceFounded := func(devId BS2_DEVICE_ID) uintptr {
		if testing.Verbose() {
			t.Logf("CALLBACK:device found(ID=%d).", devId)
		}
		return 0
	}
	deviceAccepted := func(devId BS2_DEVICE_ID) uintptr {
		if testing.Verbose() {
			t.Logf("CALLBACK:device accepted(ID=%d).", devId)
		}
		return 0
	}
	deviceConnected := func(devId BS2_DEVICE_ID) uintptr {
		if testing.Verbose() {
			t.Logf("CALLBACK:device connted(ID=%d).", devId)
		}
		return 0
	}
	deviceDisconnected := func(devId BS2_DEVICE_ID) uintptr {
		if testing.Verbose() {
			t.Logf("CALLBACK:device disconected(ID=%d).", devId)
		}
		return 0
	}

	if ret := SetDeviceEventListener(context, deviceFounded, deviceAccepted, deviceConnected, deviceDisconnected); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred on Setting device-event-listener : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if ret := DisconnectDevice(context, deviceId); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred on disconnecting a device(ID=%d) : %v", deviceId, ret)
	}
}

func TestSetKeepAliveTimeout(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if ret := SetKeepAliveTimeout(context, 360); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestSetNotificationListener(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}
	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
		return
	}
	defer DisconnectDevice(context, deviceId)

	wg := new(sync.WaitGroup)
	wg.Add(1)

	go func() {
		defer wg.Done()
		time.Sleep(15 * time.Second)
	}()

	fnOnAlarmFired := func(deviceId uint32, event *BS2Event) uintptr {
		if testing.Verbose() {
			t.Logf("Alarm-fired occurred on a device(%d):%v", deviceId, *event)
		}
		return uintptr(0)
	}
	fnOnInputDetected := func(deviceId uint32, event *BS2Event) uintptr {
		if testing.Verbose() {
			t.Logf("Input-detected occurred on a devcie(%d):%v", deviceId, *event)
		}
		return uintptr(0)
	}
	fnOnConfigChanged := func(deviceId uint32, configMask uint32) uintptr {
		if testing.Verbose() {
			t.Logf("Config-changed occurred on a device(%d):%x", deviceId, configMask)
		}
		return uintptr(0)
	}
	if ret := SetNotificationListener(context, fnOnAlarmFired, fnOnInputDetected, fnOnConfigChanged); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
	wg.Wait()
}

func TestSetServerPort(t *testing.T) {

	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if ret := SetServerPort(context, parameter.ServerPort); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

}

func TestSetSSLHandler(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if ret := SetSSLHandler(context, nil, nil, nil, nil, nil, nil); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

}

func TestDisableSSL(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}
	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
		return
	}
	defer DisconnectDevice(context, deviceId)

	if ret := DisableSSL(context, deviceId); BS_SDK_SUCCESS != ret {
		t.Errorf("can not make the device(ID=%d) SSL-disabled:error code =%d.", deviceId, ret)
	}
}

func TestGetDeviceInfo(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	deviceConnected := func(devId BS2_DEVICE_ID) uintptr {
		if testing.Verbose() {
			t.Logf("CALLBACK:device connted(ID=%d).", devId)
		}
		return 0
	}
	deviceDisconnected := func(devId BS2_DEVICE_ID) uintptr {
		if testing.Verbose() {
			t.Logf("CALLBACK:device disconected(ID=%d).", devId)
		}
		return 0
	}

	if ret := SetDeviceEventListener(context, nil, nil, deviceConnected, deviceDisconnected); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred on Setting device-event-listener : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
		return
	}
	defer DisconnectDevice(context, deviceId)

	simpleDeviceInfo := new(BS2SimpleDeviceInfo)
	if ret := GetDeviceInfo(context, deviceId, simpleDeviceInfo); BS_SDK_SUCCESS != ret {
		t.Fatalf("failedto get infos of a device(%d) : %v", deviceId, ret)
	}

	if deviceId != simpleDeviceInfo.Id {
		t.Errorf("device id missmatched : %v vs. %v", deviceId, simpleDeviceInfo.Id)
	}
	if testing.Verbose() {
		t.Logf("%v", *simpleDeviceInfo)
	}
}

func TestGetDeviceTime(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var gmtTime uint32
	if ret := GetDeviceTime(context, deviceId, &gmtTime); BS_SDK_SUCCESS != ret {
		t.Fatalf("failedto get the device-time(%d) : %v", deviceId, ret)
	}

	if 0 == gmtTime {
		t.Errorf("wrong deviceTime : %v", gmtTime)
	}
	if testing.Verbose() {
		t.Logf("Device Time = %d,%s", gmtTime, time.Unix(int64(gmtTime), 0).Format(time.RFC3339))
	}
}

func TestSetDeviceTime(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var gmtTime uint32
	if ret := GetDeviceTime(context, deviceId, &gmtTime); BS_SDK_SUCCESS != ret {
		t.Fatalf("failedto get the device-time(%d) : %v", deviceId, ret)
	}
	if testing.Verbose() {
		t.Logf("Device Time = %d", gmtTime)
	}

	var delta uint32 = 100
	settingTime := gmtTime + delta
	if ret := SetDeviceTime(context, deviceId, settingTime); BS_SDK_SUCCESS != ret {
		t.Errorf("failed to set device-time(%d) : %v", deviceId, ret)
	}

	var newGmtTime uint32
	if ret := GetDeviceTime(context, deviceId, &newGmtTime); BS_SDK_SUCCESS != ret {
		t.Errorf("failedto get the device-time(%d) : %v", deviceId, ret)
	}

	if (newGmtTime - gmtTime) < delta {
		t.Errorf("the new gmt-time should be more great than the previous(new=%d,old=%d)", newGmtTime, gmtTime)
	}
	if testing.Verbose() {
		t.Logf("New Device Time = %d", newGmtTime)
	}
}

func TestSetDeviceTimeWithCurrent(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	settingTime := uint32(time.Now().Unix())
	if ret := SetDeviceTime(context, deviceId, settingTime); BS_SDK_SUCCESS != ret {
		t.Errorf("failed to set device-time(%d) : %v", deviceId, ret)
	}
}

func TestClearDatabase(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	if ret := ClearDatabase(context, deviceId); BS_SDK_SUCCESS != ret {
		t.Errorf("failed to clear database fo the device(%d) : %v", deviceId, ret)
	}

}

func TestFactoryReset(t *testing.T) {
	t.SkipNow() // Skipped  because of running next test-cases
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	if ret := FactoryReset(context, deviceId); BS_SDK_SUCCESS != ret {
		t.Errorf("failed to factory-reset the device-time(%d) : %v", deviceId, ret)
	}
}

func TestRebootDevice(t *testing.T) {
	t.SkipNow() // Skipped, because of running next test-cases.
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	if ret := RebootDevice(context, deviceId); BS_SDK_SUCCESS != ret {
		t.Errorf("failed to reboot device(%d) : %v", deviceId, ret)
	}

}

func TestLockDevice(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	if ret := LockDevice(context, deviceId); BS_SDK_SUCCESS != ret {
		t.Errorf("failed to lock device(%d) : %v", deviceId, ret)
	}

	time.Sleep(time.Second * 1)

	if ret := UnlockDevice(context, deviceId); BS_SDK_SUCCESS != ret {
		t.Errorf("failed to unlock device(%d) : %v", deviceId, ret)
	}
}

func TestUnlockDevice(t *testing.T) {

	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	if ret := LockDevice(context, deviceId); BS_SDK_SUCCESS != ret {
		t.Errorf("failed to lock device(%d) : %v", deviceId, ret)
	}

	if ret := UnlockDevice(context, deviceId); BS_SDK_SUCCESS != ret {
		t.Errorf("failed to unlock device(%d) : %v", deviceId, ret)
	}

	if ret := UnlockDevice(context, deviceId); BS_SDK_SUCCESS != ret {
		t.Errorf("failed to unlock device(%d) : %v", deviceId, ret)
	}

}

func TestUpgradeFirmware(t *testing.T) {
	t.SkipNow()
	//Not implemented
}

func TestUpdateResource(t *testing.T) {
	t.SkipNow()
	//Not implemented
}

func TestGetLog(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var mutex = new(sync.RWMutex)
	var size int
	waitingGroup := new(sync.WaitGroup)

	waitingGroup.Add(1)
	go func(offset uint32, limit uint32, total uint32) {
		defer waitingGroup.Done()
		for i := 0; uint32(i)*limit < total; i++ {
			var eventLogs []BS2Event

			if ret := GetLog(context, deviceId, offset+limit*uint32(i), limit, &eventLogs); BS_SDK_SUCCESS != ret {
				t.Fatalf("can not retrieve a eventLogs(%s:%d)", parameter.IPAddr, parameter.Port)
			}

			waitingGroup.Add(1)

			go func(events []BS2Event, offset uint32) {
				defer waitingGroup.Done()

				// for i,e := range events {
				// 	t.Logf("%03d:%04d-%v",offset, i,e)
				// }
				mutex.Lock()
				size += len(events)
				mutex.Unlock()
			}(eventLogs, limit*uint32(i))
			mutex.RLock()
			if testing.Verbose() {
				t.Logf("%.1f%%", float32(size)/float32(total)*100.0)
			}
			mutex.RUnlock()
		}

	}(0, 10000, 100000) // A2 기준,총 100000개의 로그를 장치로 부터 가져 올 때, 10000개씩 10번 가져오는 것이 100,1000개씩 가져오는 것보다 빠름.

	waitingGroup.Wait()
	if testing.Verbose() {
		t.Logf("total size = %d", size)
	}
}

func TestGetFilteredLog(t *testing.T) {
	//anyone knows why there are nothing retrieved from device
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var eventLogBuffer []BS2Event
	if ret := GetFilteredLog(context, parameter.DeviceID, "1", 0, 0, 10000, 0, &eventLogBuffer); BS_SDK_SUCCESS != ret {
		t.Fatalf("can not retrieve even logs from the device(ID=%d)", parameter.DeviceID)
	}

	size := len(eventLogBuffer)
	if testing.Verbose() {
		for _, e := range eventLogBuffer {
			t.Logf("%v", e)
		}

		t.Logf("total size = %d", size)
	}
}

func TestClearLog(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	if ret := ClearLog(context, parameter.DeviceID); BS_SDK_SUCCESS != ret {
		t.Errorf("can not clear log in the device(ID=%d)", parameter.DeviceID)
	}

}

func TestStartMonitoringLog(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	fnOnReceivedLog := func(deviceId uint32, event *BS2Event) uintptr {
		if testing.Verbose() {
			t.Logf("%d - %v", deviceId, *event)
		}
		return 0
	}
	waitingGroup := new(sync.WaitGroup)

	waitingGroup.Add(1)
	go func() {
		defer waitingGroup.Done()

		if ret := StartMonitoringLog(context, parameter.DeviceID, fnOnReceivedLog); BS_SDK_SUCCESS != ret {
			t.Errorf("can not start monigoring log")
		}
		time.Sleep(10 * time.Second)

		if ret := StopMonitoringLog(context, parameter.DeviceID); BS_SDK_SUCCESS != ret {
			t.Errorf("can not stop monitoring log")
		}
	}()

	waitingGroup.Wait()

}

func TestStopMonitoringLog(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if ret := SearchDevices(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred on searching devices: %v", ret)
	}

	var deviceList []uint32
	if ret := GetDevices(context, &deviceList); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred on getting devices : %v", ret)
	}

	var wg sync.WaitGroup

	for _, d := range deviceList {
		wg.Add(1)
		go func(id BS2_DEVICE_ID) {
			defer wg.Done()

			if ret := ConnectDevice(context, id); BS_SDK_SUCCESS != ret {
				if testing.Verbose() {
					t.Logf("ignore the device(ID=%d) due to not being able to connect", id)
				}
				return
			}

			if ret := StartMonitoringLog(context, id, nil); BS_SDK_SUCCESS != ret {
				if testing.Verbose() {
					t.Logf("ignore the device(ID=%d) due to not being able to start monitoring", id)
				}
				return
			}
			time.Sleep(20 * time.Second)

			if ret := StopMonitoringLog(context, id); BS_SDK_SUCCESS != ret {
				t.Errorf("can not stop monitoring log")
			}
		}(d)
	}

	wg.Wait()
}

func TestGetLogBlob(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var eventLogs []BS2EventBlob

	if ret := GetLogBlob(context, deviceId, 0xFF, 0, 10, &eventLogs); BS_SDK_SUCCESS != ret {
		t.Errorf("can not retrieve a eventLogs(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	if testing.Verbose() {
		for _, e := range eventLogs {
			t.Logf("%v\n", e)

		}
	}
}

func TestGetFilteredLogSinceEventId(t *testing.T) {
	//Anyone knows why the result of eventLogs is empty
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var eventLogs []BS2Event
	timestamp := time.Now().Unix()
	if ret := GetFilteredLogSinceEventId(context, deviceId, "", 0, uint32(timestamp-10000), uint32(timestamp), 0, 0, 100, &eventLogs); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
	if testing.Verbose() {
		for _, e := range eventLogs {
			t.Logf("%v\n", e)
		}
	}
}

func TestGetUserList(t *testing.T) {

	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var userList []BS2_USER_ID
	if ret := GetUserList(context, deviceId, &userList, nil); BS_SDK_SUCCESS != ret {
		t.Fatalf("can not get user-list from the device(ID=%d)", deviceId)
	}

	if testing.Verbose() {
		for i, u := range userList {
			t.Logf("%03d-%v", i, u)
		}
	}
}

func TestGetUserInfos(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var userList []BS2_USER_ID
	if ret := GetUserList(context, deviceId, &userList, nil); BS_SDK_SUCCESS != ret {
		t.Fatalf("can not get user-list from the device(ID=%d)", deviceId)
	}

	if 0 == len(userList) {
		t.Fatalf("empty user from device.")
	}

	if testing.Verbose() {

		for i, u := range userList {
			t.Logf("%03d-%v", i, u)
		}
		t.Logf("count of user to request info was %d", len(userList))
	}

	userBlob := make([]BS2UserBlob, USER_PAGE_SIZE)
	if ret := GetUserInfos(context, deviceId, userList, &userBlob); BS_SDK_SUCCESS != ret {
		t.Fatalf("can not get users-info from the device(ID=%d)", deviceId)
	}

	if testing.Verbose() {
		t.Logf("%v", userBlob[0])
	}
}

func TestGetUserDatas(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var userList []BS2_USER_ID
	if ret := GetUserList(context, deviceId, &userList, nil); BS_SDK_SUCCESS != ret {
		t.Fatalf("can not get user-list from the device(ID=%d)", deviceId)
	}

	if 0 == len(userList) {
		t.Fatalf("empty user from device.")
	}

	if testing.Verbose() {

		for i, u := range userList {
			t.Logf("%03d-%v", i, u)
		}
		t.Logf("count of user to request info was %d", len(userList))
	}

	userBlob := make([]BS2UserBlob, USER_PAGE_SIZE)
	if ret := GetUserDatas(context, deviceId, userList, &userBlob, BS2_USER_MASK_ALL); BS_SDK_SUCCESS != ret {
		t.Fatalf("can not get users-info from the device(ID=%d)", deviceId)
	}

	if testing.Verbose() {
		t.Logf("%v", userBlob[0])
	}
}

func TestEnrolUserWithEmpryUserBlob(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var userBlob []BS2UserBlob
	if r := EnrolUser(context, deviceId, &userBlob, 0); BS_SDK_ERROR_NULL_POINTER != int16(r) {
		t.Errorf("expected -10000,but got %v", (int16)(r))
	}
}

func TestEnrolUserToOverwrite(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var userBlob []BS2UserBlob
	userBlob = append(userBlob, BS2UserBlob{User: BS2User{UserID: [32]byte{52, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}})

	if ret := RemoveAllUser(context, deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("can not remove users on the device(ID=%d)", deviceId)
	}

	if ret := EnrolUser(context, deviceId, &userBlob, 0); BS_SDK_SUCCESS != ret {
		t.Fatalf("can not enroll an user:%v", userBlob[0])
	}

	if ret := EnrolUser(context, deviceId, &userBlob, 1); BS_SDK_SUCCESS != ret {
		t.Errorf("expected 1,but got %v", ret)
	}
}

func TestEnrolUserToNonOverwrite(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var userBlob []BS2UserBlob
	userBlob = append(userBlob, BS2UserBlob{User: BS2User{UserID: [32]byte{52, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}})

	if r := EnrolUser(context, deviceId, &userBlob, 1); BS_SDK_SUCCESS != int16(r) {
		t.Errorf("expected 1,but got %v", (int16)(r))
	}

	if r := EnrolUser(context, deviceId, &userBlob, 0); BS_SDK_ERROR_DUPLICATE_ID != int16(r) {
		t.Errorf("expected -705,but got %v", (int16)(r))
	}
}
func TestRemoveUser(t *testing.T) {

	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	if ret := RemoveAllUser(context, deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("can not remove users on the device(ID=%d)", deviceId)
	}

	var userBlob []BS2UserBlob
	userBlob = append(userBlob, BS2UserBlob{User: BS2User{UserID: [32]byte{52, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}})

	if r := EnrolUser(context, deviceId, &userBlob, 1); BS_SDK_SUCCESS != int16(r) {
		t.Errorf("expected 1,but got %v", (int16)(r))
	}

	var userListToDelete []BS2_USER_ID
	userListToDelete = append(userListToDelete, BS2_USER_ID{52, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})

	if ret := RemoveUser(context, deviceId, userListToDelete); BS_SDK_SUCCESS != ret {
		t.Errorf("can not remove users on the device(ID=%d)", deviceId)
	}

	var userList []BS2_USER_ID
	if ret := GetUserList(context, deviceId, &userList, nil); BS_SDK_SUCCESS != ret {
		t.Fatalf("can not get user-list from the device(ID=%d)", deviceId)
	}

	if 0 != len(userList) {
		t.Errorf("failed to remove user: %v", userList)
	}

}

func TestRemoveAllUser(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}
	defer DisconnectDevice(context, deviceId)

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}

	if ret := RemoveAllUser(context, deviceId); BS_SDK_SUCCESS != ret {
		t.Errorf("can not remove users on the device(ID=%d)", deviceId)
	}

	var userList []BS2_USER_ID
	if ret := GetUserList(context, deviceId, &userList, nil); BS_SDK_SUCCESS != ret {
		t.Fatalf("can not get user-list from the device(ID=%d)", deviceId)
	}

	if 0 != len(userList) {
		t.Errorf("failed to remove all user: %v", userList)
	}
}

func TestGetUserInfosEx(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var userList []BS2_USER_ID
	if ret := GetUserList(context, deviceId, &userList, nil); BS_SDK_SUCCESS != ret {
		t.Fatalf("can not get user-list from the device(ID=%d)", deviceId)
	}

	if 0 == len(userList) {
		t.Fatalf("empty user from device.")
	}

	if testing.Verbose() {

		for i, u := range userList {
			t.Logf("%03d-%v", i, u)
		}
		t.Logf("count of user to request info was %d", len(userList))
	}

	userBlob := make([]BS2UserBlobEx, USER_PAGE_SIZE)
	if ret := GetUserInfosEx(context, deviceId, userList, &userBlob); BS_SDK_SUCCESS != ret {
		t.Fatalf("can not get users-info from the device(ID=%d)", deviceId)
	}

	if testing.Verbose() {
		t.Logf("%v", userBlob[0])
	}

}

func TestEnrolUserEx(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var userBlob []BS2UserBlobEx
	userBlob = append(userBlob, BS2UserBlobEx{User: BS2User{UserID: [32]byte{52, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}})

	if ret := RemoveAllUser(context, deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("can not remove users on the device(ID=%d)", deviceId)
	}

	if ret := EnrolUserEx(context, deviceId, &userBlob, 0); BS_SDK_SUCCESS != ret {
		t.Fatalf("can not enroll an user:%v", userBlob[0])
	}
}

func TestGetUserDatabaseInfo(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var numUsers, numCards, numFingers, numFaces uint32

	if ret := GetUserDatabaseInfo(context, deviceId, &numUsers, &numCards, &numFingers, &numFaces, nil); BS_SDK_SUCCESS != ret {
		t.Errorf("can not get user-database-info from the device(ID=%d)", deviceId)
	}

	if testing.Verbose() {
		t.Logf("Users : %d, Cards : %d, Fingers : %d, Faces : %d", numUsers, numCards, numFingers, numFaces)
	}
}

func TestResetConfig(t *testing.T) {
	t.SkipNow()
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	if ret := ResetConfig(context, deviceId, true); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestGetConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2Configs
	config.ConfigMask = 0xFFFFFFFF

	if ret := GetConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", config)
	}

}

func TestGetFactoryConfig(t *testing.T) {

	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2FactoryConfig
	if ret := GetFactoryConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", config)
	}

}

func TestGetSystemConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2SystemConfig
	if ret := GetSystemConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", config)
	}
}

func TestSetSystemConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2SystemConfig
	if ret := SetSystemConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", config)
	}
}

func TestGetAuthConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2AuthConfig
	if ret := GetAuthConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", config)
	}
}

func TestSetAuthConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2AuthConfig
	if ret := GetAuthConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	newConfig := config
	var newUseGlobalAPB uint8
	if 0 == config.UseGlobalAPB {
		newUseGlobalAPB = 1
	}
	newConfig.UseGlobalAPB = newUseGlobalAPB
	if ret := SetAuthConfig(context, deviceId, &newConfig); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if ret := GetAuthConfig(context, deviceId, &newConfig); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if newConfig.UseGlobalAPB != newUseGlobalAPB {
		t.Errorf("expected %d, but got %d", newUseGlobalAPB, newConfig.UseGlobalAPB)
	}

	if ret := SetAuthConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestGetStatusConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2StatusConfig
	if ret := GetStatusConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", config)
	}
}

func TestSetStatusConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2StatusConfig
	if ret := GetStatusConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	newConfig := config
	var newConfigSyncRequired uint8
	if 0 == config.ConfigSyncRequired {
		newConfigSyncRequired = 1
	}
	newConfig.ConfigSyncRequired = newConfigSyncRequired

	if ret := SetStatusConfig(context, deviceId, &newConfig); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if ret := GetStatusConfig(context, deviceId, &newConfig); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if newConfig.ConfigSyncRequired != newConfigSyncRequired {
		t.Errorf("expected %d, but got %d", newConfigSyncRequired, newConfig.ConfigSyncRequired)
	}

	if ret := SetStatusConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

}

func TestGetDisplayConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2DisplayConfig
	if ret := GetDisplayConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", config)
	}
}

func TestSetDisplayConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2DisplayConfig
	if ret := SetDisplayConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestGetIPConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2IpConfig
	if ret := GetIPConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", config)
	}
}

func TestGetIPConfigViaUDP(t *testing.T) {

	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}
	var deviceId uint32 = parameter.DeviceID
	// if ret := SearchDevices(context); BS_SDK_SUCCESS != ret {
	// 	t.Fatalf("Error occurred on searching devices: %v", ret)
	// }

	var config BS2IpConfig
	if ret := GetIPConfigViaUDP(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", config)
	}
}

/*
func TestGetIPConfigViaUDP(t *testing.T) {

	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}

	var config BS2IpConfig
	if ret := GetIPConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if ret := DisconnectDevice(context, deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if ret := SearchDevices(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred on searching devices: %v", ret)
	}

	var newUseDHCP bool = !config.UseDHCP
	config.UseDHCP = newUseDHCP

	if ret := GetIPConfigViaUDP(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var config BS2IpConfig
	if ret := GetIPConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}
	var newConfig BS2IpConfig
	if ret := GetIPConfig(context, deviceId, &newConfig); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if newUseDHCP != newConfig.UseDHCP {
		t.Errorf("missmatched between %v and %v ", newUseDHCP, newConfig.UseDHCP)
	}
	// if 0 != bytes.Compare(newServerAddr[:], newConfig.ServerAddr[:]) {
	// 	t.Errorf("missmatched between %v and %v ", newServerAddr, newConfig.ServerAddr)
	// }

}
*/

func TestSetIPConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred on connecting: %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2IpConfig
	if ret := GetIPConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred on getting ip config: %v", ret)
	}
	var newUseDHCP bool = !config.UseDHCP
	config.UseDHCP = newUseDHCP
	if ret := SetIPConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred on setting ip config: %v", ret)
	}

	var newConfig BS2IpConfig
	if ret := GetIPConfig(context, deviceId, &newConfig); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if newConfig.UseDHCP != newUseDHCP {
		t.Errorf("missmatched between %v and %v", newUseDHCP, newConfig.UseDHCP)
	}
}

func TestSetIPConfigViaUDP(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred on connecting : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}

	var config BS2IpConfig
	if ret := GetIPConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred on getting ip config: %v", ret)
	}

	if ret := DisconnectDevice(context, deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred on disconnecting : %v", ret)
	}

	var newUseDHCP bool = !config.UseDHCP
	config.UseDHCP = newUseDHCP

	if ret := SetIPConfigViaUDP(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred on setting ip-config via udp : %v", ret)
	}

	// time.Sleep(60 * time.Second)

	// if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
	// 	t.Fatalf("Error occurred  on re-connecting: %v", ret)
	// }
	// var newConfig BS2IpConfig
	// if ret := GetIPConfig(context, deviceId, &newConfig); BS_SDK_SUCCESS != ret {
	// 	t.Fatalf("Error occurred on getting ip-config to complete: %v", ret)
	// }

	// if newUseDHCP != newConfig.UseDHCP {
	// 	t.Errorf("missmatched between %v and %v ", newUseDHCP, newConfig.UseDHCP)
	// }
}

func TestGetIPConfigExt(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2IpConfigExt
	if ret := GetIPConfigExt(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", config)
	}
}

func TestSetIPConfigExt(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred on connecting : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2IpConfigExt
	if ret := GetIPConfigExt(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred on getting ip-configext : %v", ret)
	}

	if ret := SetIPConfigExt(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("error occurred on setting ip-configext : %v", ret)
	}
}

func TestGetTNAConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2TNAConfig
	if ret := GetTNAConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", config)
	}
}

func TestSetTNAConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2TNAConfig
	if ret := GetTNAConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
	if ret := SetTNAConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestGetCardConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2CardConfig
	if ret := GetCardConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", config)
	}
}

func TestSetCardConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2CardConfig
	if ret := GetCardConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
	if ret := SetCardConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestGetFingerprintConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2FingerprintConfig
	if ret := GetFingerprintConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", config)
	}
}

func TestSetFingerprintConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2FingerprintConfig
	if ret := GetFingerprintConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}
	if ret := SetFingerprintConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestGetRS485Config(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2Rs485Config
	if ret := GetRS485Config(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", config)
	}
}

func TestSetRS485Config(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2Rs485Config
	if ret := GetRS485Config(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
	if ret := SetRS485Config(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

// GetWiegandConfig

func TestGetWiegandConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2WiegandConfig
	if ret := GetWiegandConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", config)
	}
}

func TestSetWiegandConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2WiegandConfig
	if ret := GetWiegandConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}
	if ret := SetWiegandConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestGetWiegandDeviceConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2WiegandDeviceConfig
	if ret := GetWiegandDeviceConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", config)
	}
}

func TestSetWiegandDeviceConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2WiegandDeviceConfig
	if ret := GetWiegandDeviceConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}
	if ret := SetWiegandDeviceConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestGetInputConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2InputConfig
	if ret := GetInputConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", config)
	}
}

func TestSetInputConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2InputConfig
	if ret := GetInputConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}
	if ret := SetInputConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

//GetWlanConfig
func TestGetWlanConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2WlanConfig
	if ret := GetWlanConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", config)
	}
}

func TestSetWlanConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2WlanConfig
	if ret := GetWlanConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}
	if ret := SetWlanConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestGetTriggerActionConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2TriggerActionConfig
	if ret := GetTriggerActionConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", config)
	}
}

func TestSetTriggerActionConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2TriggerActionConfig
	if ret := GetTriggerActionConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if ret := SetTriggerActionConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

//GetEventConfig
func TestGetEventConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2EventConfig
	if ret := GetEventConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", config)
	}
}

func TestSetEventConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2EventConfig
	if ret := GetEventConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}
	if ret := SetEventConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestGetWiegandMultiConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2WiegandMultiConfig
	if ret := GetWiegandMultiConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", config)
	}
}

func TestSetWiegandMultiConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2WiegandMultiConfig
	if ret := GetWiegandMultiConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}
	if ret := SetWiegandMultiConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestGetCard1xConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS1CardConfig
	if ret := GetCard1xConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", config)
	}
}

func TestSetCard1xConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS1CardConfig
	if ret := GetCard1xConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}
	if ret := SetCard1xConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestGetSystemExtConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2SystemConfigExt
	if ret := GetSystemExtConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", config)
	}
}

func TestSetSystemExtConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2SystemConfigExt
	if ret := GetSystemExtConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}
	if ret := SetSystemExtConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestGetVoipConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2VoipConfig
	if ret := GetVoipConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", config)
	}
}

func TestSetVoipConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2VoipConfig
	if ret := GetVoipConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}
	if ret := SetVoipConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

//GetFaceConfig

func TestGetFaceConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2FaceConfig
	if ret := GetFaceConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", config)
	}
}

func TestSetFaceConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2FaceConfig
	if ret := GetFaceConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}
	if ret := SetFaceConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestGetRS485ConfigEx(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2Rs485ConfigEX
	if ret := GetRS485ConfigEx(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", config)
	}
}

func TestSetRS485ConfigEx(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2Rs485ConfigEX
	if ret := GetRS485ConfigEx(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if ret := SetRS485ConfigEx(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestGetCardConfigEx(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2CardConfigEx
	if ret := GetCardConfigEx(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", config)
	}
}

func TestSetCardConfigEx(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2CardConfigEx
	if ret := GetCardConfigEx(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if ret := SetCardConfigEx(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestGetDstConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2DstConfig
	if ret := GetDstConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", config)
	}
}

func TestSetDstConfig(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var config BS2DstConfig
	if ret := GetDstConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
	if ret := SetDstConfig(context, deviceId, &config); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestScanCard(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	onReadyToScan := func(deviceId uint32, sequence uint32) uintptr {
		if testing.Verbose() {
			t.Logf("recieved ready-to-scan(ID=%d, seq=%d).", deviceId, sequence)
		}
		return 0
	}

	var cardData BS2Card
	if ret := ScanCard(context, deviceId, &cardData, onReadyToScan); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", cardData)
	}
}

func TestWriteCard(t *testing.T) {
	//ERROR : INVALID PARAM (-200)
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var cardData BS2SmartCardData
	if ret := WriteCard(context, deviceId, &cardData); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestEraseCard(t *testing.T) {
	//BS_SDK_ERROR_CARD_CANNOT_WRITE_DATA
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	if ret := EraseCard(context, deviceId); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestScanFingerprint(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	onReadyToScan := func(deviceId uint32, sequence uint32) uintptr {
		if testing.Verbose() {
			t.Logf("recieved ready-to-scan(ID=%d, seq=%d).", deviceId, sequence)
		}
		return 0
	}

	var finger BS2Fingerprint
	if ret := ScanFingerprint(context, deviceId, &finger, 0 /*Standard*/, 40 /*Suprema*/, 0, onReadyToScan); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", finger)
	}
}

func TestScanFingerprintEx(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	onReadyToScan := func(deviceId uint32, sequence uint32) uintptr {
		if testing.Verbose() {
			t.Logf("recieved ready-to-scan(ID=%d, seq=%d).", deviceId, sequence)
		}
		return 0
	}

	var outQuality uint32
	var finger BS2Fingerprint
	if ret := ScanFingerprintEx(context, deviceId, &finger, 0 /*Standard*/, 40 /*Suprema*/, 0, &outQuality, onReadyToScan); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("quality:%d-%v", outQuality, finger)
	}
}

func TestVerifyFingerprint(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	if ret := RemoveAllUser(context, deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	onReadyToScan := func(deviceId uint32, sequence uint32) uintptr {
		if testing.Verbose() {
			t.Logf("recieved ready-to-scan(ID=%d, seq=%d).", deviceId, sequence)
		}
		return 0
	}

	var finger BS2Fingerprint
	if ret := ScanFingerprint(context, deviceId, &finger, 0 /*Standard*/, 40 /*Suprema*/, 0, onReadyToScan); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if ret := VerifyFingerprint(context, deviceId, &finger); BS_SDK_ERROR_NOT_SAME_FINGERPRINT != int16(ret) {
		t.Errorf("expected %d, but got %d", BS_SDK_ERROR_NOT_SAME_FINGERPRINT, ret)
	}
}

func TestGetLastFingerprintImage(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var finger BS2Fingerprint
	if ret := ScanFingerprint(context, deviceId, &finger, 0 /*Standard*/, 40 /*Suprema*/, 0, nil); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var width, height uint32
	var image []byte
	if ret := GetLastFingerprintImage(context, deviceId, &image, &width, &height); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("width:%d,height:%d,sizeofByte:%d-%v", width, height, len(image), image)
	}
}

func TestScanFace(t *testing.T) {
	t.SkipNow()
}
func TestGetAuthGroup(t *testing.T) {
	t.SkipNow()
}
func TestGetAllAuthGroup(t *testing.T) {
	t.SkipNow()
}
func TestSetAuthGroup(t *testing.T) {
	t.SkipNow()
}
func TestRemoveAuthGroup(t *testing.T) {
	t.SkipNow()
}
func TestRemoveAllAuthGroup(t *testing.T) {
	t.SkipNow()
}

func TestGetAccessGroup(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var accessGroups []BS2AccessGroup
	if ret := GetAccessGroup(context, deviceId, []uint32{1}, &accessGroups); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", accessGroups)
	}

}

func TestGetAllAccessGroup(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var accessGroups []BS2AccessGroup
	if ret := GetAllAccessGroup(context, deviceId, &accessGroups); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", accessGroups)
	}
}

func TestSetAccessGroup(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var accessGroups []BS2AccessGroup
	accessGroups = append(accessGroups, BS2AccessGroup{})

	if ret := SetAccessGroup(context, deviceId, &accessGroups); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", accessGroups)
	}
}

func TestRemoveAccessGroup(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	if ret := RemoveAccessGroup(context, deviceId, []uint32{1}); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestRemoveAllAccessGroup(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	if ret := RemoveAllAccessGroup(context, deviceId); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestGetAccessLevel(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var accessLevels []BS2AccessLevel
	if ret := GetAccessLevel(context, deviceId, []uint32{1}, &accessLevels); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", accessLevels)
	}
}

func TestSetAccessLevel(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var accessLevels []BS2AccessLevel
	accessLevels = append(accessLevels, BS2AccessLevel{})
	if ret := SetAccessLevel(context, deviceId, &accessLevels); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", accessLevels)
	}
}

func TestGetAllAccessLevel(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var accessLevels []BS2AccessLevel
	if ret := GetAllAccessLevel(context, deviceId, &accessLevels); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		t.Logf("%v", accessLevels)
	}
}

func TestRemoveAccessLevel(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	if ret := RemoveAccessLevel(context, deviceId, []uint32{1}); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestRemoveAllAccessLevel(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	if ret := RemoveAllAccessLevel(context, deviceId); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestGetAccessSchedule(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var accessSchedules []BS2Schedule
	if ret := GetAccessSchedule(context, deviceId, []uint32{1}, &accessSchedules); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		for _, a := range accessSchedules {
			t.Logf("%v", a)
		}
	}
}

func TestGetAllAccessSchedule(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var accessSchedules []BS2Schedule
	if ret := GetAllAccessSchedule(context, deviceId, &accessSchedules); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if testing.Verbose() {
		for _, a := range accessSchedules {
			t.Logf("%v", a)
		}
	}
}

func TestGetHolidayGroup(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var holidayGroupObj []BS2HolidayGroup
	if ret := GetHolidayGroup(context, deviceId, []uint32{1}, &holidayGroupObj); BS_SDK_SUCCESS != ret {
		t.Errorf("can not get holidayGroups from a device(ID=%d)", deviceId)
	}
}

func TestGetAllHolidayGroup(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var holidayGroupObj []BS2HolidayGroup

	if ret := GetAllHolidayGroup(context, deviceId, &holidayGroupObj); BS_SDK_SUCCESS != ret {
		t.Errorf("can not get holidayGroups from a device(ID=%d)", deviceId)
	}

	if testing.Verbose() {
		for _, h := range holidayGroupObj {
			t.Logf("%v\n", h)
		}
	}
}

func TestSetHolidayGroup(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var holidayGroupObj []BS2HolidayGroup
	holidayGroupObj = append(holidayGroupObj, BS2HolidayGroup{})

	if ret := SetHolidayGroup(context, deviceId, &holidayGroupObj); BS_SDK_SUCCESS != ret {
		t.Errorf("can not set holidayGroups from a device(ID=%d)", deviceId)
	}
}

func TestRemoveHolidayGroup(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	if ret := RemoveHolidayGroup(context, deviceId, []uint32{1}); BS_SDK_SUCCESS != ret {
		t.Errorf("can not set holidayGroups from a device(ID=%d)", deviceId)
	}
}

func TestRemoveAllHolidayGroup(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	if ret := RemoveAllHolidayGroup(context, deviceId); BS_SDK_SUCCESS != ret {
		t.Errorf("can not set holidayGroups from a device(ID=%d)", deviceId)
	}
}

func TestGetBlackList(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var blackListObj []BS2BlackList
	if ret := GetBlackList(context, deviceId, []BS2BlackList{{CardID: [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, IssueCount: 0}}, &blackListObj); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestGetAllBlackList(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var blackListObj []BS2BlackList
	if ret := GetAllBlackList(context, deviceId, &blackListObj); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestSetBlackList(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var blackListObj []BS2BlackList
	blackListObj = append(blackListObj, BS2BlackList{})
	if ret := SetBlackList(context, deviceId, &blackListObj); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestRemoveBlackList(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	if ret := RemoveBlackList(context, deviceId, []BS2BlackList{{CardID: [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, IssueCount: 0}}); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestRemoveAllBlackList(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	if ret := RemoveAllBlackList(context, deviceId); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestGetDoor(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var doors []BS2Door
	if ret := GetDoor(context, deviceId, []uint32{1}, &doors); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestGetAllDoor(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var doors []BS2Door
	if ret := GetAllDoor(context, deviceId, &doors); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestGetDoorStatus(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var doorStatuses []BS2DoorStatus
	if ret := GetDoorStatus(context, deviceId, []uint32{1}, &doorStatuses); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestGetAllDoorStatus(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var doorStatuses []BS2DoorStatus
	if ret := GetAllDoorStatus(context, deviceId, &doorStatuses); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestSetDoor(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	if ret := RemoveAllDoor(context, deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var doors []BS2Door
	doors = append(doors, BS2Door{DoorID: 1, EntryDeviceID: deviceId, Relay: BS2DoorRelay{DeviceID: deviceId, Port: 0}})

	if ret := SetDoor(context, deviceId, &doors); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if ret := RemoveAllDoor(context, deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

}

func TestSetDoorAlarm(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	if ret := SetDoorAlarm(context, deviceId, 1, []uint32{1}); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestRemoveDoor(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	if ret := RemoveAllDoor(context, deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var doors []BS2Door
	doors = append(doors, BS2Door{DoorID: 1, EntryDeviceID: deviceId, Relay: BS2DoorRelay{DeviceID: deviceId, Port: 0}})

	if ret := SetDoor(context, deviceId, &doors); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if ret := RemoveDoor(context, deviceId, []uint32{1}); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestRemoveAllDoor(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}

	if ret := RemoveAllDoor(context, deviceId); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestReleaseDoor(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	if ret := RemoveAllDoor(context, deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var doors []BS2Door
	doors = append(doors, BS2Door{DoorID: 1, EntryDeviceID: deviceId, Relay: BS2DoorRelay{DeviceID: deviceId, Port: 0}})

	if ret := SetDoor(context, deviceId, &doors); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if ret := ReleaseDoor(context, deviceId, 1, []uint32{1}); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if ret := RemoveAllDoor(context, deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

}

func TestLockDoor(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	if ret := RemoveAllDoor(context, deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var doors []BS2Door
	doors = append(doors, BS2Door{DoorID: 1, EntryDeviceID: deviceId, Relay: BS2DoorRelay{DeviceID: deviceId, Port: 0}})

	if ret := SetDoor(context, deviceId, &doors); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if ret := LockDoor(context, deviceId, 1, []uint32{1}); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

	if ret := RemoveAllDoor(context, deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}
}

func TestUnlockDoor(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	if ret := RemoveAllDoor(context, deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var doors []BS2Door
	doors = append(doors, BS2Door{DoorID: 1, EntryDeviceID: deviceId, Relay: BS2DoorRelay{DeviceID: deviceId, Port: 0}})

	if ret := SetDoor(context, deviceId, &doors); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if ret := UnlockDoor(context, deviceId, 1, []uint32{1}); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
	if ret := RemoveAllDoor(context, deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}
}

func TestGetSlaveDevice(t *testing.T) {
	//BS_SDK_ERROR_NOT_SUPPORTED
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}

	defer DisconnectDevice(context, deviceId)

	var slaveDeviceObj []BS2Rs485SlaveDevice
	if ret := GetSlaveDevice(context, deviceId, &slaveDeviceObj); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}

}

func TestSetSlaveDevice(t *testing.T) {
	//BS_SDK_ERROR_NOT_SUPPORTED
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var slaveDeviceObj []BS2Rs485SlaveDevice
	slaveDeviceObj = append(slaveDeviceObj, BS2Rs485SlaveDevice{})
	if ret := SetSlaveDevice(context, deviceId, &slaveDeviceObj); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestGetSlaveExDevice(t *testing.T) {
	//BS_SDK_ERROR_NOT_SUPPORTED
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var outchannelPort uint32
	var slaveDeviceObj []BS2Rs485SlaveDeviceEX
	if ret := GetSlaveExDevice(context, deviceId, 0, &slaveDeviceObj, &outchannelPort); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestSetSlaveExDevice(t *testing.T) {
	//BS_SDK_ERROR_NOT_SUPPORTED
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var slaveDeviceObj []BS2Rs485SlaveDeviceEX
	slaveDeviceObj = append(slaveDeviceObj, BS2Rs485SlaveDeviceEX{})
	if ret := SetSlaveExDevice(context, deviceId, 0, &slaveDeviceObj); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestSearchWiegandDevices(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	t.Logf("deviceId = %d", deviceId)
	var wiegandDeviceObj []uint32
	if ret := SearchWiegandDevices(context, deviceId, &wiegandDeviceObj); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestGetWiegandDevices(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var wiegandDeviceObj []uint32
	if ret := GetWiegandDevices(context, deviceId, &wiegandDeviceObj); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestAddWiegandDevices(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var wiegandDeviceObj []uint32
	if ret := SearchWiegandDevices(context, deviceId, &wiegandDeviceObj); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == len(wiegandDeviceObj) {
		t.Fatalf("had no wiegand devices from the device.")
	}

	if ret := AddWiegandDevices(context, deviceId, &wiegandDeviceObj); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestRemoveWiegandDevices(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var wiegandDeviceObj []uint32
	wiegandDeviceObj = append(wiegandDeviceObj, 1)
	if ret := RemoveWiegandDevices(context, deviceId, wiegandDeviceObj); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestSetServerMatchingHandler(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	if ret := SetServerMatchingHandler(context, nil, nil); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestVerifyUser(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var ub BS2UserBlob
	if ret := VerifyUser(context, deviceId, 1, 1, &ub); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestIdentifyUser(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var ub BS2UserBlob
	if ret := IdentifyUser(context, deviceId, 1, 1, &ub); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestVerifyUserEx(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var ub BS2UserBlobEx
	if ret := VerifyUserEx(context, deviceId, 1, 1, &ub); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}
func TestIdentifyUserEx(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var ub BS2UserBlobEx
	if ret := IdentifyUserEx(context, deviceId, 1, 1, &ub); BS_SDK_SUCCESS != ret {
		t.Errorf("Error occurred : %v", ret)
	}
}

func TestGetSupportedConfigMask(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var configMask BS2_CONFIG_MASK
	if ret := GetSupportedConfigMask(context, deviceId, &configMask); BS_SDK_SUCCESS != ret {
		t.Errorf("can not get supported config mask(%d).", ret)
	}

	if testing.Verbose() {
		t.Logf("supported config-mask : 0x%x", configMask)
	}

}

func TestGetSupportedUserMask(t *testing.T) {
	context := AllocateContext()
	defer ReleaseContext(context)

	if ret := Initialize(context); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	var deviceId uint32
	if ret := ConnectDeviceViaIP(context, parameter.IPAddr, parameter.Port, &deviceId); BS_SDK_SUCCESS != ret {
		t.Fatalf("Error occurred : %v", ret)
	}

	if 0 == deviceId {
		t.Fatalf("can not connect a device(%s:%d)", parameter.IPAddr, parameter.Port)
	}
	defer DisconnectDevice(context, deviceId)

	var configMask BS2_USER_MASK
	if ret := GetSupportedUserMask(context, deviceId, &configMask); BS_SDK_SUCCESS != ret {
		t.Errorf("can not get supported config mask(%d).", ret)
	}

	if testing.Verbose() {
		t.Logf("supported config-mask : 0x%x", configMask)
	}
}
