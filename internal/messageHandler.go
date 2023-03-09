package handlers

import (
	"encoding/json"
	"home-things/pkg/services"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"golang.org/x/exp/slices"
)

type Payload struct {
	EndDeviceIds struct {
		DeviceId string `json:"device_id"`
	} `json:"end_device_ids"`
	Message struct {
		DecodedPayload map[string]interface{} `json:"decoded_payload"`
	} `json:"uplink_message"`
}

type DevicesDataType struct {
	deviceId string
	id       int
	key      string
}

var knownDevicesDataTypes []DevicesDataType

func MessageHandler(client mqtt.Client, msg mqtt.Message) {
	var payload Payload

	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		return

	}

	idx := slices.IndexFunc(knownDevicesDataTypes, func(c DevicesDataType) bool { return c.deviceId == payload.EndDeviceIds.DeviceId })

	if idx == -1 {
		dts, err := services.GetDataTypes(payload.EndDeviceIds.DeviceId)

		if err != nil {
			return
		}

		for _, dt := range dts {
			knownDevicesDataTypes = append(knownDevicesDataTypes, DevicesDataType{
				deviceId: payload.EndDeviceIds.DeviceId,
				id:       int(dt.Id),
				key:      dt.Key,
			})
		}

	}

	for k, v := range payload.Message.DecodedPayload {
		idx := slices.IndexFunc(knownDevicesDataTypes, func(c DevicesDataType) bool { return c.deviceId == payload.EndDeviceIds.DeviceId && c.key == k })

		if idx == -1 {
			continue
		}

		dtId := knownDevicesDataTypes[idx].id

		var val float64

		switch value := v.(type) {
		case float64:
			val = value
		case int:
			val = float64(value)
		case string:
			if v == "normal" {
				val = 0
			} else {
				val = -1
			}
		}

		services.CreateData(services.Data{
			DataTypeID: int32(dtId),
			Value:      val,
		})

	}

}
