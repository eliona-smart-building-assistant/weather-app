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

package weather

import (
	"github.com/eliona-smart-building-assistant/go-eliona/assets"
	"github.com/eliona-smart-building-assistant/go-eliona/db"
	"github.com/eliona-smart-building-assistant/go-eliona/log"
	"time"
	"weather/api"
	"weather/conf"
)

// Input is a structure holds input data getting from the api endpoint. This structure corresponds
// to the input heap data in eliona.
type Input struct {
	Humidity      int     `json:"humidity"`
	Precipitation int     `json:"precipitation"`
	Wind          float64 `json:"wind"`
	Temperature   float64 `json:"temperature"`
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

// CollectData reads the defined weather location from configuration and writes the data as eliona heap
func CollectData() {

	locations := make(chan conf.Location)
	go func() {
		_ = conf.ReadLocations(locations)
	}()
	for location := range locations {

		// Reads the current weather condition for location
		condition, err := api.Today(location.Location)
		if err != nil {
			log.Error("Weather", "Error during requesting API endpoint: %v", err)
			return
		}
		log.Debug("Weather", "New condition for location '%s' found: %s", location.Location, condition.Comment)

		// Writes input data as heap
		upsertHeap(assets.InputSubtype, location.AssetId, Input{
			Temperature:   condition.Temperature,
			Wind:          condition.Wind,
			Humidity:      condition.Humidity,
			Precipitation: condition.Precipitation,
		})

		// Writes info data as heap
		upsertHeap(assets.InfoSubtype, location.AssetId, Info{
			Daytime: condition.Daytime,
		})

		// Writes status data as heap
		upsertHeap(assets.StatusSubtype, location.AssetId, Status{
			Comment: condition.Comment,
		})
	}
}

func upsertHeap[T any](subtype assets.Subtype, assetId int, data T) {
	var statusHeap assets.Heap[T]
	statusHeap.Subtype = subtype
	statusHeap.TimeStamp = time.Now()
	statusHeap.AssetId = assetId
	statusHeap.Data = data
	err := assets.UpsertHeap(db.Pool(), statusHeap)
	if err != nil {
		log.Error("Weather", "Error during writing heap: %v", err)
		return
	}
}
