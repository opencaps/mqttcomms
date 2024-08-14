package mqttcomms

type Conf struct {
	ClientID     string
	MqttUrl      string
	MqttUser     string
	MqttPass     string
	MqttCAPath   string
	MqttCertPath string
	MqttKeyPath  string
	UniqueID     string
}
