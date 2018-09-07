package devicesdk

import (
	"bytes"
	"log"
	"reflect"
	"syscall"
	"unsafe"

	//"encoding/binary"
	_ "runtime/cgo" /*https://stackoverflow.com/questions/44968397/windows-callback-function-in-golang*/
)

const (
	BS_WRAPPER_ERROR_NOT_SUPPORTED = 0x7FFF
	BS_WRAPPER_ERROR_SMALL_BUFFER  = 0x7FFE
)

type SdkContext = uintptr
type BS2_DEVICE_ID = uint32

var (
	Handle, _                         = syscall.LoadLibrary(`BS_SDK_V2.dll`)
	procBs2Version, _                 = syscall.GetProcAddress(Handle, "BS2_Version")
	procBs2AllocateContext, _         = syscall.GetProcAddress(Handle, "BS2_AllocateContext")
	procBs2ReleaseContext, _          = syscall.GetProcAddress(Handle, "BS2_ReleaseContext")
	procBs2Initialize, _              = syscall.GetProcAddress(Handle, "BS2_Initialize")
	procBs2ReleaseObject, _           = syscall.GetProcAddress(Handle, "BS2_ReleaseObject")
	procBs2MakePinCode, _             = syscall.GetProcAddress(Handle, "BS2_MakePinCode")
	procBs2SetMaxThreadCount, _       = syscall.GetProcAddress(Handle, "BS2_SetMaxThreadCount")
	procBs2ComputeCRC16CCITT, _       = syscall.GetProcAddress(Handle, "BS2_ComputeCRC16CCITT")
	procBs2GetCardModel, _            = syscall.GetProcAddress(Handle, "BS2_GetCardModel")
	procBs2GetDataEncryptKey, _       = syscall.GetProcAddress(Handle, "BS2_GetDataEncryptKey")
	procBs2SetDataEncryptKey, _       = syscall.GetProcAddress(Handle, "BS2_SetDataEncryptKey")
	procBs2RemoveDataEncryptKey, _    = syscall.GetProcAddress(Handle, "BS2_RemoveDataEncryptKey")
	procBs2SetDeviceEventListener, _  = syscall.GetProcAddress(Handle, "BS2_SetDeviceEventListener")
	procBs2SearchDevices, _           = syscall.GetProcAddress(Handle, "BS2_SearchDevices")
	procBs2SearchDevicesEx, _         = syscall.GetProcAddress(Handle, "BS2_SearchDevicesEx")
	procBs2GetDevices, _              = syscall.GetProcAddress(Handle, "BS2_GetDevices")
	procBs2ConnectDevice, _           = syscall.GetProcAddress(Handle, "BS2_ConnectDevice")
	procBs2ConnectDeviceViaIP, _      = syscall.GetProcAddress(Handle, "BS2_ConnectDeviceViaIP")
	procBs2DisconnectDevice, _        = syscall.GetProcAddress(Handle, "BS2_DisconnectDevice")
	procBs2SetKeepAliveTimeout, _     = syscall.GetProcAddress(Handle, "BS2_SetKeepAliveTimeout")
	procBs2SetNotificationListener, _ = syscall.GetProcAddress(Handle, "BS2_SetNotificationListener")
	procBs2SetServerPort, _           = syscall.GetProcAddress(Handle, "BS2_SetServerPort")

	procBs2GetDeviceTopology, _ = syscall.GetProcAddress(Handle, "BS2_GetDeviceTopology") //
	procBs2SetDeviceTopology, _ = syscall.GetProcAddress(Handle, "BS2_SetDeviceTopology") //

	procBs2SetSSLHandler, _ = syscall.GetProcAddress(Handle, "BS2_SetSSLHandler")
	procBs2DisableSSL, _    = syscall.GetProcAddress(Handle, "BS2_DisableSSL")

	procBs2GetDeviceInfo, _   = syscall.GetProcAddress(Handle, "BS2_GetDeviceInfo")
	procBs2GetDeviceTime, _   = syscall.GetProcAddress(Handle, "BS2_GetDeviceTime")
	procBs2SetDeviceTime, _   = syscall.GetProcAddress(Handle, "BS2_SetDeviceTime")
	procBs2ClearDatabase, _   = syscall.GetProcAddress(Handle, "BS2_ClearDatabase")
	procBs2FactoryReset, _    = syscall.GetProcAddress(Handle, "BS2_FactoryReset")
	procBs2RebootDevice, _    = syscall.GetProcAddress(Handle, "BS2_RebootDevice")
	procBs2LockDevice, _      = syscall.GetProcAddress(Handle, "BS2_LockDevice")
	procBs2UnlockDevice, _    = syscall.GetProcAddress(Handle, "BS2_UnlockDevice")
	procBs2UpgradeFirmware, _ = syscall.GetProcAddress(Handle, "BS2_UpgradeFirmware")
	procBs2UpdateResource, _  = syscall.GetProcAddress(Handle, "BS2_UpdateResource")

	procBs2GetLog, _                     = syscall.GetProcAddress(Handle, "BS2_GetLog")
	procBs2GetFilteredLog, _             = syscall.GetProcAddress(Handle, "BS2_GetFilteredLog")
	procBs2ClearLog, _                   = syscall.GetProcAddress(Handle, "BS2_ClearLog")
	procBs2StartMonitoringLog, _         = syscall.GetProcAddress(Handle, "BS2_StartMonitoringLog")
	procBs2StopMonitoringLog, _          = syscall.GetProcAddress(Handle, "BS2_StopMonitoringLog")
	procBs2GetLogBlob, _                 = syscall.GetProcAddress(Handle, "BS2_GetLogBlob")
	procBs2GetFilteredLogSinceEventId, _ = syscall.GetProcAddress(Handle, "BS2_GetFilteredLogSinceEventId")

	procBs2GetUserList, _         = syscall.GetProcAddress(Handle, "BS2_GetUserList")
	procBs2GetUserInfos, _        = syscall.GetProcAddress(Handle, "BS2_GetUserInfos")
	procBs2GetUserDatas, _        = syscall.GetProcAddress(Handle, "BS2_GetUserDatas")
	procBs2EnrolUser, _           = syscall.GetProcAddress(Handle, "BS2_EnrolUser")
	procBs2RemoveUser, _          = syscall.GetProcAddress(Handle, "BS2_RemoveUser")
	procBs2RemoveAllUser, _       = syscall.GetProcAddress(Handle, "BS2_RemoveAllUser")
	procBs2GetUserInfosEx, _      = syscall.GetProcAddress(Handle, "BS2_GetUserInfosEx")
	procBs2EnrolUserEx, _         = syscall.GetProcAddress(Handle, "BS2_EnrolUserEx")
	procBs2GetUserDatabaseInfo, _ = syscall.GetProcAddress(Handle, "BS2_GetUserDatabaseInfo")

	procBs2ResetConfig, _              = syscall.GetProcAddress(Handle, "BS2_ResetConfig")
	procBs2ResetConfigExceptNetInfo, _ = syscall.GetProcAddress(Handle, "BS2_ResetConfigExceptNetInfo")
	procBs2GetConfig, _                = syscall.GetProcAddress(Handle, "BS2_GetConfig")
	procBs2SetConfig, _                = syscall.GetProcAddress(Handle, "BS2_SetConfig")
	procBs2GetFactoryConfig, _         = syscall.GetProcAddress(Handle, "BS2_GetFactoryConfig")
	procBs2GetSystemConfig, _          = syscall.GetProcAddress(Handle, "BS2_GetSystemConfig")
	procBs2SetSystemConfig, _          = syscall.GetProcAddress(Handle, "BS2_SetSystemConfig")
	procBs2GetAuthConfig, _            = syscall.GetProcAddress(Handle, "BS2_GetAuthConfig")
	procBs2SetAuthConfig, _            = syscall.GetProcAddress(Handle, "BS2_SetAuthConfig")
	procBs2GetStatusConfig, _          = syscall.GetProcAddress(Handle, "BS2_GetStatusConfig")
	procBs2SetStatusConfig, _          = syscall.GetProcAddress(Handle, "BS2_SetStatusConfig")
	procBs2GetDisplayConfig, _         = syscall.GetProcAddress(Handle, "BS2_GetDisplayConfig")
	procBs2SetDisplayConfig, _         = syscall.GetProcAddress(Handle, "BS2_SetDisplayConfig")
	procBs2GetIPConfig, _              = syscall.GetProcAddress(Handle, "BS2_GetIPConfig")
	procBs2GetIPConfigViaUDP, _        = syscall.GetProcAddress(Handle, "BS2_GetIPConfigViaUDP")
	procBs2SetIPConfig, _              = syscall.GetProcAddress(Handle, "BS2_SetIPConfig")
	procBs2SetIPConfigViaUDP, _        = syscall.GetProcAddress(Handle, "BS2_SetIPConfigViaUDP")
	procBs2GetIPConfigExt, _           = syscall.GetProcAddress(Handle, "BS2_GetIPConfigExt")
	procBs2SetIPConfigExt, _           = syscall.GetProcAddress(Handle, "BS2_SetIPConfigExt")
	procBs2GetTNAConfig, _             = syscall.GetProcAddress(Handle, "BS2_GetTNAConfig")
	procBs2SetTNAConfig, _             = syscall.GetProcAddress(Handle, "BS2_SetTNAConfig")
	procBs2GetCardConfig, _            = syscall.GetProcAddress(Handle, "BS2_GetCardConfig")
	procBs2SetCardConfig, _            = syscall.GetProcAddress(Handle, "BS2_SetCardConfig")
	procBs2GetFingerprintConfig, _     = syscall.GetProcAddress(Handle, "BS2_GetFingerprintConfig")
	procBs2SetFingerprintConfig, _     = syscall.GetProcAddress(Handle, "BS2_SetFingerprintConfig")
	procBs2GetRS485Config, _           = syscall.GetProcAddress(Handle, "BS2_GetRS485Config")
	procBs2SetRS485Config, _           = syscall.GetProcAddress(Handle, "BS2_SetRS485Config")
	procBs2GetWiegandConfig, _         = syscall.GetProcAddress(Handle, "BS2_GetWiegandConfig")
	procBs2SetWiegandConfig, _         = syscall.GetProcAddress(Handle, "BS2_SetWiegandConfig")
	procBs2GetWiegandDeviceConfig, _   = syscall.GetProcAddress(Handle, "BS2_GetWiegandDeviceConfig")
	procBs2SetWiegandDeviceConfig, _   = syscall.GetProcAddress(Handle, "BS2_SetWiegandDeviceConfig")
	procBs2GetInputConfig, _           = syscall.GetProcAddress(Handle, "BS2_GetInputConfig")
	procBs2SetInputConfig, _           = syscall.GetProcAddress(Handle, "BS2_SetInputConfig")
	procBs2GetWlanConfig, _            = syscall.GetProcAddress(Handle, "BS2_GetWlanConfig")
	procBs2SetWlanConfig, _            = syscall.GetProcAddress(Handle, "BS2_SetWlanConfig")
	procBs2GetTriggerActionConfig, _   = syscall.GetProcAddress(Handle, "BS2_GetTriggerActionConfig")
	procBs2SetTriggerActionConfig, _   = syscall.GetProcAddress(Handle, "BS2_SetTriggerActionConfig")
	procBs2GetEventConfig, _           = syscall.GetProcAddress(Handle, "BS2_GetEventConfig")
	procBs2SetEventConfig, _           = syscall.GetProcAddress(Handle, "BS2_SetEventConfig")
	procBs2GetWiegandMultiConfig, _    = syscall.GetProcAddress(Handle, "BS2_GetWiegandMultiConfig")
	procBs2SetWiegandMultiConfig, _    = syscall.GetProcAddress(Handle, "BS2_SetWiegandMultiConfig")
	procBs2GetCard1xConfig, _          = syscall.GetProcAddress(Handle, "BS2_GetCard1xConfig")
	procBs2SetCard1xConfig, _          = syscall.GetProcAddress(Handle, "BS2_SetCard1xConfig")
	procBs2GetSystemExtConfig, _       = syscall.GetProcAddress(Handle, "BS2_GetSystemExtConfig")
	procBs2SetSystemExtConfig, _       = syscall.GetProcAddress(Handle, "BS2_SetSystemExtConfig")
	procBs2GetVoipConfig, _            = syscall.GetProcAddress(Handle, "BS2_GetVoipConfig")
	procBs2SetVoipConfig, _            = syscall.GetProcAddress(Handle, "BS2_SetVoipConfig")
	procBs2GetFaceConfig, _            = syscall.GetProcAddress(Handle, "BS2_GetFaceConfig")
	procBs2SetFaceConfig, _            = syscall.GetProcAddress(Handle, "BS2_SetFaceConfig")
	procBs2GetRS485ConfigEx, _         = syscall.GetProcAddress(Handle, "BS2_GetRS485ConfigEx")
	procBs2SetRS485ConfigEx, _         = syscall.GetProcAddress(Handle, "BS2_SetRS485ConfigEx")
	procBs2GetCardConfigEx, _          = syscall.GetProcAddress(Handle, "BS2_GetCardConfigEx")
	procBs2SetCardConfigEx, _          = syscall.GetProcAddress(Handle, "BS2_SetCardConfigEx")
	procBs2GetDstConfig, _             = syscall.GetProcAddress(Handle, "BS2_GetDstConfig")
	procBs2SetDstConfig, _             = syscall.GetProcAddress(Handle, "BS2_SetDstConfig")

	procBs2ScanCard, _  = syscall.GetProcAddress(Handle, "BS2_ScanCard")
	procBs2WriteCard, _ = syscall.GetProcAddress(Handle, "BS2_WriteCard")
	procBs2EraseCard, _ = syscall.GetProcAddress(Handle, "BS2_EraseCard")

	procBs2ScanFingerprint, _         = syscall.GetProcAddress(Handle, "BS2_ScanFingerprint")
	procBs2ScanFingerprintEx, _       = syscall.GetProcAddress(Handle, "BS2_ScanFingerprintEx")
	procBs2VerifyFingerprint, _       = syscall.GetProcAddress(Handle, "BS2_VerifyFingerprint")
	procBs2GetLastFingerprintImage, _ = syscall.GetProcAddress(Handle, "BS2_GetLastFingerprintImage")

	procBs2ScanFace, _           = syscall.GetProcAddress(Handle, "BS2_ScanFace")
	procBs2GetAuthGroup, _       = syscall.GetProcAddress(Handle, "BS2_GetAuthGroup")
	procBs2GetAllAuthGroup, _    = syscall.GetProcAddress(Handle, "BS2_GetAllAuthGroup")
	procBs2SetAuthGroup, _       = syscall.GetProcAddress(Handle, "BS2_SetAuthGroup")
	procBs2RemoveAuthGroup, _    = syscall.GetProcAddress(Handle, "BS2_RemoveAuthGroup")
	procBs2RemoveAllAuthGroup, _ = syscall.GetProcAddress(Handle, "BS2_RemoveAllAuthGroup")

	procBs2GetAccessGroup, _          = syscall.GetProcAddress(Handle, "BS2_GetAccessGroup")
	procBs2GetAllAccessGroup, _       = syscall.GetProcAddress(Handle, "BS2_GetAllAccessGroup")
	procBs2SetAccessGroup, _          = syscall.GetProcAddress(Handle, "BS2_SetAccessGroup")
	procBs2RemoveAccessGroup, _       = syscall.GetProcAddress(Handle, "BS2_RemoveAccessGroup")
	procBs2RemoveAllAccessGroup, _    = syscall.GetProcAddress(Handle, "BS2_RemoveAllAccessGroup")
	procBs2GetAccessLevel, _          = syscall.GetProcAddress(Handle, "BS2_GetAccessLevel")
	procBs2GetAllAccessLevel, _       = syscall.GetProcAddress(Handle, "BS2_GetAllAccessLevel")
	procBs2SetAccessLevel, _          = syscall.GetProcAddress(Handle, "BS2_SetAccessLevel")
	procBs2RemoveAccessLevel, _       = syscall.GetProcAddress(Handle, "BS2_RemoveAccessLevel")
	procBs2RemoveAllAccessLevel, _    = syscall.GetProcAddress(Handle, "BS2_RemoveAllAccessLevel")
	procBs2GetAccessSchedule, _       = syscall.GetProcAddress(Handle, "BS2_GetAccessSchedule")
	procBs2GetAllAccessSchedule, _    = syscall.GetProcAddress(Handle, "BS2_GetAllAccessSchedule")
	procBs2SetAccessSchedule, _       = syscall.GetProcAddress(Handle, "BS2_SetAccessSchedule")
	procBs2RemoveAccessSchedule, _    = syscall.GetProcAddress(Handle, "BS2_RemoveAccessSchedule")
	procBs2RemoveAllAccessSchedule, _ = syscall.GetProcAddress(Handle, "BS2_RemoveAllAccessSchedule")
	procBs2GetHolidayGroup, _         = syscall.GetProcAddress(Handle, "BS2_GetHolidayGroup")
	procBs2GetAllHolidayGroup, _      = syscall.GetProcAddress(Handle, "BS2_GetAllHolidayGroup")
	procBs2SetHolidayGroup, _         = syscall.GetProcAddress(Handle, "BS2_SetHolidayGroup")
	procBs2RemoveHolidayGroup, _      = syscall.GetProcAddress(Handle, "BS2_RemoveHolidayGroup")
	procBs2RemoveAllHolidayGroup, _   = syscall.GetProcAddress(Handle, "BS2_RemoveAllHolidayGroup")

	procBs2GetBlackList, _       = syscall.GetProcAddress(Handle, "BS2_GetBlackList")
	procBs2GetAllBlackList, _    = syscall.GetProcAddress(Handle, "BS2_GetAllBlackList")
	procBs2SetBlackList, _       = syscall.GetProcAddress(Handle, "BS2_SetBlackList")
	procBs2RemoveBlackList, _    = syscall.GetProcAddress(Handle, "BS2_RemoveBlackList")
	procBs2RemoveAllBlackList, _ = syscall.GetProcAddress(Handle, "BS2_RemoveAllBlackList")

	procBs2GetDoor, _          = syscall.GetProcAddress(Handle, "BS2_GetDoor")
	procBs2GetAllDoor, _       = syscall.GetProcAddress(Handle, "BS2_GetAllDoor")
	procBs2GetDoorStatus, _    = syscall.GetProcAddress(Handle, "BS2_GetDoorStatus")
	procBs2GetAllDoorStatus, _ = syscall.GetProcAddress(Handle, "BS2_GetAllDoorStatus")
	procBs2SetDoor, _          = syscall.GetProcAddress(Handle, "BS2_SetDoor")
	procBs2SetDoorAlarm, _     = syscall.GetProcAddress(Handle, "BS2_SetDoorAlarm")
	procBs2RemoveDoor, _       = syscall.GetProcAddress(Handle, "BS2_RemoveDoor")
	procBs2RemoveAllDoor, _    = syscall.GetProcAddress(Handle, "BS2_RemoveAllDoor")
	procBs2ReleaseDoor, _      = syscall.GetProcAddress(Handle, "BS2_ReleaseDoor")
	procBs2LockDoor, _         = syscall.GetProcAddress(Handle, "BS2_LockDoor")
	procBs2UnlockDoor, _       = syscall.GetProcAddress(Handle, "BS2_UnlockDoor")

	procBs2GetLift, _             = syscall.GetProcAddress(Handle, "BS2_GetLift")
	procBs2GetAllLift, _          = syscall.GetProcAddress(Handle, "BS2_GetAllLift")
	procBs2GetLiftStatus, _       = syscall.GetProcAddress(Handle, "BS2_GetLiftStatus")
	procBs2GetAllLiftStatus, _    = syscall.GetProcAddress(Handle, "BS2_GetAllLiftStatus")
	procBs2SetLift, _             = syscall.GetProcAddress(Handle, "BS2_SetLift")
	procBs2SetLiftAlarm, _        = syscall.GetProcAddress(Handle, "BS2_SetLiftAlarm")
	procBs2RemoveLift, _          = syscall.GetProcAddress(Handle, "BS2_RemoveLift")
	procBs2RemoveAllLift, _       = syscall.GetProcAddress(Handle, "BS2_RemoveAllLift")
	procBs2ReleaseFloor, _        = syscall.GetProcAddress(Handle, "BS2_ReleaseFloor")
	procBs2ActivateFloor, _       = syscall.GetProcAddress(Handle, "BS2_ActivateFloor")
	procBs2DeActivateFloor, _     = syscall.GetProcAddress(Handle, "BS2_DeActivateFloor")
	procBs2GetFloorLevel, _       = syscall.GetProcAddress(Handle, "BS2_GetFloorLevel")
	procBs2GetAllFloorLevel, _    = syscall.GetProcAddress(Handle, "BS2_GetAllFloorLevel")
	procBs2SetFloorLevel, _       = syscall.GetProcAddress(Handle, "BS2_SetFloorLevel")
	procBs2RemoveFloorLevel, _    = syscall.GetProcAddress(Handle, "BS2_RemoveFloorLevel")
	procBs2RemoveAllFloorLevel, _ = syscall.GetProcAddress(Handle, "BS2_RemoveAllFloorLevel")

	procBs2GetSlaveDevice, _   = syscall.GetProcAddress(Handle, "BS2_GetSlaveDevice")
	procBs2SetSlaveDevice, _   = syscall.GetProcAddress(Handle, "BS2_SetSlaveDevice")
	procBs2GetSlaveExDevice, _ = syscall.GetProcAddress(Handle, "BS2_GetSlaveExDevice")
	procBs2SetSlaveExDevice, _ = syscall.GetProcAddress(Handle, "BS2_SetSlaveExDevice")

	procBs2SearchWiegandDevices, _ = syscall.GetProcAddress(Handle, "BS2_SearchWiegandDevices")
	procBs2GetWiegandDevices, _    = syscall.GetProcAddress(Handle, "BS2_GetWiegandDevices")
	procBs2AddWiegandDevices, _    = syscall.GetProcAddress(Handle, "BS2_AddWiegandDevices")
	procBs2RemoveWiegandDevices, _ = syscall.GetProcAddress(Handle, "BS2_RemoveWiegandDevices")

	procBs2SetServerMatchingHandler, _ = syscall.GetProcAddress(Handle, "BS2_SetServerMatchingHandler")
	procBs2VerifyUser, _               = syscall.GetProcAddress(Handle, "BS2_VerifyUser")
	procBs2IdentifyUser, _             = syscall.GetProcAddress(Handle, "BS2_IdentifyUser")
	procBs2VerifyUserEx, _             = syscall.GetProcAddress(Handle, "BS2_VerifyUserEx")
	procBs2IdentifyUserEx, _           = syscall.GetProcAddress(Handle, "BS2_IdentifyUserEx")

	procBs2GetAntiPassbackZone, _               = syscall.GetProcAddress(Handle, "BS2_GetAntiPassbackZone")
	procBs2GetAllAntiPassbackZone, _            = syscall.GetProcAddress(Handle, "BS2_GetAllAntiPassbackZone")
	procBs2GetAntiPassbackZoneStatus, _         = syscall.GetProcAddress(Handle, "BS2_GetAntiPassbackZoneStatus")
	procBs2GetAllAntiPassbackZoneStatus, _      = syscall.GetProcAddress(Handle, "BS2_GetAllAntiPassbackZoneStatus")
	procBs2SetAntiPassbackZone, _               = syscall.GetProcAddress(Handle, "BS2_SetAntiPassbackZone")
	procBs2SetAntiPassbackZoneAlarm, _          = syscall.GetProcAddress(Handle, "BS2_SetAntiPassbackZoneAlarm")
	procBs2RemoveAntiPassbackZone, _            = syscall.GetProcAddress(Handle, "BS2_RemoveAntiPassbackZone")
	procBs2RemoveAllAntiPassbackZone, _         = syscall.GetProcAddress(Handle, "BS2_RemoveAllAntiPassbackZone")
	procBs2ClearAntiPassbackZoneStatus, _       = syscall.GetProcAddress(Handle, "BS2_ClearAntiPassbackZoneStatus")
	procBs2ClearAllAntiPassbackZoneStatus, _    = syscall.GetProcAddress(Handle, "BS2_ClearAllAntiPassbackZoneStatus")
	procBs2SetCheckGlobalAPBViolationHandler, _ = syscall.GetProcAddress(Handle, "BS2_SetCheckGlobalAPBViolationHandler")
	procBs2CheckGlobalAPBViolation, _           = syscall.GetProcAddress(Handle, "BS2_CheckGlobalAPBViolation")

	procBs2GetTimedAntiPassbackZone, _            = syscall.GetProcAddress(Handle, "BS2_GetTimedAntiPassbackZone")
	procBs2GetAllTimedAntiPassbackZone, _         = syscall.GetProcAddress(Handle, "BS2_GetAllTimedAntiPassbackZone")
	procBs2GetTimedAntiPassbackZoneStatus, _      = syscall.GetProcAddress(Handle, "BS2_GetTimedAntiPassbackZoneStatus")
	procBs2GetAllTimedAntiPassbackZoneStatus, _   = syscall.GetProcAddress(Handle, "BS2_GetAllTimedAntiPassbackZoneStatus")
	procBs2SetTimedAntiPassbackZone, _            = syscall.GetProcAddress(Handle, "BS2_SetTimedAntiPassbackZone")
	procBs2SetTimedAntiPassbackZoneAlarm, _       = syscall.GetProcAddress(Handle, "BS2_SetTimedAntiPassbackZoneAlarm")
	procBs2RemoveTimedAntiPassbackZone, _         = syscall.GetProcAddress(Handle, "BS2_RemoveTimedAntiPassbackZone")
	procBs2RemoveAllTimedAntiPassbackZone, _      = syscall.GetProcAddress(Handle, "BS2_RemoveAllTimedAntiPassbackZone")
	procBs2ClearTimedAntiPassbackZoneStatus, _    = syscall.GetProcAddress(Handle, "BS2_ClearTimedAntiPassbackZoneStatus")
	procBs2ClearAllTimedAntiPassbackZoneStatus, _ = syscall.GetProcAddress(Handle, "BS2_ClearAllTimedAntiPassbackZoneStatus")

	procBs2GetFireAlarmZone, _          = syscall.GetProcAddress(Handle, "BS2_GetFireAlarmZone")
	procBs2GetAllFireAlarmZone, _       = syscall.GetProcAddress(Handle, "BS2_GetAllFireAlarmZone")
	procBs2GetFireAlarmZoneStatus, _    = syscall.GetProcAddress(Handle, "BS2_GetFireAlarmZoneStatus")
	procBs2GetAllFireAlarmZoneStatus, _ = syscall.GetProcAddress(Handle, "BS2_GetAllFireAlarmZoneStatus")
	procBs2SetFireAlarmZone, _          = syscall.GetProcAddress(Handle, "BS2_SetFireAlarmZone")
	procBs2SetFireAlarmZoneAlarm, _     = syscall.GetProcAddress(Handle, "BS2_SetFireAlarmZoneAlarm")
	procBs2RemoveFireAlarmZone, _       = syscall.GetProcAddress(Handle, "BS2_RemoveFireAlarmZone")
	procBs2RemoveAllFireAlarmZone, _    = syscall.GetProcAddress(Handle, "BS2_RemoveAllFireAlarmZone")

	procBs2GetScheduledLockUnlockZone, _          = syscall.GetProcAddress(Handle, "BS2_GetScheduledLockUnlockZone")
	procBs2GetAllScheduledLockUnlockZone, _       = syscall.GetProcAddress(Handle, "BS2_GetAllScheduledLockUnlockZone")
	procBs2GetScheduledLockUnlockZoneStatus, _    = syscall.GetProcAddress(Handle, "BS2_GetScheduledLockUnlockZoneStatus")
	procBs2GetAllScheduledLockUnlockZoneStatus, _ = syscall.GetProcAddress(Handle, "BS2_GetAllScheduledLockUnlockZoneStatus")
	procBs2SetScheduledLockUnlockZone, _          = syscall.GetProcAddress(Handle, "BS2_SetScheduledLockUnlockZone")
	procBs2SetScheduledLockUnlockZoneAlarm, _     = syscall.GetProcAddress(Handle, "BS2_SetScheduledLockUnlockZoneAlarm")
	procBs2RemoveScheduledLockUnlockZone, _       = syscall.GetProcAddress(Handle, "BS2_RemoveScheduledLockUnlockZone")
	procBs2RemoveAllScheduledLockUnlockZone, _    = syscall.GetProcAddress(Handle, "BS2_RemoveAllScheduledLockUnlockZone")

	procBs2GetIntrusionAlarmZone, _          = syscall.GetProcAddress(Handle, "BS2_GetIntrusionAlarmZone")
	procBs2GetIntrusionAlarmZoneStatus, _    = syscall.GetProcAddress(Handle, "BS2_GetIntrusionAlarmZoneStatus")
	procBs2GetAllIntrusionAlarmZoneStatus, _ = syscall.GetProcAddress(Handle, "BS2_GetAllIntrusionAlarmZoneStatus")
	procBs2SetIntrusionAlarmZone, _          = syscall.GetProcAddress(Handle, "BS2_SetIntrusionAlarmZone")
	procBs2SetIntrusionAlarmZoneAlarm, _     = syscall.GetProcAddress(Handle, "BS2_SetIntrusionAlarmZoneAlarm")
	procBs2RemoveIntrusionAlarmZone, _       = syscall.GetProcAddress(Handle, "BS2_RemoveIntrusionAlarmZone")
	procBs2RemoveAllIntrusionAlarmZone, _    = syscall.GetProcAddress(Handle, "BS2_RemoveAllIntrusionAlarmZone")

	procBs2GetInterlockZone, _          = syscall.GetProcAddress(Handle, "BS2_GetInterlockZone")
	procBs2GetInterlockZoneStatus, _    = syscall.GetProcAddress(Handle, "BS2_GetInterlockZoneStatus")
	procBs2GetAllInterlockZoneStatus, _ = syscall.GetProcAddress(Handle, "BS2_GetAllInterlockZoneStatus")
	procBs2SetInterlockZone, _          = syscall.GetProcAddress(Handle, "BS2_SetInterlockZone")
	procBs2SetInterlockZoneAlarm, _     = syscall.GetProcAddress(Handle, "BS2_SetInterlockZoneAlarm")
	procBs2RemoveInterlockZone, _       = syscall.GetProcAddress(Handle, "BS2_RemoveInterlockZone")
	procBs2RemoveAllInterlockZone, _    = syscall.GetProcAddress(Handle, "BS2_RemoveAllInterlockZone")

	procBs2GetDeviceZone, _                      = syscall.GetProcAddress(Handle, "BS2_GetDeviceZone")
	procBs2GetAllDeviceZone, _                   = syscall.GetProcAddress(Handle, "BS2_GetAllDeviceZone")
	procBs2SetDeviceZone, _                      = syscall.GetProcAddress(Handle, "BS2_SetDeviceZone")
	procBs2RemoveDeviceZone, _                   = syscall.GetProcAddress(Handle, "BS2_RemoveDeviceZone")
	procBs2RemoveAllDeviceZone, _                = syscall.GetProcAddress(Handle, "BS2_RemoveAllDeviceZone")
	procBs2SetDeviceZoneAlarm, _                 = syscall.GetProcAddress(Handle, "BS2_SetDeviceZoneAlarm")
	procBs2ClearDeviceZoneAccessRecord, _        = syscall.GetProcAddress(Handle, "BS2_ClearDeviceZoneAccessRecord")
	procBs2ClearAllDeviceZoneAccessRecord, _     = syscall.GetProcAddress(Handle, "BS2_ClearAllDeviceZoneAccessRecord")
	procBs2GetAccessGroupEntranceLimit, _        = syscall.GetProcAddress(Handle, "BS2_GetAccessGroupEntranceLimit")
	procBs2GetAllAccessGroupEntranceLimit, _     = syscall.GetProcAddress(Handle, "BS2_GetAllAccessGroupEntranceLimit")
	procBs2SetAccessGroupEntranceLimit, _        = syscall.GetProcAddress(Handle, "BS2_SetAccessGroupEntranceLimit")
	procBs2RemoveAccessGroupEntranceLimit, _     = syscall.GetProcAddress(Handle, "BS2_RemoveAccessGroupEntranceLimit")
	procBs2RemoveAllAccessGroupEntranceLimit, _  = syscall.GetProcAddress(Handle, "BS2_RemoveAllAccessGroupEntranceLimit")
	procBs2GetAllDeviceZoneAGEntranceLimit, _    = syscall.GetProcAddress(Handle, "BS2_GetAllDeviceZoneAGEntranceLimit")
	procBs2SetDeviceZoneAGEntranceLimit, _       = syscall.GetProcAddress(Handle, "BS2_SetDeviceZoneAGEntranceLimit")
	procBs2RemoveDeviceZoneAGEntranceLimit, _    = syscall.GetProcAddress(Handle, "BS2_RemoveDeviceZoneAGEntranceLimit")
	procBs2RemoveAllDeviceZoneAGEntranceLimit, _ = syscall.GetProcAddress(Handle, "BS2_RemoveAllDeviceZoneAGEntranceLimit")

	procBs2GetUserDatabaseInfoFromDir, _ = syscall.GetProcAddress(Handle, "BS2_GetUserDatabaseInfoFromDir")
	procBs2GetUserListFromDir, _         = syscall.GetProcAddress(Handle, "BS2_GetUserListFromDir")
	procBs2GetUserInfosFromDir, _        = syscall.GetProcAddress(Handle, "BS2_GetUserInfosFromDir")
	procBs2GetUserDatasFromDir, _        = syscall.GetProcAddress(Handle, "BS2_GetUserDatasFromDir")
	procBs2GetUserInfosExFromDir, _      = syscall.GetProcAddress(Handle, "BS2_GetUserInfosExFromDir")
	procBs2GetUserDatasExFromDir, _      = syscall.GetProcAddress(Handle, "BS2_GetUserDatasExFromDir")
	procBs2GetLogFromDir, _              = syscall.GetProcAddress(Handle, "BS2_GetLogFromDir")
	procBs2GetFilteredLogFromDir, _      = syscall.GetProcAddress(Handle, "BS2_GetFilteredLogFromDir")
	procBs2GetLogBlobFromDir, _          = syscall.GetProcAddress(Handle, "BS2_GetLogBlobFromDir")

	procBs2GetSupportedConfigMask, _ = syscall.GetProcAddress(Handle, "BS2_GetSupportedConfigMask")
	procBs2GetSupportedUserMask, _   = syscall.GetProcAddress(Handle, "BS2_GetSupportedUserMask")
)

func Version() string {
	ret, _, err := syscall.Syscall(procBs2Version, uintptr(0), uintptr(0), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}

	return uintptrToString(ret)
}

func AllocateContext() SdkContext {
	ret, _, err := syscall.Syscall(procBs2AllocateContext, uintptr(0), uintptr(0), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	// log.Printf("AllocateContext called(result=%x).\n", ret)
	return ret
}

func ReleaseContext(context SdkContext) {
	_, _, err := syscall.Syscall(procBs2ReleaseContext, uintptr(1), context, uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	// log.Printf("ReleaseContext called(parameter=%x).\n", context)
}

func Initialize(context SdkContext) int16 {
	ret, _, err := syscall.Syscall(procBs2Initialize, uintptr(1), context, uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func ReleaseObject(object uintptr) {
	_, _, err := syscall.Syscall(procBs2ReleaseObject, uintptr(1), object, uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	//log.Printf("ReleaseObject called(parameter=%x).\n", object)
}

func MakePinCode(context SdkContext, plaintext string, ciphertext *bytes.Buffer) int16 {
	cipherTextBuffer := make([]byte, BS2_PIN_HASH_SIZE)
	ret, _, err := syscall.Syscall(procBs2MakePinCode, uintptr(3), context, uintptr(unsafe.Pointer(syscall.StringBytePtr(plaintext))), uintptr(unsafe.Pointer(&cipherTextBuffer[0])))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	ciphertext.Write(cipherTextBuffer[:])
	return int16(ret)
}

func SetMaxThreadCount(context SdkContext, maxThreadCount uint32) int16 {
	ret, _, err := syscall.Syscall(procBs2SetMaxThreadCount, uintptr(2), context, uintptr(maxThreadCount), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func ComputeCRC16CCITT(data []byte, crc *uint16) int16 {
	ret, _, err := syscall.Syscall(procBs2ComputeCRC16CCITT, uintptr(3), uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)), uintptr(unsafe.Pointer(crc)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

const (
	BS2_CARD_MODEL_OMPW = iota
	BS2_CARD_MODEL_OIPW
	BS2_CARD_MODEL_OEPW
	BS2_CARD_MODEL_OHPW
	BS2_CARD_MODEL_ODPW
	BS2_CARD_MODEL_OAPW
)

type BS2_CARD_MODEL = uint16

func GetCardModel(modelName string, cardModel *BS2_CARD_MODEL) int16 {
	ret, _, err := syscall.Syscall(procBs2GetCardModel, uintptr(2), uintptr(uintptr(unsafe.Pointer(syscall.StringBytePtr(modelName)))), uintptr(unsafe.Pointer(cardModel)), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

const (
	BS2_ENC_KEY_SIZDE = 32
)

type BS2EncryptKey struct {
	Key      [BS2_ENC_KEY_SIZDE]byte
	Reserved [32]byte
}

func GetDataEncryptKey(context SdkContext, deviceId BS2_DEVICE_ID, keyInfo *BS2EncryptKey) int16 {
	if 0 == procBs2GetDataEncryptKey {
		return BS_WRAPPER_ERROR_NOT_SUPPORTED
	}
	ret, _, err := syscall.Syscall(procBs2GetDataEncryptKey, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(keyInfo)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func SetDataEncryptKey(context SdkContext, deviceId BS2_DEVICE_ID, keyInfo *BS2EncryptKey) int16 {
	if 0 == procBs2GetDataEncryptKey {
		return BS_WRAPPER_ERROR_NOT_SUPPORTED
	}
	ret, _, err := syscall.Syscall(procBs2SetDataEncryptKey, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(keyInfo)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func RemoveDataEncryptKey(context SdkContext, deviceId BS2_DEVICE_ID) int16 {
	if 0 == procBs2GetDataEncryptKey {
		return BS_WRAPPER_ERROR_NOT_SUPPORTED
	}
	ret, _, err := syscall.Syscall(procBs2RemoveDataEncryptKey, uintptr(2), context, uintptr(deviceId), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

// type OnDeviceFound func(deviceId BS2_DEVICE_ID)
// type OnDeviceAccepted func(deviceId BS2_DEVICE_ID)
// type OnDeviceConnected func(deviceId BS2_DEVICE_ID)
// type OnDeviceDisconnected func(deviceId BS2_DEVICE_ID)

func SetDeviceEventListener(context SdkContext, fnOnDeviceFound interface{}, fnOnDeviceAccepted interface{}, fnOnDeviceConnected interface{}, fnOnDeviceDisconnected interface{}) int16 {

	var fnOnDeviceFoundPtr uintptr
	var fnOnDeviceAcceptedPtr uintptr
	var fnOnDeviceConnectedPtr uintptr
	var fnOnDeviceDisconnectedPtr uintptr

	if nil != fnOnDeviceFound {
		fnOnDeviceFoundPtr = syscall.NewCallback(fnOnDeviceFound)
	}
	if nil != fnOnDeviceAccepted {
		fnOnDeviceAcceptedPtr = syscall.NewCallback(fnOnDeviceAccepted)
	}
	if nil != fnOnDeviceConnected {
		fnOnDeviceConnectedPtr = syscall.NewCallback(fnOnDeviceConnected)
	}
	if nil != fnOnDeviceDisconnected {
		fnOnDeviceDisconnectedPtr = syscall.NewCallback(fnOnDeviceDisconnected)
	}

	ret, _, err := syscall.Syscall6(procBs2SetDeviceEventListener, uintptr(5), context, fnOnDeviceFoundPtr, fnOnDeviceAcceptedPtr, fnOnDeviceConnectedPtr, fnOnDeviceDisconnectedPtr, uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func SearchDevices(context SdkContext) int16 {
	ret, _, err := syscall.Syscall(procBs2SearchDevices, uintptr(1), context, uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func SearchDevicesEx(context SdkContext, hostipAddr string) int16 {
	ret, _, err := syscall.Syscall(procBs2SearchDevicesEx, uintptr(2), context, uintptr(uintptr(unsafe.Pointer(syscall.StringBytePtr(hostipAddr)))), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func GetDevices(context SdkContext, arr *[]uint32) int16 {
	var devicelistObj *uint32
	var numDevice uint32

	ret, _, err := syscall.Syscall(procBs2GetDevices, uintptr(3), context, uintptr(unsafe.Pointer(&devicelistObj)), uintptr(unsafe.Pointer(&numDevice)))
	defer ReleaseObject(uintptr(unsafe.Pointer(devicelistObj)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}

	temp := make([]uint32, numDevice)

	//log.Printf("recevied result on getting device-id(parameter=%x).\n",devicelistObj)

	var i uintptr = 0
	for ; i < (uintptr)(numDevice); i++ {
		temp[i] = *(*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(devicelistObj)) + i*unsafe.Sizeof(*devicelistObj)))
	}
	*arr = temp

	return int16(ret)
}

func ConnectDevice(context SdkContext, deviceId BS2_DEVICE_ID) int16 {
	ret, _, err := syscall.Syscall(procBs2ConnectDevice, uintptr(2), context, uintptr(deviceId), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func ConnectDeviceViaIP(context SdkContext, deviceAddress string, defaultDevicePort uint16, deviceId *uint32) int16 {
	ret, _, err := syscall.Syscall6(procBs2ConnectDeviceViaIP, uintptr(4), context,
		uintptr(unsafe.Pointer(syscall.StringBytePtr(deviceAddress))),
		uintptr(defaultDevicePort),
		uintptr(unsafe.Pointer(deviceId)),
		uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func DisconnectDevice(context SdkContext, deviceId BS2_DEVICE_ID) int16 {
	ret, _, err := syscall.Syscall(procBs2DisconnectDevice, uintptr(2), context, uintptr(deviceId), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func SetKeepAliveTimeout(context SdkContext, ms int) int16 {
	ret, _, err := syscall.Syscall(procBs2SetKeepAliveTimeout, uintptr(2), context, uintptr(ms), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

type BS2Event struct {
	Id       uint32
	Datetime uint32
	deviceId BS2_DEVICE_ID

	Data [32]byte

	Code  uint16
	Param uint8
	Image bool
}

// typedef void (*OnAlarmFired)(BS2_DEVICE_ID deviceId, const BS2Event* event);
// typedef void (*OnInputDetected)(BS2_DEVICE_ID deviceId, const BS2Event* event);
// typedef void (*OnConfigChanged)(BS2_DEVICE_ID deviceId, uint32_t configMask);

func SetNotificationListener(context SdkContext, fnOnAlarmFired interface{}, fnOnInputDetected interface{}, fnOnConfigChanged interface{}) int16 {

	var fnOnAlarmFiredPtr, fnOnInputDetectedPtr, fnOnConfigChangedptr uintptr

	if nil != fnOnAlarmFired {
		fnOnAlarmFiredPtr = syscall.NewCallback(fnOnAlarmFired)
	}
	if nil != fnOnInputDetected {
		fnOnInputDetectedPtr = syscall.NewCallback(fnOnInputDetected)
	}
	if nil != fnOnConfigChanged {
		fnOnConfigChangedptr = syscall.NewCallback(fnOnConfigChanged)
	}

	ret, _, err := syscall.Syscall6(procBs2SetNotificationListener, uintptr(4), context,
		fnOnAlarmFiredPtr,
		fnOnInputDetectedPtr,
		fnOnConfigChangedptr,
		uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func SetServerPort(context SdkContext, serverPort uint16) int16 {
	ret, _, err := syscall.Syscall(procBs2SetServerPort, uintptr(2), context, uintptr(serverPort), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

// typedef uint32_t (*PreferMethod)(BS2_DEVICE_ID deviceID);
// typedef const char* (*GetRootCaFilePath)(uint32_t deviceId);
// typedef const char* (*GetServerCaFilePath)(BS2_DEVICE_ID deviceId);
// typedef const char* (*GetServerPrivateKeyFilePath)(uint32_t deviceId);
// typedef const char* (*GetPassword)(uint32_t deviceId);
// typedef void (*OnErrorOccured)(BS2_DEVICE_ID deviceId, int errCode);

func SetSSLHandler(context SdkContext, fnPreferMethod interface{}, fnGetRootCaFilePath interface{}, fnGetServerCaFilePath interface{}, fnGetServerPrivateKeyFilePath interface{}, fnGetPassword interface{}, fnOnErrorOccured interface{}) int16 {

	var fnPreferMethodPtr uintptr
	if nil != fnPreferMethod {
		fnPreferMethodPtr = syscall.NewCallback(fnPreferMethod)
	}
	var fnGetRootCaFilePathPtr uintptr
	if nil != fnGetRootCaFilePath {
		fnGetRootCaFilePathPtr = syscall.NewCallback(fnGetRootCaFilePath)
	}
	var fnGetServerCaFilePathPtr uintptr
	if nil != fnGetServerCaFilePath {
		fnGetServerCaFilePathPtr = syscall.NewCallback(fnGetServerCaFilePath)
	}
	var fnGetServerPrivateKeyFilePathPtr uintptr
	if nil != fnGetServerPrivateKeyFilePath {
		fnGetServerPrivateKeyFilePathPtr = syscall.NewCallback(fnGetServerPrivateKeyFilePath)
	}
	var fnGetPasswordPtr uintptr
	if nil != fnGetPassword {
		fnGetPasswordPtr = syscall.NewCallback(fnGetPassword)
	}
	var fnOnErrorOccuredPtr uintptr
	if nil != fnOnErrorOccured {
		fnOnErrorOccuredPtr = syscall.NewCallback(fnOnErrorOccured)
	}
	ret, _, err := syscall.Syscall9(procBs2SetSSLHandler, uintptr(7), context, fnPreferMethodPtr, fnGetRootCaFilePathPtr, fnGetServerCaFilePathPtr,
		fnGetServerPrivateKeyFilePathPtr, fnGetPasswordPtr, fnOnErrorOccuredPtr, uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func DisableSSL(context SdkContext, deviceId BS2_DEVICE_ID) int16 {
	ret, _, err := syscall.Syscall(procBs2DisableSSL, uintptr(2), context, uintptr(deviceId), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

// if isError(ret) {
// 	log.Println("MakePinCode - Error no : " + itoa(int(ret)))
// }

type BS2SimpleDeviceInfo struct {
	Id                     uint32
	Type                   uint16
	ConnectionMode         uint8
	Ipv4Address            uint32
	Port                   uint16
	MaxNumOfUser           uint32
	UserNameSupported      uint8
	UserPhotoSupported     uint8
	PinSupported           uint8
	CardSupported          uint8
	FingerSupported        uint8
	FaceSupported          uint8
	WlanSupported          uint8
	TnaSupported           uint8
	TriggerActionSupported uint8
	WiegandSupported       uint8
	ImageLogSupported      uint8
	DnsSupported           uint8
	JobCodeSupported       uint8
	WiegandMultiSupported  uint8
	Rs485Mode              uint8
	SslSupported           uint8
	RootCertExist          uint8
	DualIDSupported        uint8
	UseAlphanumericID      uint8
	ConnectedIP            uint32
	PhraseCodeSupported    uint8
	Card1xSupported        uint8
	SystemExtSupported     uint8
	VoipSupported          uint8
}

func GetDeviceInfo(context SdkContext, deviceId BS2_DEVICE_ID, deviceInfo *BS2SimpleDeviceInfo) int16 {
	ret, _, err := syscall.Syscall(procBs2GetDeviceInfo, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(deviceInfo)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func GetDeviceTime(context SdkContext, deviceId BS2_DEVICE_ID, gmtTime *uint32) int16 {
	ret, _, err := syscall.Syscall(procBs2GetDeviceTime, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(gmtTime)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func SetDeviceTime(context SdkContext, deviceId BS2_DEVICE_ID, gmtTime uint32) int16 {
	ret, _, err := syscall.Syscall(procBs2SetDeviceTime, uintptr(3), context, uintptr(deviceId), uintptr(gmtTime))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func ClearDatabase(context SdkContext, deviceId BS2_DEVICE_ID) int16 {
	ret, _, err := syscall.Syscall(procBs2SetDeviceTime, uintptr(2), context, uintptr(deviceId), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func FactoryReset(context SdkContext, deviceId BS2_DEVICE_ID) int16 {
	ret, _, err := syscall.Syscall(procBs2FactoryReset, uintptr(2), context, uintptr(deviceId), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func RebootDevice(context SdkContext, deviceId BS2_DEVICE_ID) int16 {
	ret, _, err := syscall.Syscall(procBs2RebootDevice, uintptr(2), context, uintptr(deviceId), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func LockDevice(context SdkContext, deviceId BS2_DEVICE_ID) int16 {
	ret, _, err := syscall.Syscall(procBs2LockDevice, uintptr(2), context, uintptr(deviceId), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func UnlockDevice(context SdkContext, deviceId BS2_DEVICE_ID) int16 {
	ret, _, err := syscall.Syscall(procBs2UnlockDevice, uintptr(2), context, uintptr(deviceId), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func UpgradeFirmware(context SdkContext, deviceId BS2_DEVICE_ID, firmwareData []byte, keepVerifyingSlaveDevice uint8, fnOnProgressChanged interface{}) int16 {
	ret, _, err := syscall.Syscall6(procBs2UpgradeFirmware, uintptr(6), context, uintptr(deviceId), uintptr(unsafe.Pointer(&firmwareData[0])), uintptr(len(firmwareData)), uintptr(keepVerifyingSlaveDevice), syscall.NewCallback(fnOnProgressChanged))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

type InnerResData struct {
	Index   uint8
	DataLen uint32
	Data    uint8
}

type BS2ResourceElement struct {
	TypeCode   uint8
	NumResData uint32
	ResData    [128]InnerResData
}

func UpdateResource(context SdkContext, deviceId BS2_DEVICE_ID, resourceElement *BS2ResourceElement, keepVerifyingSlaveDevice uint8, fnOnProgressChanged interface{}) int16 {
	ret, _, err := syscall.Syscall6(procBs2UpdateResource, uintptr(6), context, uintptr(deviceId), uintptr(unsafe.Pointer(resourceElement)), uintptr(keepVerifyingSlaveDevice), syscall.NewCallback(fnOnProgressChanged), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func GetLog(context SdkContext, deviceId BS2_DEVICE_ID, eventId uint32, amount uint32, logs *[]BS2Event) int16 {
	var eventPtr *BS2Event
	var numLog uint32

	ret, _, err := syscall.Syscall6(procBs2GetLog, uintptr(6), context, uintptr(deviceId), uintptr(eventId), uintptr(amount), uintptr(unsafe.Pointer(&eventPtr)), uintptr(unsafe.Pointer(&numLog)))
	defer ReleaseObject(uintptr(unsafe.Pointer(eventPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}

	temp := make([]BS2Event, numLog)

	var i uintptr = 0
	for ; i < (uintptr)(numLog); i++ {
		temp[i] = *(*BS2Event)(unsafe.Pointer(uintptr(unsafe.Pointer(eventPtr)) + i*unsafe.Sizeof(*eventPtr)))
	}
	*logs = temp

	return int16(ret)
}

func GetFilteredLog(context SdkContext, deviceId BS2_DEVICE_ID, userId string, eventCode uint16, start uint32, end uint32, tnaKey uint8, logs *[]BS2Event) int16 {
	var eventPtr *BS2Event
	var numLog uint32

	ret, _, err := syscall.Syscall9(procBs2GetFilteredLog, uintptr(9), context, uintptr(deviceId), uintptr(unsafe.Pointer(syscall.StringBytePtr(userId))), uintptr(eventCode), uintptr(start), uintptr(end), uintptr(tnaKey), uintptr(unsafe.Pointer(&eventPtr)), uintptr(unsafe.Pointer(&numLog)))
	defer ReleaseObject(uintptr(unsafe.Pointer(eventPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}

	temp := make([]BS2Event, numLog)

	var i uintptr = 0
	for ; i < (uintptr)(numLog); i++ {
		temp[i] = *(*BS2Event)(unsafe.Pointer(uintptr(unsafe.Pointer(eventPtr)) + i*unsafe.Sizeof(*eventPtr)))
	}
	*logs = temp

	return int16(ret)
}

func ClearLog(context SdkContext, deviceId BS2_DEVICE_ID) int16 {
	ret, _, err := syscall.Syscall(procBs2ClearLog, uintptr(2), context, uintptr(deviceId), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func StartMonitoringLog(context SdkContext, deviceId BS2_DEVICE_ID, fnOnLogReceived interface{}) int16 {

	var fnOnLogReceivedPtr uintptr
	if nil != fnOnLogReceived {
		fnOnLogReceivedPtr = syscall.NewCallback(fnOnLogReceived)
	}

	ret, _, err := syscall.Syscall(procBs2StartMonitoringLog, uintptr(3), context, uintptr(deviceId), fnOnLogReceivedPtr)
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	log.Printf("started monitoring events on the device(ID=%d)", deviceId)
	return int16(ret)
}

func StopMonitoringLog(context SdkContext, deviceId BS2_DEVICE_ID) int16 {
	ret, _, err := syscall.Syscall(procBs2StopMonitoringLog, uintptr(2), context, uintptr(deviceId), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	log.Printf("stopped monitoring events on the device(ID=%d)", deviceId)
	return int16(ret)
}

const (
	BS2_EVENT_MAX_IMAGE_SIZE = 16 * 1024
)

type BS2EventExtInfo struct {
	DateTime uint32
	deviceId BS2_DEVICE_ID
	Code     uint16
	Reserved [2]byte
}

type BS2EventBlob struct {
	EventMask uint16
	Id        uint32
	Info      BS2EventExtInfo
	IdData    [32]byte
	TnaKey    uint8
	JobCode   uint32
	ImageSize uint16
	Image     [BS2_EVENT_MAX_IMAGE_SIZE]uint8
	Reserved  uint8
}

func GetLogBlob(context SdkContext, deviceId BS2_DEVICE_ID, eventMask uint16, eventId uint32, amount uint32, logObj *[]BS2EventBlob) int16 {
	var eventBlob *BS2EventBlob
	var numLog uint32
	ret, _, err := syscall.Syscall9(procBs2GetLogBlob, uintptr(7), context, uintptr(deviceId), uintptr(eventMask), uintptr(eventId), uintptr(amount), uintptr(unsafe.Pointer(&eventBlob)), uintptr(unsafe.Pointer(&numLog)), uintptr(0), uintptr(0))
	defer ReleaseObject(uintptr(unsafe.Pointer(eventBlob)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2EventBlob, numLog)

	var i uintptr = 0
	for ; i < (uintptr)(numLog); i++ {
		temp[i] = *(*BS2EventBlob)(unsafe.Pointer(uintptr(unsafe.Pointer(eventBlob)) + i*unsafe.Sizeof(*eventBlob)))
	}
	*logObj = temp

	return int16(ret)
}

type BS2_EVENT_CODE = uint16
type BS2_TIMESTAMP = uint32
type BS2_EVENT_ID = uint32

func GetFilteredLogSinceEventId(context SdkContext, deviceId BS2_DEVICE_ID, userId string, eventCode BS2_EVENT_CODE, start BS2_TIMESTAMP, end BS2_TIMESTAMP, tnaKey uint8, lastEventId BS2_EVENT_ID, amount uint32, logs *[]BS2Event) int16 {
	var logsObj *BS2Event
	var numLog uint32
	ret, _, err := syscall.Syscall12(
		procBs2GetFilteredLogSinceEventId,
		uintptr(11),
		context, uintptr(deviceId), uintptr(unsafe.Pointer(syscall.StringBytePtr(userId))), uintptr(eventCode),
		uintptr(start), uintptr(end), uintptr(tnaKey), uintptr(lastEventId), uintptr(amount),
		uintptr(unsafe.Pointer(&logsObj)), uintptr(unsafe.Pointer(&numLog)), uintptr(0))
	defer ReleaseObject(uintptr(unsafe.Pointer(logsObj)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2Event, numLog)

	var i uintptr = 0
	for ; i < (uintptr)(numLog); i++ {
		temp[i] = *(*BS2Event)(unsafe.Pointer(uintptr(unsafe.Pointer(logsObj)) + i*unsafe.Sizeof(*logsObj)))
	}
	*logs = temp

	return int16(ret)
}

type BS2_USER_ID = [32]byte

func GetUserList(context SdkContext, deviceId BS2_DEVICE_ID, userIds *[]BS2_USER_ID, fnIsAcceptableUserID interface{}) int16 {
	var uidsObjs *BS2_USER_ID
	var numUid uint32

	var fnIsAcceptableUserIDPtr uintptr
	if nil != fnIsAcceptableUserID {
		fnIsAcceptableUserIDPtr = syscall.NewCallback(fnIsAcceptableUserID)
	}
	ret, _, err := syscall.Syscall6(procBs2GetUserList, uintptr(5), context, uintptr(deviceId), uintptr(unsafe.Pointer(&uidsObjs)), uintptr(unsafe.Pointer(&numUid)), fnIsAcceptableUserIDPtr, uintptr(0))
	defer ReleaseObject(uintptr(unsafe.Pointer(uidsObjs)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2_USER_ID, numUid)

	var i uintptr = 0
	for ; i < (uintptr)(numUid); i++ {
		temp[i] = *(*BS2_USER_ID)(unsafe.Pointer(uintptr(unsafe.Pointer(uidsObjs)) + i*unsafe.Sizeof(*uidsObjs)))
	}
	*userIds = temp

	return int16(ret)
}

type BS2User struct {
	UserID        BS2_USER_ID
	FormatVersion uint8
	Flag          uint8
	Version       uint16
	NumCards      uint8
	NumFingers    uint8
	NumFaces      uint8
	Reserved2     uint8
	AuthGroupID   uint32
	FaceChecksum  uint32
}

type BS2UserSetting struct {
	StartTime      uint32
	EndTime        uint32
	FingerAuthMode uint8
	CardAuthMode   uint8
	IdAuthMode     uint8
	SecurityLevel  uint8
}

type BS2UserBlob struct {
	User          BS2User
	Setting       BS2UserSetting
	Name          [BS2_USER_NAME_LEN]uint8
	Photo         BS2UserPhoto
	Pin           [BS2_PIN_HASH_SIZE]uint8
	CardObjs      uint32 //*BS2CSNCard
	FingerObjs    uint32 //*BS2Fingerprint
	FaceObjs      uint32 //*BS2Face
	AccessGroupId [BS2_MAX_NUM_OF_ACCESS_GROUP_PER_USER]uint32
}

type BS2UserPhoto struct {
	Size uint32
	Data [BS2_USER_PHOTO_SIZE]byte
}

const (
	USER_PAGE_SIZE = 1024
)

func GetUserInfos(context SdkContext, deviceId BS2_DEVICE_ID, userIds []BS2_USER_ID, userBlob *[]BS2UserBlob) int16 {

	if len(userIds) > len(*userBlob) || 0 == len(*userBlob) {
		return BS_WRAPPER_ERROR_SMALL_BUFFER
	}
	var userIdsPtr uintptr
	if 0 != len(userIds) {
		userIdsPtr = uintptr(unsafe.Pointer(&userIds[0]))
	}

	ret, _, err := syscall.Syscall6(procBs2GetUserInfos, uintptr(5), context, uintptr(deviceId), userIdsPtr, uintptr(len(*userBlob)), uintptr(unsafe.Pointer(&(*userBlob)[0])), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}

	return int16(ret)
}

func GetUserDatas(context SdkContext, deviceId BS2_DEVICE_ID, userIds []BS2_USER_ID, userBlob *[]BS2UserBlob, userMask uint32) int16 {

	if len(userIds) > len(*userBlob) || 0 == len(*userBlob) {
		return BS_WRAPPER_ERROR_SMALL_BUFFER
	}
	var userIdsPtr uintptr
	if 0 != len(userIds) {
		userIdsPtr = uintptr(unsafe.Pointer(&userIds[0]))
	}

	ret, _, err := syscall.Syscall6(procBs2GetUserDatas, uintptr(6), context, uintptr(deviceId), userIdsPtr, uintptr(len(*userBlob)), uintptr(unsafe.Pointer(&(*userBlob)[0])), uintptr(userMask))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}

	return int16(ret)
}

func EnrolUser(context SdkContext, deviceId BS2_DEVICE_ID, userBlob *[]BS2UserBlob, overwrite uint8) int16 {
	var userBlobPtr uintptr
	if 0 != len(*userBlob) {
		userBlobPtr = uintptr(unsafe.Pointer(&(*userBlob)[0]))
	}

	ret, _, err := syscall.Syscall6(procBs2EnrolUser, uintptr(5), context, uintptr(deviceId), userBlobPtr, uintptr(len(*userBlob)), uintptr(overwrite), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func RemoveUser(context SdkContext, deviceId BS2_DEVICE_ID, userIds []BS2_USER_ID) int16 {
	var userIdsPtr uintptr
	if 0 != len(userIds) {
		userIdsPtr = uintptr(unsafe.Pointer(&userIds[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2RemoveUser, uintptr(4), context, uintptr(deviceId), userIdsPtr, uintptr(len(userIds)), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func RemoveAllUser(context SdkContext, deviceId BS2_DEVICE_ID) int16 {
	ret, _, err := syscall.Syscall(procBs2RemoveAllUser, uintptr(3), context, uintptr(deviceId), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

type BS2_JOB_CODE = uint32
type BS2_JOB_LABEL = uint8

type InnerJobCode struct {
	Code  BS2_JOB_CODE
	Label BS2_JOB_LABEL
}

type BS2Job struct {
	NumJobs  uint8
	Reserved [3]byte
	Jobs     [BS2_MAX_JOB_SIZE]InnerJobCode
}

type BS2_USER_PHRAZE = uint32 // check up this.
type BS2UserBlobEx struct {
	User          BS2User
	Setting       BS2UserSetting
	Nmae          [BS2_USER_NAME_SIZE]byte
	Photo         BS2UserPhoto
	Pin           [BS2_PIN_HASH_SIZE]byte
	CardObjs      *BS2CSNCard
	FingerObjs    *BS2Fingerprint
	FaceObj       *BS2Face
	Job           *BS2Job
	Phrase        *BS2_USER_PHRAZE
	AccessGroupId [BS2_MAX_NUM_OF_ACCESS_GROUP_PER_USER]uint32
}

func GetUserInfosEx(context SdkContext, deviceId BS2_DEVICE_ID, userIds []BS2_USER_ID, userBlob *[]BS2UserBlobEx) int16 {

	if len(userIds) > len(*userBlob) || 0 == len(*userBlob) {
		return BS_WRAPPER_ERROR_SMALL_BUFFER
	}
	var userIdsPtr uintptr
	if 0 != len(userIds) {
		userIdsPtr = uintptr(unsafe.Pointer(&userIds[0]))
	}

	ret, _, err := syscall.Syscall6(procBs2GetUserInfosEx, uintptr(5), context, uintptr(deviceId), userIdsPtr, uintptr(len(*userBlob)), uintptr(unsafe.Pointer(&(*userBlob)[0])), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}

	return int16(ret)
}

func EnrolUserEx(context SdkContext, deviceId BS2_DEVICE_ID, userBlob *[]BS2UserBlobEx, overwrite uint8) int16 {
	countOfUser := len(*userBlob)

	var userBlobPtr uintptr
	if 0 != len(*userBlob) {
		userBlobPtr = uintptr(unsafe.Pointer(&(*userBlob)[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2EnrolUserEx, uintptr(5), context, uintptr(deviceId), userBlobPtr, uintptr(countOfUser), uintptr(overwrite), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func GetUserDatabaseInfo(context SdkContext, deviceId BS2_DEVICE_ID, numUsers *uint32, numCards *uint32, numFingers *uint32, numFaces *uint32, fnIsAcceptableUserID interface{}) int16 {
	var fnIsAcceptableUserIDPtr uintptr
	if nil != fnIsAcceptableUserID {
		fnIsAcceptableUserIDPtr = syscall.NewCallback(fnIsAcceptableUserID)
	}

	ret, _, err := syscall.Syscall9(procBs2GetUserDatabaseInfo, uintptr(7), context, uintptr(deviceId), uintptr(unsafe.Pointer(numUsers)), uintptr(unsafe.Pointer(numCards)), uintptr(unsafe.Pointer(numFingers)), uintptr(unsafe.Pointer(numFaces)), fnIsAcceptableUserIDPtr, uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func ResetConfig(context SdkContext, deviceId BS2_DEVICE_ID, includingDB bool) int16 {
	var includingDBNum uint8
	if includingDB {
		includingDBNum = 1
	}
	ret, _, err := syscall.Syscall(procBs2ResetConfig, uintptr(3), context, uintptr(deviceId), uintptr(includingDBNum))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func ResetConfigExceptNetInfo(context SdkContext, deviceId BS2_DEVICE_ID, includingDB uint8) int16 {
	ret, _, err := syscall.Syscall(procBs2ResetConfigExceptNetInfo, uintptr(3), context, uintptr(deviceId), uintptr(includingDB))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

type BS2Configs struct {
	ConfigMask          uint32
	FactoryConfig       BS2FactoryConfig
	SystemConfig        BS2SystemConfig
	AuthConfig          BS2AuthConfig
	StatusConfig        BS2StatusConfig
	DisplayConfig       BS2DisplayConfig
	IpConfig            BS2IpConfig
	IpConfigExt         BS2IpConfigExt
	TnaConfig           BS2TNAConfig
	CardConfig          BS2CardConfig
	FingerprintConfig   BS2FingerprintConfig
	Rs485Config         BS2Rs485Config
	WiegandConfig       BS2WiegandConfig
	WiegandDeviceConfig BS2WiegandDeviceConfig
	InputConfig         BS2InputConfig
	WlanConfig          BS2WlanConfig
	TriggerActionConfig BS2TriggerActionConfig
	EventConfig         BS2EventConfig
	WiegandMultiConfig  BS2WiegandMultiConfig
	Card1xConfig        BS1CardConfig
	SystemExtConfig     BS2SystemConfigExt
	VoipConfig          BS2VoipConfig
	FaceConfig          BS2FaceConfig
}

func GetConfig(context SdkContext, deviceId BS2_DEVICE_ID, configs *BS2Configs) int16 {
	ret, _, err := syscall.Syscall(procBs2GetConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(configs)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func SetConfig(context SdkContext, deviceId BS2_DEVICE_ID, configs *BS2Configs) int16 {
	ret, _, err := syscall.Syscall(procBs2SetConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(configs)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

const (
	BS2_MAC_ADDR_LEN     = 6
	BS2_MODEL_NAME_LEN   = 32
	BS2_KERNEL_REV_LEN   = 32
	BS2_BSCORE_REV_LEN   = 32
	BS2_FIRMWARE_REV_LEN = 32
)

type Ver struct {
	major    uint8
	minor    uint8
	ext      uint8
	reserved [1]uint8
}

type BS2FactoryConfig struct {
	deviceID    uint32
	macAddr     [BS2_MAC_ADDR_LEN]uint8
	reserved    [2]uint8
	modelName   [BS2_MODEL_NAME_LEN]byte
	boardVer    Ver
	kernelVer   Ver
	bscoreVer   Ver
	firmwareVer Ver
	kernelRev   [BS2_KERNEL_REV_LEN]byte
	bscoreRev   [BS2_BSCORE_REV_LEN]byte
	firmwareRev [BS2_FIRMWARE_REV_LEN]byte
	reserved2   [32]uint8
}

func GetFactoryConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2FactoryConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2GetFactoryConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

type BS2SystemConfig struct {
	notUsed           [768]uint8
	Timezone          int32
	SyncTime          uint8
	ServerSync        uint8
	DeviceLocked      uint8
	UseInterphone     uint8
	UseUSBConnection  uint8
	KeyEncrypted      uint8
	UseJobCode        uint8
	UseAlphanumericID uint8
	CameraFrequency   uint32
	SecureTamper      bool

	tamperOn  bool //(writeprotected)
	reserved  [2]uint8
	reserved2 [20]uint8
}

func GetSystemConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2SystemConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2GetSystemConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func SetSystemConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2SystemConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2SetSystemConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

const (
	BS2_NUM_OF_AUTH_MODE = 11
	BS2_MAX_OPERATORS    = 10
)

type Oper struct {
	UserID   [BS2_USER_ID_SIZE]byte
	Level    uint8
	Reserved [3]uint8
}

type BS2AuthConfig struct {
	AuthSchedule        [BS2_NUM_OF_AUTH_MODE]uint32
	UseGlobalAPB        uint8
	GlobalAPBFailAction uint8
	UseGroupMatching    uint8
	Reserved            uint8
	Reserved2           [28]uint8
	UsePrivateAuth      uint8
	FaceDetectionLevel  uint8
	UseServerMatching   uint8
	UseFullAccess       uint8
	MatchTimeout        uint8
	AuthTimeout         uint8
	NumOperators        uint8
	Reserved3           [1]uint8
	Operators           [BS2_MAX_OPERATORS]Oper
}

func GetAuthConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2AuthConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2GetAuthConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func SetAuthConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2AuthConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2SetAuthConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

const (
	BS2_LED_SIGNAL_NUM    = 3
	BS2_BUZZER_SIGNAL_NUM = 3
	BS2_DEVICE_STATUS_NUM = 15
)

type BS2LedSignal struct {
	Color    uint8
	reserved uint8
	Duration uint16
	Delay    uint16
}

type BS2BuzzerSignal struct {
	Tone     uint8
	FadeOut  bool
	Duration uint16
	Delay    uint16
}

type Led_T struct {
	Enabled  uint8
	Reserved [1]uint8
	Count    uint16
	Signal   [BS2_LED_SIGNAL_NUM]BS2LedSignal
}
type Buzzer_T struct {
	Enabled  uint8
	Reserved [1]uint8
	Count    uint16
	Signal   [BS2_BUZZER_SIGNAL_NUM]BS2BuzzerSignal
}
type BS2StatusConfig struct {
	Led                [BS2_DEVICE_STATUS_NUM]Led_T
	Reserved1          [32]uint8
	Buzzer             [BS2_DEVICE_STATUS_NUM]Buzzer_T
	ConfigSyncRequired uint8
	Reserved2          [31]uint8
}

func GetStatusConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2StatusConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2GetStatusConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func SetStatusConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2StatusConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2SetStatusConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

const (
	BS2_MAX_SHORTCUT_HOME = 8
)

type BS2DisplayConfig struct {
	Language   uint32
	Background uint8
	Volume     uint8
	BgTheme    uint8
	DateFormat uint8

	MenuTimeout      uint16
	MsgTimeout       uint16
	BacklightTimeout uint16
	DisplayDateTime  bool
	UseVoice         bool

	TimeFormat    uint8
	HomeFormation uint8
	UseUserPhrase bool
	Reserved      uint8

	ShortcutHome [BS2_MAX_SHORTCUT_HOME]uint8
	TnaIcon      [BS2_MAX_TNA_KEY]uint8

	Reserved1 [32]uint8
}

func GetDisplayConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2DisplayConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2GetDisplayConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func SetDisplayConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2DisplayConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2SetDisplayConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

type BS2IpConfig struct {
	ConnectionMode uint8    ///< 1 byte
	UseDHCP        bool     ///< 1 byte
	UseDNS         bool     ///< 1 byte
	Reserved       uint8    ///< 1 byte (packing)
	IpAddress      [16]byte ///< 16 bytes
	Gateway        [16]byte ///< 16 bytes
	SubnetMask     [16]byte ///< 16 bytes
	ServerAddr     [16]byte ///< 16 bytes
	Port           uint16   ///< 2 bytes
	ServerPor      uint16   ///< 2 bytes
	MtuSize        uint16   ///< 2 bytes
	Baseband       uint8    ///< 1 byte
	Reserved2      uint8    ///< 1 byte (packing)

	SslServerPort uint16    ///< 2 bytes
	Reserved3     [30]uint8 ///< 30 bytes (reserved)
}

func GetIPConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2IpConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2GetIPConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func GetIPConfigViaUDP(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2IpConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2GetIPConfigViaUDP, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func SetIPConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2IpConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2SetIPConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func SetIPConfigViaUDP(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2IpConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2SetIPConfigViaUDP, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

type BS2IpConfigExt struct {
	dnsAddr   [16]byte
	serverUrl [256]byte
	reserved  [32]uint8
}

func GetIPConfigExt(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2IpConfigExt) int16 {
	ret, _, err := syscall.Syscall(procBs2GetIPConfigExt, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func SetIPConfigExt(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2IpConfigExt) int16 {
	ret, _, err := syscall.Syscall(procBs2SetIPConfigExt, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

const (
	BS2_MAX_TNA_LABEL_LEN = 16 * 3
)

const (
	// tnaKey
	BS2_TNA_UNSPECIFIED = 0

	BS2_TNA_KEY_1
	BS2_TNA_KEY_2
	BS2_TNA_KEY_3
	BS2_TNA_KEY_4
	BS2_TNA_KEY_5
	BS2_TNA_KEY_6
	BS2_TNA_KEY_7
	BS2_TNA_KEY_8
	BS2_TNA_KEY_9
	BS2_TNA_KEY_10
	BS2_TNA_KEY_11
	BS2_TNA_KEY_12
	BS2_TNA_KEY_13
	BS2_TNA_KEY_14
	BS2_TNA_KEY_15
	BS2_TNA_KEY_16

	BS2_MAX_TNA_KEY = 16
)

type BS2TNAInfo struct {
	TnaMode     uint8
	TnaKey      uint8
	TnaRequired bool
	Reserved    uint8

	TnaSchedule [BS2_MAX_TNA_KEY]uint32
	Unused      [BS2_MAX_TNA_KEY]uint8
}

type BS2TNAExtInfo struct {
	TnaLabel [BS2_MAX_TNA_KEY][BS2_MAX_TNA_LABEL_LEN]byte
	Unused   [BS2_MAX_TNA_KEY]uint8
}

type BS2TNAConfig struct {
	TnaInfo    BS2TNAInfo
	TnaExtInfo BS2TNAExtInfo

	Reserved2 [32]byte
}

func GetTNAConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2TNAConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2GetTNAConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func SetTNAConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2TNAConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2SetTNAConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

/**
*  BS2MifareCard
 */
type BS2MifareCard struct {
	PrimaryKey      [6]uint8
	Reserved1       [2]uint8
	SecondaryKey    [6]uint8
	Reserved2       [2]uint8
	StartBlockIndex uint16
	Reserved        [6]uint8
} //24 Bytes

type BS2IClassCard struct {
	PrimaryKey      [8]uint8
	SecondaryKey    [8]uint8
	StartBlockIndex uint16
	Reserved        [6]uint8
} //24 Bytes

type BS2DesFireCard struct {
	PrimaryKey     [16]uint8
	SecondaryKey   [16]uint8
	AppID          [3]uint8
	FileID         uint8
	EncryptionType uint8 //for DesFire DES/3DES or AES. AES will be provided at future(TBD).
	Reserved       [3]uint8
} //40 Bytes

type BS2SeosCard struct {
	Oid_ADF          [13]uint8
	Size_ADF         uint8
	Reserved1        [2]uint8  ///< 16 bytes
	Oid_DataObjectID [8]uint8  ///< 24 bytes
	Size_DataObject  [8]uint16 ///< 40 bytes
	PrimaryKeyAuth   [16]uint8
	SecondaryKeyAuth [16]uint8 ///< 72 bytes
	Reserved2        [24]uint8
} ///< 96 bytes

type BS2CardConfig struct {
	ByteOrder        uint8 ///< 1 byte
	UseWiegandFormat bool  ///< 1 byte

	DataType        uint8 ///< 1 byte
	UseSecondaryKey bool  ///< 1 byte

	Mifare  BS2MifareCard  ///< 24 bytes
	Iclass  BS2IClassCard  ///< 24 bytes
	Desfire BS2DesFireCard ///< 40 bytes

	FormatID uint32    ///< 4 bytes (card format ID / use only application)
	Reserved [24]uint8 ///< 24 bytes (packing)
}

func GetCardConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2CardConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2GetCardConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func SetCardConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2CardConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2SetCardConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

type BS2_FINGER_TEMPLATE_FORMAT = uint8
type BS2_FINGER_SECURITY_LEVEL = uint8
type BS2_FINGER_FAST_MODE = uint8
type BS2_FINGER_SENSITIVITY = uint8
type BS2_FINGER_SENSOR_MODE = uint8

type BS2FingerprintConfig struct {
	SecurityLevel BS2_FINGER_SECURITY_LEVEL
	FastMode      BS2_FINGER_FAST_MODE
	Sensitivity   BS2_FINGER_SENSITIVITY
	SensorMode    BS2_FINGER_SENSOR_MODE

	TemplateFormat BS2_FINGER_TEMPLATE_FORMAT
	Reserved       uint8
	ScanTimeout    uint16

	SuccessiveScan     bool
	AdvancedEnrollment bool
	ShowImage          bool
	LfdLevel           uint8

	Reserved3 [32]uint8
}

func GetFingerprintConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2FingerprintConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2GetFingerprintConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func SetFingerprintConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2FingerprintConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2SetFingerprintConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

const (
	BS2_RS485_MAX_CHANNELS           = 4
	BS2_RS485_MAX_SLAVES_PER_CHANNEL = 32
)

type BS2_RS485_MODE = uint8
type BS2_DEVICE_TYPE = uint16

type BS2Rs485SlaveDevice struct {
	DeviceID   BS2_DEVICE_ID
	DeviceType BS2_DEVICE_TYPE
	EnableOSDP bool
	Connected  bool
}

type BS2Rs485Channel struct {
	BaudRate      uint32
	ChannelIndex  uint8
	UseRegistance uint8
	NumOfDevices  uint8
	Reserved      uint8
	SlaveDevices  [BS2_RS485_MAX_SLAVES_PER_CHANNEL]BS2Rs485SlaveDevice
}

type BS2Rs485Config struct {
	Mode          BS2_RS485_MODE
	NumOfChannels uint8
	Reserved      [2]uint8

	Reserved1 [32]uint8

	Channels [BS2_RS485_MAX_CHANNELS]BS2Rs485Channel
}

func GetRS485Config(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2Rs485Config) int16 {
	ret, _, err := syscall.Syscall(procBs2GetRS485Config, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func SetRS485Config(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2Rs485Config) int16 {
	ret, _, err := syscall.Syscall(procBs2SetRS485Config, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

const (
	BS2_WIEGAND_FIELD_SIZE      = 32 ///< 256 bits
	BS2_WIEGAND_FIELD_MASK_SIZE = BS2_WIEGAND_FIELD_SIZE * 8
	BS2_WIEGAND_MAX_FIELDS      = 4
	BS2_WIEGAND_MAX_PARITIES    = 4
)

type BS2_WIEGAND_PARITY uint8

type BS2WiegandFormat struct {
	length       uint32
	idFields     [BS2_WIEGAND_MAX_FIELDS][BS2_WIEGAND_FIELD_SIZE]uint8
	parityFields [BS2_WIEGAND_MAX_PARITIES][BS2_WIEGAND_FIELD_SIZE]uint8
	parityType   [BS2_WIEGAND_MAX_PARITIES]BS2_WIEGAND_PARITY
	parityPos    [BS2_WIEGAND_MAX_PARITIES]uint8
}

type BS2WiegandConfig struct {
	InOut            uint8
	UseWiegandBypass bool  ///< 1 byte
	UseFailCode      bool  ///< 1 byte
	FailCode         uint8 ///< 1 byte

	outPulseWidth    uint16 ///< 2 bytes (20 ~ 100 us, default = 40)
	outPulseInterval uint16 ///< 2 bytes (200 ~ 20000 us, default = 10000)

	FormatID uint32 ///< 4 bytes (wiegand format ID)
	Format   BS2WiegandFormat

	WiegandInputMask uint16   ///< 2 Bytes (bitmask , no use 0 postion bit, 1~15 bit)
	WiegandCardMask  uint16   ///< 2 Bytes (bitmask , no use 0 postion bit, 1~15 bit)
	WiegandCSNIndex  uint8    ///< 1 Bytes (1~15)
	Reserved         [27]byte ///< 27 bytes (reserved)
}

func GetWiegandConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2WiegandConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2GetWiegandConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func SetWiegandConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2WiegandConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2SetWiegandConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

const (
	BS2_WIEGAND_STATUS_NUM = 2
)

type BS2WiegandTamperInput struct {
	DeviceID   uint32
	Port       uint16
	SwitchType uint8
	Reserved   uint8
}

type BS2WiegandBuzzerOutput struct {
	DeviceID uint32
	Port     uint16

	Reserved [34]byte
}

type BS2WiegandLedOutput struct {
	DeviceID uint32
	port     uint16
	reserved [10]byte
}

type BS2WiegandDeviceConfig struct {
	Tamper   BS2WiegandTamperInput
	Led      [BS2_WIEGAND_STATUS_NUM]BS2WiegandLedOutput
	Buzzer   BS2WiegandBuzzerOutput
	Reserved [32]uint32
}

func GetWiegandDeviceConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2WiegandDeviceConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2GetWiegandDeviceConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func SetWiegandDeviceConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2WiegandDeviceConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2SetWiegandDeviceConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

const (
	BS2_MAX_INPUT_NUM = 8
)

const (
	BS2_INPUT_TYPE_NORMAL = iota
	BS2_INPUT_TYPE_SUPERVISERD
)

const (
	SUPERVISED_REG_1K = 0
	SUPERVISED_REG_2_2K
	SUPERVISED_REG_4_7K
	SUPERVISED_REG_10K

	SUPERVISED_REG_CUSTOM = 255
)

type BS2SVInputRange struct {
	MinValue uint16
	MaxValue uint16
}

type BS2SupervisedInputConfig struct {
	ShortInput BS2SVInputRange
	OpenInput  BS2SVInputRange
	OnInput    BS2SVInputRange
	OffInput   BS2SVInputRange
}

type SuprevisedInput struct {
	PortIndex        uint8
	Enabled          bool
	Supervised_index uint8

	Reserved [5]uint8

	Config BS2SupervisedInputConfig
}

type BS2InputConfig struct {
	numInputs         uint8
	numSupervised     uint8
	reseved           uint16
	supervised_inputs [BS2_MAX_INPUT_NUM]SuprevisedInput
}

func GetInputConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2InputConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2GetInputConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func SetInputConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2InputConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2SetInputConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

const (
	BS2_WLAN_SSID_SIZE = 32
	BS2_WLAN_KEY_SIZE  = 64
)

const (
	BS2_WLAN_OPMODE_MANAGED = iota
	BS2_WLAN_OPMODE_ADHOC

	BS2_WLAN_OPMODE_DEFAULT = BS2_WLAN_OPMODE_MANAGED
)

const (
	BS2_WLAN_AUTH_OPEN = iota
	BS2_WLAN_AUTH_SHARED
	BS2_WLAN_AUTH_WPA_PSK
	BS2_WLAN_AUTH_WPA2_PSK

	BS2_WLAN_AUTH_DEFAULT = BS2_WLAN_AUTH_OPEN
)

const (
	BS2_WLAN_ENC_NONE = iota
	BS2_WLAN_ENC_WEP
	BS2_WLAN_ENC_TKIP_AES
	BS2_WLAN_ENC_AES
	BS2_WLAN_ENC_TKIP

	BS2_WLAN_ENC_DEFAULT = BS2_WLAN_ENC_NONE
)

type BS2WlanConfig struct {
	Enabled        bool
	OperationMode  uint8
	AuthType       uint8
	EncryptionType uint8

	Essid   [BS2_WLAN_SSID_SIZE]byte
	AuthKey [BS2_WLAN_KEY_SIZE]byte

	Reserved2 [32]uint8
}

func GetWlanConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2WlanConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2GetWlanConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func SetWlanConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2WlanConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2SetWlanConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

type BS2_TRIGGER_TYPE = uint8

const (
	BS2_TRIGGER_NONE = iota

	BS2_TRIGGER_EVENT
	BS2_TRIGGER_INPUT
	BS2_TRIGGER_SCHEDULE
)

const (
	BS2_SCHEDULE_TRIGGER_ON_STAR = iota
	BS2_SCHEDULE_TRIGGER_ON_END
)

type BS2EventTrigger struct {
	Code     uint16
	Reserved [2]byte
}

type BS2Trigger struct {
	DeviceID BS2_DEVICE_ID
	Type     BS2_TRIGGER_TYPE
	Reserved [3]uint8

	//Event BS2EventTrigger
	TriggerUnionData [8]byte
}

type BS2Action struct {
	DeviceID        uint32
	Type            uint8
	StopFlag        uint8
	Delay           uint16
	ActionUnionData [24]byte
}

const (
	BS2_MAX_TRIGGER_ACTION = 128
)

type BS2TriggerAction struct {
	trigger BS2Trigger
	action  BS2Action
}

/**
*	BS2TriggerActionConfig
 */
type BS2TriggerActionConfig struct {
	NumItems uint8    ///< 1 byte
	Reserved [3]uint8 ///< 3 bytes

	Items [BS2_MAX_TRIGGER_ACTION]BS2TriggerAction

	Reserved1 [32]uint8 ///< 32 bytes (reserved)
}

func GetTriggerActionConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2TriggerActionConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2GetTriggerActionConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func SetTriggerActionConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2TriggerActionConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2SetTriggerActionConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

const (
	BS2_EVENT_MAX_IMAGE_CODE_COUNT = 32
)

type ImageEventFilter struct {
	MainEventCode uint8
	Reserved      [3]uint8
	ScheduleID    uint32
}

type BS2EventConfig struct {
	NumImageEventFilter uint32

	ImageEventFilter [BS2_EVENT_MAX_IMAGE_CODE_COUNT]ImageEventFilter

	Reserved [32]uint8
}

func GetEventConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2EventConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2GetEventConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func SetEventConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2EventConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2SetEventConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

type BS2WiegandInConfig struct {
	FormatID uint32
	Format   BS2WiegandFormat

	Reserved [32]byte
}

const (
	MAX_WIEGAND_IN_COUNT = 15
)

type BS2WiegandMultiConfig struct {
	Formats  [MAX_WIEGAND_IN_COUNT]BS2WiegandInConfig
	Reserved [32]uint8
}

func GetWiegandMultiConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2WiegandMultiConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2GetWiegandMultiConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func SetWiegandMultiConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2WiegandMultiConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2SetWiegandMultiConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

const (
	MIFARE_KEY_SIZE     = 6
	MIFARE_MAX_TEMPLATE = 4

	VALID_MAGIC_NO = 0x1f1f1f1f
)

type BS1CardConfig struct {
	// Options
	MagicNo            uint32
	Disabled           uint32
	UseCSNOnly         uint32
	BioentryCompatible uint32

	// Keys
	UseSecondaryKey uint32
	Reserved1       uint32
	PrimaryKey      [MIFARE_KEY_SIZE]uint8
	Reserved2       [2]uint8
	SecondaryKey    [MIFARE_KEY_SIZE]uint8
	Reserved3       [2]uint8

	// Layout
	CisIndex           uint32
	NumOfTemplate      uint32
	TemplateSize       uint32
	TemplateStartBlock [MIFARE_MAX_TEMPLATE]uint32

	Reserve4 [15]uint32
}

func GetCard1xConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS1CardConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2GetCard1xConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func SetCard1xConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS1CardConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2SetCard1xConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

const (
	SEC_KEY_SIZE = 16
)

type BS2SystemConfigExt struct {
	PrimarySecureKey   [SEC_KEY_SIZE]byte
	SecondarySecureKey [SEC_KEY_SIZE]byte

	Reserved3 [32]uint8
}

func GetSystemExtConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2SystemConfigExt) int16 {
	ret, _, err := syscall.Syscall(procBs2GetSystemExtConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func SetSystemExtConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2SystemConfigExt) int16 {
	ret, _, err := syscall.Syscall(procBs2SetSystemExtConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

const (
	BS2_MAX_DESCRIPTION_NAME_LEN = 144
	BS2_VOIP_MAX_PHONEBOOK       = 32
)

type BS2_URL = [256]byte
type BS2_PORT = uint16

type BS2UserPhoneItem struct {
	PhoneNumber BS2_USER_ID
	Descript    [BS2_MAX_DESCRIPTION_NAME_LEN]byte

	Reserved2 [32]uint8
}

type BS2VoipConfig struct {
	ServerUrl  BS2_URL
	ServerPort BS2_PORT
	UserID     BS2_USER_ID
	UserPW     BS2_USER_ID

	ExitButton uint8
	DtmfMode   uint8
	BUse       bool
	Reseverd   [1]uint8

	NumPhonBook uint32
	Phonebook   [BS2_VOIP_MAX_PHONEBOOK]BS2UserPhoneItem

	Reserved2 [32]uint8
}

func GetVoipConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2VoipConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2GetVoipConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func SetVoipConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2VoipConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2SetVoipConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

type BS2_FACE_SECURITY_LEVEL = uint8
type BS2_FACE_LIGHT_CONDITON = uint8
type BS2_FACE_ENROLL_THRESHOLD = uint8
type BS2_FACE_DETECT_SENSITIVITY = uint8

type BS2FaceConfig struct {
	SecurityLevel     BS2_FACE_SECURITY_LEVEL
	LightCondition    BS2_FACE_LIGHT_CONDITON
	EnrollThreshold   BS2_FACE_ENROLL_THRESHOLD
	DetectSensitivity BS2_FACE_DETECT_SENSITIVITY

	EnrollTimeout uint16

	Reserved3 [32]uint8
}

func GetFaceConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2FaceConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2GetFaceConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func SetFaceConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2FaceConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2SetFaceConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

const (
	BS2_RS485_MAX_CHANNELS_EX = 8
)

const (
	BS2_RS485_INVALID_BAUD_RATE = -1
)

type BS2Rs485SlaveDeviceEX struct {
	DeviceID    uint32
	DeviceType  uint16
	EnableOSDP  bool
	Connected   bool
	ChannelInfo uint8
	Reserved    [3]byte
}

type BS2Rs485ChannelEX struct {
	BaudRate      uint32
	ChannelIndex  uint8
	UseRegistance uint8
	NumOfDevices  uint8
	Reserved      byte
	SlaveDevices  [BS2_RS485_MAX_SLAVES_PER_CHANNEL]BS2Rs485SlaveDeviceEX
}

type BS2Rs485ConfigEX struct {
	Mode [BS2_RS485_MAX_CHANNELS_EX]uint8

	NumOfChannels uint16
	Reserved      [2]uint8

	Reserved1 [32]uint8

	Channels [BS2_RS485_MAX_CHANNELS_EX]BS2Rs485ChannelEX
}

func GetRS485ConfigEx(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2Rs485ConfigEX) int16 {
	ret, _, err := syscall.Syscall(procBs2GetRS485ConfigEx, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func SetRS485ConfigEx(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2Rs485ConfigEX) int16 {
	ret, _, err := syscall.Syscall(procBs2SetRS485ConfigEx, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

type BS2CardConfigEx struct {
	Seos     BS2SeosCard ///< 96 bytes
	Reserved [24]uint8
}

func GetCardConfigEx(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2CardConfigEx) int16 {
	ret, _, err := syscall.Syscall(procBs2GetCardConfigEx, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func SetCardConfigEx(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2CardConfigEx) int16 {
	ret, _, err := syscall.Syscall(procBs2SetCardConfigEx, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

type BS2_YEAR = uint16
type BS2_MONTH = uint8
type BS2_WEEKDAY = uint8
type BS2_ORDINAL = int8

const (
	BS2_MAX_DST_SCHEDULE = 2
)

type BS2WeekTime struct {
	Year BS2_YEAR

	Month   BS2_MONTH
	Ordinal BS2_ORDINAL
	WeekDay BS2_WEEKDAY
	Hour    uint8
	Minute  uint8
	Second  uint8
}

type BS2DstSchedule struct {
	StartTime  BS2WeekTime
	EndTime    BS2WeekTime
	TimeOffset int32
	Reserved   [4]uint8
}

type BS2DstConfig struct {
	numSchedules uint8
	reserved     [31]uint8

	schedules [BS2_MAX_DST_SCHEDULE]BS2DstSchedule
}

func GetDstConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2DstConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2GetDstConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func SetDstConfig(context SdkContext, deviceId BS2_DEVICE_ID, config *BS2DstConfig) int16 {
	ret, _, err := syscall.Syscall(procBs2SetDstConfig, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(config)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

const (
	BS2_CARD_DATA_SIZE = 32
	BS2_CARD_ID_SIZE   = 24
)

const (
	BS2_CARD_TYPE_UNKNOWN     = 0x00
	BS2_CARD_TYPE_CSN         = 0x01
	BS2_CARD_TYPE_SECURE      = 0x02
	BS2_CARD_TYPE_ACCESS      = 0x03
	BS2_CARD_TYPE_WIEGAND     = 0x0A
	BS2_CARD_TYPE_CONFIG_CARD = 0x0B
)

type BS2Card struct {
	isSmartCard uint8
	Data        [1656]byte
}

type BS2SmartCardHeader struct {
	HdrCRC        uint16
	CardCRC       uint16
	CardType      BS2_CARD_TYPE
	NumOfTemplate uint8
	TemplateSize  uint16
	IssueCount    uint16
	DuressMask    uint8
	Reserved      [5]uint8
}

const (
	BS2_SMART_CARD_MAX_TEMPLATE_COUNT     = 4
	BS2_SMART_CARD_MAX_ACCESS_GROUP_COUNT = 16
)

type BS2SmartCardCredentials struct {
	Pin          [BS2_PIN_HASH_SIZE]uint8
	templateData [BS2_SMART_CARD_MAX_TEMPLATE_COUNT * BS2_FINGER_TEMPLATE_SIZE]uint8
}

type BS2_DATETIME = uint32
type BS2AccessOnCardData struct {
	AccessGroupID [BS2_SMART_CARD_MAX_ACCESS_GROUP_COUNT]uint16
	StartTime     BS2_DATETIME
	EndTime       BS2_DATETIME
}

type BS2SmartCardData struct {
	Header       BS2SmartCardHeader
	CardID       [BS2_CARD_DATA_SIZE]uint8
	Credentials  BS2SmartCardCredentials
	AccessOnData BS2AccessOnCardData
}

func ScanCard(context SdkContext, deviceId BS2_DEVICE_ID, card *BS2Card, fnOnReadyToScan interface{}) int16 {
	var fnOnReadyToScanPtr uintptr
	if nil != fnOnReadyToScan {
		fnOnReadyToScanPtr = uintptr(syscall.NewCallback(fnOnReadyToScan))
	}
	ret, _, err := syscall.Syscall6(procBs2ScanCard, uintptr(4), context, uintptr(deviceId), uintptr(unsafe.Pointer(card)), fnOnReadyToScanPtr, uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func WriteCard(context SdkContext, deviceId BS2_DEVICE_ID, smartCard *BS2SmartCardData) int16 {
	ret, _, err := syscall.Syscall(procBs2WriteCard, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(smartCard)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func EraseCard(context SdkContext, deviceId BS2_DEVICE_ID) int16 {
	ret, _, err := syscall.Syscall(procBs2EraseCard, uintptr(2), context, uintptr(deviceId), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

type BS2Fingerprint struct {
	Index    uint8
	Flag     uint8
	Reserved [2]uint8
	Data     [BS2_TEMPLATE_PER_FINGER][BS2_FINGER_TEMPLATE_SIZE]uint8
}

func ScanFingerprint(context SdkContext, deviceId BS2_DEVICE_ID, finger *BS2Fingerprint, templateIndex uint32, quality uint32, templateFormat uint8, fnOnReadyToScan interface{}) int16 {
	var fnOnReadyToScanPtr uintptr
	if nil != fnOnReadyToScan {
		fnOnReadyToScanPtr = uintptr(syscall.NewCallback(fnOnReadyToScan))
	}
	ret, _, err := syscall.Syscall9(procBs2ScanFingerprint, uintptr(7), context, uintptr(deviceId), uintptr(unsafe.Pointer(finger)), uintptr(templateIndex), uintptr(quality), uintptr(templateFormat), fnOnReadyToScanPtr, uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func ScanFingerprintEx(context SdkContext, deviceId BS2_DEVICE_ID, finger *BS2Fingerprint, templateIndex uint32, quality uint32, templateFormat uint8, outQuality *uint32, fnOnReadyToScan interface{}) int16 {
	if 0 == procBs2ScanFingerprintEx {
		return BS_WRAPPER_ERROR_NOT_SUPPORTED
	}
	ret, _, err := syscall.Syscall9(procBs2ScanFingerprintEx, uintptr(8), context, uintptr(deviceId), uintptr(unsafe.Pointer(finger)), uintptr(templateIndex), uintptr(quality), uintptr(templateFormat), uintptr(unsafe.Pointer(outQuality)), uintptr(syscall.NewCallback(fnOnReadyToScan)), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func VerifyFingerprint(context SdkContext, deviceId BS2_DEVICE_ID, finger *BS2Fingerprint) int16 {
	ret, _, err := syscall.Syscall(procBs2VerifyFingerprint, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(finger)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func GetLastFingerprintImage(context SdkContext, deviceId BS2_DEVICE_ID, imageObj *[]byte, imageWidth *uint32, imageHeight *uint32) int16 {
	var buffer *byte
	ret, _, err := syscall.Syscall6(procBs2GetLastFingerprintImage, uintptr(5), context, uintptr(deviceId), uintptr(unsafe.Pointer(&buffer)), uintptr(unsafe.Pointer(imageWidth)), uintptr(unsafe.Pointer(imageHeight)), uintptr(0))
	defer ReleaseObject(uintptr(unsafe.Pointer(buffer)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}

	temp := make([]byte, 2048)

	var i uintptr = 0
	for ; i < 1024; i++ {
		temp[i] = *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(buffer)) + i*unsafe.Sizeof(*buffer)))
	}
	*imageObj = temp

	return int16(ret)
}

type BS2Face struct {
	FaceIndex     uint8
	NumOfTemplate uint8
	Flag          uint8
	Reserved      uint8

	ImageLen  uint16
	Reserved2 [2]uint8

	ImageData    [BS2_FACE_IMAGE_SIZE]uint8
	TemplateData [BS2_TEMPLATE_PER_FACE][BS2_FACE_TEMPLATE_LENGTH]uint8
}

func ScanFace(context SdkContext, deviceId BS2_DEVICE_ID, face *BS2Face, enrollmentThreshold uint8, fnOnReadyToScan interface{}) int16 {
	ret, _, err := syscall.Syscall6(procBs2ScanFace, uintptr(5), context, uintptr(deviceId), uintptr(unsafe.Pointer(face)), uintptr(enrollmentThreshold), uintptr(syscall.NewCallback(fnOnReadyToScan)), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

type BS2_AUTH_GROUP_ID = uint32

const (
	BS2_MAX_AUTH_GROUP_NAME_LEN = 144
)

type BS2AuthGroup struct {
	Id       BS2_AUTH_GROUP_ID
	Name     [BS2_MAX_AUTH_GROUP_NAME_LEN]byte
	Reserved [32]uint8
}

func GetAuthGroup(context SdkContext, deviceId BS2_DEVICE_ID, authGroupdIds []uint32, authGroupObj *[]BS2AuthGroup) int16 {
	var authGroupBuffer *BS2AuthGroup
	var numAuthGroup uint32

	var authGroupdIdsPtr uintptr
	if 0 != len(authGroupdIds) {
		authGroupdIdsPtr = uintptr(unsafe.Pointer(&authGroupdIds[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2GetAuthGroup, uintptr(6), context, uintptr(deviceId), authGroupdIdsPtr, uintptr(len(authGroupdIds)), uintptr(unsafe.Pointer(&authGroupBuffer)), uintptr(numAuthGroup))
	defer ReleaseObject(uintptr(unsafe.Pointer(authGroupBuffer)))

	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}

	temp := make([]BS2AuthGroup, numAuthGroup)

	var i uintptr = 0
	for ; i < (uintptr)(numAuthGroup); i++ {
		temp[i] = *(*BS2AuthGroup)(unsafe.Pointer(uintptr(unsafe.Pointer(authGroupBuffer)) + i*unsafe.Sizeof(*authGroupBuffer)))
	}
	*authGroupObj = temp

	return int16(ret)
}
func GetAllAuthGroup(context SdkContext, deviceId BS2_DEVICE_ID, authGroupObj *[]BS2AuthGroup) int16 {
	var authGroupBuffer *BS2AuthGroup
	var numAuthGroup uint32
	ret, _, err := syscall.Syscall6(procBs2GetAllAuthGroup, uintptr(4), context, uintptr(deviceId), uintptr(unsafe.Pointer(&authGroupBuffer)), uintptr(numAuthGroup), uintptr(0), uintptr(0))
	defer ReleaseObject(uintptr(unsafe.Pointer(authGroupBuffer)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}

	temp := make([]BS2AuthGroup, numAuthGroup)

	var i uintptr = 0
	for ; i < (uintptr)(numAuthGroup); i++ {
		temp[i] = *(*BS2AuthGroup)(unsafe.Pointer(uintptr(unsafe.Pointer(authGroupBuffer)) + i*unsafe.Sizeof(*authGroupBuffer)))
	}
	*authGroupObj = temp

	return int16(ret)
}
func SetAuthGroup(context SdkContext, deviceId BS2_DEVICE_ID, authGroups *[]BS2AuthGroup) int16 {
	ret, _, err := syscall.Syscall6(procBs2SetAuthGroup, uintptr(4), context, uintptr(deviceId), uintptr(unsafe.Pointer(authGroups)), uintptr(len(*authGroups)), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveAuthGroup(context SdkContext, deviceId BS2_DEVICE_ID, authGroupIds []uint32) int16 {
	var authGroupIdsPtr uintptr
	if 0 != len(authGroupIds) {
		authGroupIdsPtr = uintptr(unsafe.Pointer(&authGroupIds[0]))
	}

	ret, _, err := syscall.Syscall6(procBs2RemoveAuthGroup, uintptr(4), context, uintptr(deviceId), authGroupIdsPtr, uintptr(len(authGroupIds)), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveAllAuthGroup(context SdkContext, deviceId BS2_DEVICE_ID) int16 {
	ret, _, err := syscall.Syscall(procBs2RemoveAllAuthGroup, uintptr(2), context, uintptr(deviceId), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

const (
	BS2_MAX_ACCESS_LEVEL_PER_ACCESS_GROUP = 128 // #accessLevels + #floorLevels
	BS2_MAX_ACCESS_GROUP_PER_USER         = 16
	BS2_MAX_ACCESS_GROUP_NAME_LEN         = 48 * 3

	BS2_INVALID_ACCESS_GROUP_ID = 0
)

type BS2AccessGroupUsers struct {
	AccessGroupID uint32
	NumUsers      uint32
	UserID        []byte
}

type BS2UserAccessGroups struct {
	NumAccessGroups uint8
	Reserved        [3]byte
	AccessGroupID   [BS2_MAX_ACCESS_GROUP_PER_USER]uint32
}

type BS2AccessGroup struct {
	Id              uint32
	Name            [BS2_MAX_ACCESS_GROUP_NAME_LEN]byte
	NumAccessLevels uint8
	Reserved        [3]byte
	AccessLevels    [BS2_MAX_ACCESS_LEVEL_PER_ACCESS_GROUP]uint32
}

func GetAccessGroup(context SdkContext, deviceId BS2_DEVICE_ID, accessGroupIds []uint32, accessGroupObj *[]BS2AccessGroup) int16 {
	var accessGroupBuffer *BS2AccessGroup
	var numAccessGroup uint32

	var accessGroupIdsPtr uintptr
	if 0 != len(accessGroupIds) {
		accessGroupIdsPtr = uintptr(unsafe.Pointer(&accessGroupIds[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2GetAccessGroup, uintptr(6), context, uintptr(deviceId), accessGroupIdsPtr, uintptr(len(accessGroupIds)), uintptr(unsafe.Pointer(&accessGroupBuffer)), uintptr(unsafe.Pointer(&numAccessGroup)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}

	temp := make([]BS2AccessGroup, numAccessGroup)

	var i uintptr = 0
	for ; i < (uintptr)(numAccessGroup); i++ {
		temp[i] = *(*BS2AccessGroup)(unsafe.Pointer(uintptr(unsafe.Pointer(accessGroupBuffer)) + i*unsafe.Sizeof(*accessGroupBuffer)))
	}
	*accessGroupObj = temp
	ReleaseObject(uintptr(unsafe.Pointer(accessGroupBuffer)))

	return int16(ret)
}
func GetAllAccessGroup(context SdkContext, deviceId BS2_DEVICE_ID, accessGroupObj *[]BS2AccessGroup) int16 {
	var accessGroupBuffer *BS2AccessGroup
	var numAccessGroup uint32
	ret, _, err := syscall.Syscall6(procBs2GetAllAccessGroup, uintptr(4), context, uintptr(deviceId), uintptr(unsafe.Pointer(&accessGroupBuffer)), uintptr(unsafe.Pointer(&numAccessGroup)), uintptr(0), uintptr(0))
	defer ReleaseObject(uintptr(unsafe.Pointer(accessGroupBuffer)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2AccessGroup, numAccessGroup)

	var i uintptr = 0
	for ; i < (uintptr)(numAccessGroup); i++ {
		temp[i] = *(*BS2AccessGroup)(unsafe.Pointer(uintptr(unsafe.Pointer(accessGroupBuffer)) + i*unsafe.Sizeof(*accessGroupBuffer)))
	}
	*accessGroupObj = temp
	return int16(ret)
}
func SetAccessGroup(context SdkContext, deviceId BS2_DEVICE_ID, accessGroups *[]BS2AccessGroup) int16 {
	var accessGroupsPtr uintptr
	if 0 != len(*accessGroups) {
		accessGroupsPtr = uintptr(unsafe.Pointer(&(*accessGroups)[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2SetAccessGroup, uintptr(4), context, uintptr(deviceId), accessGroupsPtr, uintptr(len(*accessGroups)), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveAccessGroup(context SdkContext, deviceId BS2_DEVICE_ID, accessGroupIds []uint32) int16 {
	var accessGroupIdsPtr uintptr
	if 0 != len(accessGroupIds) {
		accessGroupIdsPtr = uintptr(unsafe.Pointer(&accessGroupIds[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2RemoveAccessGroup, uintptr(4), context, uintptr(deviceId), accessGroupIdsPtr, uintptr(len(accessGroupIds)), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveAllAccessGroup(context SdkContext, deviceId BS2_DEVICE_ID) int16 {
	ret, _, err := syscall.Syscall(procBs2RemoveAllAccessGroup, uintptr(2), context, uintptr(deviceId), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

const (
	BS2_MAX_ACCESS_LEVEL_ITEMS    = 128
	BS2_MAX_ACCESS_LEVEL_NAME_LEN = 48 * 3

	BS2_INVALID_ACCESS_LEVEL_ID = 0
)

type BS2DoorSchedule struct {
	DoorID     uint32
	ScheduleID uint32
}

type BS2AccessLevel struct {
	Id           uint32
	Name         [BS2_MAX_ACCESS_LEVEL_NAME_LEN]byte
	Reserved     [3]byte
	DooSchedules [BS2_MAX_ACCESS_LEVEL_ITEMS]BS2DoorSchedule
}

func GetAccessLevel(context SdkContext, deviceId BS2_DEVICE_ID, accessLevelIds []uint32, accessLevelObj *[]BS2AccessLevel) int16 {
	var accessLevelBuffer *BS2AccessLevel
	var numAccessLevel uint32

	var accessLevelIdsPtr uintptr
	if 0 != len(accessLevelIds) {
		accessLevelIdsPtr = uintptr(unsafe.Pointer(&accessLevelIds[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2GetAccessLevel, uintptr(6), context, uintptr(deviceId), accessLevelIdsPtr, uintptr(len(accessLevelIds)), uintptr(unsafe.Pointer(&accessLevelBuffer)), uintptr(unsafe.Pointer(&numAccessLevel)))
	defer ReleaseObject(uintptr(unsafe.Pointer(accessLevelBuffer)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2AccessLevel, numAccessLevel)

	var i uintptr = 0
	for ; i < (uintptr)(numAccessLevel); i++ {
		temp[i] = *(*BS2AccessLevel)(unsafe.Pointer(uintptr(unsafe.Pointer(accessLevelBuffer)) + i*unsafe.Sizeof(*accessLevelBuffer)))
	}
	*accessLevelObj = temp
	return int16(ret)
}
func GetAllAccessLevel(context SdkContext, deviceId BS2_DEVICE_ID, accessLevelObj *[]BS2AccessLevel) int16 {
	var accessLevelBuffer *BS2AccessLevel
	var numAccessLevel uint32
	ret, _, err := syscall.Syscall6(procBs2GetAllAccessLevel, uintptr(4), context, uintptr(deviceId), uintptr(unsafe.Pointer(&accessLevelBuffer)), uintptr(unsafe.Pointer(&numAccessLevel)), uintptr(0), uintptr(0))
	defer ReleaseObject(uintptr(unsafe.Pointer(accessLevelBuffer)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2AccessLevel, numAccessLevel)

	var i uintptr = 0
	for ; i < (uintptr)(numAccessLevel); i++ {
		temp[i] = *(*BS2AccessLevel)(unsafe.Pointer(uintptr(unsafe.Pointer(accessLevelBuffer)) + i*unsafe.Sizeof(*accessLevelBuffer)))
	}
	*accessLevelObj = temp
	return int16(ret)
}
func SetAccessLevel(context SdkContext, deviceId BS2_DEVICE_ID, accessLevels *[]BS2AccessLevel) int16 {
	var accessLevelsPtr uintptr
	if 0 != len(*accessLevels) {
		accessLevelsPtr = uintptr(unsafe.Pointer(&(*accessLevels)[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2SetAccessLevel, uintptr(4), context, uintptr(deviceId), accessLevelsPtr, uintptr(len(*accessLevels)), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveAccessLevel(context SdkContext, deviceId BS2_DEVICE_ID, accessLevelIds []uint32) int16 {
	var accessLevelIdsPtr uintptr
	if 0 != len(accessLevelIds) {
		accessLevelIdsPtr = uintptr(unsafe.Pointer(&accessLevelIds[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2RemoveAccessLevel, uintptr(4), context, uintptr(deviceId), accessLevelIdsPtr, uintptr(len(accessLevelIds)), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveAllAccessLevel(context SdkContext, deviceId BS2_DEVICE_ID) int16 {
	ret, _, err := syscall.Syscall(procBs2RemoveAllAccessLevel, uintptr(2), context, uintptr(deviceId), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

const (
	BS2_MAX_HOLIDAY_GROUPS_PER_SCHEDULE = 4
	BS2_MAX_DAYS_PER_DAILY_SCHEDULE     = 90

	BS2_MAX_SCHEDULE_NAME_LEN = 48 * 3

	BS2_INVALID_SCHEDULE_ID = 0

	BS2_SCHEDULE_NEVER_ID  = 0
	BS2_SCHEDULE_ALWAYS_ID = 1
)
const (
	BS2_MAX_TIME_PERIODS_PER_DAY = 5
)

type BS2TimePeriod struct {
	StartTime int16
	EndTime   int16
}

type BS2DaySchedule struct {
	NumPeriods uint8
	Reserved   byte
	Periods    [BS2_MAX_TIME_PERIODS_PER_DAY]BS2TimePeriod
}

type BS2HolidaySchedule struct {
	Id        uint32
	Schedules BS2DaySchedule
}

type BS2Schedule struct {
	Id   uint32
	Name [BS2_MAX_SCHEDULE_NAME_LEN]byte

	IsDaily             bool
	NumHolidaySchedules uint8
	Reserved            [2]byte

	//sizeof(BS2DaySchedule) is 24.
	//sizeof(BS2Weeklyschedule) is 168.
	//sizeof(BS2DailySchedule) is 90*24+8
	ScheduleUnionData [90*24 + 8]byte

	HolidaySchedules [BS2_MAX_HOLIDAY_GROUPS_PER_SCHEDULE]BS2HolidaySchedule
}

func GetAccessSchedule(context SdkContext, deviceId BS2_DEVICE_ID, accessSheduleIds []uint32, accessScheduleObj *[]BS2Schedule) int16 {
	var scheduleBuffer *BS2Schedule
	var numSchedule uint32

	var accessSheduleIdsPtr uintptr
	if 0 != len(accessSheduleIds) {
		accessSheduleIdsPtr = uintptr(unsafe.Pointer(&accessSheduleIds[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2GetAccessSchedule, uintptr(6), context, uintptr(deviceId), accessSheduleIdsPtr, uintptr(len(accessSheduleIds)), uintptr(unsafe.Pointer(&scheduleBuffer)), uintptr(unsafe.Pointer(&numSchedule)))
	defer ReleaseObject(uintptr(unsafe.Pointer(scheduleBuffer)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2Schedule, numSchedule)

	var i uintptr = 0
	for ; i < (uintptr)(numSchedule); i++ {
		temp[i] = *(*BS2Schedule)(unsafe.Pointer(uintptr(unsafe.Pointer(scheduleBuffer)) + i*unsafe.Sizeof(*scheduleBuffer)))
	}
	*accessScheduleObj = temp
	return int16(ret)
}
func GetAllAccessSchedule(context SdkContext, deviceId BS2_DEVICE_ID, accessScheduleObj *[]BS2Schedule) int16 {
	var scheduleBuffer *BS2Schedule
	var numSchedule uint32
	ret, _, err := syscall.Syscall6(procBs2GetAllAccessSchedule, uintptr(4), context, uintptr(deviceId), uintptr(unsafe.Pointer(&scheduleBuffer)), uintptr(unsafe.Pointer(&numSchedule)), uintptr(0), uintptr(0))
	defer ReleaseObject(uintptr(unsafe.Pointer(scheduleBuffer)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2Schedule, numSchedule)

	var i uintptr = 0
	for ; i < (uintptr)(numSchedule); i++ {
		temp[i] = *(*BS2Schedule)(unsafe.Pointer(uintptr(unsafe.Pointer(scheduleBuffer)) + i*unsafe.Sizeof(*scheduleBuffer)))
	}
	*accessScheduleObj = temp
	return int16(ret)
}
func SetAccessSchedule(context SdkContext, deviceId BS2_DEVICE_ID, accessSchedules *[]BS2Schedule) int16 {
	var accessSchedulesPtr uintptr
	if 0 != len(*accessSchedules) {
		accessSchedulesPtr = uintptr(unsafe.Pointer(&(*accessSchedules)[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2SetAccessSchedule, uintptr(4), context, uintptr(deviceId), accessSchedulesPtr, uintptr(len(*accessSchedules)), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveAccessSchedule(context SdkContext, deviceId BS2_DEVICE_ID, accessSheduleIds []uint32) int16 {
	var accessSheduleIdsPtr uintptr
	if 0 != len(accessSheduleIds) {
		accessSheduleIdsPtr = uintptr(unsafe.Pointer(&accessSheduleIds[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2RemoveAccessSchedule, uintptr(4), context, uintptr(deviceId), accessSheduleIdsPtr, uintptr(len(accessSheduleIds)), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveAllAccessSchedule(context SdkContext, deviceId BS2_DEVICE_ID) int16 {
	ret, _, err := syscall.Syscall(procBs2RemoveAllAccessSchedule, uintptr(2), context, uintptr(deviceId), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

const (
	BS2_MAX_HOLIDAYS_PER_GROUP = 128
)

type BS2Holiday struct {
	Date       uint32
	Recurrence uint8
}

type BS2HolidayGroup struct {
	Id          uint32
	Name        [BS2_MAX_SCHEDULE_NAME_LEN]byte
	NumHolidays uint8
	Reserved    [3]byte
	Holidays    [BS2_MAX_HOLIDAYS_PER_GROUP]BS2Holiday
}

func GetHolidayGroup(context SdkContext, deviceId BS2_DEVICE_ID, holidayGroupIds []uint32, holidayGroupObj *[]BS2HolidayGroup) int16 {
	var holidayGroupBuffer *BS2HolidayGroup
	var numHolidayGroup uint32

	var holidayGroupIdsPtr uintptr
	if 0 != len(holidayGroupIds) {
		holidayGroupIdsPtr = uintptr(unsafe.Pointer(&holidayGroupIds[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2GetHolidayGroup, uintptr(6), context, uintptr(deviceId), holidayGroupIdsPtr, uintptr(len(holidayGroupIds)), uintptr(unsafe.Pointer(&holidayGroupBuffer)), uintptr(unsafe.Pointer(&numHolidayGroup)))
	defer ReleaseObject(uintptr(unsafe.Pointer(holidayGroupBuffer)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2HolidayGroup, numHolidayGroup)

	var i uintptr = 0
	for ; i < (uintptr)(numHolidayGroup); i++ {
		temp[i] = *(*BS2HolidayGroup)(unsafe.Pointer(uintptr(unsafe.Pointer(holidayGroupBuffer)) + i*unsafe.Sizeof(*holidayGroupBuffer)))
	}
	*holidayGroupObj = temp
	return int16(ret)
}
func GetAllHolidayGroup(context SdkContext, deviceId BS2_DEVICE_ID, holidayGroupObj *[]BS2HolidayGroup) int16 {
	var holidayGroupBuffer *BS2HolidayGroup
	var numHolidayGroup uint32
	ret, _, err := syscall.Syscall6(procBs2GetAllHolidayGroup, uintptr(4), context, uintptr(deviceId), uintptr(unsafe.Pointer(&holidayGroupBuffer)), uintptr(unsafe.Pointer(&numHolidayGroup)), uintptr(0), uintptr(0))
	defer ReleaseObject(uintptr(unsafe.Pointer(holidayGroupBuffer)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2HolidayGroup, numHolidayGroup)

	var i uintptr = 0
	for ; i < (uintptr)(numHolidayGroup); i++ {
		temp[i] = *(*BS2HolidayGroup)(unsafe.Pointer(uintptr(unsafe.Pointer(holidayGroupBuffer)) + i*unsafe.Sizeof(*holidayGroupBuffer)))
	}
	*holidayGroupObj = temp
	return int16(ret)
}
func SetHolidayGroup(context SdkContext, deviceId BS2_DEVICE_ID, holidayGroupObj *[]BS2HolidayGroup) int16 {
	var holidayGroupObjPtr uintptr
	if 0 != len(*holidayGroupObj) {
		holidayGroupObjPtr = uintptr(unsafe.Pointer(&(*holidayGroupObj)[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2SetHolidayGroup, uintptr(4), context, uintptr(deviceId), holidayGroupObjPtr, uintptr(len(*holidayGroupObj)), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveHolidayGroup(context SdkContext, deviceId BS2_DEVICE_ID, holidayGroupIds []uint32) int16 {
	var holidayGroupIdsPtr uintptr
	if 0 != len(holidayGroupIds) {
		holidayGroupIdsPtr = uintptr(unsafe.Pointer(&holidayGroupIds[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2RemoveHolidayGroup, uintptr(4), context, uintptr(deviceId), holidayGroupIdsPtr, uintptr(len(holidayGroupIds)), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveAllHolidayGroup(context SdkContext, deviceId BS2_DEVICE_ID) int16 {
	ret, _, err := syscall.Syscall(procBs2RemoveAllHolidayGroup, uintptr(2), context, uintptr(deviceId), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

type BS2BlackList struct {
	CardID     [BS2_CARD_DATA_SIZE]byte
	IssueCount uint16
}

func GetBlackList(context SdkContext, deviceId BS2_DEVICE_ID, blackLists []BS2BlackList, blacklistObj *[]BS2BlackList) int16 {
	var blacklistBuffer *BS2HolidayGroup
	var numBlacklist uint32
	var blackListsPtr uintptr
	if 0 != len(blackLists) {
		blackListsPtr = uintptr(unsafe.Pointer(&blackLists[0]))
	}

	ret, _, err := syscall.Syscall6(procBs2GetBlackList, uintptr(6), context, uintptr(deviceId), blackListsPtr, uintptr(len(blackLists)), uintptr(unsafe.Pointer(&blacklistBuffer)), uintptr(unsafe.Pointer(&numBlacklist)))
	defer ReleaseObject(uintptr(unsafe.Pointer(blacklistBuffer)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2BlackList, numBlacklist)

	var i uintptr = 0
	for ; i < (uintptr)(numBlacklist); i++ {
		temp[i] = *(*BS2BlackList)(unsafe.Pointer(uintptr(unsafe.Pointer(blacklistBuffer)) + i*unsafe.Sizeof(*blacklistBuffer)))
	}
	*blacklistObj = temp
	return int16(ret)
}
func GetAllBlackList(context SdkContext, deviceId BS2_DEVICE_ID, blacklistObj *[]BS2BlackList) int16 {
	var blacklistBuffer *BS2HolidayGroup
	var numBlacklist uint32
	ret, _, err := syscall.Syscall6(procBs2GetAllBlackList, uintptr(4), context, uintptr(deviceId), uintptr(unsafe.Pointer(&blacklistBuffer)), uintptr(unsafe.Pointer(&numBlacklist)), uintptr(0), uintptr(0))
	defer ReleaseObject(uintptr(unsafe.Pointer(blacklistBuffer)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2BlackList, numBlacklist)

	var i uintptr = 0
	for ; i < (uintptr)(numBlacklist); i++ {
		temp[i] = *(*BS2BlackList)(unsafe.Pointer(uintptr(unsafe.Pointer(blacklistBuffer)) + i*unsafe.Sizeof(*blacklistBuffer)))
	}
	*blacklistObj = temp
	return int16(ret)
}
func SetBlackList(context SdkContext, deviceId BS2_DEVICE_ID, blacklists *[]BS2BlackList) int16 {
	var blacklistsPtr uintptr
	if 0 != len(*blacklists) {
		blacklistsPtr = uintptr(unsafe.Pointer(&(*blacklists)[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2SetBlackList, uintptr(4), context, uintptr(deviceId), blacklistsPtr, uintptr(len(*blacklists)), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveBlackList(context SdkContext, deviceId BS2_DEVICE_ID, blacklists []BS2BlackList) int16 {
	var blacklistsPtr uintptr
	if 0 != len(blacklists) {
		blacklistsPtr = uintptr(unsafe.Pointer(&blacklists[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2RemoveBlackList, uintptr(4), context, uintptr(deviceId), blacklistsPtr, uintptr(len(blacklists)), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveAllBlackList(context SdkContext, deviceId BS2_DEVICE_ID) int16 {
	ret, _, err := syscall.Syscall(procBs2RemoveAllBlackList, uintptr(2), context, uintptr(deviceId), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

const (
	BS2_MAX_DOOR_NAME_LEN = 48 * 3

	BS2_MAX_DUAL_AUTH_APPROVAL_GROUP = 16
	BS2_DEFAULT_AUTO_LOCK_TIMEOUT    = 3  ///< in seconds
	BS2_DEFAULT_HELD_OPEN_TIMEOUT    = 10 ///< in seconds
	BS2_DEFAULT_DUAL_AUTH_TIMEOUT    = 15 ///< in seconds
	BS2_INVALID_DOOR_ID              = 0

	BS2_MAX_HELD_OPEN_ALARM_ACTION   = 5
	BS2_MAX_FORCED_OPEN_ALARM_ACTION = 5
)

const (
	BS2_DUAL_AUTH_APPROVAL_NONE = 0
	BS2_DUAL_AUTH_APPROVAL_LAST = 1
)

const (
	BS2_DUAL_AUTH_NO_DEVICE         = 0
	BS2_DUAL_AUTH_ENTRY_DEVICE_ONLY = 1
	BS2_DUAL_AUTH_EXIT_DEVICE_ONLY  = 2
	BS2_DUAL_AUTH_BOTH_DEVICE       = 3
)
const (
	BS2_DOOR_FLAG_NONE      = 0x00
	BS2_DOOR_FLAG_SCHEDULE  = 0x01
	BS2_DOOR_FLAG_OPERATOR  = 0x04
	BS2_DOOR_FLAG_EMERGENCY = 0x02
)

/**
 *	BS2_DOOR_ALARM_FLAG
 */
const (
	BS2_DOOR_ALARM_FLAG_NONE        = 0x00
	BS2_DOOR_ALARM_FLAG_FORCED_OPEN = 0x01
	BS2_DOOR_ALARM_FLAG_HELD_OPEN   = 0x02
	BS2_DOOR_ALARM_FLAG_APB         = 0x04
)

type BS2Door struct {
	DoorID uint32
	Name   [BS2_MAX_DOOR_NAME_LEN]byte

	EntryDeviceID uint32
	ExitDeviceID  uint32

	Relay  BS2DoorRelay
	Sensor BS2DoorSensor
	Button BS2ExitButton

	AutoLockTimeout uint32
	HeldOpenTimeout uint32

	InstantLock       bool
	UnlockFlags       uint8
	LockFlags         uint8
	UnconditionalLock bool

	ForcedOpenAlarm [BS2_MAX_FORCED_OPEN_ALARM_ACTION]BS2Action
	HeldOpenAlarm   [BS2_MAX_HELD_OPEN_ALARM_ACTION]BS2Action

	DualAuthScheduleID        uint32
	DualAuthDevice            uint8
	DualAuthApprovalType      uint8
	DualAuthTimeout           uint32
	NumDualAuthApprovalGroups uint8
	Reserved2                 byte
	DualAuthApprovalGroupID   [BS2_MAX_DUAL_AUTH_APPROVAL_GROUP]uint32

	ApbZone BS2AntiPassbackZone
}
type BS2DoorRelay struct {
	DeviceID uint32
	Port     uint8
	Reserved [3]uint8
}
type BS2DoorSensor struct {
	DeviceID uint32

	Port       uint8
	SwitchType uint8
	Reserved   [2]byte
}

type BS2DoorStatus struct {
	Id           uint32
	Opened       uint8
	Unlocked     uint8
	HeldOpened   uint8
	UnlockFlags  uint8
	LockFlags    uint8
	AlarmFlags   uint8
	Reserved     [2]byte
	LastOpenTime uint32
}
type BS2ExitButton struct {
	DeviceID   uint32
	Port       uint8
	SwitchType uint8
	Reserved   [2]uint
}

const (
	BS2_MAX_READERS_PER_APB_ZONE       = 64
	BS2_MAX_BYPASS_GROUPS_PER_APB_ZONE = 16
	BS2_MAX_APB_ALARM_ACTION           = 5

	BS2_RESET_DURATION_DEFAULT = 86400
)

const (
	BS2_APB_ZONE_HARD = 0x00
	BS2_APB_ZONE_SOFT = 0x01
)

const (
	BS2_APB_ZONE_READER_NONE  = -1
	BS2_APB_ZONE_READER_ENTRY = 0
	BS2_APB_ZONE_READER_EXIT  = 1
)

type BS2ApbMember struct {
	DeviceID uint32
	Type     uint8
	Reserved [3]uint8
}

const (
	BS2_MAX_ZONE_NAME_LEN = 144
)

type BS2AntiPassbackZone struct {
	ZoneID uint32
	Name   [BS2_MAX_ZONE_NAME_LEN]byte

	Type            uint8
	NumReaders      uint8
	NumBypassGroups uint8
	Disabled        bool

	Alarmed  bool
	Reserved [3]byte

	ResetDuration uint32

	Alarm [BS2_MAX_APB_ALARM_ACTION]BS2Action

	Readers   [BS2_MAX_READERS_PER_APB_ZONE]BS2ApbMember
	Reserved2 [8 * 64]byte

	BypassGroupIDs [BS2_MAX_BYPASS_GROUPS_PER_APB_ZONE]uint32
}

func GetDoor(context SdkContext, deviceId BS2_DEVICE_ID, doorIds []uint32, doorObj *[]BS2Door) int16 {
	var doorBuffer *BS2Door
	var numDoor uint32

	var doorIdsPtr uintptr
	if 0 != len(doorIds) {
		doorIdsPtr = uintptr(unsafe.Pointer(&doorIds[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2GetDoor, uintptr(6), context, uintptr(deviceId), doorIdsPtr, uintptr(len(doorIds)), uintptr(unsafe.Pointer(&doorBuffer)), uintptr(unsafe.Pointer(&numDoor)))
	defer ReleaseObject(uintptr(unsafe.Pointer(doorBuffer)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2Door, numDoor)

	var i uintptr = 0
	for ; i < (uintptr)(numDoor); i++ {
		temp[i] = *(*BS2Door)(unsafe.Pointer(uintptr(unsafe.Pointer(doorBuffer)) + i*unsafe.Sizeof(*doorBuffer)))
	}
	*doorObj = temp
	return int16(ret)
}
func GetAllDoor(context SdkContext, deviceId BS2_DEVICE_ID, doorObj *[]BS2Door) int16 {
	var doorBuffer *BS2Door
	var numDoor uint32
	ret, _, err := syscall.Syscall6(procBs2GetAllDoor, uintptr(4), context, uintptr(deviceId), uintptr(unsafe.Pointer(&doorBuffer)), uintptr(unsafe.Pointer(&numDoor)), uintptr(0), uintptr(0))
	defer ReleaseObject(uintptr(unsafe.Pointer(doorBuffer)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2Door, numDoor)

	var i uintptr = 0
	for ; i < (uintptr)(numDoor); i++ {
		temp[i] = *(*BS2Door)(unsafe.Pointer(uintptr(unsafe.Pointer(doorBuffer)) + i*unsafe.Sizeof(*doorBuffer)))
	}
	*doorObj = temp
	return int16(ret)
}
func GetDoorStatus(context SdkContext, deviceId BS2_DEVICE_ID, doorIds []uint32, doorStatusObje *[]BS2DoorStatus) int16 {
	var doorBuffer *BS2DoorStatus
	var numDoor uint32

	var doorIdsPtr uintptr
	if 0 != len(doorIds) {
		doorIdsPtr = uintptr(unsafe.Pointer(&doorIds[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2GetDoorStatus, uintptr(6), context, uintptr(deviceId), doorIdsPtr, uintptr(len(doorIds)), uintptr(unsafe.Pointer(&doorBuffer)), uintptr(unsafe.Pointer(&numDoor)))
	defer ReleaseObject(uintptr(unsafe.Pointer(doorBuffer)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2DoorStatus, numDoor)

	var i uintptr = 0
	for ; i < (uintptr)(numDoor); i++ {
		temp[i] = *(*BS2DoorStatus)(unsafe.Pointer(uintptr(unsafe.Pointer(doorBuffer)) + i*unsafe.Sizeof(*doorBuffer)))
	}
	*doorStatusObje = temp
	return int16(ret)
}
func GetAllDoorStatus(context SdkContext, deviceId BS2_DEVICE_ID, doorStatusObje *[]BS2DoorStatus) int16 {
	var doorBuffer *BS2DoorStatus
	var numDoor uint32
	ret, _, err := syscall.Syscall6(procBs2GetAllDoorStatus, uintptr(4), context, uintptr(deviceId), uintptr(unsafe.Pointer(&doorBuffer)), uintptr(unsafe.Pointer(&numDoor)), uintptr(0), uintptr(0))
	defer ReleaseObject(uintptr(unsafe.Pointer(doorBuffer)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2DoorStatus, numDoor)

	var i uintptr = 0
	for ; i < (uintptr)(numDoor); i++ {
		temp[i] = *(*BS2DoorStatus)(unsafe.Pointer(uintptr(unsafe.Pointer(doorBuffer)) + i*unsafe.Sizeof(*doorBuffer)))
	}
	*doorStatusObje = temp
	return int16(ret)
}
func SetDoor(context SdkContext, deviceId BS2_DEVICE_ID, doors *[]BS2Door) int16 {

	var doorsPtr uintptr
	if 0 != len(*doors) {
		doorsPtr = uintptr(unsafe.Pointer(&(*doors)[0]))
	}

	ret, _, err := syscall.Syscall6(procBs2SetDoor, uintptr(4), context, uintptr(deviceId), doorsPtr, uintptr(len(*doors)), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func SetDoorAlarm(context SdkContext, deviceId BS2_DEVICE_ID, flag uint8, doorIds []uint32) int16 {
	var doorIdsPtr uintptr
	if 0 != len(doorIds) {
		doorIdsPtr = uintptr(unsafe.Pointer(&doorIds[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2SetDoorAlarm, uintptr(5), context, uintptr(deviceId), uintptr(flag), doorIdsPtr, uintptr(len(doorIds)), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveDoor(context SdkContext, deviceId BS2_DEVICE_ID, doorIds []uint32) int16 {
	var doorIdsPtr uintptr
	if 0 != len(doorIds) {
		doorIdsPtr = uintptr(unsafe.Pointer(&doorIds[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2RemoveDoor, uintptr(4), context, uintptr(deviceId), doorIdsPtr, uintptr(len(doorIds)), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveAllDoor(context SdkContext, deviceId BS2_DEVICE_ID) int16 {
	ret, _, err := syscall.Syscall(procBs2RemoveAllDoor, uintptr(2), context, uintptr(deviceId), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func ReleaseDoor(context SdkContext, deviceId BS2_DEVICE_ID, flag uint8, doorIds []uint32) int16 {
	var doorIdsPtr uintptr
	if 0 != len(doorIds) {
		doorIdsPtr = uintptr(unsafe.Pointer(&doorIds[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2ReleaseDoor, uintptr(5), context, uintptr(deviceId), uintptr(flag), doorIdsPtr, uintptr(len(doorIds)), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func LockDoor(context SdkContext, deviceId BS2_DEVICE_ID, flag uint8, doorIds []uint32) int16 {
	var doorIdsPtr uintptr
	if 0 != len(doorIds) {
		doorIdsPtr = uintptr(unsafe.Pointer(&doorIds[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2LockDoor, uintptr(5), context, uintptr(deviceId), uintptr(flag), doorIdsPtr, uintptr(len(doorIds)), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func UnlockDoor(context SdkContext, deviceId BS2_DEVICE_ID, flag uint8, doorIds []uint32) int16 {
	var doorIdsPtr uintptr
	if 0 != len(doorIds) {
		doorIdsPtr = uintptr(unsafe.Pointer(&doorIds[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2UnlockDoor, uintptr(5), context, uintptr(deviceId), uintptr(flag), doorIdsPtr, uintptr(len(doorIds)), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func GetSlaveDevice(context SdkContext, deviceId BS2_DEVICE_ID, slaveDeviceObj *[]BS2Rs485SlaveDevice) int16 {
	var slaveDeviceBuffer *BS2Rs485SlaveDevice
	var numSlaveDevice uint32
	ret, _, err := syscall.Syscall6(procBs2GetSlaveDevice, uintptr(4), context, uintptr(deviceId), uintptr(unsafe.Pointer(&slaveDeviceBuffer)), uintptr(unsafe.Pointer(&numSlaveDevice)), uintptr(0), uintptr(0))
	defer ReleaseObject(uintptr(unsafe.Pointer(slaveDeviceBuffer)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2Rs485SlaveDevice, numSlaveDevice)

	var i uintptr = 0
	for ; i < (uintptr)(numSlaveDevice); i++ {
		temp[i] = *(*BS2Rs485SlaveDevice)(unsafe.Pointer(uintptr(unsafe.Pointer(slaveDeviceBuffer)) + i*unsafe.Sizeof(*slaveDeviceBuffer)))
	}
	*slaveDeviceObj = temp
	return int16(ret)
}
func SetSlaveDevice(context SdkContext, deviceId BS2_DEVICE_ID, slaveDeviceObj *[]BS2Rs485SlaveDevice) int16 {

	var slaveDeviceObjPtr uintptr
	if 0 != len(*slaveDeviceObj) {
		slaveDeviceObjPtr = uintptr(unsafe.Pointer(&(*slaveDeviceObj)[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2SetSlaveDevice, uintptr(4), context, uintptr(deviceId), slaveDeviceObjPtr, uintptr(len(*slaveDeviceObj)), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func GetSlaveExDevice(context SdkContext, deviceId BS2_DEVICE_ID, channelPort uint32, slaveDevices *[]BS2Rs485SlaveDeviceEX, outchannelPort *uint32) int16 {
	var slaveDeviceBuffer *BS2Rs485SlaveDeviceEX
	var numSlaveDevice uint32
	ret, _, err := syscall.Syscall6(procBs2GetSlaveExDevice, uintptr(6), context, uintptr(deviceId), uintptr(channelPort), uintptr(unsafe.Pointer(&slaveDeviceBuffer)), uintptr(unsafe.Pointer(outchannelPort)), uintptr(unsafe.Pointer(&numSlaveDevice)))
	defer ReleaseObject(uintptr(unsafe.Pointer(slaveDeviceBuffer)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2Rs485SlaveDeviceEX, numSlaveDevice)

	var i uintptr = 0
	for ; i < (uintptr)(numSlaveDevice); i++ {
		temp[i] = *(*BS2Rs485SlaveDeviceEX)(unsafe.Pointer(uintptr(unsafe.Pointer(slaveDeviceBuffer)) + i*unsafe.Sizeof(*slaveDeviceBuffer)))
	}
	*slaveDevices = temp
	return int16(ret)
}
func SetSlaveExDevice(context SdkContext, deviceId BS2_DEVICE_ID, channelPort uint32, slaveDevices *[]BS2Rs485SlaveDeviceEX) int16 {
	var slaveDevicesPtr uintptr
	if 0 != len(*slaveDevices) {
		slaveDevicesPtr = uintptr(unsafe.Pointer(&(*slaveDevices)[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2SetSlaveExDevice, uintptr(5), context, uintptr(deviceId), uintptr(channelPort), slaveDevicesPtr, uintptr(len(*slaveDevices)), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func SearchWiegandDevices(context SdkContext, deviceId BS2_DEVICE_ID, wiegandDeviceObj *[]uint32) int16 {
	var buffer *uint32
	var numWiegandDevice uint32
	ret, _, err := syscall.Syscall6(procBs2SearchWiegandDevices, uintptr(4), context, uintptr(deviceId),
		uintptr(unsafe.Pointer(&buffer)), uintptr(unsafe.Pointer(&numWiegandDevice)), uintptr(0), uintptr(0))
	defer ReleaseObject(uintptr(unsafe.Pointer(buffer)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]uint32, numWiegandDevice)

	var i uintptr = 0
	for ; i < (uintptr)(numWiegandDevice); i++ {
		temp[i] = *(*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(buffer)) + i*unsafe.Sizeof(*buffer)))
	}
	*wiegandDeviceObj = temp
	return int16(ret)
}
func GetWiegandDevices(context SdkContext, deviceId BS2_DEVICE_ID, wiegandDeviceObj *[]uint32) int16 {
	var buffer *uint32
	var numWiegandDevice uint32
	ret, _, err := syscall.Syscall6(procBs2GetWiegandDevices, uintptr(4), context, uintptr(deviceId), uintptr(unsafe.Pointer(&buffer)), uintptr(unsafe.Pointer(&numWiegandDevice)), uintptr(0), uintptr(0))
	defer ReleaseObject(uintptr(unsafe.Pointer(buffer)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]uint32, numWiegandDevice)

	var i uintptr = 0
	for ; i < (uintptr)(numWiegandDevice); i++ {
		temp[i] = *(*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(buffer)) + i*unsafe.Sizeof(*buffer)))
	}
	*wiegandDeviceObj = temp
	return int16(ret)
}
func AddWiegandDevices(context SdkContext, deviceId BS2_DEVICE_ID, wiegandDeviceObj *[]uint32) int16 {
	log.Printf("%v, %v, %v\n", context, deviceId, *wiegandDeviceObj)
	var wiegandDeviceObjPtr uintptr
	if 0 != len(*wiegandDeviceObj) {
		wiegandDeviceObjPtr = uintptr(unsafe.Pointer(&(*wiegandDeviceObj)[0]))
	}
	var numWiegandDevice uint32 = (uint32)(len(*wiegandDeviceObj))
	ret, _, err := syscall.Syscall6(procBs2AddWiegandDevices, uintptr(4), context, uintptr(deviceId), wiegandDeviceObjPtr, uintptr(numWiegandDevice), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveWiegandDevices(context SdkContext, deviceId BS2_DEVICE_ID, wiegandDevice []uint32) int16 {
	var wiegandDevicePtr uintptr
	if 0 != len(wiegandDevice) {
		wiegandDevicePtr = uintptr(unsafe.Pointer(&wiegandDevice[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2RemoveWiegandDevices, uintptr(4), context, uintptr(deviceId), wiegandDevicePtr, uintptr(len(wiegandDevice)), uintptr(0), uintptr(0))

	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func SetServerMatchingHandler(context SdkContext, fnOnVerifyUser interface{}, fnOnIdentifyUser interface{}) int16 {

	var fnOnVerifyUserPtr, fnOnIdentifyUserPtr uintptr
	if nil != fnOnVerifyUser {
		fnOnVerifyUserPtr = uintptr(syscall.NewCallback(fnOnVerifyUser))
	}
	if nil != fnOnIdentifyUser {
		fnOnIdentifyUserPtr = uintptr(syscall.NewCallback(fnOnIdentifyUser))
	}
	ret, _, err := syscall.Syscall(procBs2SetServerMatchingHandler, uintptr(3), context, fnOnVerifyUserPtr, fnOnIdentifyUserPtr)
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func VerifyUser(context SdkContext, deviceId BS2_DEVICE_ID, seq uint16, handleResult int, userBlob *BS2UserBlob) int16 {
	ret, _, err := syscall.Syscall6(procBs2VerifyUser, uintptr(5), context, uintptr(deviceId), uintptr(seq), uintptr(handleResult), uintptr(unsafe.Pointer(userBlob)), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func IdentifyUser(context SdkContext, deviceId BS2_DEVICE_ID, seq uint16, handleResult int, userBlob *BS2UserBlob) int16 {
	ret, _, err := syscall.Syscall6(procBs2IdentifyUser, uintptr(5), context, uintptr(deviceId), uintptr(seq), uintptr(handleResult), uintptr(unsafe.Pointer(userBlob)), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func VerifyUserEx(context SdkContext, deviceId BS2_DEVICE_ID, seq uint16, handleResult int, userBlob *BS2UserBlobEx) int16 {
	ret, _, err := syscall.Syscall6(procBs2VerifyUserEx, uintptr(5), context, uintptr(deviceId), uintptr(seq), uintptr(handleResult), uintptr(unsafe.Pointer(userBlob)), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func IdentifyUserEx(context SdkContext, deviceId BS2_DEVICE_ID, seq uint16, handleResult int, userBlob *BS2UserBlobEx) int16 {
	ret, _, err := syscall.Syscall6(procBs2IdentifyUserEx, uintptr(5), context, uintptr(deviceId), uintptr(seq), uintptr(handleResult), uintptr(unsafe.Pointer(userBlob)), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

type BS2_CONFIG_MASK = uint32

func GetSupportedConfigMask(context SdkContext, deviceId BS2_DEVICE_ID, configMask *BS2_CONFIG_MASK) int16 {
	ret, _, err := syscall.Syscall(procBs2GetSupportedConfigMask, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(configMask)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

type BS2_USER_MASK = uint16

const (
	BS2_USER_MASK_ID_ONLY      = 0      // fill only user id in BS2User
	BS2_USER_MASK_DATA         = 0x0001 // BS2User
	BS2_USER_MASK_SETTING      = 0x0002 // BS2UserSetting
	BS2_USER_MASK_NAME         = 0x0004 // BS2_USER_NAME
	BS2_USER_MASK_PHOTO        = 0x0008 // BS2UserPhoto
	BS2_USER_MASK_PIN          = 0x0010 // BS2_HASH256
	BS2_USER_MASK_CARD         = 0x0020 // BS2CSNCard
	BS2_USER_MASK_FINGER       = 0x0040 // BS2FingerTemplate
	BS2_USER_MASK_FACE         = 0x0080 // BS2FaceTemplate
	BS2_USER_MASK_ACCESS_GROUP = 0x0100 // BS2_ACCESS_GROUP_ID
	BS2_USER_MASK_JOB          = 0x0200 // BS2Job
	BS2_USER_MASK_ALL          = 0xFFFF // 4 bytes

)

func GetSupportedUserMask(context SdkContext, deviceId BS2_DEVICE_ID, configMask *BS2_USER_MASK) int16 {
	ret, _, err := syscall.Syscall(procBs2GetSupportedUserMask, uintptr(3), context, uintptr(deviceId), uintptr(unsafe.Pointer(configMask)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

type BS2_LIFT_ID = uint32

/**
 *  Lift Constants
 */
const (
	BS2_MAX_LIFT_NAME_LEN = 48 * 3

	BS2_MAX_DEVICES_ON_LIFT = 4
	BS2_MAX_FLOORS_ON_LIFT  = 255

	BS2_MAX_ALARMS_ON_LIFT                   = 2
	BS2_MAX_DUAL_AUTH_APPROVAL_GROUP_ON_LIFT = 16

	BS2_DEFAULT_ACTIVATE_TIMEOUT_ON_LIFT  = 10 ///< in seconds
	BS2_DEFAULT_DUAL_AUTH_TIMEOUT_ON_LIFT = 15 ///< in seconds
	BS2_INVALID_LIFT_ID                   = 0
)

/**
 *	BS2_DUAL_AUTH_APPROVAL_TYPE
 */
const (
	BS2_DUAL_AUTH_APPROVAL_NONE_ON_LIFT = 0
	BS2_DUAL_AUTH_APPROVAL_LAST_ON_LIFT = 1

//	BS2_DUAL_AUTH_APPROVAL_BOTH			= 2
)

type BS2_DUAL_AUTH_APPROVAL = uint8

/**
 *	BS2_FLOOR_FLAG
 */
const (
	BS2_FLOOR_FLAG_NONE      = 0x00
	BS2_FLOOR_FLAG_SCHEDULE  = 0x01
	BS2_FLOOR_FLAG_OPERATOR  = 0x04
	BS2_FLOOR_FLAG_EMERGENCY = 0x02
)

type BS2_FLOOR_FLAG = uint8

/**
 *  BS2FloorStatus
 */
type BS2FloorStatus struct {
	Activated       bool           ///< 1 byte
	ActivateFlags   BS2_FLOOR_FLAG ///< 1 byte
	DeactivateFlags BS2_FLOOR_FLAG ///< 1 byte
}

/**
 *	BS2LiftFloor
 */
type BS2LiftFloor struct {
	DeviceID BS2_DEVICE_ID  ///< 4 bytes
	Port     uint8          ///< 1 byte : 1 ~ 16
	Status   BS2FloorStatus ///< 3 bytes
}

type BS2_SWITCH_TYPE byte
type BS2_SCHEDULE_ID uint32

/**
 *	BS2LiftSensor
 */
type BS2LiftSensor struct {
	DeviceID   BS2_DEVICE_ID   ///< 4 bytes
	Port       uint8           ///< 1 byte
	SwitchType BS2_SWITCH_TYPE ///< 1 byte
	Duration   uint16          ///< 2 bytes
	ScheduleID BS2_SCHEDULE_ID ///< 4 bytes
}

type BS2LiftAlarm struct {
	Sensor BS2LiftSensor
	Action BS2Action
}

/**
 * BS2_LIFT_ALARM_FLAG
 */
const (
	BS2_LIFT_ALARM_FLAG_NONE   = 0x00
	BS2_LIFT_ALARM_FLAG_FIRST  = 0x01
	BS2_LIFT_ALARM_FLAG_SECOND = 0x02
	BS2_LIFT_ALARM_FLAG_TAMPER = 0x04
)

type BS2_LIFT_ALARM_FLAG = uint8

/**
 *  BS2LiftStatus
 */
type BS2LiftStatus struct {
	LiftID     BS2_LIFT_ID                            ///< 4 bytes
	NumFloors  uint16                                 ///< 2 bytes
	AlarmFlags BS2_LIFT_ALARM_FLAG                    ///< 1 byte
	TamperOn   bool                                   ///< 1 byte
	Floors     [BS2_MAX_FLOORS_ON_LIFT]BS2FloorStatus ///< 3 * 255 bytes
}

type BS2_ACCESS_GROUP_ID uint32

/**
 *	BS2Lift
 */
type BS2Lift struct {
	LiftID BS2_LIFT_ID ///< 4 bytes
	Name   [BS2_MAX_LIFT_NAME_LEN]byte

	DeviceID [BS2_MAX_DEVICES_ON_LIFT]BS2_DEVICE_ID ///< 4 * 4 bytes

	ActivateTimeout uint32 ///< 4 bytes (in seconds)
	DualAuthTimeout uint32 ///< 4 bytes

	NumFloors                 uint8                  ///< 1 byte
	NumDualAuthApprovalGroups uint8                  ///< 1 byte
	DualAuthApprovalType      BS2_DUAL_AUTH_APPROVAL ///< 1 byte
	TamperOn                  bool                   ///< 1 byte

	DualAuthRequired   [BS2_MAX_DEVICES_ON_LIFT]bool ///< 4 * 1 byte
	DualAuthScheduleID BS2_SCHEDULE_ID               ///< 4 bytes

	Floor                   [BS2_MAX_FLOORS_ON_LIFT]BS2LiftFloor                          ///< 8 * 255 bytes
	DualAuthApprovalGroupID [BS2_MAX_DUAL_AUTH_APPROVAL_GROUP_ON_LIFT]BS2_ACCESS_GROUP_ID ///< 4 * 16 bytes

	Alarm  [BS2_MAX_ALARMS_ON_LIFT]BS2LiftAlarm
	Tamper BS2LiftAlarm

	AlarmFlags BS2_LIFT_ALARM_FLAG ///< 1 byte
	Reserved   [3]uint8            ///< 3 bytes (packing)
}

func GetLift(context SdkContext, deviceId BS2_DEVICE_ID, liftIds []uint32, liftObj *[]BS2Lift) int16 {

	var liftIdsPtr uintptr
	var liftObjPtr *BS2Lift
	var numLift uint32

	if 0 != len(liftIds) {
		liftIdsPtr = uintptr(unsafe.Pointer(&liftIds[0]))
	}

	ret, _, err := syscall.Syscall6(procBs2GetLift, uintptr(6), context, uintptr(deviceId), liftIdsPtr, uintptr(len(liftIds)), uintptr(unsafe.Pointer(&liftObjPtr)), uintptr(unsafe.Pointer(&numLift)))

	defer ReleaseObject(uintptr(unsafe.Pointer(liftObjPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2Lift, numLift)

	var i uintptr = 0
	for ; i < (uintptr)(numLift); i++ {
		temp[i] = *(*BS2Lift)(unsafe.Pointer(uintptr(unsafe.Pointer(liftObjPtr)) + i*unsafe.Sizeof(*liftObjPtr)))
	}
	*liftObj = temp

	return int16(ret)
}
func GetAllLift(context SdkContext, deviceId BS2_DEVICE_ID, liftObj *[]BS2Lift) int16 {
	var liftObjPtr *BS2Lift
	var numLift uint32
	ret, _, err := syscall.Syscall6(procBs2GetAllLift, uintptr(4), context, uintptr(deviceId), uintptr(unsafe.Pointer(&liftObjPtr)), uintptr(unsafe.Pointer(&numLift)), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}

	return int16(ret)
}

func GetLiftStatus(context SdkContext, deviceId BS2_DEVICE_ID, liftIds []BS2_LIFT_ID, liftStatusObj *[]BS2LiftStatus) int16 {
	var liftIdsPtr uintptr
	var liftStatusObjPtr *BS2LiftStatus
	var numLift uint32

	if 0 != len(liftIds) {
		liftIdsPtr = uintptr(unsafe.Pointer(&liftIds[0]))
	}

	ret, _, err := syscall.Syscall6(procBs2GetLiftStatus, uintptr(6), context, liftIdsPtr, uintptr(len(liftIds)), uintptr(unsafe.Pointer(&liftStatusObjPtr)), uintptr(unsafe.Pointer(&numLift)), uintptr(0))

	defer ReleaseObject(uintptr(unsafe.Pointer(liftStatusObjPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2LiftStatus, numLift)

	var i uintptr = 0
	for ; i < (uintptr)(numLift); i++ {
		temp[i] = *(*BS2LiftStatus)(unsafe.Pointer(uintptr(unsafe.Pointer(liftStatusObjPtr)) + i*unsafe.Sizeof(*liftStatusObjPtr)))
	}
	*liftStatusObj = temp

	return int16(ret)
}
func GetAllLiftStatus(context SdkContext, deviceId BS2_DEVICE_ID, liftStatusObj *[]BS2LiftStatus) int16 {
	var liftStatusObjPtr *BS2LiftStatus
	var numLift uint32
	ret, _, err := syscall.Syscall6(procBs2GetAllLiftStatus, uintptr(4), context, uintptr(deviceId), uintptr(unsafe.Pointer(&liftStatusObjPtr)), uintptr(unsafe.Pointer(&numLift)), uintptr(0), uintptr(0))

	defer ReleaseObject(uintptr(unsafe.Pointer(liftStatusObjPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2LiftStatus, numLift)

	var i uintptr = 0
	for ; i < (uintptr)(numLift); i++ {
		temp[i] = *(*BS2LiftStatus)(unsafe.Pointer(uintptr(unsafe.Pointer(liftStatusObjPtr)) + i*unsafe.Sizeof(*liftStatusObjPtr)))
	}
	*liftStatusObj = temp
	return int16(ret)
}
func SetLift(context SdkContext, deviceId BS2_DEVICE_ID, lifts []BS2Lift) int16 {
	var liftsPtr uintptr
	if 0 != len(lifts) {
		liftsPtr = uintptr(unsafe.Pointer(&lifts[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2SetLift, uintptr(4), context, uintptr(deviceId), liftsPtr, uintptr(len(lifts)), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func SetLiftAlarm(context SdkContext, deviceId BS2_DEVICE_ID, flag BS2_LIFT_ALARM_FLAG, liftIds []BS2_LIFT_ID) int16 {
	var liftIdsPtr uintptr
	if 0 != len(liftIds) {
		liftIdsPtr = uintptr(unsafe.Pointer(&liftIds[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2SetLiftAlarm, uintptr(5), context, uintptr(deviceId), uintptr(flag), liftIdsPtr, uintptr(len(liftIds)), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveLift(context SdkContext, deviceId BS2_DEVICE_ID, liftIds []uint32) int16 {
	var liftIdsPtr uintptr
	if 0 != len(liftIds) {
		liftIdsPtr = uintptr(unsafe.Pointer(&liftIds[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2RemoveLift, uintptr(4), context, uintptr(deviceId), liftIdsPtr, uintptr(len(liftIds)), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveAllLift(context SdkContext, deviceId BS2_DEVICE_ID) int16 {
	ret, _, err := syscall.Syscall(procBs2RemoveAllLift, uintptr(2), context, uintptr(deviceId), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func ReleaseFloor(context SdkContext, deviceId BS2_DEVICE_ID, flag BS2_FLOOR_FLAG, liftId BS2_LIFT_ID, floorIndex []uint16) int16 {
	var floorIndicesPtr uintptr
	if 0 != len(floorIndex) {
		floorIndicesPtr = uintptr(unsafe.Pointer(&floorIndex[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2ReleaseFloor, uintptr(6), context, uintptr(deviceId), uintptr(flag), uintptr(liftId), floorIndicesPtr, uintptr(len(floorIndex)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func ActivateFloor(context SdkContext, deviceId BS2_DEVICE_ID, flag BS2_FLOOR_FLAG, liftId BS2_LIFT_ID, floorIndex []uint16) int16 {
	var floorIndicesPtr uintptr
	if 0 != len(floorIndex) {
		floorIndicesPtr = uintptr(unsafe.Pointer(&floorIndex[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2ActivateFloor, uintptr(6), context, uintptr(deviceId), uintptr(flag), uintptr(liftId), floorIndicesPtr, uintptr(len(floorIndex)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func DeActivateFloor(context SdkContext, deviceId BS2_DEVICE_ID, flag BS2_FLOOR_FLAG, liftId BS2_LIFT_ID, floorIndex []uint16) int16 {
	var floorIndicesPtr uintptr
	if 0 != len(floorIndex) {
		floorIndicesPtr = uintptr(unsafe.Pointer(&floorIndex[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2DeActivateFloor, uintptr(6), context, uintptr(deviceId), uintptr(flag), uintptr(liftId), floorIndicesPtr, uintptr(len(floorIndex)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

const (
	BS2_MAX_FLOOR_LEVEL_ITEMS    = 128
	BS2_MAX_FLOOR_LEVEL_NAME_LEN = 48 * 3

	BS2_INVALID_FLOOR_LEVEL_ID = 0
)

/**
 *	BS2FloorSchedule
 */
type BS2FloorSchedule struct {
	LiftID     BS2_LIFT_ID
	FloorIndex uint16
	Reserved   [2]uint8
	ScheduleID BS2_SCHEDULE_ID
}

type BS2_FLOOR_LEVEL_ID = uint32

/**
 *	BS2FloorLevel
 */
type BS2FloorLevel struct {
	Id                BS2_FLOOR_LEVEL_ID // id >= 32768 (BS2_ACCESS_LEVEL_ID < 32768)
	Name              [BS2_MAX_FLOOR_LEVEL_NAME_LEN]byte
	NumFloorSchedules uint8
	Reserved          [3]uint8
	FloorSchedules    [BS2_MAX_FLOOR_LEVEL_ITEMS]BS2FloorSchedule
}

func GetFloorLevel(context SdkContext, deviceId BS2_DEVICE_ID, floorLevelIds []uint32, floorLevelObj *[]BS2FloorLevel) int16 {
	var floorLevelIdsPtr uintptr
	if 0 != len(floorLevelIds) {
		floorLevelIdsPtr = uintptr(unsafe.Pointer(&floorLevelIds[0]))
	}

	var floorLevelObjPtr *BS2FloorLevel
	var numFloorLevel uint32

	ret, _, err := syscall.Syscall6(procBs2GetFloorLevel, uintptr(6), context, uintptr(deviceId), floorLevelIdsPtr, uintptr(len(floorLevelIds)), uintptr(unsafe.Pointer(&floorLevelObjPtr)), uintptr(unsafe.Pointer(&numFloorLevel)))
	defer ReleaseObject(uintptr(unsafe.Pointer(floorLevelObjPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2FloorLevel, numFloorLevel)

	var i uintptr = 0
	for ; i < (uintptr)(numFloorLevel); i++ {
		temp[i] = *(*BS2FloorLevel)(unsafe.Pointer(uintptr(unsafe.Pointer(floorLevelObjPtr)) + i*unsafe.Sizeof(*floorLevelObjPtr)))
	}
	*floorLevelObj = temp

	return int16(ret)
}
func GetAllFloorLevel(context SdkContext, deviceId BS2_DEVICE_ID, floorLevelObj *[]BS2FloorLevel) int16 {
	var floorLevelObjPtr *BS2FloorLevel
	var numFloorLevel uint32

	ret, _, err := syscall.Syscall6(procBs2GetAllFloorLevel, uintptr(4), context, uintptr(deviceId), uintptr(unsafe.Pointer(&floorLevelObjPtr)), uintptr(unsafe.Pointer(&numFloorLevel)), uintptr(0), uintptr(0))
	defer ReleaseObject(uintptr(unsafe.Pointer(floorLevelObjPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2FloorLevel, numFloorLevel)

	var i uintptr = 0
	for ; i < (uintptr)(numFloorLevel); i++ {
		temp[i] = *(*BS2FloorLevel)(unsafe.Pointer(uintptr(unsafe.Pointer(floorLevelObjPtr)) + i*unsafe.Sizeof(*floorLevelObjPtr)))
	}
	*floorLevelObj = temp

	return int16(ret)
}
func SetFloorLevel(context SdkContext, deviceId BS2_DEVICE_ID, floorLevels []BS2FloorLevel) int16 {
	var floorLevelsPtr uintptr
	if 0 != len(floorLevels) {
		floorLevelsPtr = uintptr(unsafe.Pointer(&floorLevels[0]))
	}

	ret, _, err := syscall.Syscall6(procBs2SetFloorLevel, uintptr(4), context, uintptr(deviceId), floorLevelsPtr, uintptr(len(floorLevels)), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveFloorLevel(context SdkContext, deviceId BS2_DEVICE_ID, floorLevelIds []uint32) int16 {
	var floorLevelIdsPtr uintptr
	if 0 != len(floorLevelIds) {
		floorLevelIdsPtr = uintptr(unsafe.Pointer(&floorLevelIds[0]))
	}

	ret, _, err := syscall.Syscall6(procBs2RemoveFloorLevel, uintptr(4), context, uintptr(deviceId), floorLevelIdsPtr, uintptr(len(floorLevelIds)), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveAllFloorLevel(context SdkContext, deviceId BS2_DEVICE_ID) int16 {
	ret, _, err := syscall.Syscall(procBs2RemoveAllFloorLevel, uintptr(2), context, uintptr(deviceId), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

type BS2_ZONE_ID = uint32

func GetAntiPassbackZone(context SdkContext, deviceId BS2_DEVICE_ID, zoneIds []BS2_ZONE_ID, zoneObj *[]BS2AntiPassbackZone) int16 {
	var zoneIdsPtr uintptr
	if 0 != len(zoneIds) {
		zoneIdsPtr = uintptr(unsafe.Pointer(&zoneIds[0]))
	}

	var zoneObjPtr *BS2AntiPassbackZone
	var numZone uint32

	ret, _, err := syscall.Syscall6(procBs2GetAntiPassbackZone, uintptr(6), context, uintptr(deviceId), zoneIdsPtr, uintptr(len(zoneIds)), uintptr(unsafe.Pointer(&zoneObjPtr)), uintptr(unsafe.Pointer(&numZone)))
	defer ReleaseObject(uintptr(unsafe.Pointer(zoneObjPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}

	temp := make([]BS2AntiPassbackZone, numZone)

	var i uintptr = 0
	for ; i < (uintptr)(numZone); i++ {
		temp[i] = *(*BS2AntiPassbackZone)(unsafe.Pointer(uintptr(unsafe.Pointer(zoneObjPtr)) + i*unsafe.Sizeof(*zoneObjPtr)))
	}
	*zoneObj = temp

	return int16(ret)
}
func GetAllAntiPassbackZone(context SdkContext, deviceId BS2_DEVICE_ID, zoneObj *[]BS2AntiPassbackZone) int16 {

	var zoneObjPtr *BS2AntiPassbackZone
	var numZone uint32
	ret, _, err := syscall.Syscall6(procBs2GetAllAntiPassbackZone, uintptr(4), context, uintptr(deviceId), uintptr(unsafe.Pointer(&zoneObjPtr)), uintptr(unsafe.Pointer(&numZone)), uintptr(0), uintptr(0))
	defer ReleaseObject(uintptr(unsafe.Pointer(zoneObjPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2AntiPassbackZone, numZone)

	var i uintptr = 0
	for ; i < (uintptr)(numZone); i++ {
		temp[i] = *(*BS2AntiPassbackZone)(unsafe.Pointer(uintptr(unsafe.Pointer(zoneObjPtr)) + i*unsafe.Sizeof(*zoneObjPtr)))
	}
	*zoneObj = temp
	return int16(ret)
}

/**
 *  Zone Constants
 */
const (
	BS2_INVALID_ZONE_ID = 0
)

/**
 * BS2_ZONE_TYPE
 */
const (
	BS2_ZONE_APB = iota
	BS2_ZONE_TIMED_APB
	BS2_ZONE_FIRE_ALARM
	BS2_ZONE_SCHEDULED_LOCK_UNLOCK                                  // Defined old (Until BioStar SDK V2.3.3)
	BS2_ZONE_FORCED_LOCK_UNLOCK    = BS2_ZONE_SCHEDULED_LOCK_UNLOCK // Defined newly - Recommend use this (From BioStar SDK V2.4.0)
	BS2_ZONE_INTRUSION_ALARM
)

type BS2_ZONE_TYPE uint8

/**
*	BS2_ZONE_STATUS
 */
const (
	BS2_ZONE_STATUS_NORMAL = 0x00
	BS2_ZONE_STATUS_ALARM  = 0x01
	// Defined old (Until BioStar SDK V2.3.3)
	BS2_ZONE_STATUS_SCHEDULED_LOCKED   = 0x02
	BS2_ZONE_STATUS_SCHEDULED_UNLOCKED = 0x04
	// Defined newly - Recommend use this (From BioStar SDK V2.4.0)
	BS2_ZONE_STATUS_FORCED_LOCKED   = 0x02
	BS2_ZONE_STATUS_FORCED_UNLOCKED = 0x04
	BS2_ZONE_STATUS_ARM             = 0x08
	BS2_ZONE_STATUS_DISARM          = BS2_ZONE_STATUS_NORMAL
)

type BS2_ZONE_STATUS uint8

/**
 *  BS2ZoneStatus
 */
type BS2ZoneStatus struct {
	Id       BS2_ZONE_ID     ///< 4 bytes
	Status   BS2_ZONE_STATUS ///< 1 byte
	Disabled bool            ///< 1 byte
	Reserved [6]uint8        ///< 6 bytes (packing)
}

func GetAntiPassbackZoneStatus(context SdkContext, deviceId BS2_DEVICE_ID, zoneIds []BS2_ZONE_ID, zoneStatusObj *[]BS2ZoneStatus) int16 {
	var zoneIdsPtr uintptr
	if 0 != len(zoneIds) {
		zoneIdsPtr = uintptr(unsafe.Pointer(&zoneIds[0]))
	}

	var zoneStatusObjPtr *BS2ZoneStatus
	var numZone uint32

	ret, _, err := syscall.Syscall6(procBs2GetAntiPassbackZoneStatus, uintptr(6), context, uintptr(deviceId), zoneIdsPtr, uintptr(len(zoneIds)), uintptr(unsafe.Pointer(&zoneStatusObjPtr)), uintptr(unsafe.Pointer(&numZone)))
	defer ReleaseObject(uintptr(unsafe.Pointer(zoneStatusObjPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2ZoneStatus, numZone)

	var i uintptr = 0
	for ; i < (uintptr)(numZone); i++ {
		temp[i] = *(*BS2ZoneStatus)(unsafe.Pointer(uintptr(unsafe.Pointer(zoneStatusObjPtr)) + i*unsafe.Sizeof(*zoneStatusObjPtr)))
	}
	*zoneStatusObj = temp
	return int16(ret)
}
func GetAllAntiPassbackZoneStatus(context SdkContext, deviceId BS2_DEVICE_ID, zoneStatusObj *[]BS2ZoneStatus) int16 {

	var zoneStatusObjPtr *BS2ZoneStatus
	var numZone uint32

	ret, _, err := syscall.Syscall6(procBs2GetAllAntiPassbackZoneStatus, uintptr(4), context, uintptr(deviceId), uintptr(unsafe.Pointer(&zoneStatusObjPtr)), uintptr(unsafe.Pointer(&numZone)), uintptr(0), uintptr(0))
	defer ReleaseObject(uintptr(unsafe.Pointer(zoneStatusObjPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2ZoneStatus, numZone)

	var i uintptr = 0
	for ; i < (uintptr)(numZone); i++ {
		temp[i] = *(*BS2ZoneStatus)(unsafe.Pointer(uintptr(unsafe.Pointer(zoneStatusObjPtr)) + i*unsafe.Sizeof(*zoneStatusObjPtr)))
	}
	*zoneStatusObj = temp
	return int16(ret)
}
func SetAntiPassbackZone(context SdkContext, deviceId BS2_DEVICE_ID, zoneStatusObj *[]BS2AntiPassbackZone) int16 {
	var sizeOfZoneStaus uint32 = uint32(len(*zoneStatusObj))
	var zoneStatusObjPtr uintptr
	if 0 != sizeOfZoneStaus {
		zoneStatusObjPtr = uintptr(unsafe.Pointer(&(*zoneStatusObj)[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2SetAntiPassbackZone, uintptr(4), context, uintptr(deviceId), zoneStatusObjPtr, uintptr(sizeOfZoneStaus), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func SetAntiPassbackZoneAlarm(context SdkContext, deviceId BS2_DEVICE_ID, alarmed uint8, zoneIds []BS2_ZONE_ID) int16 {
	var sizeOfZoneIds uint32 = uint32(len(zoneIds))
	var zoneIdsPtr uintptr
	if 0 != sizeOfZoneIds {
		zoneIdsPtr = uintptr(unsafe.Pointer(&zoneIds[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2SetAntiPassbackZoneAlarm, uintptr(5), context, uintptr(deviceId), uintptr(alarmed), zoneIdsPtr, uintptr(sizeOfZoneIds), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveAntiPassbackZone(context SdkContext, deviceId BS2_DEVICE_ID, zoneIds []BS2_ZONE_ID) int16 {
	var sizeOfZoneIds uint32 = uint32(len(zoneIds))
	var zoneIdsPtr uintptr
	if 0 != sizeOfZoneIds {
		zoneIdsPtr = uintptr(unsafe.Pointer(&zoneIds[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2RemoveAntiPassbackZone, uintptr(4), context, uintptr(deviceId), zoneIdsPtr, uintptr(sizeOfZoneIds), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveAllAntiPassbackZone(context SdkContext, deviceId BS2_DEVICE_ID) int16 {
	ret, _, err := syscall.Syscall(procBs2RemoveAllAntiPassbackZone, uintptr(2), context, uintptr(deviceId), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func ClearAntiPassbackZoneStatus(context SdkContext, deviceId BS2_DEVICE_ID, zoneId BS2_ZONE_ID, userIds []BS2_USER_ID) int16 {
	var sizeOfUserIds uint32 = uint32(len(userIds))
	var userIdsPtr uintptr
	if 0 != sizeOfUserIds {
		userIdsPtr = uintptr(unsafe.Pointer(&userIds[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2ClearAntiPassbackZoneStatus, uintptr(5), context, uintptr(deviceId), uintptr(zoneId), userIdsPtr, uintptr(sizeOfUserIds), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func ClearAllAntiPassbackZoneStatus(context SdkContext, deviceId BS2_DEVICE_ID, zoneId BS2_ZONE_ID) int16 {
	ret, _, err := syscall.Syscall(procBs2ClearAllAntiPassbackZoneStatus, uintptr(3), context, uintptr(deviceId), uintptr(zoneId))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func SetCheckGlobalAPBViolationHandler(context SdkContext, fnOnCheckGlobalAPBViolation interface{}) int16 {
	var fnOnCheckGlobalAPBViolationPtr uintptr
	if nil != fnOnCheckGlobalAPBViolation {
		fnOnCheckGlobalAPBViolationPtr = syscall.NewCallback(fnOnCheckGlobalAPBViolation)
	}
	ret, _, err := syscall.Syscall(procBs2SetCheckGlobalAPBViolationHandler, uintptr(2), context, fnOnCheckGlobalAPBViolationPtr, uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func CheckGlobalAPBViolation(context SdkContext, deviceId BS2_DEVICE_ID, seq uint16, handleResult int, zoneId uint32) int16 {
	ret, _, err := syscall.Syscall6(procBs2CheckGlobalAPBViolation, uintptr(5), context, uintptr(deviceId), uintptr(seq), uintptr(handleResult), uintptr(zoneId), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

const (
	BS2_MAX_READERS_PER_TIMED_APB_ZONE       = 64
	BS2_MAX_BYPASS_GROUPS_PER_TIMED_APB_ZONE = 16
	BS2_MAX_TIMED_APB_ALARM_ACTION           = 5
)

/**
 *  BS2_TIMED_APB_ZONE_TYPE
 */
const (
	BS2_TIMED_APB_ZONE_HARD = 0x00
	BS2_TIMED_APB_ZONE_SOFT = 0x01
)

type BS2_TIMED_APB_ZONE_TYPE = uint8

type BS2TimedApbMember struct {
	DeviceID BS2_DEVICE_ID ///< 4 bytes
	Reserved [4]uint8      ///< 4 bytes (packing)
}

type BS2TimedAntiPassbackZone struct {
	ZoneID BS2_ZONE_ID                 ///< 4 bytes
	Name   [BS2_MAX_ZONE_NAME_LEN]byte ///< 48 * 3 bytes

	Type            BS2_TIMED_APB_ZONE_TYPE ///< 1 byte
	NumReaders      uint8                   ///< 1 byte
	NumBypassGroups uint8                   ///< 1 byte
	Disabled        bool                    ///< 1 byte

	Alarmed  bool     ///< 1 byte
	Reserved [3]uint8 ///< 3 bytes (packing)

	ResetDuration uint32 ///< 4 bytes: in seconds 0: no reset

	Alarm [BS2_MAX_TIMED_APB_ALARM_ACTION]BS2Action ///< 32 * 5 bytes

	Readers   [BS2_MAX_READERS_PER_TIMED_APB_ZONE]BS2TimedApbMember ///< 8 * 64 bytes
	Reserved2 [8 * 40]uint8                                         ///< 8 * 40 bytes (packing)

	BypassGroupIDs [BS2_MAX_BYPASS_GROUPS_PER_TIMED_APB_ZONE]BS2_ACCESS_GROUP_ID ///< 4 * 16 bytes
}

func GetTimedAntiPassbackZone(context SdkContext, deviceId BS2_DEVICE_ID, zoneIds []uint32, zoneObj *[]BS2TimedAntiPassbackZone) int16 {
	var zoneIdsPtr uintptr
	if 0 != len(zoneIds) {
		zoneIdsPtr = uintptr(unsafe.Pointer(&zoneIds[0]))
	}

	var zoneObjPtr *BS2TimedAntiPassbackZone
	var numZone uint32

	ret, _, err := syscall.Syscall6(procBs2GetTimedAntiPassbackZone, uintptr(6), context, uintptr(deviceId), zoneIdsPtr, uintptr(len(zoneIds)), uintptr(unsafe.Pointer(&zoneObjPtr)), uintptr(unsafe.Pointer(&numZone)))
	defer ReleaseObject(uintptr(unsafe.Pointer(zoneObjPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2TimedAntiPassbackZone, numZone)

	var i uintptr = 0
	for ; i < (uintptr)(numZone); i++ {
		temp[i] = *(*BS2TimedAntiPassbackZone)(unsafe.Pointer(uintptr(unsafe.Pointer(zoneObjPtr)) + i*unsafe.Sizeof(*zoneObjPtr)))
	}
	*zoneObj = temp
	return int16(ret)
}
func GetAllTimedAntiPassbackZone(context SdkContext, deviceId BS2_DEVICE_ID, zoneObj *[]BS2TimedAntiPassbackZone) int16 {

	var zoneObjPtr *BS2TimedAntiPassbackZone
	var numZone uint32
	ret, _, err := syscall.Syscall6(procBs2GetAllTimedAntiPassbackZone, uintptr(4), context, uintptr(deviceId), uintptr(unsafe.Pointer(&zoneObjPtr)), uintptr(unsafe.Pointer(&numZone)), uintptr(0), uintptr(0))
	defer ReleaseObject(uintptr(unsafe.Pointer(zoneObjPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2TimedAntiPassbackZone, numZone)

	var i uintptr = 0
	for ; i < (uintptr)(numZone); i++ {
		temp[i] = *(*BS2TimedAntiPassbackZone)(unsafe.Pointer(uintptr(unsafe.Pointer(zoneObjPtr)) + i*unsafe.Sizeof(*zoneObjPtr)))
	}
	*zoneObj = temp
	return int16(ret)
}
func GetTimedAntiPassbackZoneStatus(context SdkContext, deviceId BS2_DEVICE_ID, zoneIds []uint32, zoneStatusObj *[]BS2ZoneStatus) int16 {

	var zoneIdsPtr uintptr
	if 0 != len(zoneIds) {
		zoneIdsPtr = uintptr(unsafe.Pointer(&zoneIds[0]))
	}

	var zoneObjPtr *BS2ZoneStatus
	var numZone uint32

	ret, _, err := syscall.Syscall6(procBs2GetTimedAntiPassbackZoneStatus, uintptr(6), context, uintptr(deviceId), zoneIdsPtr, uintptr(len(zoneIds)), uintptr(unsafe.Pointer(&zoneObjPtr)), uintptr(unsafe.Pointer(&numZone)))
	defer ReleaseObject(uintptr(unsafe.Pointer(zoneObjPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2ZoneStatus, numZone)

	var i uintptr = 0
	for ; i < (uintptr)(numZone); i++ {
		temp[i] = *(*BS2ZoneStatus)(unsafe.Pointer(uintptr(unsafe.Pointer(zoneObjPtr)) + i*unsafe.Sizeof(*zoneObjPtr)))
	}
	*zoneStatusObj = temp
	return int16(ret)
}
func GetAllTimedAntiPassbackZoneStatus(context SdkContext, deviceId BS2_DEVICE_ID, zoneStatusObj *[]BS2ZoneStatus) int16 {

	var zoneObjPtr *BS2ZoneStatus
	var numZone uint32
	ret, _, err := syscall.Syscall6(procBs2GetAllTimedAntiPassbackZoneStatus, uintptr(4), context, uintptr(deviceId), uintptr(unsafe.Pointer(&zoneObjPtr)), uintptr(unsafe.Pointer(&numZone)), uintptr(0), uintptr(0))
	defer ReleaseObject(uintptr(unsafe.Pointer(zoneObjPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2ZoneStatus, numZone)

	var i uintptr = 0
	for ; i < (uintptr)(numZone); i++ {
		temp[i] = *(*BS2ZoneStatus)(unsafe.Pointer(uintptr(unsafe.Pointer(zoneObjPtr)) + i*unsafe.Sizeof(*zoneObjPtr)))
	}
	*zoneStatusObj = temp
	return int16(ret)
}
func SetTimedAntiPassbackZone(context SdkContext, deviceId BS2_DEVICE_ID, zones *[]BS2TimedAntiPassbackZone) int16 {
	var sizeOfZones uint32 = uint32(len(*zones))
	var zonePtr uintptr
	if 0 != sizeOfZones {
		zonePtr = uintptr(unsafe.Pointer(&(*zones)[0]))
	}

	ret, _, err := syscall.Syscall6(procBs2SetTimedAntiPassbackZone, uintptr(4), context, uintptr(deviceId), zonePtr, uintptr(sizeOfZones), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func SetTimedAntiPassbackZoneAlarm(context SdkContext, deviceId BS2_DEVICE_ID, alarmed uint8, zoneIds []BS2_ZONE_ID) int16 {

	var zoneIdsPtr uintptr
	if 0 != len(zoneIds) {
		zoneIdsPtr = uintptr(unsafe.Pointer(&zoneIds[0]))
	}

	ret, _, err := syscall.Syscall6(procBs2SetTimedAntiPassbackZoneAlarm, uintptr(5), context, uintptr(deviceId), uintptr(alarmed), uintptr(zoneIdsPtr), uintptr(len(zoneIds)), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveTimedAntiPassbackZone(context SdkContext, deviceId BS2_DEVICE_ID, zoneIds []uint32) int16 {
	var zoneIdsPtr uintptr
	if 0 != len(zoneIds) {
		zoneIdsPtr = uintptr(unsafe.Pointer(&zoneIds[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2RemoveTimedAntiPassbackZone, uintptr(4), context, uintptr(deviceId), zoneIdsPtr, uintptr(len(zoneIds)), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveAllTimedAntiPassbackZone(context SdkContext, deviceId BS2_DEVICE_ID) int16 {
	ret, _, err := syscall.Syscall(procBs2RemoveAllTimedAntiPassbackZone, uintptr(2), context, uintptr(deviceId), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func ClearTimedAntiPassbackZoneStatus(context SdkContext, deviceId BS2_DEVICE_ID, zoneId uint32, userIds []BS2_USER_ID) int16 {
	var userIdsPtr uintptr
	if 0 != len(userIds) {
		userIdsPtr = uintptr(unsafe.Pointer(&userIds[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2ClearTimedAntiPassbackZoneStatus, uintptr(5), context, uintptr(deviceId), uintptr(zoneId), userIdsPtr, uintptr(len(userIds)), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func ClearAllTimedAntiPassbackZoneStatus(context SdkContext, deviceId BS2_DEVICE_ID, zoneId uint32) int16 {
	ret, _, err := syscall.Syscall(procBs2ClearAllTimedAntiPassbackZoneStatus, uintptr(3), context, uintptr(deviceId), uintptr(zoneId))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

type BS2_FIRE_SWITCH = uint8

type BS2FireSensor struct {
	DeviceID   BS2_DEVICE_ID   ///< 4 bytes
	Port       uint8           ///< 1 byte
	SwitchType BS2_SWITCH_TYPE ///< 1 byte
	Duration   uint16          ///< 2 bytes
}

const (
	BS2_MAX_FIRE_SENSORS_PER_FIRE_ALARM_ZONE = 8
	BS2_MAX_FIRE_ALARM_ACTION                = 5
	BS2_MAX_DOORS_PER_FIRE_ALARM_ZONE        = 32
)

type BS2FireAlarmZone struct {
	ZoneID BS2_ZONE_ID                 ///< 4 bytes
	Name   [BS2_MAX_ZONE_NAME_LEN]byte ///< 48 * 3 bytes

	NumSensors uint8 ///< 1 byte

	NumMembers uint8

	Alarmed  bool ///< 1 byte
	Disabled bool ///< 1 byte

	Reserved [8]uint8 ///< 8 bytes (packing)

	Sensor [BS2_MAX_FIRE_SENSORS_PER_FIRE_ALARM_ZONE]BS2FireSensor ///< 8 * 8 bytes
	Alarm  [BS2_MAX_FIRE_ALARM_ACTION]BS2Action                    ///< 32 * 5 bytes

	Reserved2 [32]uint8 ///< 32 bytes (packing)

	MemberIDs [BS2_MAX_DOORS_PER_FIRE_ALARM_ZONE]uint32 ///< 4 * 32 bytes

}

func GetFireAlarmZone(context SdkContext, deviceId BS2_DEVICE_ID, zoneIds []uint32, zoneObj *[]BS2FireAlarmZone) int16 {
	var zoneIdsPtr uintptr
	if 0 != len(zoneIds) {
		zoneIdsPtr = uintptr(unsafe.Pointer(&zoneIds[0]))
	}

	var zoneObjPtr *BS2FireAlarmZone
	var numZone uint32
	ret, _, err := syscall.Syscall6(procBs2GetFireAlarmZone, uintptr(6), context, uintptr(deviceId), zoneIdsPtr, uintptr(len(zoneIds)), uintptr(unsafe.Pointer(&zoneObjPtr)), uintptr(unsafe.Pointer(&numZone)))
	defer ReleaseObject(uintptr(unsafe.Pointer(zoneObjPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2FireAlarmZone, numZone)

	var i uintptr = 0
	for ; i < (uintptr)(numZone); i++ {
		temp[i] = *(*BS2FireAlarmZone)(unsafe.Pointer(uintptr(unsafe.Pointer(zoneObjPtr)) + i*unsafe.Sizeof(*zoneObjPtr)))
	}
	*zoneObj = temp
	return int16(ret)
}
func GetAllFireAlarmZone(context SdkContext, deviceId BS2_DEVICE_ID, zoneObj *[]BS2FireAlarmZone) int16 {

	var zoneObjPtr *BS2TimedAntiPassbackZone
	var numZone uint32

	ret, _, err := syscall.Syscall6(procBs2GetAllFireAlarmZone, uintptr(4), context, uintptr(deviceId), uintptr(unsafe.Pointer(&zoneObjPtr)), uintptr(unsafe.Pointer(&numZone)), uintptr(0), uintptr(0))
	defer ReleaseObject(uintptr(unsafe.Pointer(zoneObjPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2FireAlarmZone, numZone)

	var i uintptr = 0
	for ; i < (uintptr)(numZone); i++ {
		temp[i] = *(*BS2FireAlarmZone)(unsafe.Pointer(uintptr(unsafe.Pointer(zoneObjPtr)) + i*unsafe.Sizeof(*zoneObjPtr)))
	}
	*zoneObj = temp
	return int16(ret)
}
func GetFireAlarmZoneStatus(context SdkContext, deviceId BS2_DEVICE_ID, zoneIds []uint32, zoneStatusObj *[]BS2ZoneStatus) int16 {
	var zoneIdsPtr uintptr
	if 0 != len(zoneIds) {
		zoneIdsPtr = uintptr(unsafe.Pointer(&zoneIds[0]))
	}

	var zoneObjPtr *BS2ZoneStatus
	var numZone uint32
	ret, _, err := syscall.Syscall6(procBs2GetFireAlarmZoneStatus, uintptr(6), context, uintptr(deviceId), zoneIdsPtr, uintptr(len(zoneIds)), uintptr(unsafe.Pointer(&zoneObjPtr)), uintptr(unsafe.Pointer(&numZone)))
	defer ReleaseObject(uintptr(unsafe.Pointer(zoneObjPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2ZoneStatus, numZone)

	var i uintptr = 0
	for ; i < (uintptr)(numZone); i++ {
		temp[i] = *(*BS2ZoneStatus)(unsafe.Pointer(uintptr(unsafe.Pointer(zoneObjPtr)) + i*unsafe.Sizeof(*zoneObjPtr)))
	}
	*zoneStatusObj = temp
	return int16(ret)
}
func GetAllFireAlarmZoneStatus(context SdkContext, deviceId BS2_DEVICE_ID, zoneStatusObj *[]BS2ZoneStatus) int16 {
	var zoneObjPtr *BS2ZoneStatus
	var numZone uint32
	ret, _, err := syscall.Syscall6(procBs2GetAllFireAlarmZoneStatus, uintptr(4), context, uintptr(deviceId), uintptr(unsafe.Pointer(&zoneObjPtr)), uintptr(unsafe.Pointer(&numZone)), uintptr(0), uintptr(0))
	defer ReleaseObject(uintptr(unsafe.Pointer(zoneObjPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2ZoneStatus, numZone)

	var i uintptr = 0
	for ; i < (uintptr)(numZone); i++ {
		temp[i] = *(*BS2ZoneStatus)(unsafe.Pointer(uintptr(unsafe.Pointer(zoneObjPtr)) + i*unsafe.Sizeof(*zoneObjPtr)))
	}
	*zoneStatusObj = temp
	return int16(ret)
}
func SetFireAlarmZone(context SdkContext, deviceId BS2_DEVICE_ID, zones *[]BS2FireAlarmZone) int16 {
	var sizeOfZones uint32 = uint32(len(*zones))
	var zonesPtr uintptr
	if 0 != sizeOfZones {
		zonesPtr = uintptr(unsafe.Pointer(&(*zones)[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2SetFireAlarmZone, uintptr(4), context, uintptr(deviceId), zonesPtr, uintptr(sizeOfZones), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func SetFireAlarmZoneAlarm(context SdkContext, deviceId BS2_DEVICE_ID, alarmed uint8, zoneIds []uint32) int16 {
	var zoneIdsPtr uintptr
	if 0 != len(zoneIds) {
		zoneIdsPtr = uintptr(unsafe.Pointer(&zoneIds[0]))
	}

	ret, _, err := syscall.Syscall6(procBs2SetFireAlarmZoneAlarm, uintptr(4), context, uintptr(deviceId), uintptr(alarmed), zoneIdsPtr, uintptr(len(zoneIds)), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveFireAlarmZone(context SdkContext, deviceId BS2_DEVICE_ID, zoneIds []uint32) int16 {
	var zoneIdsPtr uintptr
	if 0 != len(zoneIds) {
		zoneIdsPtr = uintptr(unsafe.Pointer(&zoneIds[0]))
	}

	ret, _, err := syscall.Syscall6(procBs2RemoveFireAlarmZone, uintptr(4), context, uintptr(deviceId), zoneIdsPtr, uintptr(len(zoneIds)), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveAllFireAlarmZone(context SdkContext, deviceId BS2_DEVICE_ID) int16 {
	ret, _, err := syscall.Syscall(procBs2RemoveAllFireAlarmZone, uintptr(2), context, uintptr(deviceId), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

const (
	BS2_MAX_SCHEDULED_LOCK_UNLOCK_ALARM_ACTION          = 5
	BS2_MAX_DOORS_IN_SCHEDULED_LOCK_UNLOCK_ZONE         = 32
	BS2_MAX_BYPASS_GROUPS_IN_SCHEDULED_LOCK_UNLOCK_ZONE = 16
	BS2_MAX_UNLOCK_GROUPS_IN_SCHEDULED_LOCK_UNLOCK_ZONE = 16
)

type BS2_DOOR_ID = uint32

type BS2ScheduledLockUnlockZone struct {
	ZoneID BS2_ZONE_ID                 ///< 4 bytes
	Name   [BS2_MAX_ZONE_NAME_LEN]byte ///< 48 * 3 bytes

	LockScheduleID   BS2_SCHEDULE_ID ///< 4 bytes
	UnlockScheduleID BS2_SCHEDULE_ID ///< 4 bytes

	NumDoors          uint8 ///< 1 byte
	NumBypassGroups   uint8 ///< 1 byte
	NumUnlockGroups   uint8 ///< 1 byte
	BidirectionalLock bool  ///< 1 byte

	Disabled bool     ///< 1 byte
	Alarmed  bool     ///< 1 byte
	Reserved [6]uint8 ///< 6 bytes (packing)

	Alarm [BS2_MAX_SCHEDULED_LOCK_UNLOCK_ALARM_ACTION]BS2Action ///< 32 * 5 bytes

	Reserved2 [32]uint8 ///< 32 bytes (packing)

	DoorIDs        [BS2_MAX_DOORS_IN_SCHEDULED_LOCK_UNLOCK_ZONE]BS2_DOOR_ID                 ///< 4 * 32 bytes
	BypassGroupIDs [BS2_MAX_BYPASS_GROUPS_IN_SCHEDULED_LOCK_UNLOCK_ZONE]BS2_ACCESS_GROUP_ID ///< 4 * 16 bytes
	UnlockGroupIDs [BS2_MAX_UNLOCK_GROUPS_IN_SCHEDULED_LOCK_UNLOCK_ZONE]BS2_ACCESS_GROUP_ID ///< 4 * 16 bytes
}

func GetScheduledLockUnlockZone(context SdkContext, deviceId BS2_DEVICE_ID, zoneIds []uint32, zoneObj *[]BS2ScheduledLockUnlockZone) int16 {
	var zoneIdsPtr uintptr
	if 0 != len(zoneIds) {
		zoneIdsPtr = uintptr(unsafe.Pointer(&zoneIds[0]))
	}

	var zoneObjPtr *BS2ScheduledLockUnlockZone
	var numZone uint32

	ret, _, err := syscall.Syscall6(procBs2GetScheduledLockUnlockZone, uintptr(6), context, uintptr(deviceId), zoneIdsPtr, uintptr(len(zoneIds)), uintptr(unsafe.Pointer(&zoneObjPtr)), uintptr(unsafe.Pointer(&numZone)))
	defer ReleaseObject(uintptr(unsafe.Pointer(zoneObjPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2ScheduledLockUnlockZone, numZone)

	var i uintptr = 0
	for ; i < (uintptr)(numZone); i++ {
		temp[i] = *(*BS2ScheduledLockUnlockZone)(unsafe.Pointer(uintptr(unsafe.Pointer(zoneObjPtr)) + i*unsafe.Sizeof(*zoneObjPtr)))
	}
	*zoneObj = temp
	return int16(ret)
}
func GetAllScheduledLockUnlockZone(context SdkContext, deviceId BS2_DEVICE_ID, zoneObj *[]BS2ZoneStatus) int16 {

	var zoneObjPtr *BS2ZoneStatus
	var numZone uint32
	ret, _, err := syscall.Syscall6(procBs2GetAllScheduledLockUnlockZone, uintptr(4), context, uintptr(deviceId), uintptr(unsafe.Pointer(&zoneObjPtr)), uintptr(unsafe.Pointer(&numZone)), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2ZoneStatus, numZone)

	var i uintptr = 0
	for ; i < (uintptr)(numZone); i++ {
		temp[i] = *(*BS2ZoneStatus)(unsafe.Pointer(uintptr(unsafe.Pointer(zoneObjPtr)) + i*unsafe.Sizeof(*zoneObjPtr)))
	}
	*zoneObj = temp
	return int16(ret)
}
func GetScheduledLockUnlockZoneStatus(context SdkContext, deviceId BS2_DEVICE_ID, zoneIds []uint32, zoneStatusObj *[]BS2ZoneStatus) int16 {
	var zoneIdsPtr uintptr
	if 0 != len(zoneIds) {
		zoneIdsPtr = uintptr(unsafe.Pointer(&zoneIds[0]))
	}

	var zoneObjPtr *BS2ZoneStatus
	var numZone uint32

	ret, _, err := syscall.Syscall6(procBs2GetScheduledLockUnlockZoneStatus, uintptr(6), context, uintptr(deviceId), zoneIdsPtr, uintptr(len(zoneIds)), uintptr(unsafe.Pointer(&zoneObjPtr)), uintptr(unsafe.Pointer(&numZone)))
	defer ReleaseObject(uintptr(unsafe.Pointer(zoneObjPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2ZoneStatus, numZone)

	var i uintptr = 0
	for ; i < (uintptr)(numZone); i++ {
		temp[i] = *(*BS2ZoneStatus)(unsafe.Pointer(uintptr(unsafe.Pointer(zoneObjPtr)) + i*unsafe.Sizeof(*zoneObjPtr)))
	}
	*zoneStatusObj = temp
	return int16(ret)
}
func GetAllScheduledLockUnlockZoneStatus(context SdkContext, deviceId BS2_DEVICE_ID, zoneStatusObj *[]BS2ZoneStatus) int16 {
	var zoneObjPtr *BS2ZoneStatus
	var numZone uint32

	ret, _, err := syscall.Syscall6(procBs2GetAllScheduledLockUnlockZoneStatus, uintptr(4), context, uintptr(deviceId), uintptr(unsafe.Pointer(&zoneObjPtr)), uintptr(unsafe.Pointer(&numZone)), uintptr(0), uintptr(0))
	defer ReleaseObject(uintptr(unsafe.Pointer(zoneObjPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2ZoneStatus, numZone)

	var i uintptr = 0
	for ; i < (uintptr)(numZone); i++ {
		temp[i] = *(*BS2ZoneStatus)(unsafe.Pointer(uintptr(unsafe.Pointer(zoneObjPtr)) + i*unsafe.Sizeof(*zoneObjPtr)))
	}
	*zoneStatusObj = temp
	return int16(ret)
}
func SetScheduledLockUnlockZone(context SdkContext, deviceId BS2_DEVICE_ID, zoneObj *[]BS2ScheduledLockUnlockZone) int16 {
	var sizeOfZoneStatus uint32 = uint32(len(*zoneObj))
	var zoneObjPtr uintptr
	if 0 != sizeOfZoneStatus {
		zoneObjPtr = uintptr(unsafe.Pointer(&(*zoneObj)[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2SetScheduledLockUnlockZone, uintptr(4), context, uintptr(deviceId), zoneObjPtr, uintptr(sizeOfZoneStatus), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func SetScheduledLockUnlockZoneAlarm(context SdkContext, deviceId BS2_DEVICE_ID, alarmed uint8, zoneIds []uint32) int16 {
	var zoneIdsPtr uintptr
	if 0 != len(zoneIds) {
		zoneIdsPtr = uintptr(unsafe.Pointer(&zoneIds[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2SetScheduledLockUnlockZoneAlarm, uintptr(5), context, uintptr(deviceId), uintptr(alarmed), zoneIdsPtr, uintptr(len(zoneIds)), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveScheduledLockUnlockZone(context SdkContext, deviceId BS2_DEVICE_ID, zoneIds []uint32) int16 {
	var zoneIdsPtr uintptr
	if 0 != len(zoneIds) {
		zoneIdsPtr = uintptr(unsafe.Pointer(&zoneIds[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2RemoveScheduledLockUnlockZone, uintptr(4), context, uintptr(deviceId), zoneIdsPtr, uintptr(len(zoneIds)), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveAllScheduledLockUnlockZone(context SdkContext, deviceId BS2_DEVICE_ID) int16 {
	ret, _, err := syscall.Syscall(procBs2RemoveAllScheduledLockUnlockZone, uintptr(2), context, uintptr(deviceId), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

const (
	BS2_MAX_READERS_IN_INTRUSION_ALARM_ZONE = 128
	BS2_MAX_INPUTS_IN_INTRUSION_ALARM_ZONE  = 128
	BS2_MAX_OUTPUTS_IN_INTRUSION_ALARM_ZONE = 128
	BS2_MAX_CARDS_IN_INTRUSION_ALARM_ZONE   = 128
	BS2_MAX_DOORS_IN_INTRUSION_ALARM_ZONE   = 128 / 2
	BS2_MAX_GROUPS_IN_INTRUSION_ALARM_ZONE  = 128
)

const (
	INTRUSION_ALARM_ZONE_INPUT_MASK_NONE = 0x00
	INTRUSION_ALARM_ZONE_INPUT_MASK_CARD = 0x01
	INTRUSION_ALARM_ZONE_INPUT_MASK_KEY  = 0x02
	INTRUSION_ALARM_ZONE_INPUT_MASK_ALL  = 0xFF
)

const (
	INTRUSION_ALARM_ZONE_OPERATION_MASK_NONE        = 0x00
	INTRUSION_ALARM_ZONE_OPERATION_MASK_ARM         = 0x01
	INTRUSION_ALARM_ZONE_OPERATION_MASK_DISARM      = 0x02
	INTRUSION_ALARM_ZONE_OPERATION_MASK_TOGGLE      = INTRUSION_ALARM_ZONE_OPERATION_MASK_ARM | INTRUSION_ALARM_ZONE_OPERATION_MASK_DISARM
	INTRUSION_ALARM_ZONE_OPERATION_MASK_ALARM       = 0x04
	INTRUSION_ALARM_ZONE_OPERATION_MASK_ALARM_CLEAR = 0x08
)

type BS2AlarmZoneMember struct {
	DeviceID      BS2_DEVICE_ID ///< 4 bytes
	InputType     uint8         ///< 1 byte  - INTRUSION_ALARM_ZONE_INPUT_MASK_CARD | INTRUSION_ALARM_ZONE_INPUT_MASK_KEY
	OperationType uint8         ///< 1 byte - INTRUSION_ALARM_ZONE_OPERATION_MASK_[ARM|DISARM|TOGGLE]
	Reserved      [2]uint8      ///< 2 bytes (packing)
}

type BS2AlarmZoneInput struct {
	DeviceID BS2_DEVICE_ID ///< 4 bytes

	Port       uint8           ///< 1 byte
	SwitchType BS2_SWITCH_TYPE ///< 1 byte
	Duration   uint16          ///< 2 bytes

	OperationType uint8    ///< 1 byte - INTRUSION_ALARM_ZONE_OPERATION_MASK_*
	Reserved      [3]uint8 /// 3 bytes (packing)
}

type BS2AlarmZoneOutput struct {
	Event    BS2_EVENT_CODE ///< 2 byte - BS2_EVENT_ZONE_INTRUSION_ALARM_[VIOLATION|(DIS)ARM(ED|_FAIL)|ALARM(_[INPUT|CLEAR])]
	Reserved [2]uint8       /// 2 bytes (packing)
	Action   BS2Action
}

type BS2IntrusionAlarmZone struct {
	ZoneID BS2_ZONE_ID                 ///< 4 bytes
	Name   [BS2_MAX_ZONE_NAME_LEN]byte ///< 48 * 3 bytes

	ArmDelay   uint8 ///< 1 byte
	AlarmDelay uint8 ///< 1 byte
	Disabled   bool  ///< 1 byte
	Reserved   uint8 ///< 1 byte (packing)

	NumReaders uint8 ///< 1 byte
	NumInputs  uint8 ///< 1 byte
	NumOutputs uint8 ///< 1 byte
	NumCards   uint8 ///< 1 byte
	NumDoors   uint8 ///< 1 byte
	NumGroups  uint8 ///< 1 byte

	Reserved2 [10]uint8 ///< 10 bytes (packing)

}

type BS2IntrusionAlarmZoneBlob struct {
	IntrusionAlarmZone BS2IntrusionAlarmZone
	MemberObjs         uintptr
	InputObjs          uintptr
	OutputObjs         uintptr
	CardObjs           uintptr
	DoorIDs            uintptr
	GroupIDs           uintptr
}

func GetIntrusionAlarmZone(context SdkContext, deviceId BS2_DEVICE_ID, zoneBlob *[]BS2IntrusionAlarmZoneBlob) int16 {
	var zoneObjPtr *BS2IntrusionAlarmZoneBlob
	var numZone uint32
	ret, _, err := syscall.Syscall6(procBs2GetIntrusionAlarmZone, uintptr(4), context, uintptr(deviceId), uintptr(unsafe.Pointer(&zoneObjPtr)), uintptr(unsafe.Pointer(&numZone)), uintptr(0), uintptr(0))
	defer ReleaseObject(uintptr(unsafe.Pointer(zoneObjPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2IntrusionAlarmZoneBlob, numZone)

	var i uintptr = 0
	for ; i < (uintptr)(numZone); i++ {
		temp[i] = *(*BS2IntrusionAlarmZoneBlob)(unsafe.Pointer(uintptr(unsafe.Pointer(zoneObjPtr)) + i*unsafe.Sizeof(*zoneObjPtr)))
	}
	*zoneBlob = temp
	return int16(ret)
}
func GetIntrusionAlarmZoneStatus(context SdkContext, deviceId BS2_DEVICE_ID, zoneIds []uint32, zoneStatusObj *[]BS2ZoneStatus) int16 {
	var zoneIdsPtr uintptr
	if 0 != len(zoneIds) {
		zoneIdsPtr = uintptr(unsafe.Pointer(&zoneIds[0]))
	}

	var zoneObjPtr *BS2ZoneStatus
	var numZone uint32
	ret, _, err := syscall.Syscall6(procBs2GetIntrusionAlarmZoneStatus, uintptr(6), context, uintptr(deviceId), zoneIdsPtr, uintptr(len(zoneIds)), uintptr(unsafe.Pointer(&zoneObjPtr)), uintptr(unsafe.Pointer(&numZone)))
	defer ReleaseObject(uintptr(unsafe.Pointer(zoneObjPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2ZoneStatus, numZone)

	var i uintptr = 0
	for ; i < (uintptr)(numZone); i++ {
		temp[i] = *(*BS2ZoneStatus)(unsafe.Pointer(uintptr(unsafe.Pointer(zoneObjPtr)) + i*unsafe.Sizeof(*zoneObjPtr)))
	}
	*zoneStatusObj = temp
	return int16(ret)
}
func GetAllIntrusionAlarmZoneStatus(context SdkContext, deviceId BS2_DEVICE_ID, zoneStatusObj *[]BS2ZoneStatus) int16 {

	var zoneObjPtr *BS2ZoneStatus
	var numZone uint32

	ret, _, err := syscall.Syscall6(procBs2GetAllIntrusionAlarmZoneStatus, uintptr(4), context, uintptr(deviceId), uintptr(unsafe.Pointer(&zoneObjPtr)), uintptr(unsafe.Pointer(&numZone)), uintptr(0), uintptr(0))

	defer ReleaseObject(uintptr(unsafe.Pointer(zoneObjPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2ZoneStatus, numZone)

	var i uintptr = 0
	for ; i < (uintptr)(numZone); i++ {
		temp[i] = *(*BS2ZoneStatus)(unsafe.Pointer(uintptr(unsafe.Pointer(zoneObjPtr)) + i*unsafe.Sizeof(*zoneObjPtr)))
	}
	*zoneStatusObj = temp
	return int16(ret)
}
func SetIntrusionAlarmZone(context SdkContext, deviceId BS2_DEVICE_ID, zoneObj *[]BS2IntrusionAlarmZoneBlob) int16 {
	var sizeOfZoneObj uint32 = uint32(len(*zoneObj))
	var zoneObjPtr uintptr
	if 0 != sizeOfZoneObj {
		zoneObjPtr = uintptr(unsafe.Pointer(&(*zoneObj)[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2SetIntrusionAlarmZone, uintptr(4), context, uintptr(deviceId), zoneObjPtr, uintptr(sizeOfZoneObj), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func SetIntrusionAlarmZoneAlarm(context SdkContext, deviceId BS2_DEVICE_ID, alarmed uint8, zoneIds []uint32) int16 {
	var zoneIdsPtr uintptr
	if 0 != len(zoneIds) {
		zoneIdsPtr = uintptr(unsafe.Pointer(&zoneIds[0]))
	}

	ret, _, err := syscall.Syscall6(procBs2SetIntrusionAlarmZoneAlarm, uintptr(5), context, uintptr(deviceId), uintptr(alarmed), zoneIdsPtr, uintptr(len(zoneIds)), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveIntrusionAlarmZone(context SdkContext, deviceId BS2_DEVICE_ID, zoneIds []uint32) int16 {
	var zoneIdsPtr uintptr
	if 0 != len(zoneIds) {
		zoneIdsPtr = uintptr(unsafe.Pointer(&zoneIds[0]))
	}

	ret, _, err := syscall.Syscall6(procBs2RemoveIntrusionAlarmZone, uintptr(4), context, uintptr(deviceId), zoneIdsPtr, uintptr(len(zoneIds)), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveAllIntrusionAlarmZone(context SdkContext, deviceId BS2_DEVICE_ID) int16 {
	ret, _, err := syscall.Syscall(procBs2RemoveAllIntrusionAlarmZone, uintptr(2), context, uintptr(deviceId), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func GetInterlockZone(context SdkContext, deviceId BS2_DEVICE_ID, zoneBlob *[]BS2InterlockZoneBlob) int16 {
	var zoneObjPtr *BS2InterlockZoneBlob
	var numZone uint32

	ret, _, err := syscall.Syscall6(procBs2GetInterlockZone, uintptr(4), context, uintptr(deviceId), uintptr(unsafe.Pointer(&zoneObjPtr)), uintptr(unsafe.Pointer(&numZone)), uintptr(0), uintptr(0))
	defer ReleaseObject(uintptr(unsafe.Pointer(zoneObjPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2InterlockZoneBlob, numZone)

	var i uintptr = 0
	for ; i < (uintptr)(numZone); i++ {
		temp[i] = *(*BS2InterlockZoneBlob)(unsafe.Pointer(uintptr(unsafe.Pointer(zoneObjPtr)) + i*unsafe.Sizeof(*zoneObjPtr)))
	}
	*zoneBlob = temp
	return int16(ret)
}
func GetInterlockZoneStatus(context SdkContext, deviceId BS2_DEVICE_ID, zoneIds []uint32, zoneStatusObj *[]BS2ZoneStatus) int16 {
	var zoneIdsPtr uintptr
	if 0 != len(zoneIds) {
		zoneIdsPtr = uintptr(unsafe.Pointer(&zoneIds[0]))
	}

	var zoneObjPtr *BS2ZoneStatus
	var numZone uint32

	ret, _, err := syscall.Syscall6(procBs2GetInterlockZoneStatus, uintptr(6), context, uintptr(deviceId), zoneIdsPtr, uintptr(len(zoneIds)), uintptr(unsafe.Pointer(&zoneObjPtr)), uintptr(unsafe.Pointer(&numZone)))

	defer ReleaseObject(uintptr(unsafe.Pointer(zoneObjPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2ZoneStatus, numZone)

	var i uintptr = 0
	for ; i < (uintptr)(numZone); i++ {
		temp[i] = *(*BS2ZoneStatus)(unsafe.Pointer(uintptr(unsafe.Pointer(zoneObjPtr)) + i*unsafe.Sizeof(*zoneObjPtr)))
	}
	*zoneStatusObj = temp
	return int16(ret)
}
func GetAllInterlockZoneStatus(context SdkContext, deviceId BS2_DEVICE_ID, zoneStatusObj *[]BS2ZoneStatus) int16 {
	var zoneObjPtr *BS2ZoneStatus
	var numZone uint32

	ret, _, err := syscall.Syscall6(procBs2GetAllInterlockZoneStatus, uintptr(4), context, uintptr(deviceId), uintptr(unsafe.Pointer(&zoneObjPtr)), uintptr(unsafe.Pointer(&numZone)), uintptr(0), uintptr(0))
	defer ReleaseObject(uintptr(unsafe.Pointer(zoneObjPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2ZoneStatus, numZone)

	var i uintptr = 0
	for ; i < (uintptr)(numZone); i++ {
		temp[i] = *(*BS2ZoneStatus)(unsafe.Pointer(uintptr(unsafe.Pointer(zoneObjPtr)) + i*unsafe.Sizeof(*zoneObjPtr)))
	}
	*zoneStatusObj = temp
	return int16(ret)
}

type BS2InterlockZone struct {
	ZoneID     BS2_ZONE_ID
	Name       [BS2_MAX_ZONE_NAME_LEN]byte
	Disabled   bool
	NumInputs  uint8
	NumOutputs uint8
	NumDoors   uint8
	reserved   [8]byte
}

type BS2InterlockZoneBlob struct {
	InterlockZone BS2InterlockZone
	inputObjs     uintptr
	outputObjs    uintptr
	doorIDs       uintptr
}

func SetInterlockZone(context SdkContext, deviceId BS2_DEVICE_ID, zoneObj *[]BS2InterlockZoneBlob) int16 {
	var sizeOfZoneObj uint32 = uint32(len(*zoneObj))
	var zoneObjPtr uintptr
	if 0 != sizeOfZoneObj {
		zoneObjPtr = uintptr(unsafe.Pointer(&(*zoneObj)[0]))
	}

	ret, _, err := syscall.Syscall6(procBs2SetInterlockZone, uintptr(4), context, uintptr(deviceId), zoneObjPtr, uintptr(sizeOfZoneObj), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func SetInterlockZoneAlarm(context SdkContext, deviceId BS2_DEVICE_ID, alarmed uint8, zoneIds []uint32) int16 {
	var zoneIdsPtr uintptr
	if 0 != len(zoneIds) {
		zoneIdsPtr = uintptr(unsafe.Pointer(&zoneIds[0]))
	}

	ret, _, err := syscall.Syscall6(procBs2SetInterlockZoneAlarm, uintptr(5), context, uintptr(deviceId), uintptr(alarmed), zoneIdsPtr, uintptr(len(zoneIds)), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveInterlockZone(context SdkContext, deviceId BS2_DEVICE_ID, zoneIds []uint32) int16 {
	var zoneIdsPtr uintptr
	if 0 != len(zoneIds) {
		zoneIdsPtr = uintptr(unsafe.Pointer(&zoneIds[0]))
	}

	ret, _, err := syscall.Syscall6(procBs2RemoveInterlockZone, uintptr(4), context, uintptr(deviceId), zoneIdsPtr, uintptr(len(zoneIds)), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveAllInterlockZone(context SdkContext, deviceId BS2_DEVICE_ID) int16 {
	ret, _, err := syscall.Syscall(procBs2RemoveAllInterlockZone, uintptr(2), context, uintptr(deviceId), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

type BS2DeviceZone struct {
	ZoneID   uint32
	ZoneType uint8
	NodeType uint8
	Enable   uint8
	Reserved uint8

	MemberData [884]byte
}

func GetDeviceZone(context SdkContext, deviceId BS2_DEVICE_ID, ids []BS2_DEVICE_ZONE_TABLE_ID, zoneObj *[]BS2DeviceZone) int16 {
	var zoneIdsPtr uintptr
	if 0 != len(ids) {
		zoneIdsPtr = uintptr(unsafe.Pointer(&ids[0]))
	}

	var zoneObjPtr *BS2DeviceZone
	var numZone uint32

	ret, _, err := syscall.Syscall6(procBs2GetDeviceZone, uintptr(6), context, uintptr(deviceId), zoneIdsPtr, uintptr(len(ids)), uintptr(unsafe.Pointer(&zoneObjPtr)), uintptr(unsafe.Pointer(&numZone)))
	defer ReleaseObject(uintptr(unsafe.Pointer(zoneObjPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2DeviceZone, numZone)

	var i uintptr = 0
	for ; i < (uintptr)(numZone); i++ {
		temp[i] = *(*BS2DeviceZone)(unsafe.Pointer(uintptr(unsafe.Pointer(zoneObjPtr)) + i*unsafe.Sizeof(*zoneObjPtr)))
	}
	*zoneObj = temp
	return int16(ret)
}
func GetAllDeviceZone(context SdkContext, deviceId BS2_DEVICE_ID, zoneObj *[]BS2DeviceZone) int16 {

	var zoneObjPtr *BS2DeviceZone
	var numZone uint32
	ret, _, err := syscall.Syscall6(procBs2GetAllDeviceZone, uintptr(4), context, uintptr(deviceId), uintptr(unsafe.Pointer(&zoneObjPtr)), uintptr(unsafe.Pointer(&numZone)), uintptr(0), uintptr(0))
	defer ReleaseObject(uintptr(unsafe.Pointer(zoneObjPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2DeviceZone, numZone)

	var i uintptr = 0
	for ; i < (uintptr)(numZone); i++ {
		temp[i] = *(*BS2DeviceZone)(unsafe.Pointer(uintptr(unsafe.Pointer(zoneObjPtr)) + i*unsafe.Sizeof(*zoneObjPtr)))
	}
	*zoneObj = temp
	return int16(ret)
}
func SetDeviceZone(context SdkContext, deviceId BS2_DEVICE_ID, zoneObj *[]BS2DeviceZone) int16 {
	var sizeOfZoneObj uint32 = uint32(len(*zoneObj))
	var zoneObjPtr uintptr
	if 0 != sizeOfZoneObj {
		zoneObjPtr = uintptr(unsafe.Pointer(&(*zoneObj)[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2SetDeviceZone, uintptr(4), context, uintptr(deviceId), zoneObjPtr, uintptr(sizeOfZoneObj), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

type BS2_DEVICE_ZONE_TABLE_ID struct {
	ZoneID   BS2_ZONE_ID
	NodeType uint32
}

func RemoveDeviceZone(context SdkContext, deviceId BS2_DEVICE_ID, ids []BS2_DEVICE_ZONE_TABLE_ID) int16 {
	var zoneIdsPtr uintptr
	if 0 != len(ids) {
		zoneIdsPtr = uintptr(unsafe.Pointer(&ids[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2RemoveDeviceZone, uintptr(4), context, uintptr(deviceId), zoneIdsPtr, uintptr(len(ids)), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveAllDeviceZone(context SdkContext, deviceId BS2_DEVICE_ID) int16 {
	ret, _, err := syscall.Syscall(procBs2RemoveAllDeviceZone, uintptr(2), context, uintptr(deviceId), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func SetDeviceZoneAlarm(context SdkContext, deviceId BS2_DEVICE_ID, alarmed uint8, zoneIds []BS2_ZONE_ID) int16 {
	var sizeOfZoneIds uint32 = uint32(len(zoneIds))
	var zoneIdsPtr uintptr
	if 0 != len(zoneIds) {
		zoneIdsPtr = uintptr(unsafe.Pointer(&zoneIds[0]))
	}

	ret, _, err := syscall.Syscall6(procBs2SetDeviceZoneAlarm, uintptr(5), context, uintptr(deviceId), uintptr(alarmed), zoneIdsPtr, uintptr(sizeOfZoneIds), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func ClearDeviceZoneAccessRecord(context SdkContext, deviceId BS2_DEVICE_ID, zoneId uint32, userIds *[]BS2_USER_ID) int16 {
	var sizeOfuserIds uint32 = uint32(len(*userIds))
	var userIdsPtr uintptr
	if 0 != len(*userIds) {
		userIdsPtr = uintptr(unsafe.Pointer(&(*userIds)[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2ClearDeviceZoneAccessRecord, uintptr(5), context, uintptr(deviceId), uintptr(zoneId), userIdsPtr, uintptr(sizeOfuserIds), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func ClearAllDeviceZoneAccessRecord(context SdkContext, deviceId BS2_DEVICE_ID, zoneId uint32) int16 {
	ret, _, err := syscall.Syscall(procBs2ClearAllDeviceZoneAccessRecord, uintptr(3), context, uintptr(deviceId), uintptr(zoneId))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

/**
 *  Constants
 */
const (
	BS2_MAX_READERS_PER_DEVICE_ZONE_ENTRANCE_LIMIT        = 64
	BS2_MAX_BYPASS_GROUPS_PER_DEVICE_ZONE_ENTRANCE_LIMIT  = 16
	BS2_MAX_DEVICE_ZONE_ENTRANCE_LIMIT_ALARM_ACTION       = 5
	BS2_MAX_ENTRANCE_LIMIT_PER_ZONE                       = 24
	BS2_MAX_ACCESS_GROUP_ENTRANCE_LIMIT_PER_ENTRACE_LIMIT = 16
	BS2_ENTRY_COUNT_FOR_ACCESS_GROUP_ENTRANCE_LIMIT       = -2
	BS2_OTHERWISE_ACCESS_GROUP_ID                         = -1
	BS2_ENTRY_COUNT_NO_LIMIT                              = -1
)

/**
 *  BS2_DEVICE_ZONE_ENTRANCE_LIMIT_TYPE
 */
const (
	BS2_DEVICE_ZONE_ENTRANCE_LIMIT_SOFT = 0x01
	BS2_DEVICE_ZONE_ENTRANCE_LIMIT_HARD = 0x02
)

type BS2_DEVICE_ZONE_ENTRANCE_LIMIT_TYPE = uint8

/**
 *  BS2_DEVICE_ZONE_ENTRANCE_LIMIT_DISCONNECTED_ACTION_TYPE
 */
const (
	BS2_DEVICE_ZONE_ENTRANCE_LIMIT_DISCONNECTED_ACTION_SOFT = 0x01
	BS2_DEVICE_ZONE_ENTRANCE_LIMIT_DISCONNECTED_ACTION_HARD = 0x02
)

type BS2_DEVICE_ZONE_ENTRANCE_LIMIT_DISCONNECTED_ACTION_TYPE = uint8

type BS2DeviceZoneEntranceLimitMemberInfo struct {
	ReaderID BS2_DEVICE_ID ///< 4 bytes
}

/**
 *  BS2DeviceZoneEntranceLimitMaster
 */

type BS2DeviceZoneEntranceLimitMaster struct {
	Name [BS2_MAX_ZONE_NAME_LEN]byte ///< 48 * 3 bytes

	Type      BS2_DEVICE_ZONE_ENTRANCE_LIMIT_TYPE ///< 1 byte
	Reserved1 [3]uint8                            ///< 3 bytes (packing)

	EntryLimitInterval_s uint32 ///< 4 bytes: in seconds 0: no limit

	NumEntranceLimit uint8 ///< 1 byte
	NumReaders       uint8 ///< 1 byte
	NumAlarm         uint8 ///< 1 byte
	NumBypassGroups  uint8 ///< 1 byte

	MaxEntry      [BS2_MAX_ENTRANCE_LIMIT_PER_ZONE]uint8  ///< 1 * 24 bytes // 0 (always limit)
	PeriodStart_s [BS2_MAX_ENTRANCE_LIMIT_PER_ZONE]uint32 ///< 4 * 24 bytes: in seconds
	PeriodEnd_s   [BS2_MAX_ENTRANCE_LIMIT_PER_ZONE]uint32 ///< 4 * 24 bytes: in seconds

	Readers        [BS2_MAX_READERS_PER_DEVICE_ZONE_ENTRANCE_LIMIT]BS2DeviceZoneEntranceLimitMemberInfo ///< 4 * 64 bytes
	Alarm          [BS2_MAX_DEVICE_ZONE_ENTRANCE_LIMIT_ALARM_ACTION]BS2Action                           ///< 32 * 5 bytes
	BypassGroupIDs [BS2_MAX_BYPASS_GROUPS_PER_DEVICE_ZONE_ENTRANCE_LIMIT]BS2_ACCESS_GROUP_ID            ///< 4 * 16 bytes

	Reserved3 [8 * 4]uint8 ///< 8 * 4 bytes
} ///884 bytes

type BS2DeviceZoneAGEntranceLimit struct {
	ZoneID             BS2_ZONE_ID
	NumAGEntranceLimit uint16
	Reserved1          uint16
	PeriodStart_s      [BS2_MAX_ENTRANCE_LIMIT_PER_ZONE]uint32
	PeriodEnd_s        [BS2_MAX_ENTRANCE_LIMIT_PER_ZONE]uint32
	NumEntry           [BS2_MAX_ENTRANCE_LIMIT_PER_ZONE]uint16
	MaxEntry           [BS2_MAX_ENTRANCE_LIMIT_PER_ZONE][BS2_MAX_ACCESS_GROUP_ENTRANCE_LIMIT_PER_ENTRACE_LIMIT]uint16
	AccessGroupID      [BS2_MAX_ENTRANCE_LIMIT_PER_ZONE][BS2_MAX_ACCESS_GROUP_ENTRANCE_LIMIT_PER_ENTRACE_LIMIT]BS2_ACCESS_GROUP_ID
}

type BS2_IPV4_ADDR = [16]byte

/**
 *  BS2DeviceZoneEntranceLimitMember
 */
type BS2DeviceZoneEntranceLimitMember struct {
	MasterPort         BS2_PORT                                                ///< 2 bytes
	ActionInDisconnect BS2_DEVICE_ZONE_ENTRANCE_LIMIT_DISCONNECTED_ACTION_TYPE ///< 1 byte
	Reserved1          uint8                                                   ///< 1 byte (packing)
	MasterIP           BS2_IPV4_ADDR                                           ///< 16 bytes
} ///20 bytes

func GetAccessGroupEntranceLimit(context SdkContext, deviceId BS2_DEVICE_ID, zoneIds []BS2_ZONE_ID, zoneObj *[]BS2DeviceZoneAGEntranceLimit) int16 {
	var zoneIdsPtr uintptr
	if 0 != len(zoneIds) {
		zoneIdsPtr = uintptr(unsafe.Pointer(&zoneIds[0]))
	}

	var zoneObjPtr *BS2DeviceZoneAGEntranceLimit
	var numZone uint32

	ret, _, err := syscall.Syscall6(procBs2GetAccessGroupEntranceLimit, uintptr(6), context, uintptr(deviceId), zoneIdsPtr, uintptr(len(zoneIds)), uintptr(unsafe.Pointer(&zoneObjPtr)), uintptr(unsafe.Pointer(&numZone)))
	defer ReleaseObject(uintptr(unsafe.Pointer(zoneObjPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2DeviceZoneAGEntranceLimit, numZone)

	var i uintptr = 0
	for ; i < (uintptr)(numZone); i++ {
		temp[i] = *(*BS2DeviceZoneAGEntranceLimit)(unsafe.Pointer(uintptr(unsafe.Pointer(zoneObjPtr)) + i*unsafe.Sizeof(*zoneObjPtr)))
	}
	*zoneObj = temp
	return int16(ret)
}
func GetAllAccessGroupEntranceLimit(context SdkContext, deviceId BS2_DEVICE_ID, zoneObj *[]BS2DeviceZoneAGEntranceLimit) int16 {
	var zoneObjPtr *BS2DeviceZoneAGEntranceLimit
	var numZone uint32
	ret, _, err := syscall.Syscall6(procBs2GetAllAccessGroupEntranceLimit, uintptr(4), context, uintptr(deviceId), uintptr(unsafe.Pointer(&zoneObjPtr)), uintptr(unsafe.Pointer(&numZone)), uintptr(0), uintptr(0))
	defer ReleaseObject(uintptr(unsafe.Pointer(zoneObjPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2DeviceZoneAGEntranceLimit, numZone)

	var i uintptr = 0
	for ; i < (uintptr)(numZone); i++ {
		temp[i] = *(*BS2DeviceZoneAGEntranceLimit)(unsafe.Pointer(uintptr(unsafe.Pointer(zoneObjPtr)) + i*unsafe.Sizeof(*zoneObjPtr)))
	}
	*zoneObj = temp
	return int16(ret)
}
func SetAccessGroupEntranceLimit(context SdkContext, deviceId BS2_DEVICE_ID, zoneObj *[]BS2DeviceZoneAGEntranceLimit) int16 {
	var sizeOfZoneObj uint32 = uint32(len(*zoneObj))
	var zoneObjPtr uintptr
	if 0 != sizeOfZoneObj {
		zoneObjPtr = uintptr(unsafe.Pointer(&(*zoneObj)[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2SetAccessGroupEntranceLimit, uintptr(4), context, uintptr(deviceId), zoneObjPtr, uintptr(sizeOfZoneObj), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveAccessGroupEntranceLimit(context SdkContext, deviceId BS2_DEVICE_ID, zoneIds []BS2_ZONE_ID) int16 {
	var zoneIdsPtr uintptr
	if 0 != len(zoneIds) {
		zoneIdsPtr = uintptr(unsafe.Pointer(&zoneIds[0]))
	}

	ret, _, err := syscall.Syscall6(procBs2RemoveAccessGroupEntranceLimit, uintptr(4), context, uintptr(deviceId), zoneIdsPtr, uintptr(len(zoneIds)), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveAllAccessGroupEntranceLimit(context SdkContext, deviceId BS2_DEVICE_ID) int16 {
	ret, _, err := syscall.Syscall(procBs2RemoveAllAccessGroupEntranceLimit, uintptr(2), context, uintptr(deviceId), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func GetAllDeviceZoneAGEntranceLimit(context SdkContext, deviceId BS2_DEVICE_ID, zoneObj *[]BS2DeviceZoneAGEntranceLimit) int16 {
	var zoneObjPtr *BS2DeviceZoneAGEntranceLimit
	var numZone uint32

	ret, _, err := syscall.Syscall6(procBs2GetAllDeviceZoneAGEntranceLimit, uintptr(4), context, uintptr(deviceId), uintptr(unsafe.Pointer(&zoneObjPtr)), uintptr(unsafe.Pointer(&numZone)), uintptr(0), uintptr(0))
	defer ReleaseObject(uintptr(unsafe.Pointer(zoneObjPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2DeviceZoneAGEntranceLimit, numZone)

	var i uintptr = 0
	for ; i < (uintptr)(numZone); i++ {
		temp[i] = *(*BS2DeviceZoneAGEntranceLimit)(unsafe.Pointer(uintptr(unsafe.Pointer(zoneObjPtr)) + i*unsafe.Sizeof(*zoneObjPtr)))
	}
	*zoneObj = temp
	return int16(ret)
}
func SetDeviceZoneAGEntranceLimit(context SdkContext, deviceId BS2_DEVICE_ID, zoneObj *[]BS2DeviceZoneAGEntranceLimit) int16 {
	var sizeOfZoneObj uint32 = uint32(len(*zoneObj))
	var zoneObjPtr uintptr
	if 0 != sizeOfZoneObj {
		zoneObjPtr = uintptr(unsafe.Pointer(&(*zoneObj)[0]))
	}
	ret, _, err := syscall.Syscall6(procBs2SetDeviceZoneAGEntranceLimit, uintptr(4), context, uintptr(deviceId), zoneObjPtr, uintptr(sizeOfZoneObj), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveDeviceZoneAGEntranceLimit(context SdkContext, deviceId BS2_DEVICE_ID, zoneIds []BS2_ZONE_ID) int16 {
	var zoneIdsPtr uintptr
	if 0 != len(zoneIds) {
		zoneIdsPtr = uintptr(unsafe.Pointer(&zoneIds[0]))
	}

	ret, _, err := syscall.Syscall6(procBs2RemoveDeviceZoneAGEntranceLimit, uintptr(4), context, uintptr(deviceId), zoneIdsPtr, uintptr(len(zoneIds)), uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func RemoveAllDeviceZoneAGEntranceLimit(context SdkContext, deviceId BS2_DEVICE_ID) int16 {
	ret, _, err := syscall.Syscall(procBs2RemoveAllDeviceZoneAGEntranceLimit, uintptr(2), context, uintptr(deviceId), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func GetUserDatabaseInfoFromDir(context SdkContext, filePathName string, numUsers *uint32, numCards *uint32, numFingers *uint32, numFaces *uint32, fnIsAcceptableUserID interface{}) int16 {
	var fnIsAcceptableUserIDPtr uintptr
	if nil != fnIsAcceptableUserID {
		fnIsAcceptableUserIDPtr = syscall.NewCallback(fnIsAcceptableUserID)
	}

	ret, _, err := syscall.Syscall9(procBs2GetUserDatabaseInfoFromDir, uintptr(7), context, uintptr(unsafe.Pointer(syscall.StringBytePtr(filePathName))), uintptr(unsafe.Pointer(numUsers)), uintptr(unsafe.Pointer(numCards)), uintptr(unsafe.Pointer(numFingers)), uintptr(unsafe.Pointer(numFaces)), fnIsAcceptableUserIDPtr, uintptr(0), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}

func GetUserListFromDir(context SdkContext, filePathName string, userIds *[]BS2_USER_ID, fnIsAcceptableUserID interface{}) int16 {
	var fnIsAcceptableUserIDPtr uintptr
	if nil != fnIsAcceptableUserID {
		fnIsAcceptableUserIDPtr = syscall.NewCallback(fnIsAcceptableUserID)
	}

	var userIdsPtr *BS2_USER_ID
	var numUserId uint32
	ret, _, err := syscall.Syscall6(procBs2GetUserDatabaseInfoFromDir, uintptr(5), context, uintptr(unsafe.Pointer(syscall.StringBytePtr(filePathName))), uintptr(unsafe.Pointer(&userIdsPtr)), uintptr(unsafe.Pointer(&numUserId)), fnIsAcceptableUserIDPtr, uintptr(0))

	defer ReleaseObject(uintptr(unsafe.Pointer(userIdsPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2_USER_ID, numUserId)

	var i uintptr = 0
	for ; i < (uintptr)(numUserId); i++ {
		temp[i] = *(*BS2_USER_ID)(unsafe.Pointer(uintptr(unsafe.Pointer(userIdsPtr)) + i*unsafe.Sizeof(*userIdsPtr)))
	}
	*userIds = temp
	return int16(ret)
}

func GetUserInfosFromDir(context SdkContext, filePathName string, userIds []BS2_USER_ID, userObj *[]BS2UserBlob) int16 {
	var userIdsPtr uintptr
	if 0 != len(userIds) {
		userIdsPtr = uintptr(unsafe.Pointer(&userIds[0]))
	}

	var numUsers uint32 = uint32(len(*userObj))

	ret, _, err := syscall.Syscall6(procBs2GetUserDatabaseInfoFromDir, uintptr(5), context, uintptr(unsafe.Pointer(syscall.StringBytePtr(filePathName))), userIdsPtr, uintptr(numUsers), uintptr(unsafe.Pointer(&(*userObj)[0])), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func GetUserDatasFromDir(context SdkContext, filePathName string, userIds []BS2_USER_ID, userObj *[]BS2UserBlob, mask uint32) int16 {
	var userIdsPtr uintptr
	if 0 != len(userIds) {
		userIdsPtr = uintptr(unsafe.Pointer(&userIds[0]))
	}

	var numUsers uint32 = uint32(len(*userObj))

	ret, _, err := syscall.Syscall6(procBs2GetUserDatasFromDir, uintptr(6), context, uintptr(unsafe.Pointer(syscall.StringBytePtr(filePathName))), userIdsPtr, uintptr(numUsers), uintptr(unsafe.Pointer(&(*userObj)[0])), uintptr(mask))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func GetUserInfosExFromDir(context SdkContext, filePathName string, userIds []BS2_USER_ID, userObj *[]BS2UserBlobEx) int16 {
	var userIdsPtr uintptr
	if 0 != len(userIds) {
		userIdsPtr = uintptr(unsafe.Pointer(&userIds[0]))
	}

	var numUsers uint32 = uint32(len(*userObj))

	ret, _, err := syscall.Syscall6(procBs2GetUserInfosExFromDir, uintptr(5), context, uintptr(unsafe.Pointer(syscall.StringBytePtr(filePathName))), userIdsPtr, uintptr(numUsers), uintptr(unsafe.Pointer(&(*userObj)[0])), uintptr(0))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func GetUserDatasExFromDir(context SdkContext, filePathName string, userIds []BS2_USER_ID, userObj *[]BS2UserBlobEx, mask uint32) int16 {
	var userIdsPtr uintptr
	if 0 != len(userIds) {
		userIdsPtr = uintptr(unsafe.Pointer(&userIds[0]))
	}

	var numUsers uint32 = uint32(len(*userObj))

	ret, _, err := syscall.Syscall6(procBs2GetUserDatasExFromDir, uintptr(6), context, uintptr(unsafe.Pointer(syscall.StringBytePtr(filePathName))), userIdsPtr, uintptr(numUsers), uintptr(unsafe.Pointer(&(*userObj)[0])), uintptr(mask))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	return int16(ret)
}
func GetLogFromDir(context SdkContext, filePathName string, eventId BS2_EVENT_ID, amount uint32, eventLogs *[]BS2Event) int16 {

	var eventLogsPtr *BS2Event
	var numLogs uint32

	ret, _, err := syscall.Syscall6(procBs2GetLogFromDir, uintptr(6), context, uintptr(unsafe.Pointer(syscall.StringBytePtr(filePathName))), uintptr(eventId), uintptr(amount), uintptr(unsafe.Pointer(&eventLogsPtr)), uintptr(unsafe.Pointer(&numLogs)))
	defer ReleaseObject(uintptr(unsafe.Pointer(eventLogsPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2Event, numLogs)

	var i uintptr = 0
	for ; i < (uintptr)(numLogs); i++ {
		temp[i] = *(*BS2Event)(unsafe.Pointer(uintptr(unsafe.Pointer(eventLogsPtr)) + i*unsafe.Sizeof(*eventLogsPtr)))
	}
	*eventLogs = temp
	return int16(ret)
}

func GetFilteredLogFromDir(context SdkContext, filePathName string, userId BS2_USER_ID, eventCode BS2_EVENT_CODE, startTimestamp uint32, endTimestamp uint32, tnaKey uint8, eventLogs *[]BS2Event) int16 {

	var eventLogsPtr *BS2Event
	var numLogs uint32
	ret, _, err := syscall.Syscall9(procBs2GetFilteredLogFromDir, uintptr(9), context, uintptr(unsafe.Pointer(syscall.StringBytePtr(filePathName))), uintptr(unsafe.Pointer(&userId[0])), uintptr(eventCode), uintptr(startTimestamp), uintptr(endTimestamp), uintptr(tnaKey), uintptr(unsafe.Pointer(&eventLogsPtr)), uintptr(unsafe.Pointer(&numLogs)))
	defer ReleaseObject(uintptr(unsafe.Pointer(eventLogsPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2Event, numLogs)

	var i uintptr = 0
	for ; i < (uintptr)(numLogs); i++ {
		temp[i] = *(*BS2Event)(unsafe.Pointer(uintptr(unsafe.Pointer(eventLogsPtr)) + i*unsafe.Sizeof(*eventLogsPtr)))
	}
	*eventLogs = temp
	return int16(ret)
}

func GetLogBlobFromDir(context SdkContext, filePathName string, eventMask uint16, lastEventId uint32, amount uint32, eventLogs *[]BS2EventBlob) int16 {

	var eventLogsPtr *BS2EventBlob
	var numLogs uint32
	ret, _, err := syscall.Syscall9(procBs2GetLogBlobFromDir, uintptr(7), context, uintptr(unsafe.Pointer(syscall.StringBytePtr(filePathName))), uintptr(eventMask), uintptr(lastEventId), uintptr(amount), uintptr(unsafe.Pointer(&eventLogsPtr)), uintptr(unsafe.Pointer(&numLogs)), uintptr(0), uintptr(0))
	defer ReleaseObject(uintptr(unsafe.Pointer(eventLogsPtr)))
	if 0 != err {
		log.Fatal(syscall.Errno(err))
	}
	temp := make([]BS2EventBlob, numLogs)

	var i uintptr = 0
	for ; i < (uintptr)(numLogs); i++ {
		temp[i] = *(*BS2EventBlob)(unsafe.Pointer(uintptr(unsafe.Pointer(eventLogsPtr)) + i*unsafe.Sizeof(*eventLogsPtr)))
	}
	*eventLogs = temp
	return int16(ret)
}

const (
	BS_SDK_SUCCESS             = 1
	BS_SDK_DURESS_SUCCESS      = 2
	BS_SDK_FIRST_AUTH_SUCCESS  = 3
	BS_SDK_SECOND_AUTH_SUCCESS = 4
	BS_SDK_DUAL_AUTH_SUCCESS   = 5

	// Communication errors
	BS_SDK_ERROR_CANNOT_OPEN_SOCKET      = -101
	BS_SDK_ERROR_CANNOT_CONNECT_SOCKET   = -102
	BS_SDK_ERROR_CANNOT_LISTEN_SOCKET    = -103
	BS_SDK_ERROR_CANNOT_ACCEPT_SOCKET    = -104
	BS_SDK_ERROR_CANNOT_READ_SOCKET      = -105
	BS_SDK_ERROR_CANNOT_WRITE_SOCKET     = -106
	BS_SDK_ERROR_SOCKET_IS_NOT_CONNECTED = -107
	BS_SDK_ERROR_SOCKET_IS_NOT_OPEN      = -108
	BS_SDK_ERROR_SOCKET_IS_NOT_LISTENED  = -109
	BS_SDK_ERROR_SOCKET_IN_PROGRESS      = -110

	// Packet errors
	BS_SDK_ERROR_INVALID_PARAM       = -200
	BS_SDK_ERROR_INVALID_PACKET      = -201
	BS_SDK_ERROR_INVALID_DEVICE_ID   = -202
	BS_SDK_ERROR_INVALID_DEVICE_TYPE = -203
	BS_SDK_ERROR_PACKET_CHECKSUM     = -204
	BS_SDK_ERROR_PACKET_INDEX        = -205
	BS_SDK_ERROR_PACKET_COMMAND      = -206
	BS_SDK_ERROR_PACKET_SEQUENCE     = -207
	BS_SDK_ERROR_NO_PACKET           = -209

	//Fingerprint errors
	BS_SDK_ERROR_EXTRACTION_FAIL            = -300
	BS_SDK_ERROR_VERIFY_FAIL                = -301
	BS_SDK_ERROR_IDENTIFY_FAIL              = -302
	BS_SDK_ERROR_IDENTIFY_TIMEOUT           = -303
	BS_SDK_ERROR_FINGERPRINT_CAPTURE_FAIL   = -304
	BS_SDK_ERROR_FINGERPRINT_SCAN_TIMEOUT   = -305
	BS_SDK_ERROR_FINGERPRINT_SCAN_CANCELLED = -306
	BS_SDK_ERROR_NOT_SAME_FINGERPRINT       = -307
	BS_SDK_ERROR_EXTRACTION_LOW_QUALITY     = -308
	BS_SDK_ERROR_CAPTURE_LOW_QUALITY        = -309
	BS_SDK_ERROR_CANNOT_FIND_FINGERPRINT    = -310
	BS_SDK_ERROR_FAKE_FINGER_DETECTED       = -311

	//File I/O errors
	BS_SDK_ERROR_CANNOT_OPEN_DIR    = -400
	BS_SDK_ERROR_CANNOT_OPEN_FILE   = -401
	BS_SDK_ERROR_CANNOT_WRITE_FILE  = -402
	BS_SDK_ERROR_CANNOT_SEEK_FILE   = -403
	BS_SDK_ERROR_CANNOT_READ_FILE   = -404
	BS_SDK_ERROR_CANNOT_GET_STAT    = -405
	BS_SDK_ERROR_CANNOT_GET_SYSINFO = -406
	BS_SDK_ERROR_DATA_MISMATCH      = -407

	// I/O errors
	BS_SDK_ERROR_INVALID_RELAY             = -500
	BS_SDK_ERROR_CANNOT_WRITE_IO_PACKET    = -501
	BS_SDK_ERROR_CANNOT_READ_IO_PACKET     = -502
	BS_SDK_ERROR_CANNOT_READ_INPUT         = -503
	BS_SDK_ERROR_READ_INPUT_TIMEOUT        = -504
	BS_SDK_ERROR_CANNOT_ENABLE_INPUT       = -505
	BS_SDK_ERROR_CANNOT_SET_INPUT_DURATION = -506
	BS_SDK_ERROR_INVALID_PORT              = -507
	BS_SDK_ERROR_INVALID_INTERPHONE_TYPE   = -508
	BS_SDK_ERROR_INVALID_LCD_PARAM         = -510
	BS_SDK_ERROR_CANNOT_WRITE_LCD_PACKET   = -511
	BS_SDK_ERROR_CANNOT_READ_LCD_PACKET    = -512
	BS_SDK_ERROR_INVALID_LCD_PACKET        = -513
	BS_SDK_ERROR_INPUT_QUEUE_FULL          = -520
	BS_SDK_ERROR_WIEGAND_QUEUE_FULL        = -521
	BS_SDK_ERROR_MISC_INPUT_QUEUE_FULL     = -522
	BS_SDK_ERROR_WIEGAND_DATA_QUEUE_FULL   = -523
	BS_SDK_ERROR_WIEGAND_DATA_QUEUE_EMPTY  = -524

	//Util errors
	BS_SDK_ERROR_NOT_SUPPORTED = -600
	BS_SDK_ERROR_TIMEOUT       = -601

	//Database errors
	BS_SDK_ERROR_INVALID_DATA_FILE           = -700
	BS_SDK_ERROR_TOO_LARGE_DATA_FOR_SLOT     = -701
	BS_SDK_ERROR_INVALID_SLOT_NO             = -702
	BS_SDK_ERROR_INVALID_SLOT_DATA           = -703
	BS_SDK_ERROR_CANNOT_INIT_DB              = -704
	BS_SDK_ERROR_DUPLICATE_ID                = -705
	BS_SDK_ERROR_USER_FULL                   = -706
	BS_SDK_ERROR_DUPLICATE_TEMPLATE          = -707
	BS_SDK_ERROR_FINGERPRINT_FULL            = -708
	BS_SDK_ERROR_DUPLICATE_CARD              = -709
	BS_SDK_ERROR_CARD_FULL                   = -710
	BS_SDK_ERROR_NO_VALID_HDR_FILE           = -711
	BS_SDK_ERROR_INVALID_LOG_FILE            = -712
	BS_SDK_ERROR_CANNOT_FIND_USER            = -714
	BS_SDK_ERROR_ACCESS_LEVEL_FULL           = -715
	BS_SDK_ERROR_INVALID_USER_ID             = -716
	BS_SDK_ERROR_BLACKLIST_FULL              = -717
	BS_SDK_ERROR_USER_NAME_FULL              = -718
	BS_SDK_ERROR_USER_IMAGE_FULL             = -719
	BS_SDK_ERROR_USER_IMAGE_SIZE_TOO_BIG     = -720
	BS_SDK_ERROR_SLOT_DATA_CHECKSUM          = -721
	BS_SDK_ERROR_CANNOT_UPDATE_FINGERPRINT   = -722
	BS_SDK_ERROR_TEMPLATE_FORMAT_MISMATCH    = -723
	BS_SDK_ERROR_NO_ADMIN_USER               = -724
	BS_SDK_ERROR_CANNOT_FIND_LOG             = -725
	BS_SDK_ERROR_DOOR_SCHEDULE_FULL          = -726
	BS_SDK_ERROR_DB_SLOT_FULL                = -727
	BS_SDK_ERROR_ACCESS_GROUP_FULL           = -728
	BS_SDK_ERROR_ACCESS_SCHEDULE_FULL        = -730
	BS_SDK_ERROR_HOLIDAY_GROUP_FULL          = -731
	BS_SDK_ERROR_HOLIDAY_FULL                = -732
	BS_SDK_ERROR_TIME_PERIOD_FULL            = -733
	BS_SDK_ERROR_NO_CREDENTIAL               = -734
	BS_SDK_ERROR_NO_BIOMETRIC_CREDENTIAL     = -735
	BS_SDK_ERROR_NO_CARD_CREDENTIAL          = -736
	BS_SDK_ERROR_NO_PIN_CREDENTIAL           = -737
	BS_SDK_ERROR_NO_BIOMETRIC_PIN_CREDENTIAL = -738
	BS_SDK_ERROR_NO_USER_NAME                = -739
	BS_SDK_ERROR_NO_USER_IMAGE               = -740
	BS_SDK_ERROR_READER_FULL                 = -741
	BS_SDK_ERROR_CACHE_MISSED                = -742
	BS_SDK_ERROR_OPERATOR_FULL               = -743
	BS_SDK_ERROR_INVALID_LINK_ID             = -744
	BS_SDK_ERROR_TIMER_CANCELED              = -745
	BS_SDK_ERROR_USER_JOB_FULL               = -746

	//Config errors
	BS_SDK_ERROR_INVALID_CONFIG           = -800
	BS_SDK_ERROR_CANNOT_OPEN_CONFIG_FILE  = -801
	BS_SDK_ERROR_CANNOT_READ_CONFIG_FILE  = -802
	BS_SDK_ERROR_INVALID_CONFIG_FILE      = -803
	BS_SDK_ERROR_INVALID_CONFIG_DATA      = -804
	BS_SDK_ERROR_CANNOT_WRITE_CONFIG_FILE = -805
	BS_SDK_ERROR_INVALID_CONFIG_INDEX     = -806

	//Device errors
	BS_SDK_ERROR_CANNOT_SCAN_FINGER        = -900
	BS_SDK_ERROR_CANNOT_SCAN_CARD          = -901
	BS_SDK_ERROR_CANNOT_OPEN_RTC           = -902
	BS_SDK_ERROR_CANNOT_SET_RTC            = -903
	BS_SDK_ERROR_CANNOT_GET_RTC            = -904
	BS_SDK_ERROR_CANNOT_SET_LED            = -905
	BS_SDK_ERROR_CANNOT_OPEN_DEVICE_DRIVER = -906
	BS_SDK_ERROR_CANNOT_FIND_DEVICE        = -907

	//Door errors
	BS_SDK_ERROR_CANNOT_FIND_DOOR    = -1000
	BS_SDK_ERROR_DOOR_FULL           = -1001
	BS_SDK_ERROR_CANNOT_LOCK_DOOR    = -1002
	BS_SDK_ERROR_CANNOT_UNLOCK_DOOR  = -1003
	BS_SDK_ERROR_CANNOT_RELEASE_DOOR = -1004

	//Access control errors
	BS_SDK_ERROR_ACCESS_RULE_VIOLATION          = -1100
	BS_SDK_ERROR_DISABLED                       = -1101
	BS_SDK_ERROR_NOT_YET_VALID                  = -1102
	BS_SDK_ERROR_EXPIRED                        = -1103
	BS_SDK_ERROR_BLACKLIST                      = -1104
	BS_SDK_ERROR_CANNOT_FIND_ACCESS_GROUP       = -1105
	BS_SDK_ERROR_CANNOT_FIND_ACCESS_LEVEL       = -1106
	BS_SDK_ERROR_CANNOT_FIND_ACCESS_SCHEDULE    = -1107
	BS_SDK_ERROR_CANNOT_FIND_HOLIDAY_GROUP      = -1108
	BS_SDK_ERROR_CANNOT_FIND_BLACKLIST          = -1109
	BS_SDK_ERROR_AUTH_TIMEOUT                   = -1110
	BS_SDK_ERROR_DUAL_AUTH_TIMEOUT              = -1111
	BS_SDK_ERROR_INVALID_AUTH_MODE              = -1112
	BS_SDK_ERROR_AUTH_UNEXPECTED_USER           = -1113
	BS_SDK_ERROR_AUTH_UNEXPECTED_CREDENTIAL     = -1114
	BS_SDK_ERROR_DUAL_AUTH_FAIL                 = -1115
	BS_SDK_ERROR_BIOMETRIC_AUTH_REQUIRED        = -1116
	BS_SDK_ERROR_CARD_AUTH_REQUIRED             = -1117
	BS_SDK_ERROR_PIN_AUTH_REQUIRED              = -1118
	BS_SDK_ERROR_BIOMETRIC_OR_PIN_AUTH_REQUIRED = -1119
	BS_SDK_ERROR_TNA_CODE_REQUIRED              = -1120
	BS_SDK_ERROR_AUTH_SERVER_MATCH_REFUSAL      = -1121

	//Zone errors
	BS_SDK_ERROR_CANNOT_FIND_ZONE                = -1200
	BS_SDK_ERROR_ZONE_FULL                       = -1201
	BS_SDK_ERROR_HARD_APB_VIOLATION              = -1202
	BS_SDK_ERROR_SOFT_APB_VIOLATION              = -1203
	BS_SDK_ERROR_HARD_TIMED_APB_VIOLATION        = -1204
	BS_SDK_ERROR_SOFT_TIMED_APB_VIOLATION        = -1205
	BS_SDK_ERROR_SCHEDULED_LOCK_VIOLATION        = -1206
	BS_SDK_ERROR_SCHEDULED_UNLOCK_VIOLATION      = -1207
	BS_SDK_ERROR_SET_FIRE_ALARM                  = -1208
	BS_SDK_ERROR_TIMED_APB_ZONE_FULL             = -1209
	BS_SDK_ERROR_FIRE_ALARM_ZONE_FULL            = -1210
	BS_SDK_ERROR_SCHEDULED_LOCK_UNLOCK_ZONE_FULL = -1211
	BS_SDK_ERROR_INACTIVE_ZONE                   = -1212

	//Card errors
	BS_SDK_ERROR_CARD_IO                = -1300
	BS_SDK_ERROR_CARD_INIT_FAIL         = -1301
	BS_SDK_ERROR_CARD_NOT_ACTIVATED     = -1302
	BS_SDK_ERROR_CARD_CANNOT_READ_DATA  = -1303
	BS_SDK_ERROR_CARD_CIS_CRC           = -1304
	BS_SDK_ERROR_CARD_CANNOT_WRITE_DATA = -1305
	BS_SDK_ERROR_CARD_READ_TIMEOUT      = -1306
	BS_SDK_ERROR_CARD_READ_CANCELLED    = -1307
	BS_SDK_ERROR_CARD_CANNOT_SEND_DATA  = -1308
	BS_SDK_ERROR_CANNOT_FIND_CARD       = -1310

	// Operation
	BS_SDK_ERROR_INVALID_PASSWORD = -1400

	// System
	BS_SDK_ERROR_CAMERA_INIT_FAIL             = -1500
	BS_SDK_ERROR_JPEG_ENCODER_INIT_FAIL       = -1501
	BS_SDK_ERROR_CANNOT_ENCODE_JPEG           = -1502
	BS_SDK_ERROR_JPEG_ENCODER_NOT_INITIALIZED = -1503
	BS_SDK_ERROR_JPEG_ENCODER_DEINIT_FAIL     = -1504
	BS_SDK_ERROR_CAMERA_CAPTURE_FAIL          = -1505
	BS_SDK_ERROR_CANNOT_DETECT_FACE           = -1506

	//ETC.
	BS_SDK_ERROR_FILE_IO               = -2000
	BS_SDK_ERROR_ALLOC_MEM             = -2002
	BS_SDK_ERROR_CANNOT_UPGRADE        = -2003
	BS_SDK_ERROR_DEVICE_LOCKED         = -2004
	BS_SDK_ERROR_CANNOT_SEND_TO_SERVER = -2005

	//SSL
	BS_SDK_ERROR_SSL_INIT              = -3000
	BS_SDK_ERROR_SSL_EXIST             = -3001
	BS_SDK_ERROR_SSL_IS_NOT_CONNECTED  = -3002
	BS_SDK_ERROR_SSL_ALREADY_CONNECTED = -3003
	BS_SDK_ERROR_SSL_INVALID_CA        = -3004
	BS_SDK_ERROR_SSL_VERIFY_CA         = -3005
	BS_SDK_ERROR_SSL_INVALID_KEY       = -3006
	BS_SDK_ERROR_SSL_VERIFY_KEY        = -3007

	BS_SDK_ERROR_NULL_POINTER        = -10000
	BS_SDK_ERROR_UNINITIALIZED       = -10001
	BS_SDK_ERROR_CANNOT_RUN_SERVICE  = -10002
	BS_SDK_ERROR_CANCELED            = -10003
	BS_SDK_ERROR_EXIST               = -10004
	BS_SDK_ERROR_ENCRYPT             = -10005
	BS_SDK_ERROR_DECRYPT             = -10006
	BS_SDK_ERROR_DEVICE_BUSY         = -10007
	BS_SDK_ERROR_INTERNAL            = -10008
	BS_SDK_ERROR_INVALID_FILE_FORMAT = -10009
	BS_SDK_ERROR_INVALID_SCHEDULE_ID = -10010
)

const (
	BS2_IPV4_ADDR_SIZE = 16
	BS2_URL_SIZE       = 256
	BS2_USER_ID_SIZE   = 32 ///< Alpha-numeric

	BS2_USER_NAME_LEN   = 48 * 4 ///< UTF-8 Encoding
	BS2_USER_IMAGE_SIZE = 16 * 1024
	BS2_PIN_HASH_SIZE   = 32
)

const (
	BS2_USER_PIN_SIZE    = 32     ///< 16 byte -> 32 byte hash value
	BS2_USER_NAME_SIZE   = 48 * 4 ///< UTF-8 Encoding
	BS2_USER_PHOTO_SIZE  = 16 * 1024
	BS2_MAX_JOB_SIZE     = 16
	BS2_MAX_JOBLABEL_LEN = 16 * 3
	BS2_USER_PHRASE_SIZE = 32 * 4 ///< UTF-8 Encoding

	BS2_INVALID_USER_ID = 0
)

const (
	BS2_MAX_NUM_OF_ACCESS_GROUP_PER_USER = 16
)

//  const (
// 	BS2_CARD_DATA_SIZE		= 32
//  )

type BS2_CARD_TYPE = uint8

type BS2CSNCard struct {
	Type BS2_CARD_TYPE
	Size uint8
	Data [BS2_CARD_DATA_SIZE]byte
}

const (
	BS2_FINGER_TEMPLATE_SIZE = 384
	BS2_TEMPLATE_PER_FINGER  = 2
)

const (
	BS2_FACE_TEMPLATE_LENGTH = 3008 // assert(BS2_FACE_TEMPLATE_LENGTH * 30 % 16 == 0)
	BS2_TEMPLATE_PER_FACE    = 30

	BS2_FACE_IMAGE_SIZE = 16 * 1024
)

func isError(err uintptr) bool {
	return BS_SDK_SUCCESS != (int)(err)
}

func itoa(val int) string { // do it here rather than with fmt to avoid dependency
	if val < 0 {
		return "-" + uitoa(uint(-val))
	}
	return uitoa(uint(val))
}

func uitoa(val uint) string {
	var buf [32]byte // big enough for int64
	i := len(buf) - 1
	for val >= 10 {
		buf[i] = byte(val%10 + '0')
		i--
		val /= 10
	}
	buf[i] = byte(val + '0')
	return string(buf[i:])
}

func uintptrToString(data uintptr) string {
	s := unsafe.Sizeof(data)
	//log.Printf("length = %v\n", s)
	buffer := make([]byte, s)
	p := data
	for i := 0; i < (int)(s); i += 1 {
		//pb := unsafe.Pointer(p)
		u := *(*byte)(unsafe.Pointer(p))
		//log.Printf("%02d = %d\n", i, u)
		buffer[i] = u
		p += 1
	}
	return string(buffer)
}

func uintptrToStringNewVersion(i interface{}) string {
	size := reflect.TypeOf(i).Size()
	buffer := bytes.NewBuffer(make([]byte, size))
	ptr := unsafe.Pointer(&i)

	startAddr := uintptr(ptr)
	endAddr := startAddr + size

	for i := startAddr; i < endAddr; i++ {
		bytePtr := unsafe.Pointer(i)
		b := *(*byte)(bytePtr)
		buffer.WriteByte(b)
	}

	return string(buffer.Bytes())

}
