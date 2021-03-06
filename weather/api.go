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
	"encoding/json"
	"fmt"
	"github.com/eliona-smart-building-assistant/go-utils/http"
	"time"
	"weather/conf"
)

// Conditions holds the current weather condition for the current day
type Conditions struct {
	Humidity      int     `json:"humidity"`
	Precipitation int     `json:"precipitation"`
	Wind          float64 `json:"wind"`
	Temperature   float64 `json:"temperature"`
	Comment       string  `json:"comment"`
	Daytime       string  `json:"daytime"`
}

// Today returns the current weather conditions for the given latitude and longitude.
func Today(location conf.Location) (Conditions, error) {
	var conditions Conditions

	// Request the API to get current conditions conditions
	url, payload, err := request(location)
	if err != nil {
		return conditions, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(payload, &result)
	if err != nil {
		return conditions, err
	}
	status, ok := result["status"].(string)
	if ok && status == "fail" {
		return conditions, fmt.Errorf("error requesting api %s: %s", url, result["message"].(string))
	}
	conditions.Temperature, _ = result["dataseries"].([]interface{})[0].(map[string]interface{})["temp2m"].(map[string]interface{})["max"].(float64)
	conditions.Wind, _ = result["dataseries"].([]interface{})[0].(map[string]interface{})["wind10m_max"].(float64)
	conditions.Comment, _ = result["dataseries"].([]interface{})[0].(map[string]interface{})["weather"].(string)
	date, _ := result["dataseries"].([]interface{})[0].(map[string]interface{})["date"].(float64)
	conditions.Daytime = fmt.Sprintf("%0.0f", date)
	return conditions, nil
}

// request calls the api to get structured weather data
func request(location conf.Location) (string, []byte, error) {
	url := fmt.Sprintf(conf.Endpoint()+"&lon=%f&lat=%f", location.Longitude, location.Latitude)
	request, err := http.NewRequest(url)
	if err != nil {
		return url, nil, err
	}
	payload, err := http.Do(request, time.Second*10, true)
	return url, payload, err
}
