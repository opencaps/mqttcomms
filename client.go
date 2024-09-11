package mqttcomms

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"
	"sync"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/op/go-logging"
)

// Client structure for mqtt client
type Client struct {
	client      MQTT.Client
	url         string
	Log         *logging.Logger
	wg          sync.WaitGroup
	IsConnected chan bool
	sync.Mutex
}

func (c *Client) newTLSConfig(conf *Conf) *tls.Config {
	cert, err := tls.LoadX509KeyPair(conf.MqttCertPath, conf.MqttKeyPath)
	if err != nil {
		c.Log.Fatal("Error loading client certificate", err)
	}

	caCert, err := os.ReadFile(conf.MqttCAPath)
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
	c.wg.Done()
	c.IsConnected <- true
}

func (c *Client) onConnectionLostHandler(client MQTT.Client, err error) {
	c.Log.Error("MQTT connection lost", err)
	c.IsConnected <- false
	c.wg.Add(1)
}

// InitMqtt initialze mqtt client
func (c *Client) InitMqtt(conf *Conf) {
	c.IsConnected = make(chan bool)
	c.wg.Add(1)

	tlsconfig := c.newTLSConfig(conf)

	opts := MQTT.NewClientOptions()
	opts.SetClientID(conf.ClientID + conf.UniqueID)
	opts.SetUsername(conf.MqttUser)
	opts.SetPassword(conf.MqttPass)
	opts.SetTLSConfig(tlsconfig)
	opts.SetOnConnectHandler(c.onConnectHandler)
	opts.SetConnectionLostHandler(c.onConnectionLostHandler)
	opts.AddBroker("mqtts://" + conf.MqttUrl)
	c.url = "mqtts://" + conf.MqttUrl
	c.client = MQTT.NewClient(opts)
}

func (c *Client) Connect() {
	c.Log.Info("Trying to connect to the broker", c.url)
	retry := time.NewTicker(5 * time.Second)
	for range retry.C {
		if token := c.client.Connect(); token.Wait() && token.Error() == nil {
			return
		}
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
