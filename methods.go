package mqttcomms

const (
	// Device methods
	DEV_ADD_NEW      = 0x00
	DEV_ADD          = 0x01
	DEV_REMOVE       = 0x02
	DEV_GET_ACTUATOR = 0x03
	DEV_GET_SENSOR   = 0x04
	DEV_UPDATE_ZONE  = 0x05
	DEV_ADD_NEW_ERL  = 0x06
	ERL_REMOVE       = 0x07

	// Admin methods
	AD_INIT_UPDATE     = 0x00
	AD_SEND_CHUNK      = 0x01
	AD_UPDATE_FINISHED = 0x02
	AD_GET_INFO        = 0x03
	AD_GET_STATUS      = 0x04
	AD_RESET           = 0x05
	AD_ALIVE           = 0x06
	AD_BLE_COMFIRM     = 0x07

	AD_INIT_ESP_OTA       = 0x08
	AD_SEND_ESP_OTA_CHUNK = 0x09
	AD_ESP_OTA_FINISHED   = 0x0A

	AD_SET_ENTITY_ID = 0x11
	AD_FACTORY_RESET = 0x12
	AD_GOODBYE       = 0xFF

	// Time methods
	TM_GET_TIME = 0x00
	TM_PUBLISH  = 0x01

	// Calendar methods
	CL_SET_ZONE_PLANNING  = 0x00
	CL_SET_ZONE_EXCEPTION = 0x01
	CL_GET_ZONE_PLANNING  = 0x02
	CL_GET_ZONE_EXCEPTION = 0x03

	// Telemetry methods
	TELEMETRY = 0x00
)
