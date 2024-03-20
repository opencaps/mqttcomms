package mqttcomms

const (
	AD_INIT_UPDATE     = 0x00
	AD_SEND_CHUNK      = 0x01
	AD_UPDATE_FINISHED = 0x02
	AD_GET_INFO        = 0x03
	AD_GET_STATUS      = 0x04
	AD_RESET           = 0x05
	AD_GOODBYE         = 0xFF
)
