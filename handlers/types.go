package handlers

type MeasureType struct {
	name string
	unit string
}

type Device struct {
	id            string
	name          string
	location      string
	MeasureTypes  []MeasureType
	handlePayload func(payload DecodedPayload)
}

type TDevices []Device

type Payload struct {
	EndDeviceIds EndDeviceIds `json:"end_device_ids"`
	Message      Message      `json:"uplink_message"`
}

type Message struct {
	DecodedPayload DecodedPayload `json:"decoded_payload"`
}

type DecodedPayload struct {
	EM310DecodedPayload
}

type EndDeviceIds struct {
	DeviceId string `json:"device_id"`
}

type IHandler interface {
	Init()
	GetDeviceID() string
}
