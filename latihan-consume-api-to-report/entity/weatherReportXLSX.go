package entity

type WeatherReportXLSX struct {
	No                 int
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	Timezone           string  `json:"timezone"`
	Time               string  `json:"time"`
	Temperature2m      float64 `json:"temperature_2m"`
	RelativeHumidity2m int     `json:"relative_humidity_2m"`
	WindSpeed10m       float64 `json:"wind_speed_10m"`
}
