package entity

type WeatherReportPDF struct {
	Latitude    float64
	Longitude   float64
	CurrentTime string
	TimeZone    string
	Data        []WeatherData
}

type WeatherData struct {
	No                 int
	Time               string
	Temperature2m      float64
	RelativeHumidity2m int
	WindSpeed10m       float64
}
