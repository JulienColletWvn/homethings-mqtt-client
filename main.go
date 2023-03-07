package main

import (
	"fmt"
	"home-things/handlers"
	db "home-things/utils"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	db.Connect()

	options := mqtt.NewClientOptions().AddBroker(os.Getenv("MQTT_BROKER_HOST"))
	options.SetUsername(os.Getenv("MQTT_USER_NAME"))
	options.SetPassword(os.Getenv("MQTT_TOKEN"))
	options.SetDefaultPublishHandler(handlers.MessageHandler)
	options.AutoReconnect = true

	options.OnConnect = func(c mqtt.Client) {
		fmt.Println("Connected to TTN Broker")
	}

	client := mqtt.NewClient(options)

	forever := make(chan bool)

	for _, d := range handlers.Devices {

		d.Init()

		go func(topic string) {

			if token := client.Connect(); token.Wait() && token.Error() != nil {
				fmt.Println("Error: ", token.Error())
			}

			if token := client.Subscribe(topic, 0, nil); token.Wait() {
				fmt.Println("Subscribed to topic: ", topic)
				if token.Error() != nil {
					fmt.Println(token.Error())
				}
			}

		}(handlers.GetDeviceTopic(&d))

	}

	<-forever

}
