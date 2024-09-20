package mqttcomms

type Conf struct {
	ClientID     string
	MqttUrl      string
	MqttCredPath string
	MqttCAPath   string
	MqttCertPath string
	MqttKeyPath  string
	UniqueID     string
}
