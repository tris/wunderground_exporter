package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	defaultPort        = "9122"
	weatherAPIEndpoint = "https://api.weather.com/v2/pws/observations/current?stationId=%s&format=json&apiKey=%s&units=m"
)

var (
	apiKey = os.Getenv("WU_API_KEY")
)

func newWeatherMetrics() map[string]*prometheus.GaugeVec {
	return map[string]*prometheus.GaugeVec{
		"temperature": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "wunderground_temp",
				Help: "Air temperature in degrees Celsius",
			},
			[]string{"station_id", "neighborhood", "software_type", "country"},
		),
		"dewpoint": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "wunderground_dewpt",
				Help: "Dew point temperature in degrees Celsius",
			},
			[]string{"station_id", "neighborhood", "software_type", "country"},
		),
		"humidity": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "wunderground_humidity",
				Help: "Relative humidity in percentage",
			},
			[]string{"station_id", "neighborhood", "software_type", "country"},
		),
		"pressure": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "wunderground_pressure",
				Help: "Atmospheric pressure at sea level in hectopascals",
			},
			[]string{"station_id", "neighborhood", "software_type", "country"},
		),
		"windspeed": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "wunderground_windSpeed",
				Help: "Wind speed in meters per second",
			},
			[]string{"station_id", "neighborhood", "software_type", "country"},
		),
		"winddirection": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "wunderground_windDir",
				Help: "Wind direction in degrees",
			},
			[]string{"station_id", "neighborhood", "software_type", "country"},
		),
		"windgust": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "wunderground_windGust",
				Help: "Wind gust speed in meters per second",
			},
			[]string{"station_id", "neighborhood", "software_type", "country"},
		),
		"precipitation_rate": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "wunderground_precipRate",
				Help: "Precipitation rate in millimeters per hour",
			},
			[]string{"station_id", "neighborhood", "software_type", "country"},
		),
		"precipitation_total": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "wunderground_precipTotal",
				Help: "Total accumulated precipitation in millimeters",
			},
			[]string{"station_id", "neighborhood", "software_type", "country"},
		),
		"uv_index": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "wunderground_uv",
				Help: "Ultraviolet Index",
			},
			[]string{"station_id", "neighborhood", "software_type", "country"},
		),
		"solar_radiation": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "wunderground_solarRadiation",
				Help: "Solar radiation in watts per square meter",
			},
			[]string{"station_id", "neighborhood", "software_type", "country"},
		),
		"epoch": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "wunderground_epoch",
				Help: "Epoch time in seconds",
			},
			[]string{"station_id", "neighborhood", "software_type", "country"},
		),
		"visibility": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "wunderground_visibility",
				Help: "Visibility in meters",
			},
			[]string{"station_id", "neighborhood", "software_type", "country"},
		),
		"soil_temperature": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "wunderground_soilTemp",
				Help: "Soil temperature in degrees Celsius",
			},
			[]string{"station_id", "neighborhood", "software_type", "country"},
		),
		"soil_moisture": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "wunderground_soilMoisture",
				Help: "Soil moisture in percentage",
			},
			[]string{"station_id", "neighborhood", "software_type", "country"},
		),
		"windchill": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "wunderground_windChill",
				Help: "Wind chill temperature in degrees Celsius",
			},
			[]string{"station_id", "neighborhood", "software_type", "country"},
		),
		"elevation": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "wunderground_elevation",
				Help: "Elevation in meters",
			},
			[]string{"station_id", "neighborhood", "software_type", "country"},
		),
		"latitude": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "wunderground_latitude",
				Help: "Latitude",
			},
			[]string{"station_id", "neighborhood", "software_type", "country"},
		),
		"longitude": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "wunderground_longitude",
				Help: "Longitude",
			},
			[]string{"station_id", "neighborhood", "software_type", "country"},
		),
	}
}

type WeatherObservation struct {
	Observations []struct {
		StationID         string      `json:"stationID"`
		ObsTimeUTC        string      `json:"obsTimeUtc"`
		ObsTimeLocal      string      `json:"obsTimeLocal"`
		Neighborhood      string      `json:"neighborhood"`
		SoftwareType      string      `json:"softwareType"`
		Country           string      `json:"country"`
		SolarRadiation    float64     `json:"solarRadiation"`
		Lat               float64     `json:"lat"`
		Lon               float64     `json:"lon"`
		RealtimeFrequency interface{} `json:"realtimeFrequency"`
		Epoch             int         `json:"epoch"`
		UV                float64     `json:"uv"`
		WindDir           int         `json:"winddir"`
		Humidity          int         `json:"humidity"`
		QCStatus          int         `json:"qcStatus"`
		Metric            struct {
			Temp        int     `json:"temp"`
			HeatIndex   int     `json:"heatIndex"`
			DewPt       int     `json:"dewpt"`
			WindChill   int     `json:"windChill"`
			WindSpeed   int     `json:"windSpeed"`
			WindGust    int     `json:"windGust"`
			Pressure    float64 `json:"pressure"`
			PrecipRate  float64 `json:"precipRate"`
			PrecipTotal float64 `json:"precipTotal"`
			Elev        int     `json:"elev"`
		} `json:"metric"`
	} `json:"observations"`
}

type WeatherData struct {
	StationID    string
	Epoch        int
	Latitude     float64
	Longitude    float64
	Elevation    int
	Neighborhood string
	SoftwareType string
	Country      string
	Sensors      map[string]float64
}

func fetchWeatherData(stationID string) (WeatherData, error) {
	url := fmt.Sprintf(weatherAPIEndpoint, stationID, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return WeatherData{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return WeatherData{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return WeatherData{}, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var weatherObservation WeatherObservation
	err = json.Unmarshal(body, &weatherObservation)
	if err != nil {
		return WeatherData{}, err
	}

	obs := weatherObservation.Observations[0]

	data := WeatherData{
		StationID:    stationID,
		Epoch:        obs.Epoch,
		Latitude:     obs.Lat,
		Longitude:    obs.Lon,
		Elevation:    obs.Metric.Elev,
		Neighborhood: obs.Neighborhood,
		SoftwareType: obs.SoftwareType,
		Country:      obs.Country,
		Sensors: map[string]float64{
			"temperature":         float64(obs.Metric.Temp),
			"dewpoint":            float64(obs.Metric.DewPt),
			"humidity":            float64(obs.Humidity),
			"pressure":            obs.Metric.Pressure,
			"windspeed":           float64(obs.Metric.WindSpeed),
			"winddirection":       float64(obs.WindDir),
			"windgust":            float64(obs.Metric.WindGust),
			"precipitation_rate":  obs.Metric.PrecipRate,
			"precipitation_total": obs.Metric.PrecipTotal,
			"uv_index":            obs.UV,
			"solar_radiation":     obs.SolarRadiation,
		},
	}

	return data, nil
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}).ServeHTTP)
	router.HandleFunc("/scrape", func(w http.ResponseWriter, r *http.Request) {
		stationID := r.URL.Query().Get("station_id")
		if stationID == "" {
			http.Error(w, "station_id query parameter is required", http.StatusBadRequest)
			return
		}

		weatherData, err := fetchWeatherData(stationID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch weather data: %s", err), http.StatusInternalServerError)
			return
		}

		registry := prometheus.NewRegistry()
		weatherMetrics := newWeatherMetrics()
		for _, metric := range weatherMetrics {
			registry.MustRegister(metric)
		}

		for sensor, value := range weatherData.Sensors {
			if metric, ok := weatherMetrics[sensor]; ok {
				metric.WithLabelValues(stationID, weatherData.Neighborhood, weatherData.SoftwareType, weatherData.Country).Set(value)
			}
		}

		promhttp.HandlerFor(registry, promhttp.HandlerOpts{}).ServeHTTP(w, r)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
