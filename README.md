# wrapper-for-golang

This is Suprema Device SDK Wrapper for Go-lang

## Getting Started

### MS Windows 
* make a request for the BioStar 2 Device SDK from http://kb.supremainc.com/bs2sdk/
* download the wrapper-for-golang  
```
$ cd %GOPATH%

$ go get github.com/suprema-device-sdk/wrapper-for-golang.git

```
* make sure that this wrapper is imported 
```
import (
  devicesdk "github.com/suprema-device-sdk/wrapper-for-golang.git/windows-wrapper"
)
```
* you can make physical access contorl systems or utilities using wrapper-for-golang
* enjoy it!

### Linux/Mac OS

TBD. 

## Sea also

* [BioStar 2 Device SDK](http://kb.supremainc.com/bs2sdk/)


## Examples 

### How to run the device-sdk-cli example 

* copy SDK files to %GOPATH%/bin
* run the Command-Prompt and then execute below

```
$ cd %GOPATH%

$ go get github.com/suprema-device-sdk/wrapper-for-golang.git

$ cd ./src/github.com/suprema-device-sdk/wrapper-for-golang/example/device-sdk-cli

$ go install

$ cd %GOPATH%/bin
```

* As you can see you can execute commands 
  
```
$ device-sdk-cli.exe
CLI using Suprema Device-SDK [2.6.1.20]
> help

Usage :
        > Command parameter-1 parameter-2

Commands:
        exit - quit this CLI

        help - show commands and comments

        Init - initialize Deivce SDK

        Discover - discover devices via UDP and then print the devices

        Connect - connect a device(parameter#1)
                example : Connect 541530986

        ConnectViaIP - connect a device using the address(parameter#1) and port(parameter#2) of the device
                example : ConnectViaIP 192.168.1.1 51211

        Disconnect - disconnect a device(parameter#1)
                example : Disconnect 541530986

        StartMonitoring - start event-monitoring from a device(parameter#1)
                example : StartMonitoring 541530986

        StopMonitoring - stop event-monitoring from a device(parameter#1)
                example : StopMonitoring 541530986

        GetDeviceInfo - get device-info
                example : GetDeviceInfo 541530986
```

