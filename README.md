# MqttComms: MQTT Protocol Handling for Opencaps Devices

[![AsyncAPI](https://img.shields.io/badge/AsyncAPI-3.0.0-blue.svg)](https://studio.asyncapi.com/?url=https://raw.githubusercontent.com/opencaps/mqttcomms/main/asyncapi.yaml)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

## Overview

MqttComms is a Go library designed to facilitate seamless communication between Opencaps devices and services using the MQTT protocol. It provides robust tools for encoding and decoding custom protocol messages, secure communication using TLS, and managing MQTT topics and methods.

---

## Features

- **Secure MQTT Communication**: Supports TLS for secure connections.
- **Automatic Certificate Renewal**: Handles certificate expiration and reconnects automatically.
- **Topic Management**: Predefined topics for various use cases.
- **Custom Protocol Support**: Encode and decode messages with CRC validation.
- **Signal-Based Reconnection**: Reconnects on receiving `SIGHUP` signals.
- **Flexible Configuration**: Supports file-based or inline credentials.
- **AsyncAPI Specification**: Complete API documentation in AsyncAPI 3.0 format.

---

## API Documentation

The MQTT API is fully documented using the [AsyncAPI](https://www.asyncapi.com/) specification.

### View Documentation

- **Interactive Studio**: [Open in AsyncAPI Studio](https://studio.asyncapi.com/?url=https://raw.githubusercontent.com/opencaps/mqttcomms/main/asyncapi.yaml)
- **Raw Specification**: [asyncapi.yaml](asyncapi.yaml)

### Generate Documentation Locally

```bash
# Install AsyncAPI CLI
npm install -g @asyncapi/cli

# Validate the specification
asyncapi validate asyncapi.yaml

# Generate HTML documentation
asyncapi generate fromTemplate asyncapi.yaml @asyncapi/html-template -o docs
```

---

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/opencaps/mqttcomms.git
   cd mqttcomms
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

---

## Configuration

The library uses the `Conf` struct (defined in [`conf.go`](conf.go)) for configuration. Below are the key fields:

- **Client Identification**:
  - `ClientID`: Unique identifier for the client.
  - `UniqueID`: Additional unique identifier.

- **Credentials**:
  - `CredFromFile`: Boolean to indicate if credentials are read from a file.
  - `MqttUser` and `MqttPass`: Inline username and password.
  - `MqttCredPath`: Path to the credentials file.

- **TLS Configuration**:
  - `MqttCAPath`: Path to the CA certificate.
  - `MqttCertPath`: Path to the client certificate.
  - `MqttKeyPath`: Path to the client key.

- **Broker Details**:
  - `MqttUrl`: URL of the MQTT broker.

- **Renewal Mode**:
  - `RenewalMode`: Mode for renewing the connection (`expiracy` or `signal`).

---

## Usage

### 1. Initialize the Client

Create a new `Client` instance and initialize it with a configuration:

```go
client := &mqttcomms.Client{}
conf := &mqttcomms.Conf{
    ClientID:     "example-client",
    MqttUrl:      "broker.example.com",
    MqttCAPath:   "/path/to/ca.crt",
    MqttCertPath: "/path/to/client.crt",
    MqttKeyPath:  "/path/to/client.key",
    RenewalMode:  mqttcomms.RenewalModeExpiracy,
}
client.InitMqtt(conf)
```

### 2. Connect to the Broker

Establish a connection to the MQTT broker:

```go
client.Connect()
```

### 3. Subscribe to Topics

Subscribe to a topic and handle incoming messages:

```go
client.Subscribe("example/topic", func(client MQTT.Client, msg MQTT.Message) {
    fmt.Println("Received message:", string(msg.Payload()))
})
```

### 4. Publish Messages

Publish a message to a topic:

```go
err := client.WriteMQTT([]byte("Hello, MQTT!"), "example/topic")
if err != nil {
    fmt.Println("Error publishing message:", err)
}
```

### 5. Unsubscribe from Topics

Unsubscribe from a topic:

```go
client.Unsubscribe("example/topic")
```

---

## Topics

Predefined topics are available in [`topics.go`](topics.go):

- **Admin**:
  - `TopAdminEvt`: Admin events.
  - `TopAdminOrd`: Admin orders.

- **Device**:
  - `TopDevEvt`: Device events.
  - `TopDevOrd`: Device orders.

- **EnOcean**:
  - `TopEnoEvt`: EnOcean events.
  - `TopEnoOrd`: EnOcean orders.

- **Telemetry**:
  - `TopTelemetryEvt`: Telemetry events.
  - `TopTelemetryOrd`: Telemetry orders.

- **Time**:
  - `TimeEvt`: Time events.
  - `TimeOrd`: Time orders.

- **Calendar**:
  - `TopCalEvt`: Calendar events.
  - `TopCalOrd`: Calendar orders.

- **UI**:
  - `TopUIEvt`: UI events.
  - `TopUIOrd`: UI orders.

---

## Methods

Predefined methods are available in [`methods.go`](methods.go):

- **Device Methods**:
  - `DEV_ADD_NEW`, `DEV_ADD`, `DEV_REMOVE`, etc.

- **Admin Methods**:
  - `AD_INIT_UPDATE`, `AD_SEND_CHUNK`, `AD_RESET`, etc.

- **Time Methods**:
  - `TM_GET_TIME`, `TM_PUBLISH`.

- **Calendar Methods**:
  - `CL_SET_ZONE_PLANNING`, `CL_GET_ZONE_EXCEPTION`.

- **Telemetry Methods**:
  - `TELEMETRY`.

- **UI Methods**:
  - `UI_ZONE_LIST`: Request/send zone list.
  - `UI_BUILDING_NAME`: Request/send building name.
  - `UI_FULL_INFO`: Request/send full UI information.

---

## Custom Protocol

The library supports encoding and decoding custom protocol messages with CRC validation. See [`mqttcomms.go`](mqttcomms.go) for details.

### Encode a Message

Use `GenerateMsg` to encode a message:

```go
msg := &mqttcomms.Msg{
    Seq:    1,
    Op:     mqttcomms.REQRESP,
    Method: mqttcomms.DEV_ADD,
    Body:   []byte("payload"),
}
encoded, err := mqttcomms.GenerateMsg(msg)
if err != nil {
    fmt.Println("Error encoding message:", err)
}
```

### Decode a Message

Use `DecodeMsg` to decode a message:

```go
decoded, err := mqttcomms.DecodeMsg(encoded)
if err != nil {
    fmt.Println("Error decoding message:", err)
}
fmt.Println("Decoded message:", decoded)
```

---

## License

This project is licensed under the GNU General Public License v3. See the [LICENSE](LICENSE) file for details.