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

package init

import (
	"github.com/eliona-smart-building-assistant/go-eliona/assets"
	"github.com/eliona-smart-building-assistant/go-eliona/db"
	"weather/conf"
)

// Assets adds example assets for weather locations in Switzerland.
// This should be made editable by eliona frontend.
func Assets(connection db.Connection) error {
	id := assets.UpsertAsset(connection, assets.Asset{ProjectId: defaultProjectId(connection), GlobalAssetIdentifier: "Winterthur, Schweiz", Description: "Winterthur, Schweiz", Name: "Winterthur", Latitude: 47.5056400, Longitude: 8.7241300, AssetTypeId: "weather_location"})
	_ = insertLocation(connection, id, "winterthur")
	id = assets.UpsertAsset(connection, assets.Asset{ProjectId: defaultProjectId(connection), GlobalAssetIdentifier: "Zürich, Schweiz", Description: "Zürich, Schweiz", Name: "Zürich", Latitude: 47.3666700, Longitude: 8.5500000, AssetTypeId: "weather_location"})
	_ = insertLocation(connection, id, "zurich")
	id = assets.UpsertAsset(connection, assets.Asset{ProjectId: defaultProjectId(connection), GlobalAssetIdentifier: "Bern, Schweiz", Description: "Bern, Schweiz", Name: "Bern", Latitude: 47.5056400, Longitude: 8.7241300, AssetTypeId: "weather_location"})
	_ = insertLocation(connection, id, "bern")
	return nil
}

func Configuration(connection db.Connection) error {
	_ = conf.Set("endpoint", "https://weatherdbi.herokuapp.com/data/weather/")
	_ = conf.Set("polling_interval", "10")
	return nil
}

func insertLocation(connection db.Connection, assetId int, location string) error {
	return db.Exec(connection, "insert into weather.locations (asset_id, location) values ($1, $2)", assetId, location)
}

func defaultProjectId(connection db.Connection) string {
	id, _ := db.QuerySingleRow[string](connection, "select min(proj_id) from public.eliona_project")
	return id
}
