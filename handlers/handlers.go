package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	db "home-things/db/sqlc"
	"home-things/utils"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func (d *Device) Init() {
	ctx := context.Background()
	queries := db.New(utils.Database)

	device, _ := queries.GetDevice(ctx, d.id)

	if device.ID == "" {
		_, err := queries.CreateDevice(ctx, db.CreateDeviceParams{
			ID:       d.id,
			Name:     d.name,
			Location: d.location,
		})
		if err != nil {
			panic(err)
		}
	}

	for _, m := range d.MeasureTypes {
		measure, _ := queries.GetDataType(ctx, db.GetDataTypeParams{
			Name: m.name,
			Unit: m.unit,
		})

		if measure.ID == 0 {
			_, err := queries.CreateDataType(ctx, db.CreateDataTypeParams{
				Name: m.name,
				Unit: m.unit,
			})
			if err != nil {
				panic(err)
			}
		}

	}
}
func (d *Device) GetDeviceID() string {
	return d.id
}

// https://stackoverflow.com/questions/40823315/x-does-not-implement-y-method-has-a-pointer-receiver
var Devices TDevices = []Device{EM310}

func GetDeviceTopic(d IHandler) string {
	return fmt.Sprintf("v3/the-home-things@ttn/devices/%v/up", d.GetDeviceID())
}

func MessageHandler(client mqtt.Client, msg mqtt.Message) {
	var payload Payload
	err := json.Unmarshal([]byte(msg.Payload()), &payload)

	if err != nil {
		fmt.Println("Can't decode payload")
		return
	}

	for _, d := range Devices {
		d.handlePayload(payload.Message.DecodedPayload)
	}

}
