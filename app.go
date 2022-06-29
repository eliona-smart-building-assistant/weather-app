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

package main

import (
	"fmt"
	"github.com/eliona-smart-building-assistant/go-eliona/api"
	"github.com/eliona-smart-building-assistant/go-eliona/asset"
	"github.com/eliona-smart-building-assistant/go-eliona/log"
	"weather/conf"
	"weather/eliona"
	"weather/weather"
)

// CollectData reads the defined weather location from configuration and writes the data as eliona heap
func CollectData() {

	locations := make(chan conf.Location)
	go func() {
		_ = conf.ReadLocations(locations)
	}()
	for location := range locations {

		// Check if eliona project is defined; if not, ignore this location
		if len(location.ProjectId) == 0 {
			log.Warn("Weather", "Ignoring location, because no project id is defined: %s", location.Location)
			continue
		}

		// Create or update asset and get ID
		assetId, err := asset.UpsertAsset(api.Asset{
			ProjectId:             location.ProjectId,
			GlobalAssetIdentifier: fmt.Sprintf("%s %f %f", location.Location, location.Latitude, location.Longitude),
			Description:           &location.Location,
			Name:                  &location.Location,
			Latitude:              &location.Latitude,
			Longitude:             &location.Longitude,
			AssetType:             "weather_location",
		})
		log.Debug("Weather", "Determining asset id %d for location '%s'", *assetId, location.Location)

		// Reads the current weather condition for location
		condition, err := weather.Today(location)
		if err != nil {
			log.Error("Weather", "Error during requesting API endpoint: %v", err)
			return
		}
		log.Debug("Weather", "New condition for location '%s' found: %s", location.Location, condition.Comment)

		// Writes input data as heap
		eliona.UpsertHeap(api.INPUT, *assetId, eliona.Input{
			Temperature:   condition.Temperature,
			Wind:          condition.Wind,
			Humidity:      condition.Humidity,
			Precipitation: condition.Precipitation,
		})

		// Writes info data as heap
		eliona.UpsertHeap(api.INFO, *assetId, eliona.Info{
			Daytime: condition.Daytime,
		})

		// Writes status data as heap
		eliona.UpsertHeap(api.STATUS, *assetId, eliona.Status{
			Comment: condition.Comment,
		})
	}
}
