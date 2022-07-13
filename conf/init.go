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
)

// InitConfiguration adds example configuration
// This should be made editable by eliona frontend.
func InitConfiguration(connection db.Connection) error {
	err := Set(connection, "endpoint", "https://www.7timer.info/bin/civillight.php?ac=0&unit=metric&output=json&tzshift=0")
	if err != nil {
		return err
	}
	err = Set(connection, "polling_interval", "10")
	if err != nil {
		return err
	}
	return nil
}

// InitLocations adds example assets for weather locations in Switzerland.
// This should be made editable by eliona frontend.
func InitLocations(connection db.Connection) error {

	err := InsertLocation(connection, Location{
		Location:  "Winterthur, Schweiz",
		Latitude:  47.5056400,
		Longitude: 8.7241300,
	})
	if err != nil {
		return err
	}

	err = InsertLocation(connection, Location{
		Location:  "Zürich, Schweiz",
		Latitude:  47.3666700,
		Longitude: 8.5500000,
	})
	if err != nil {
		return err
	}

	err = InsertLocation(connection, Location{
		Location:  "Bern, Schweiz",
		Latitude:  47.5056400,
		Longitude: 8.7241300,
	})
	if err != nil {
		return err
	}

	return nil
}
