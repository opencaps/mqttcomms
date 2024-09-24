package mqttcomms

const (
	// RenewalModeExpiracy is the mode where the client renews the connection on expiracy
	RenewalModeExpiracy = "expiracy"
	// RenewalModeSignal is the mode where the client renews the connection on SIGHUP
	RenewalModeSignal = "signal"
)

// Conf represents the configuration settings required for MQTT communication.
// It includes fields for client identification, credential sources, and MQTT connection details.
//
// Fields:
// - ClientID: A unique identifier for the MQTT client.
// - CredFromFile: A boolean indicating whether credentials should be loaded from a file.
// - MqttUrl: The URL of the MQTT broker.
// - MqttUser: The username for MQTT authentication.
// - MqttPass: The password for MQTT authentication.
// - MqttCredPath: The file path to the MQTT credentials.
// - MqttCAPath: The file path to the Certificate Authority (CA) certificate.
// - MqttCertPath: The file path to the client certificate.
// - MqttKeyPath: The file path to the client key.
// - UniqueID: A unique identifier for the configuration instance.
// - RenewalMode: The mode for renewing the MQTT connection or credentials.
type Conf struct {
	ClientID     string
	CredFromFile bool
	MqttUrl      string
	MqttUser     string
	MqttPass     string
	MqttCredPath string
	MqttCAPath   string
	MqttCertPath string
	MqttKeyPath  string
	UniqueID     string
	RenewalMode  string
}
