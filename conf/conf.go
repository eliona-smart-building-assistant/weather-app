//  This file is part of the eliona project.
//  Copyright © 2022 LEICOM iTEC AG. All Rights Reserved.
//  ______ _ _
// |  ____| (_)
// | |__  | |_  ___  _ __   __ _
// |  __| | | |/ _ \| '_ \ / _` |
// | |____| | | (_) | | | | (_| |
// |______|_|_|\___/|_| |_|\__,_|
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
//  BUT NOT LIMITED  TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//  NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
//  DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package conf

import (
	"github.com/eliona-smart-building-assistant/go-utils/db"
	"strconv"
	"time"
)

type Location struct {
	Location  string
	Latitude  float64
	Longitude float64
	ProjectId string
}

// ReadLocations reads all configured weather locations and send each location to the given channel
func ReadLocations(locations chan Location) error {
	return db.Query(db.Pool(), "select location, latitude, longitude, coalesce(proj_id,'') from weather.locations", locations)
}

// PollingInterval returns the interval polling the weather api
func PollingInterval() time.Duration {
	interval, _ := strconv.Atoi(get("polling_interval", "10"))
	return time.Duration(interval) * time.Second
}

// Endpoint returns the configured API endpoint to get weather data.
func Endpoint() string {
	return get("endpoint", "https://www.7timer.info/bin/civillight.php?ac=0&unit=metric&output=json&tzshift=0")
}

// Value returns the configuration string referenced by key. The configuration is stored in the init
// table weather.configuration. This table should be configurable via the eliona frontend.
func get(name string, fallback string) string {
	valueChan := make(chan string)
	go func() {
		_ = db.Query(db.Pool(), "select value from weather.configuration where name = $1", valueChan, name)
	}()
	value := <-valueChan
	if len(value) == 0 {
		return fallback
	}
	return value
}

// Set sets the value of configuration
func Set(connection db.Connection, name string, value string) error {
	return db.Exec(connection, "insert into weather.configuration (name, value) values ($1, $2) on conflict (name) do update set value = excluded.value", name, value)
}

func InsertLocation(connection db.Connection, location Location) error {
	return db.Exec(connection, "insert into weather.locations (location, latitude, longitude) values ($1, $2, $3)", location.Location, location.Latitude, location.Longitude)
}
