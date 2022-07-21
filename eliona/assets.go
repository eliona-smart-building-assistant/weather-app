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

package eliona

import (
	api "github.com/eliona-smart-building-assistant/go-eliona-api-client"
	"github.com/eliona-smart-building-assistant/go-eliona/asset"
	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/db"
	"github.com/eliona-smart-building-assistant/go-utils/log"
	"time"
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

// InitAssetType creates asset type for weather locations
func InitAssetType(db.Connection) error {
	err := asset.UpsertAssetType(api.AssetType{
		Name:        "weather_location",
		Custom:      true,
		Vendor:      common.Ptr("7Timer!"),
		Translation: &api.Translation{De: common.Ptr("Wetterstandort"), En: common.Ptr("Weather location")},
		Urldoc:      common.Ptr("https://www.7timer.info/doc.php?lang=en#api"),
		Icon:        common.Ptr("weather"),
		Attributes: []api.AssetTypeAttribute{
			{
				Type:        common.Ptr("humidity"),
				Name:        "humidity",
				Subtype:     api.INPUT,
				Translation: &api.Translation{De: common.Ptr("Luftfeuchte"), En: common.Ptr("Humidity")},
				Enable:      common.Ptr(true),
				Unit:        common.Ptr("%"),
			},
			{
				Type:        common.Ptr("weather"),
				Name:        "precipitation",
				Subtype:     api.INPUT,
				Translation: &api.Translation{De: common.Ptr("Niederschlag"), En: common.Ptr("Precipitation")},
				Enable:      common.Ptr(true),
				Unit:        common.Ptr("%"),
			},
			{
				Type:        common.Ptr("weather"),
				Name:        "wind",
				Subtype:     api.INPUT,
				Translation: &api.Translation{De: common.Ptr("Wind"), En: common.Ptr("Wind")},
				Enable:      common.Ptr(true),
				Unit:        common.Ptr("km/h"),
			},
			{
				Type:        common.Ptr("temperature"),
				Name:        "temperature",
				Subtype:     api.INPUT,
				Translation: &api.Translation{De: common.Ptr("Temperatur"), En: common.Ptr("Temperature")},
				Enable:      common.Ptr(true),
				Unit:        common.Ptr("°C"),
			},
			{
				Type:        common.Ptr("weather"),
				Name:        "comment",
				Subtype:     api.STATUS,
				Translation: &api.Translation{De: common.Ptr("Kommentar"), En: common.Ptr("Comment")},
				Enable:      common.Ptr(true),
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func defaultProjectId(connection db.Connection) string {
	id, _ := db.QuerySingleRow[string](connection, "select min(proj_id) from public.eliona_project")
	return id
}

func UpsertHeap(subtype api.HeapSubtype, assetId int32, data any) {
	var statusHeap api.Heap
	statusHeap.Subtype = subtype
	statusHeap.Timestamp = common.Ptr(time.Now())
	statusHeap.AssetId = assetId
	statusHeap.Data = common.StructToMap(data)
	err := asset.UpsertHeap(statusHeap)
	if err != nil {
		log.Error("Weather", "Error during writing heap: %v", err)
		return
	}
}
