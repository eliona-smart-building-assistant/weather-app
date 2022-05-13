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

package weather

import (
	"context"
	"github.com/eliona-smart-building-assistant/go-eliona/log"
	"time"
	"weather/api"
	"weather/conf"
)

// StartCollectingData starts an infinite loop to periodically collect data for weather locations.
func StartCollectingData() {
	log.Info("Weather", "Start service for collecting weather locations.")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {

		// Collecting data for all configured weather locations. This method is called periodically
		// until the context is done.
		collectData()

		select {
		case <-time.After(conf.PollingInterval()):
		case <-ctx.Done():
			log.Info("Weather", "Finish collecting data for weather locations.")
			return
		}
	}
}

// Input is a structure holds input data getting from the api endpoint. This structure corresponds
// to the input heap data in eliona.
type Input struct {
	Humidity      int `json:"daytime"`
	Precipitation int `json:"precipitation"`
	Wind          int `json:"wind"`
	Temperature   int `json:"temperature"`
}

// Info is a structure holds informational data getting from the api endpoint. This structure corresponds
// to the info heap data in eliona.
type Info struct {
	Daytime string `json:"daytime"`
}

// Status is a structure holds data getting from the api endpoint related to state of the weather location.
// This structure corresponds to the status heap data in eliona.
type Status struct {
	Comment string `json:"comment"`
}

func collectData() {
	locations := make(chan conf.Location)
	go conf.ReadLocations(locations)
	for location := range locations {
		condition, err := api.Today(location.Location)
		if err != nil {
			log.Error("Weather", "Error during requesting API endpoint: %v", err)
		} else {
			log.Info("Weather", "%s: %f °C", location.Location, condition.Temperature)
		}
	}
}
