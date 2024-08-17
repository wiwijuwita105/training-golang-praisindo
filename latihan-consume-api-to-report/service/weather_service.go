package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/xuri/excelize/v2"
	template2 "html/template"
	"latihan-consume-api-to-report/config"
	"latihan-consume-api-to-report/entity"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

type IWeatherService interface {
	ReportWeatherToXLSX(context.Context, entity.LocationRequest) (entity.GenerateReportResponse, error)
	ReportWeatherToPDF(context.Context, entity.LocationRequest) (entity.GenerateReportResponse, error)
}

type WeatherService struct{}

func NewWeatherService() *WeatherService {
	return &WeatherService{}
}

var (
	httpClient *http.Client
	logger     *slog.Logger
)

const baseDirOutputXLSX = "output/xlsx"
const baseDirOutputPDF = "output/pdf"
const htmlTemplate = "static/template.html"

func GetWeather(ctx context.Context, endpoint string) (entity.ForecastCurrent, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		slog.WarnContext(ctx, "error when hit http.NewRequestWithContext", slog.Any("error", err))
		return entity.ForecastCurrent{}, err
	}

	res, err := httpClient.Do(httpReq)
	if err != nil {
		slog.WarnContext(ctx, "error when hit httpClient.Do", slog.Any("error", err))
		return entity.ForecastCurrent{}, err
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		errStatusCode := errors.New("not receiving status OK when hit API")
		slog.WarnContext(ctx, errStatusCode.Error(), slog.Any("error", errStatusCode), slog.Any("res.StatusCode", res.StatusCode))
		return entity.ForecastCurrent{}, errStatusCode
	}

	var foreCast entity.ForecastCurrent
	if err = json.NewDecoder(res.Body).Decode(&foreCast); err != nil {
		slog.WarnContext(ctx, "error when hit Decode(&users)", slog.Any("error", err))
		return entity.ForecastCurrent{}, err
	}

	return foreCast, nil
}

func (s WeatherService) ReportWeatherToXLSX(ctx context.Context, request entity.LocationRequest) (entity.GenerateReportResponse, error) {
	url := fmt.Sprintf(config.URL_FORECAST+"?latitude=%f&longitude=%f&current=temperature_2m,wind_speed_10m&hourly=temperature_2m,relative_humidity_2m,wind_speed_10m", request.Latitude, request.Longitude)
	startTime := time.Now()
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	httpClient = &http.Client{Timeout: 10 * time.Second}

	slog.InfoContext(ctx, "Start retrieving users data")
	foreCast, err := GetWeather(ctx, url)
	if err != nil {
		return entity.GenerateReportResponse{}, err
	}
	slog.InfoContext(ctx, fmt.Sprintf("Finished retrieving weather"))

	dataReport, err := config.ConvertResponseAPIWeatherToReportXLSX(foreCast)
	if err != nil {
		return entity.GenerateReportResponse{}, err
	}

	slog.InfoContext(ctx, "Start generating XLSX weather")

	xlsx := excelize.NewFile()
	sheet1Name := "Sheet1"
	xlsx.SetSheetName(xlsx.GetSheetName(1), sheet1Name)

	xlsx.SetCellValue(sheet1Name, "A1", "No")
	xlsx.SetCellValue(sheet1Name, "B1", "Latitude")
	xlsx.SetCellValue(sheet1Name, "C1", "Longitude")
	xlsx.SetCellValue(sheet1Name, "D1", "Timezone")
	xlsx.SetCellValue(sheet1Name, "E1", "Time")
	xlsx.SetCellValue(sheet1Name, "F1", "Temperature_2m")
	xlsx.SetCellValue(sheet1Name, "G1", "Relative_Humidity_2m")
	xlsx.SetCellValue(sheet1Name, "H1", "Wind_Speed_10m")

	for i, row := range dataReport {
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("A%d", i+2), row.No)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("B%d", i+2), row.Latitude)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("C%d", i+2), row.Longitude)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("D%d", i+2), row.Timezone)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("E%d", i+2), row.Time)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("F%d", i+2), row.Temperature2m)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("G%d", i+2), row.RelativeHumidity2m)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("H%d", i+2), row.WindSpeed10m)
	}

	timeString := foreCast.Current.Time
	timeString = strings.Replace(timeString, ":", "", -1)
	timeString = strings.Replace(timeString, "-", "", -1)
	nameFile := "Report-Weather-" + timeString + ".xlsx"

	err = xlsx.SaveAs("./" + baseDirOutputXLSX + "/" + nameFile)
	if err != nil {
		slog.Error(err.Error())
		return entity.GenerateReportResponse{}, err
	}
	slog.InfoContext(ctx, fmt.Sprintf("Finish generating xlsx weathe. Elapsed Time: %d ms", time.Since(startTime).Milliseconds()))
	return entity.GenerateReportResponse{
		Message: "Success to Generate XLSX weather.",
	}, nil
}

func (s WeatherService) ReportWeatherToPDF(ctx context.Context, request entity.LocationRequest) (entity.GenerateReportResponse, error) {
	url := fmt.Sprintf(config.URL_FORECAST+"?latitude=%f&longitude=%f&current=temperature_2m,wind_speed_10m&hourly=temperature_2m,relative_humidity_2m,wind_speed_10m", request.Latitude, request.Longitude)
	//startTime := time.Now()
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	httpClient = &http.Client{Timeout: 10 * time.Second}

	slog.InfoContext(ctx, "Start retrieving users data")
	foreCast, err := GetWeather(ctx, url)
	if err != nil {
		return entity.GenerateReportResponse{}, err
	}
	slog.InfoContext(ctx, fmt.Sprintf("Finished retrieving weather"))

	dataReport, err := config.ConvertResponseAPIWeatherToReportPDF(foreCast)
	if err != nil {
		return entity.GenerateReportResponse{}, err
	}

	template, err := template2.ParseFiles("./" + htmlTemplate)
	if err != nil {
		log.Fatalln(err)
	}

	buf := new(bytes.Buffer)

	if err := template.Execute(buf, dataReport); err != nil {
		log.Println(err)
	}

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatalln(err)
	}

	page := wkhtmltopdf.NewPageReader(strings.NewReader(buf.String()))
	// Add to document
	pdfg.AddPage(page)
	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
	}
	// Write buffer contents to file on disk
	timeString := foreCast.Current.Time
	timeString = strings.Replace(timeString, ":", "", -1)
	timeString = strings.Replace(timeString, "-", "", -1)
	nameFile := "Report-Weather-" + timeString + ".pdf"
	err = pdfg.WriteFile("./" + baseDirOutputPDF + "/" + nameFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done")

	return entity.GenerateReportResponse{
		Message: "Success to Generate PDF weather.",
	}, nil
}
