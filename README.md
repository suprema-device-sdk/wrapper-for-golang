# wrapper-for-golang

This is Suprema Device SDK Wrapper for Go-lang

## Getting Started

### MS Windows 
* install BioStar 2 Device SDK from http://kb.supremainc.com/bs2sdk/
* download wrapper-for-golang  
```
$ cd %GOPATH%

$ go get github.com/suprema-device-sdk/wrapper-for-golang.git

```
* make sure that this wrapper is imported 
```
import (
  devicesdk "github.com/suprema-device-sdk/wrapper-for-golang.git/windows-wrapper
)
```
* you can make physical access contorl systems or utilities using wrapper-for-golang
* enjoy it!

### Linux/Mac OS

TBD. 

## Sea also

* [BioStar 2 Device SDK](http://kb.supremainc.com/bs2sdk/)


## Examples 

### How to run device-sdk-cli example 

* copy SDK files to %GOPATH%/bin
* run Command-Prompt and then execute below

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
```

