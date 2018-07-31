package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"

	devicesdk "github.com/suprema-device-sdk/wrapper-for-golang/windows-wrapper"
)

func main() {

	commandReceived := make(chan string)
	connectionStatusChanged := make(chan ConnectionStatus, 1)
	deviceDiscovered := make(chan []uint32, 1)
	eventReceived := make(chan devicesdk.BS2Event, 1)

	deviceInterface := NewDeviceInterface(connectionStatusChanged, deviceDiscovered, eventReceived)
	deviceInterface.Init()

	go ScanUserCommand(deviceInterface.Version(), commandReceived)

EXIT:
	for {
		select {
		case cmd := <-commandReceived:
			switch {
			case "exit" == cmd:
				close(commandReceived)
				close(connectionStatusChanged)
				close(eventReceived)
				break EXIT
			case "help" == cmd:
				ShowHelp()
				continue
			}
			if result := ExecuteUserCommand(deviceInterface, cmd); !result {
				if len(cmd) > 0 {
					ShowMessageln(os.Stderr, "Unknown command or invalid parameters : ", cmd)
					continue
				}
			}
			ShowPrompt(false)
		case cs := <-connectionStatusChanged:
			ShowLogf("ConnectionStatusChanged: %v", cs)
		case er := <-eventReceived:
			ShowLogf("EventReceived : %v", er)
		case dd := <-deviceDiscovered:
			ShowDevicesDiscovered(dd)
		}
	}
}

func ScanUserCommand(version string, ch chan<- string) {
	scanner := bufio.NewScanner(os.Stdin)
	ShowMessagef(false, "CLI using Suprema Device-SDK %s", version)
	for scanner.Scan() {
		line := scanner.Text()
		ch <- line
	}
	if err := scanner.Err(); err != nil {
		ShowMessageln(os.Stderr, "reading standard input:", err)
	}
}

func ExecuteUserCommand(device *DeviceInterface, UserInputedCommand string) bool {
	arguments := strings.Split(UserInputedCommand, " ")

	v := reflect.ValueOf(device)
	m := v.MethodByName(arguments[0])
	if m.Kind() != reflect.Func {
		return false
	}

	argv, successfullyParsed := ParseUserCommand(&m, &arguments)

	if successfullyParsed {
		m.Call(argv)
	}

	return successfullyParsed && true
}

func ParseUserCommand(m *reflect.Value /*, t *reflect.Type*/, arguments *[]string) ([]reflect.Value, bool) {
	t := m.Type()
	argv := make([]reflect.Value, t.NumIn()) /*https://github.com/a8m/reflect-examples*/

	if len(*arguments)-1 < t.NumIn() {
		ShowMessageln(os.Stderr, "lack of arguments")
		return argv, false
	}

	for i := range argv {

		switch {
		case t.In(i).Kind() == reflect.Uint32:
			val, err := strconv.ParseUint((*arguments)[i+1], 10, 32)
			if nil != err {
				ShowMessageln(os.Stderr, "could not parse %d parameter to uint32 value : %s", i, (*arguments)[i+1])
			}
			argv[i] = reflect.ValueOf(uint32(val))
		case t.In(i).Kind() == reflect.Uint16:
			val, err := strconv.ParseUint((*arguments)[i+1], 10, 32)
			if nil != err {
				ShowMessageln(os.Stderr, "could not parse %d parameter to uint16 value : %s", i, (*arguments)[i+1])
			}
			argv[i] = reflect.ValueOf(uint16(val))
		case t.In(i).Kind() == reflect.String:
			argv[i] = reflect.ValueOf((*arguments)[i+1])
		case t.In(i).Kind() == reflect.Bool:
			val, err := strconv.ParseBool((*arguments)[i+1])
			if nil != err {
				ShowMessageln(os.Stderr, "could not parse %d parameter to boolean value : %s", i, (*arguments)[i+1])
			}
			argv[i] = reflect.ValueOf(val)
			//TODO : append other types
		default:
			ShowMessageln(os.Stderr, "could not parse %d parameter to %v value\n", i, t.In(i).Kind())
		}
	}
	return argv, true
}

func ShowDevicesDiscovered(dd []uint32) {

	buffer := fmt.Sprintf("%d deivces were founded", len(dd))
	if 0 != len(dd) {
		buffer += fmt.Sprint(", See below.")
		for i, deviceId := range dd {
			buffer += fmt.Sprintf("\n%03d Device ID - %d", i, deviceId)
		}
		buffer += fmt.Sprintf("\n\ntry to put down a command like below.\n Connect %d", dd[0])
	} else {
		buffer += fmt.Sprint(".")
	}

	ShowMessage(true, buffer)
}

func ShowMessage(newLine bool, message string) {
	if newLine {
		fmt.Print("\n")
	}
	fmt.Print(message)
	ShowPrompt(true)
}

func ShowLogf(format string, a ...interface{}) {
	fmt.Print("\n")
	fmt.Printf(format, a)
	ShowPrompt(true)
}

func ShowMessagef(newLine bool, format string, a ...interface{}) {
	if newLine {
		fmt.Print("\n")
	}
	fmt.Printf(format, a)
	ShowPrompt(true)
}

func ShowMessageln(w io.Writer, a ...interface{}) {
	fmt.Fprintln(w, a)
	ShowPrompt(false)
}

func ShowPrompt(newLine bool) {
	if newLine {
		fmt.Print("\n")
	}
	fmt.Printf("> ")
}

func ShowHelp() {
	fmt.Print(`
Usage : 
	> Command parameter-1 parameter-2

Commands:	
	exit - quit this CLI

	help - show commands and comments

	Init - initialize Deivce SDK 
	
	Discover - discover devices via UDP and then print the devices 
						
	Connect - connect the device(parameter#1)
		example : Connect 541530986

	ConnectViaIP - connect the address(parameter#1) of device via the port(parameter#2)
		example : ConnectViaIP 192.168.1.1 51211

	Disconnect - disconnect the device(parameter#1) 
		example : Disconnect 541530986

	StartMonitoring - start event-monitoring from the device(parameter#1)
		example : StartMonitoring 541530986

	StopMonitoring - stop event-monitoring from the device(parameter#1)
		example : StopMonitoring 541530986

	GetDeviceInfo - get device-info
		example : GetDeviceInfo 541530986
`)
	ShowPrompt(true)
}
