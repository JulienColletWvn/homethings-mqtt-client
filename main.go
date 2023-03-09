package main

import (
	"fmt"
	handlers "home-things/internal"
	services "home-things/pkg/services"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var Devices []services.Device

func main() {

	d, err := services.GetDevices()

	if err != nil {
		panic(err)
	}

	Devices = d

	options := mqtt.NewClientOptions().AddBroker(os.Getenv("MQTT_BROKER_HOST"))
	options.SetUsername(os.Getenv("MQTT_USER_NAME"))
	options.SetPassword(os.Getenv("MQTT_TOKEN"))
	options.SetDefaultPublishHandler(handlers.MessageHandler)
	options.AutoReconnect = true

	options.OnConnect = func(c mqtt.Client) {
		fmt.Println("Connected to TTN Broker")
		for _, d := range Devices {
			go func(device services.Device) {

				if conn := c.Subscribe(fmt.Sprintf("v3/the-home-things@ttn/devices/%v/up", device.Id), 0, nil); conn.Wait() {
					fmt.Println("Subscribed to up-messages of device", device.Name, device.Location)
					if conn.Error() != nil {
						fmt.Println(conn.Error())
					}
				}

			}(d)

		}
	}

	client := mqtt.NewClient(options)

	if conn := client.Connect(); conn.Wait() && conn.Error() != nil {
		panic(conn.Error())
	}

	select {}

}
