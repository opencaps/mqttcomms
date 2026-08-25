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
	// Network-core DFU (multi-image) — matches firmware init_dfu_net/dfu_pkt_net/fin_dfu_net.
	AD_INIT_UPDATE_NET     = 0x0B
	AD_SEND_CHUNK_NET      = 0x0C
	AD_UPDATE_FINISHED_NET = 0x0D

	AD_SET_ENTITY_ID           = 0x11
	// NOTE: factory_reset is 0x12 (firmware bumped it when net-core DFU took 0x0B-0x0D).
	// It was previously (incorrectly) 0x0B here, colliding with net-core DFU init.
	AD_FACTORY_RESET           = 0x12
	AD_SET_ENOCEAN_SENDER_ID   = 0x13
	AD_CHECK_ENTITY_ID_READY   = 0x14
	AD_ENTER_TEST_MODE         = 0x15

	// ERL OTA over Thread — ADR 0007 step 4. The bridge PULLS: the cloud offers an
	// image with BEGIN, the bridge asks for each 512 B slice by offset as its 2x512 B
	// window drains, and reports the outcome with END. One transfer per bridge at a
	// time; a second BEGIN while one is in flight is refused.
	//
	// Bodies are big-endian, like every other body on this channel. EUI-64 travels in
	// CANONICAL (display) order, not the little-endian form the container embeds — the
	// cloud reverses it when personalising, per D4.
	//
	//	BEGIN     eui64[8] fileVersion[4] fileSize[4] channel[1]
	//	SLICE_REQ offset[4] length[2]
	//	SLICE     offset[4] data[<=512]
	//	END       status[1] fileVersion[4] reason[...]  (reason is UTF-8, may be empty)
	AD_ERL_OTA_BEGIN     = 0x16
	AD_ERL_OTA_SLICE_REQ = 0x17
	AD_ERL_OTA_SLICE     = 0x18
	AD_ERL_OTA_END       = 0x19

	AD_GOODBYE = 0xFF

	// Time methods
	TM_GET_TIME = 0x00
	TM_PUBLISH  = 0x01

	// Calendar methods
	CL_SET_ZONE_PLANNING                  = 0x00 // Weekly setpoints + optional exceptions for one zone
	CL_SET_ZONE_EXCEPTION                 = 0x01 // Exception-only update for one zone
	CL_GET_ZONE_PLANNING                  = 0x02 // Bridge requests schedule (inbound)
	CL_GET_ZONE_EXCEPTION                 = 0x03 // Bridge requests exceptions (inbound)
	CL_SEND_PLANNING_GROUP_EVT            = 0x06 // Full group metadata with embedded plannings (UI)
	CL_PLANNING_GROUP_ACTIVATE_RECEIVED   = 0x08 // User activated group on ESP (inbound)
	CL_PLANNING_EXCEPTION_EVT_RECEIVED    = 0x09 // User created exception on ESP (inbound)
	CL_PLANNING_GROUP_DEACTIVATE_RECEIVED = 0x0A // User deactivated group on ESP (inbound)
	CL_CLEAR_ALL_PLANNINGS                = 0x0B // Clear all zone plannings + exceptions (optional hold targets in body)
	CL_SYNC_GROUPS_HASH                   = 0x0C // Bridge sends CRC32 hash of group list for sync check

	// Telemetry methods
	TELEMETRY = 0x00

	// UI methods
	UI_ZONE_LIST     = 0x00
	UI_BUILDING_NAME = 0x01
	UI_FULL_INFO     = 0x02
	UI_WEATHER       = 0x03
	UI_SUN_TIMES     = 0x04
)
