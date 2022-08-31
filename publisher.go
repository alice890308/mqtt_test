package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func main() {
	// var broker = "192.168.1.49:403/mosquitto"
	// var port = 1883
	opts := mqtt.NewClientOptions()
	//opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.AddBroker("tcp://192.168.1.49:1803")
	opts.SetClientID("client1_subscriber")
	opts.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	if token := client.Subscribe("ac/mqtt2ac", 0, callBack); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	opts.SetClientID("client2_publisher")
	client2 := mqtt.NewClient(opts)
	if token := client2.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter json text: ")
		text, _ := reader.ReadString('\n')
		fmt.Println(text)
		publish(client2, text)
	}
}

func publish(client mqtt.Client, text string) {
	// tmp := &msg{
	// 	DeviceID: 123,
	// }
	// text, _ := json.Marshal(tmp)
	raw := []byte(text)
	token := client.Publish("ac/ac2mqtt", 2, false, raw)
	token.Wait()
	time.Sleep(time.Second)
}

func callBack(client mqtt.Client, msg mqtt.Message) {
	// fmt.Printf("TOPIC: %s\n", msg.Topic())
	// fmt.Printf("MSG: %s\n", msg.Payload())
	fmt.Println(string(msg.Payload()))
}
