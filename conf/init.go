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
	"github.com/eliona-smart-building-assistant/go-eliona/assets"
	"github.com/eliona-smart-building-assistant/go-eliona/db"
)

// InitAssetType creates asset type for weather locations
func InitAssetType(connection db.Connection) error {
	err := assets.UpsertAssetType(connection, assets.AssetType{
		Id:               "weather_location",
		Vendor:           "weatherDB by Dron Bhattacharya & Rituraj Datta",
		Translation:      assets.Translation{German: "Wetterstandort", English: "Weather location"},
		DocumentationUrl: "https://weatherdbi.herokuapp.com/documentation/v1",
		Icon:             "weather",
	})
	if err != nil {
		return err
	}

	err = assets.UpsertAssetTypeAttribute(connection, assets.AssetTypeAttribute{
		AssetTypeId:   "weather_location",
		AttributeType: "humidity",
		Id:            "humidity",
		Subtype:       assets.InputSubtype,
		Translation:   assets.Translation{German: "Luftfeuchte", English: "Humidity"},
		Enable:        true,
		Unit:          "%",
	})
	if err != nil {
		return err
	}

	err = assets.UpsertAssetTypeAttribute(connection, assets.AssetTypeAttribute{
		AssetTypeId:   "weather_location",
		AttributeType: "weather",
		Id:            "precipitation",
		Subtype:       assets.InputSubtype,
		Translation:   assets.Translation{German: "Niederschlag", English: "Precipitation"},
		Enable:        true,
		Unit:          "%",
	})
	if err != nil {
		return err
	}

	err = assets.UpsertAssetTypeAttribute(connection, assets.AssetTypeAttribute{
		AssetTypeId:   "weather_location",
		AttributeType: "weather",
		Id:            "wind",
		Subtype:       assets.InputSubtype,
		Translation:   assets.Translation{German: "Wind", English: "Wind"},
		Enable:        true,
		Unit:          "km/h",
	})
	if err != nil {
		return err
	}

	err = assets.UpsertAssetTypeAttribute(connection, assets.AssetTypeAttribute{
		AssetTypeId:   "weather_location",
		AttributeType: "temperature",
		Id:            "temperature",
		Subtype:       assets.InputSubtype,
		Translation:   assets.Translation{German: "Temperatur", English: "Temperature"},
		Enable:        true,
		Unit:          "°C",
	})
	if err != nil {
		return err
	}

	err = assets.UpsertAssetTypeAttribute(connection, assets.AssetTypeAttribute{
		AssetTypeId:   "weather_location",
		AttributeType: "weather",
		Id:            "comment",
		Subtype:       assets.StatusSubtype,
		Translation:   assets.Translation{German: "Kommentar", English: "Comment"},
		Enable:        true,
	})
	if err != nil {
		return err
	}

	return nil
}

// InitAssets adds example assets for weather locations in Switzerland.
// This should be made editable by eliona frontend.
func InitAssets(connection db.Connection) error {
	id, err := assets.UpsertAsset(connection, assets.Asset{
		ProjectId:             defaultProjectId(connection),
		GlobalAssetIdentifier: "Winterthur, Schweiz",
		Description:           "Winterthur, Schweiz",
		Name:                  "Winterthur",
		Latitude:              47.5056400,
		Longitude:             8.7241300,
		AssetTypeId:           "weather_location",
	})
	if err != nil {
		return err
	}
	err = insertLocation(connection, id, "winterthur")
	if err != nil {
		return err
	}

	id, err = assets.UpsertAsset(connection, assets.Asset{
		ProjectId:             defaultProjectId(connection),
		GlobalAssetIdentifier: "Zürich, Schweiz",
		Description:           "Zürich, Schweiz",
		Name:                  "Zürich",
		Latitude:              47.3666700,
		Longitude:             8.5500000,
		AssetTypeId:           "weather_location",
	})
	if err != nil {
		return err
	}
	err = insertLocation(connection, id, "zurich")
	if err != nil {
		return err
	}

	id, err = assets.UpsertAsset(connection, assets.Asset{
		ProjectId:             defaultProjectId(connection),
		GlobalAssetIdentifier: "Bern, Schweiz",
		Description:           "Bern, Schweiz",
		Name:                  "Bern",
		Latitude:              47.5056400,
		Longitude:             8.7241300,
		AssetTypeId:           "weather_location",
	})
	if err != nil {
		return err
	}
	err = insertLocation(connection, id, "bern")
	if err != nil {
		return err
	}

	return nil
}

func InitConfiguration(connection db.Connection) error {
	err := Set("endpoint", "https://weatherdbi.herokuapp.com/data/weather/")
	if err != nil {
		return err
	}
	err = Set("polling_interval", "10")
	if err != nil {
		return err
	}
	return nil
}

func insertLocation(connection db.Connection, assetId int, location string) error {
	return db.Exec(connection, "insert into weather.locations (asset_id, location) values ($1, $2)", assetId, location)
}

func defaultProjectId(connection db.Connection) string {
	id, _ := db.QuerySingleRow[string](connection, "select min(proj_id) from public.eliona_project")
	return id
}
