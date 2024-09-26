package mqttcomms

const (
	// Device methods
	DEV_ADD_NEW = 0x00
	DEV_ADD     = 0x01
	DEV_REMOVE  = 0x02

	// Admin methods
	AD_INIT_UPDATE     = 0x00
	AD_SEND_CHUNK      = 0x01
	AD_UPDATE_FINISHED = 0x02
	AD_GET_INFO        = 0x03
	AD_GET_STATUS      = 0x04
	AD_RESET           = 0x05
	AD_ALIVE           = 0x06
	AD_BLE_COMFIRM     = 0x07
	AD_SET_ENTITY_ID   = 0x11
	AD_GOODBYE         = 0xFF

	// Time methods
	TM_GET_TIME = 0x00
	TM_PUBLISH  = 0x01
)
