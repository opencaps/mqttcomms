package mqttcomms

const (
	// Device methods
	DEV_ADD_NEW                = 0x00
	DEV_ADD                    = 0x01
	DEV_REMOVE                 = 0x02
	DEV_GET_ACTUATOR           = 0x03
	DEV_GET_SENSOR             = 0x04
	DEV_UPDATE_ZONE            = 0x05
	DEV_ADD_NEW_ERL            = 0x06
	ERL_REMOVE                 = 0x07
	DEV_ADD_INTERNAL_TMP_SENSOR = 0x08
	DEV_ADD_INTERNAL_HUM_SENSOR = 0x09
	DEV_ADD_INTERNAL_COV_SENSOR = 0x0A
	DEV_ADD_INTERNAL_NOX_SENSOR          = 0x0B
	DEV_UPDATE_ZONE_INTERNAL_SENSORS     = 0x0C
	DEV_REMOVE_INTERNAL_SENSORS          = 0x0E

	// Admin methods
	AD_INIT_UPDATE        = 0x00
	AD_SEND_CHUNK         = 0x01
	AD_UPDATE_FINISHED    = 0x02
	AD_GET_INFO           = 0x03
	AD_GET_STATUS         = 0x04
	AD_RESET              = 0x05
	AD_ALIVE              = 0x06
	AD_BLE_COMFIRM        = 0x07
	AD_INIT_ESP_OTA       = 0x08
	AD_SEND_ESP_OTA_CHUNK = 0x09
	AD_ESP_OTA_FINISHED   = 0x0A
	AD_FACTORY_RESET      = 0x0B
	AD_SET_ENTITY_ID           = 0x11
	AD_SET_ENOCEAN_SENDER_ID   = 0x13

	AD_GOODBYE = 0xFF

	// Time methods
	TM_GET_TIME = 0x00
	TM_PUBLISH  = 0x01

	// Calendar methods
	CL_SET_ZONE_PLANNING                  = 0x00
	CL_SET_ZONE_EXCEPTION                 = 0x01
	CL_GET_ZONE_PLANNING                  = 0x02
	CL_GET_ZONE_EXCEPTION                 = 0x03
	CL_DELETE_ZONE_PLANNING               = 0x04
	CL_DELETE_ZONE_EXCEPTION              = 0x05
	CL_SEND_PLANNING_GROUP_EVT            = 0x06
	CL_SEND_PLANNING_EVT                  = 0x07
	CL_PLANNING_GROUP_ACTIVATE_RECEIVED   = 0x08
	CL_PLANNING_EXCEPTION_EVT_RECEIVED    = 0x09
	CL_PLANNING_GROUP_DEACTIVATE_RECEIVED = 0x0A

	// Telemetry methods
	TELEMETRY = 0x00

	// UI methods
	UI_ZONE_LIST     = 0x00
	UI_BUILDING_NAME = 0x01
	UI_FULL_INFO     = 0x02
	UI_WEATHER       = 0x03
	UI_SUN_TIMES     = 0x04
)
