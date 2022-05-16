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

package main

import (
	"github.com/eliona-smart-building-assistant/go-eliona/apps"
	"github.com/eliona-smart-building-assistant/go-eliona/common"
	"github.com/eliona-smart-building-assistant/go-eliona/db"
	"github.com/eliona-smart-building-assistant/go-eliona/log"
	"sync"
	"weather/weather"
)

// The main function starts the app by starting all services necessary for this app and waits
// until all services are finished. In most cases the services run infinite, except the app is stopped
// externally, e.g. during a shut-down of the eliona environment.
func main() {
	log.Info("Weather", "Starting the app.")

	// Init the app for the first run.
	apps.Init(db.Pool(), common.AppName(),
		func(connection db.Connection) error {
			return db.ExecFile(connection, "database/init.sql")
		},
		func(connection db.Connection) error {
			return db.ExecFile(connection, "database/defaults.sql")
		},
	)

	// Patch the app 010100
	apps.Patch(db.Pool(), common.AppName(), "010100",
		func(connection db.Connection) error {
			return db.ExecFile(connection, "database/patches/010100.sql")
		},
	)

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	defer db.ClosePool()

	// Starting the service for the weather app. Normally one app has only one service. In case of the
	// weather app, the service reads weather data for configurable locations and write this data as heap
	// back to the eliona environment.
	go func() {
		weather.StartCollectingData()
		waitGroup.Done()
	}()

	waitGroup.Wait()
	log.Info("Weather", "Terminate the app.")
}
