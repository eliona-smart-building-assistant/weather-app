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
	"github.com/eliona-smart-building-assistant/go-eliona/apps"
	"github.com/eliona-smart-building-assistant/go-eliona/common"
	"github.com/eliona-smart-building-assistant/go-eliona/db"
	"github.com/eliona-smart-building-assistant/go-eliona/log"
	"weather/conf"
	"weather/weather"
)

// The main function starts the app by starting all services necessary for this app and waits
// until all services are finished. In most cases the services run infinite, except the app is stopped
// externally, e.g. during a shut-down of the eliona environment.
func main() {
	log.Info("Weather", "Starting the app.")

	// Necessary to close used init resources, because db.Pool() is used in this app.
	defer db.ClosePool()

	// Init the app before the first run.
	apps.Init(db.Pool(), common.AppName(),
		apps.ExecSqlFile("conf/init.sql"),
		conf.InitConfiguration,
		conf.InitAssetType,
		conf.InitAssets,
	)

	// Patch the app v1.1.0
	apps.Patch(db.Pool(), common.AppName(), "010100",
		conf.AddDaytimeAttribute,
	)

	// Starting the service for the weather app. Normally one app has only one service. In case of the
	// weather app, the service reads weather data for configurable locations and write this data as heap
	// back to the eliona environment.
	apps.WaitFor(
		apps.Loop(weather.CollectData, conf.PollingInterval()),
	)

	log.Info("Weather", "Terminate the app.")
}
