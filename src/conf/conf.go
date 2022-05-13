//  This file is part of the eliona project.
//  Copyright © 2022 LEICOM iTEC AG. All Rights Reserved.
//  Authors: Adam Lange, et al.
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
	"github.com/eliona-smart-building-assistant/go-eliona/db"
	"strconv"
	"time"
)

type Location struct {
	Location string
	AssetId  int
}

// ReadLocations reads all configured weather locations and send each location to the given channel
func ReadLocations(locations chan Location) {
	db.Query(db.Pool(), "select location, asset_id from weather.locations", locations)
}

// PollingInterval returns the interval polling the weather api
func PollingInterval() time.Duration {
	interval, _ := strconv.Atoi(value("polling_interval", "10"))
	return time.Duration(interval) * time.Second
}

// Endpoint returns the configured API endpoint to get weather data.
func Endpoint() string {
	return value("endpoint", "https://weatherdbi.herokuapp.com/data/weather/")
}

// Value returns the configuration string referenced by key. The configuration is stored in the database
// table weather.configuration. This table should be configurable via the eliona frontend.
func value(name string, fallback string) string {
	valueChan := make(chan string)
	go db.Query(db.Pool(), "select value from weather.configuration where name = $1", valueChan, name)
	value := <-valueChan
	if len(value) == 0 {
		return fallback
	}
	return value
}
