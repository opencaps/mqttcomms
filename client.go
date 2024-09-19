package mqttcomms

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
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

	user, err := os.ReadFile(c.conf.MqttUserPath)
	if err != nil {
		c.Log.Fatal("Error reading username file", err)
	}

	pass, err := os.ReadFile(c.conf.MqttPassPath)
	if err != nil {
		c.Log.Fatal("Error reading password file", err)
	}

	opts := MQTT.NewClientOptions()
	opts.SetUsername(string(user))
	opts.SetPassword(string(pass))
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
