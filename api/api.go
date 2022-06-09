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

package api

import (
	"encoding/json"
	"fmt"
	"github.com/eliona-smart-building-assistant/go-eliona/http"
	"strconv"
	"strings"
	"time"
	"weather/conf"
)

// Conditions holds the current weather condition for the current day
type Conditions struct {
	Humidity      int
	Precipitation int
	Wind          float64
	Temperature   float64
	Comment       string
	Daytime       string
}

// Today returns the current weather conditions for the given latitude and longitude.
func Today(location string) (Conditions, error) {
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
	conditions.Temperature, _ = result["currentConditions"].(map[string]interface{})["temp"].(map[string]interface{})["c"].(float64)
	conditions.Precipitation, _ = strconv.Atoi(strings.Replace(result["currentConditions"].(map[string]interface{})["precip"].(string), "%", "", -1))
	conditions.Humidity, _ = strconv.Atoi(strings.Replace(result["currentConditions"].(map[string]interface{})["humidity"].(string), "%", "", -1))
	conditions.Wind, _ = result["currentConditions"].(map[string]interface{})["wind"].(map[string]interface{})["km"].(float64)
	conditions.Comment, _ = result["currentConditions"].(map[string]interface{})["comment"].(string)
	conditions.Daytime, _ = result["currentConditions"].(map[string]interface{})["dayhour"].(string)
	return conditions, nil
}

// request calls the api to get structured weather data
func request(location string) (string, []byte, error) {
	url := fmt.Sprintf(conf.Endpoint()+"%s", location)
	request, err := http.NewRequest(url)
	if err != nil {
		return url, nil, err
	}
	payload, err := http.Do(request, time.Second*10, true)
	return url, payload, err
}
