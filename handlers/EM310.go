package handlers

import (
	"database/sql"
	"io/ioutil"
	"net/http"
	"os"
)

type EM310DecodedPayload struct {
	Battery  int    `json:"battery"`
	Distance int    `json:"distance"`
	Position string `json:"position"`
}

var measureTypes = []MeasureType{
	{
		name: "battery",
		unit: "%",
	},
	{
		name: "distance",
		unit: "mm",
	},
	{
		name: "position",
		unit: "position",
	},
}

var id = "eui-24e124713b498841"
var name = "em310"
var location = "rain_water_tank"

var EM310 = Device{
	id,
	name,
	location,
	measureTypes,
	func(payload DecodedPayload) {
		if payload.Battery == 0 && payload.Distance == 0 && payload.Position == "" {
			return
		}

		for _, m := range measureTypes {
			var value float64

			resp, err := http.Get(os.Getenv("API_BASE_URL" + "/api/data-types?name=" + m.name + "&unit=" + m.unit))
			if err != nil {
				continue
			}

			body, err := ioutil.ReadAll(resp.Body)
			sb := string(body)

			if m.name == "position" {
				if payload.Position == "normal" {
					value = 1
				} else {
					value = 2
				}
			} else if m.name == "battery" {
				value = float64(payload.Battery)
			} else if m.name == "distance" {
				value = float64(payload.Distance)
			} else {
				value = 0
			}

			queries.CreateData(ctx, db.CreateDataParams{
				DataTypeID: sql.NullInt32{
					Int32: dt.ID,
					Valid: true,
				},
				DeviceID: sql.NullString{
					String: id,
					Valid:  true,
				},
				Measure: value,
			})
		}

	},
}
