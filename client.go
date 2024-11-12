package mqttcomms

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/op/go-logging"
)

// Client structure for mqtt client
// Client represents an MQTT client with configuration, logging, and synchronization capabilities.
// It includes the following fields:
// - client: The underlying MQTT client instance.
// - conf: Configuration settings for the client.
// - Log: Logger for logging client activities.
// - wg: WaitGroup to manage goroutines.
// - IsConnected: Channel to signal connection status.
// - sighup: Channel to handle OS signals.
// - Mutex: Embedded mutex for synchronization.
type Client struct {
	client      MQTT.Client
	conf        *Conf
	Log         *logging.Logger
	wg          sync.WaitGroup
	IsConnected chan bool
	sighup      chan os.Signal
	sync.Mutex
}

func (c *Client) newTLSConfig() *tls.Config {
	cert, err := tls.LoadX509KeyPair(c.conf.MqttCertPath, c.conf.MqttKeyPath)
	if err != nil {
		c.Log.Fatal("Error loading client certificate", err)
	}

	caCert, err := os.ReadFile(c.conf.MqttCAPath)
	if err != nil {
		c.Log.Fatal("Error loading CA certificate", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
	}

	return tlsConfig
}

func (c *Client) readCredentials() (string, string, error) {
	file, err := os.Open(c.conf.MqttCredPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to open credentials file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var username, password string

	if scanner.Scan() {
		username = scanner.Text()
	} else {
		return "", "", fmt.Errorf("username not found in credentials file")
	}

	if scanner.Scan() {
		password = scanner.Text()
	} else {
		return "", "", fmt.Errorf("password not found in credentials file")
	}

	if err := scanner.Err(); err != nil {
		return "", "", fmt.Errorf("error reading credentials: %v", err)
	}

	return username, password, nil
}

func (c *Client) onConnectHandler(client MQTT.Client) {
	c.Log.Info("MQTT connected")
	c.IsConnected <- true
}

func (c *Client) onConnectionLostHandler(client MQTT.Client, err error) {
	c.Log.Error("MQTT connection lost", err)
	c.IsConnected <- false
}

// InitMqtt initializes the MQTT client with the provided configuration.
// It sets up channels for connection status and signal handling, and starts
// the appropriate renewal handling goroutine based on the configuration.
//
// Parameters:
//   - conf: A pointer to a Conf struct containing the client's configuration.
func (c *Client) InitMqtt(conf *Conf) {
	c.IsConnected = make(chan bool)
	c.sighup = make(chan os.Signal, 1)
	signal.Notify(c.sighup, os.Interrupt, syscall.SIGHUP)
	c.conf = conf

	switch c.conf.RenewalMode {
	case RenewalModeExpiracy:
		go c.handleExpiracy()
	case RenewalModeSignal:
		go c.handleSighup()
	}
}

// Connect establishes a connection to the MQTT broker using the client's configuration.
// It sets up TLS configuration, reads credentials from a file if specified, and configures
// the MQTT client options including username, password, broker URL, client ID, and handlers
// for connection events. The function attempts to connect to the broker and retries every
// 5 seconds until a successful connection is made.
func (c *Client) Connect() {
	tlsConfig := c.newTLSConfig()
	var user, pass string
	var err error

	if c.conf.CredFromFile {
		user, pass, err = c.readCredentials()
		if err != nil {
			c.Log.Fatal("Error reading credentials", err)
		}
	} else {
		user = c.conf.MqttUser
		pass = c.conf.MqttPass
	}

	opts := MQTT.NewClientOptions()
	opts.SetUsername(user)
	opts.SetPassword(pass)
	opts.AddBroker("mqtts://" + c.conf.MqttUrl)
	opts.SetClientID(c.conf.ClientID + c.conf.UniqueID)
	opts.SetTLSConfig(tlsConfig)
	opts.SetOnConnectHandler(c.onConnectHandler)
	opts.SetConnectionLostHandler(c.onConnectionLostHandler)
	c.client = MQTT.NewClient(opts)

	c.Log.Info("Trying to connect to the broker mqtts://" + c.conf.MqttUrl)

	retry := time.NewTicker(5 * time.Second)
	for range retry.C {
		if token := c.client.Connect(); token.Wait() && token.Error() == nil {
			return
		}
	}
}

func (c *Client) handleSighup() {
	for range c.sighup {
		c.Log.Info("SIGHUP received, reconnecting")
		c.client.Disconnect(1000)
		c.Connect()
	}
}

func (c *Client) handleExpiracy() {
	certPEM, err := os.ReadFile(c.conf.MqttCertPath)
	if err != nil {
		c.Log.Fatal("Error reading certificate", err)
	}

	block, _ := pem.Decode(certPEM)
	if block == nil {
		c.Log.Fatal("failed to parse certificate PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		c.Log.Fatal("Error parsing certificate", err)
	}

	// Calculate 80% of the certificate's validity period
	validityPeriod := cert.NotAfter.Sub(cert.NotBefore)
	timeToWait := validityPeriod * 95 / 100

	if timeToWait < 0 {
		time.Sleep(time.Hour)
		go c.handleExpiracy()
		return
	}

	expirationTime := cert.NotBefore.Add(timeToWait)
	timeUntilExpiration := time.Until(expirationTime).Round(time.Second)

	c.Log.Info("Certificate will be renewed in", timeUntilExpiration)

	timer := time.NewTimer(timeUntilExpiration)
	<-timer.C

	// Reconnect
	c.Log.Info("Certificate expired, reconnecting")
	c.client.Disconnect(1000)
	c.Connect()

	// Restart the expiration handler
	go c.handleExpiracy()
}

// Subscribe to a topic, the messages of this topic will be passed through the method callback
// Subscribe subscribes the client to a given MQTT topic and registers a callback
// function to handle incoming messages on that topic.
//
// Parameters:
//   - topic: The MQTT topic to subscribe to.//   - callback: A function of type MQTT.MessageHandler that will be called//	data  - The byte slice containing the data to be published.
//     topic - The MQTT topic to which the data should be published.
//
// Returns:
//
//	error - An error if the MQTT client is not instantiated or if the publish operation fails, otherwise nil.
func (c *Client) Subscribe(topic string, callback MQTT.MessageHandler) {
	c.client.Subscribe(topic, 0, func(client MQTT.Client, msg MQTT.Message) {
		callback(client, msg)
	})
}

// Unsubscribe unsubscribes the client from the specified MQTT topic.
// It takes a single parameter:
// - topic: the topic string from which the client should unsubscribe.
func (c *Client) Unsubscribe(topic string) {
	c.client.Unsubscribe(topic)
}

// WriteMQTT send data to a topic
// WriteMQTT publishes data to a specified MQTT topic.
// It waits for any ongoing operations to complete, locks the client for exclusive access,
// and then attempts to publish the data. If the MQTT client is not instantiated, it logs a warning
// and returns an error. If the publish operation fails, it logs a warning and returns the error.
//
// Parameters:
//
//	data  - The byte slice containing the data to be published.
//	topic - The MQTT topic to which the data should be published.
//
// Returns:
//
//	error - An error if the MQTT client is not instantiated or if the publish operation fails, otherwise nil.
func (c *Client) WriteMQTT(data []byte, topic string) error {
	c.wg.Wait()
	c.Lock()
	defer c.Unlock()

	if c.client == nil {
		c.Log.Warning("Mqtt client not instanciate")
		return errors.New("mqtt client not instanciate")
	}

	c.Log.Debug("Sending data on", topic)

	if token := c.client.Publish(topic, 0, false, data); token.Wait() && token.Error() != nil {
		c.Log.Warning("Cannot send a command", token.Error().Error())
		return token.Error()
	}
	return nil
}
