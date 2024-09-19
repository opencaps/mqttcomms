package mqttcomms

type Conf struct {
	ClientID     string
	MqttUrl      string
	MqttUserPath string
	MqttPassPath string
	MqttCAPath   string
	MqttCertPath string
	MqttKeyPath  string
	UniqueID     string
}
