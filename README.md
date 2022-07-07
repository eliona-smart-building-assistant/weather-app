# Weather app
The [weather app](https://github.com/eliona-smart-building-assistant/weather-app) is an exemplary app build for running within an [eliona](https://www.eliona.io/) environment. Apps like this add new features to eliona like supporting additional sensors or open new communication channels.

This app grabs weather data from [7timer](https://www.7timer.info/bin/civillight.php?ac=0&unit=metric&output=json&tzshift=0) web service this data to eliona. To do this, you can configure locations for which the weather data is to be read. Within eliona these locations are handled as assets and can be used to show on dashboards, trigger alarms and so on.

[<img src="weather-app.png" width="350"/>](weather-app.png)

## Configuration ##

The app needs environment variables and database tables for configuration.

### Eliona ###

To start and initialize an app in an Eliona environment, the app have to registered in Eliona. For this, an entry in the database table `public.eliona_app` is necessary.

### Environment variables ###

The `APPNAME` MUST be set to `weather`. Some resources use this name to identify the app inside an Eliona environment. For running as a Docker container inside an Eliona environment, the `APPNAME` have to set in the [Dockerfile](Dockerfile). If the app runs outside an Eliona environment the `APPNAME` must be set explicitly.

```bash
export APPNAME=weather
```

#### CONNECTION_STRING

The `CONNECTION_STRING` variable configures the [Eliona database](https://github.com/eliona-smart-building-assistant/go-eliona/tree/main/db). If the app runs as a Docker container inside an Eliona environment, the environment must provide this variable. If you run the app standalone you must set this variable. Otherwise, the app can't be initialized and started.

```bash
export CONNECTION_STRING=postgres://user:pass@localhost:5432/iot
```

#### API_ENDPOINT

The `API_ENDPOINT` variable configures the endpoint to access the [Eliona API](https://github.com/eliona-smart-building-assistant/eliona-api). If the app runs as a Docker container inside an Eliona environment, the environment must provide this variable. If you run the app standalone you must set this variable. Otherwise, the app can't be initialized and started.

```bash
export API_ENDPOINT=http://localhost:8082/v2
```

#### DEBUG_LEVEL (optional)

The `DEBUG_LEVEL` variable defines the minimum level that should be [logged](https://github.com/eliona-smart-building-assistant/go-eliona/tree/main/log). Not defined the default level is `info`.

```bash
export LOG_LEVEL=debug # optionally, default is 'info'
```

### Configuration tables ###

The app requires some configuration data that remains in the database. To do this, the app creates its own database schema `weather` during initialization. The data in this schema should be made editable by eliona frontend. This allows the app to be configured by the user without direct database access.

A good practice is to initialize the app configuration with [default values](sql/defaults.sql). This allows the user to see how what needs to be configured.

In detail, you need the following configuration data in table `weather.configuration (name, value)`.

```sql
-- weather.configuration (name, value)
('endpoint', 'https://www.7timer.info/bin/civillight.php?ac=0&unit=metric&output=json&tzshift=0') -- where is the API located
('polling_interval', '10') -- with interval in seconds is used to poll the API 
```
    
In order to define the weather locations for which conditions are to be read, an entry in the table `weather.locations (location, latitude, longitude, proj_id)` is required. The location is the name of the location, e.g. `winterthur,swizerland` or if clearly just `winterthur`. Each location is later mapped with an asset in eliona.

## API Reference

The weather app grabs weather conditions from [7timer](https://www.7timer.info/bin/civillight.php?ac=0&unit=metric&output=json&tzshift=0) web service and writes these data to eliona as heap data of assets. The heap data is separated in `weather.Input`, `weather.Info` and `weather.Status` heaps. These structures are used to write the heap data.

```json
{"wind": 6, "humidity": 97, "temperature": 18, "precipitation": 15} 
{"daytime": "Monday 8:00 AM"}
{"comment": "Mostly cloudy"}
```

In eliona these heaps are handled as `weather_location` asset type with appropriate attributes created during the [initialization](init/init.sql).


