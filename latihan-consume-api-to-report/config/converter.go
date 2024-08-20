package config

import (
	"latihan-consume-api-to-report/entity"
)

func ConvertResponseAPIWeatherToReportXLSX(forecastResponse entity.ForecastCurrent) ([]entity.WeatherReportXLSX, error) {
	var dataReport []entity.WeatherReportXLSX

	for i, row := range forecastResponse.Hourly.Time {
		Humidity := forecastResponse.Hourly.RelativeHumidity2m[i]
		WindSpeed := forecastResponse.Hourly.WindSpeed10m[i]
		Temperatur := forecastResponse.Hourly.Temperature2m[i]
		rowData := entity.WeatherReportXLSX{
			No:                 i + 1,
			Latitude:           forecastResponse.Latitude,
			Longitude:          forecastResponse.Longitude,
			Timezone:           forecastResponse.Timezone,
			Time:               row,
			Temperature2m:      Temperatur,
			RelativeHumidity2m: Humidity,
			WindSpeed10m:       WindSpeed,
		}

		dataReport = append(dataReport, rowData)
	}

	return dataReport, nil
}

func ConvertResponseAPIWeatherToReportPDF(forecastResponse entity.ForecastCurrent) (entity.WeatherReportPDF, error) {
	dataReport := entity.WeatherReportPDF{}
	dataReport.Latitude = forecastResponse.Latitude
	dataReport.Longitude = forecastResponse.Longitude
	dataReport.CurrentTime = forecastResponse.Current.Time
	dataReport.Timezone = forecastResponse.Timezone

	for i, row := range forecastResponse.Hourly.Time {
		Humidity := forecastResponse.Hourly.RelativeHumidity2m[i]
		WindSpeed := forecastResponse.Hourly.WindSpeed10m[i]
		Temperatur := forecastResponse.Hourly.Temperature2m[i]
		rowData := entity.WeatherData{
			No:                 i + 1,
			Time:               row,
			Temperature2m:      Temperatur,
			RelativeHumidity2m: Humidity,
			WindSpeed10m:       WindSpeed,
		}

		dataReport.Data = append(dataReport.Data, rowData)
	}

	return dataReport, nil
}
