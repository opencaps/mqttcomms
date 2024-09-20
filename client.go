package mqttcomms

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
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

// InitMqtt initialze mqtt client
func (c *Client) InitMqtt(conf *Conf) {
	c.IsConnected = make(chan bool)
	c.sighup = make(chan os.Signal, 1)
	signal.Notify(c.sighup, os.Interrupt, syscall.SIGHUP)
	c.conf = conf

	go c.handleSighup()
}

func (c *Client) Connect() {
	tlsConfig := c.newTLSConfig()

	user, pass, err := c.readCredentials()
	if err != nil {
		c.Log.Fatal("Error reading credentials", err)
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

// Handle sighup signal
func (c *Client) handleSighup() {
	for range c.sighup {
		c.Log.Info("SIGHUP received, reconnecting")
		c.client.Disconnect(1000)
		c.Connect()
	}
}

// Subscribe to a topic, the messages of this topic will be passed through the method callback
func (c *Client) Subscribe(topic string, callback MQTT.MessageHandler) {
	c.client.Subscribe(topic, 0, func(client MQTT.Client, msg MQTT.Message) {
		callback(client, msg)
	})
}

// Unsubscribe from a topic
func (c *Client) Unsubscribe(topic string) {
	c.client.Unsubscribe(topic)
}

// WriteMQTT send data to a topic
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
