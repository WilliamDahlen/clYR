/*
Copyright, no.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const URL = "https://api.met.no/weatherapi/locationforecast/2.0/compact?"

type Compact struct {
	Type     string `json:"type,omitempty"`
	Geometry struct {
		Type        string    `json:"type,omitempty"`
		Coordinates []float64 `json:"coordinates,omitempty"`
	} `json:"geometry,omitempty"`
	Properties struct {
		Meta struct {
			UpdatedAt time.Time `json:"updated_at,omitempty"`
			Units     struct {
				AirPressureAtSeaLevel string `json:"air_pressure_at_sea_level,omitempty"`
				AirTemperature        string `json:"air_temperature,omitempty"`
				CloudAreaFraction     string `json:"cloud_area_fraction,omitempty"`
				PrecipitationAmount   string `json:"precipitation_amount,omitempty"`
				RelativeHumidity      string `json:"relative_humidity,omitempty"`
				WindFromDirection     string `json:"wind_from_direction,omitempty"`
				WindSpeed             string `json:"wind_speed,omitempty"`
			} `json:"units,omitempty"`
		} `json:"meta,omitempty"`
		Timeseries []struct {
			Time time.Time `json:"time,omitempty"`
			Data struct {
				Instant struct {
					Details struct {
						AirPressureAtSeaLevel float64 `json:"air_pressure_at_sea_level,omitempty"`
						AirTemperature        float64 `json:"air_temperature,omitempty"`
						CloudAreaFraction     float64 `json:"cloud_area_fraction,omitempty"`
						RelativeHumidity      float64 `json:"relative_humidity,omitempty"`
						WindFromDirection     float64 `json:"wind_from_direction,omitempty"`
						WindSpeed             float64 `json:"wind_speed,omitempty"`
					} `json:"details,omitempty"`
				} `json:"instant,omitempty"`
				Next12Hours struct {
					Summary struct {
						SymbolCode string `json:"symbol_code,omitempty"`
					} `json:"summary,omitempty"`
				} `json:"next_12_hours,omitempty"`
				Next1Hours struct {
					Summary struct {
						SymbolCode string `json:"symbol_code,omitempty"`
					} `json:"summary,omitempty"`
					Details struct {
						PrecipitationAmount float64 `json:"precipitation_amount,omitempty"`
					} `json:"details,omitempty"`
				} `json:"next_1_hours,omitempty"`
				Next6Hours struct {
					Summary struct {
						SymbolCode string `json:"symbol_code,omitempty"`
					} `json:"summary,omitempty"`
					Details struct {
						PrecipitationAmount float64 `json:"precipitation_amount,omitempty"`
					} `json:"details,omitempty"`
				} `json:"next_6_hours,omitempty"`
			} `json:"data,omitempty"`
		}
	}
}

type Cities struct {
	City       string `csv:"city"` // .csv column headers
	City_ascii string `csv:"city_ascii"`
	Lat        string `csv:"lat"`
	Lng        string `csv:"lng"`
	Country    string `csv:"country"`
	Iso2       string `csv:"iso2"`
	Iso3       string `csv:"iso3"`
	Admin_name string `csv:"admin_name"`
	Capital    string `csv:"capital"`
	Population int    `csv:"population"`
	Id         int    `csv:"id"`
}

type DefaultCity struct {
	City string `json:"city,omitempty"`
}

func getCords(cityname string) string {
	in, err := os.OpenFile("lib/worldcities.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer in.Close()

	citylist := []*Cities{}

	if err := gocsv.UnmarshalFile(in, &citylist); err != nil {
		panic(err)
	}
	for _, ci := range citylist {
		if ci.City == cityname {
			return "lat=" + ci.Lat + "&lon=" + ci.Lng
		}
	}

	return ""
}

func loadCache(cityname string) (Compact, error) {
	cacheFile, _ := filepath.Glob("cache/" + city + "*.json")

	var cache Compact
	var err error

	if len(cacheFile) != 0 {
		in, err := os.OpenFile(cacheFile[0], os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			panic(err)
		}
		defer in.Close()
		byteValue, _ := io.ReadAll(in)
		if err := json.Unmarshal(byteValue, &cache); err != nil {
			panic(err)
		}
		return cache, err
	}
	if err != nil {
		panic(err)
	}
	return cache, err
}

func getDefaultCity() DefaultCity {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://ip-api.com/json/?fields=city", nil)
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var out DefaultCity
	json.Unmarshal(body, &out)

	return out
}

func currentTimeFormatted() string {
	currentTime := time.Now()

	return currentTime.Weekday().String()[0:3] + ", " + fmt.Sprint(currentTime.Day()) + " " + currentTime.Month().String()[0:3] + " " + fmt.Sprint(currentTime.Year()) + " " + fmt.Sprint(currentTime.Hour()) + ":" + fmt.Sprint(currentTime.Minute()) + ":" + fmt.Sprint(currentTime.Second()) + " GMT"
}

func getClosestTimeIndex(timestamps Compact) int {
	// Get the current time
	now := time.Now()

	// Initialize variables to store the index and difference
	var closestIndex int
	var closestDiff time.Duration

	// Loop through the timestamps and find the one that is closest to the current time
	for i, t := range timestamps.Properties.Timeseries {
		diff := now.Sub(t.Time)
		if diff < 0 {
			diff = -diff
		}
		if closestDiff == 0 || diff < closestDiff {
			closestIndex = i
			closestDiff = diff
		}
	}

	// Print the index of the closest timestamp
	return closestIndex
}

func getCompact(coordinates string) (object Compact, expireDate string) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", URL+coordinates, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", "clYR/0.1 willam@dahlen.dev")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// Retrieve the "Expires" header value as a string
	expiresHeader := resp.Header.Get("Expires")

	// Parse the header value into a time.Time value
	expires, err := time.Parse(time.RFC1123, expiresHeader)
	if err != nil {
		log.Fatalln(err)
	}

	// Format the time.Time value in the HTTP time format
	expiresFormatted := expires.Format(http.TimeFormat)

	var out Compact
	json.Unmarshal(body, &out)

	return out, expiresFormatted
}

func checkAndUpdateCache(city string) Compact {
	cacheFile, _ := filepath.Glob("cache/" + city + "*.json")
	if len(cacheFile) != 0 {
		cached, _ := loadCache(city)
		localCopyTimeStamp, _ := time.Parse(time.RFC1123, strings.Split(cacheFile[0], ".")[1])
		localTimeStamp, _ := time.Parse(time.RFC1123, time.Now().UTC().Format(http.TimeFormat))
		if localCopyTimeStamp.Before(localTimeStamp) {
			err := os.Remove(cacheFile[0])
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println("Stale Cache ")
			out, expireDate := getCompact(getCords(city))
			file, _ := json.MarshalIndent(out, "", " ")
			_ = os.WriteFile("cache/"+city+"."+expireDate+".json", file, 0644)
			return out
		} else {
			fmt.Println("Valid Cache ")
			return cached
		}
	} else {
		fmt.Println("No local cache exist, GET " + city)
		out, expireDate := getCompact(getCords(city))
		file, _ := json.MarshalIndent(out, "", " ")
		os.Create("cache/" + city + "." + expireDate + ".json")
		_ = os.WriteFile("cache/"+city+"."+expireDate+".json", file, 0644)
		return out
	}
}

var cfgFile string

var city string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "clyr",
	Short: "Command line YR",
	Long:  `This application uses the met API to display weather information in the terminal.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Long)
		if city == "" {
			defaultCity := getDefaultCity().City
			cache := checkAndUpdateCache(defaultCity)
			temp := cache.Properties.Timeseries[getClosestTimeIndex(cache)].Data.Instant.Details.AirTemperature
			timeStamp := cache.Properties.Timeseries[getClosestTimeIndex(cache)].Time.Add(1 * time.Hour).Format(http.TimeFormat)
			fmt.Printf("Current Air Temp @ "+timeStamp+": %v\n", temp)
		} else {
			cache := checkAndUpdateCache(city)
			temp := cache.Properties.Timeseries[getClosestTimeIndex(cache)].Data.Instant.Details.AirTemperature
			timeStamp := cache.Properties.Timeseries[getClosestTimeIndex(cache)].Time.Add(1 * time.Hour).Format(http.TimeFormat)
			fmt.Printf("Current Air Temp @ "+timeStamp+": %v\n", temp)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.clyr.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringVarP(&city, "city", "c", "", "The target city")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".clyr" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".clyr")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
