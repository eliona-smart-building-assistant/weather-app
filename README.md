# Weather app
The [weather app](https://github.com/eliona-smart-building-assistant/weather-app) is an exemplary app build for running within an [eliona](https://www.eliona.io/) environment. Apps like this add new features to eliona like supporting additional sensors or open new communication channels.

This app grabs weather data from [WeatherDB](https://weatherdbi.herokuapp.com/) web service this data to eliona. To do this, you can configure locations for which the weather data is to be read. Within eliona these locations are handled as assets and can be used to show on dashboards, trigger alarms and so on.

[<img src="weather-app.png" width="350"/>](weather-app.png)

## Configuration ##

The app needs environment variables and database tables for configuration.

### Environment variables ###

The `APPNAME` MUST be set to `weather`. Some resources use this name to identify the app inside an eliona environment. For running as a Docker container inside an eliona environment, the `APPNAME` have to set in the [Dockerfile](Dockerfile). If the app runs outside an eliona environment the `APPNAME` must be set explicitly.

```bash
export APPNAME=weather # For running in eliona environment set app name in Dockerfile
```

The `CONNECTION_STRING` variable configures the [eliona database](https://github.com/eliona-smart-building-assistant/go-eliona/tree/main/db). If the app runs as a Docker container inside an eliona environment this variable is already set by the environment. If you run the app outside you must provide this variable. Otherwise the app can't be initialized and started. 

```bash
export CONNECTION_STRING=postgres://user:pass@localhost::5432/iot # only if run outside eliona environment
```

The `DEBUG_LEVEL` variable defines the minimum level that should be [logged](https://github.com/eliona-smart-building-assistant/go-eliona/tree/main/log). Not defined the default level is `info`.

```bash
export LOG_LEVEL=debug # This is optionally, default is info
```

### Configuration tables ###

The app requires some configuration data that remains in the database. To do this, the app creates its own database schema `weather` during initialization. The data in this schema should be made editable by eliona frontend. This allows the app to be configured by the user without direct database access.

A good practice is to initialize the app configuration with [default values](database/defaults.sql). This allows the user to see how what needs to be configured.

In detail, you need the following configuration data in table `weather.configuration (name, value)`.

```sql
-- weather.configuration (name, value)
('endpoint', 'https://weatherdbi.herokuapp.com/data/weather/') -- where is the weatherDB located
('polling_interval', '10') -- with interval in seconds is used to poll the weatherDB 
```

In order to define the weather locations for which conditions are to be read, an entry in the table `weather.locations (location, asset_id)` is required. The location is the name of the location, e.g. `winterthur,swizerland` or if clearly just `winterthur`. Each location is mapped with an asset in eliona, represented by the `asset_id`.

Before the locations can be inserted the corresponding asset have to create. The id of this can now use to configure the location.

```sql
-- weather.locations (location, asset_id)
('winterthur', 4711) -- define Winterthur as location and map with eliona asset 4711
```

## API Reference

The weather app grabs weather conditions from [WeatherDB](https://weatherdbi.herokuapp.com/) web service and writes these data to eliona as heap data of assets. The heap data is separated in `weather.Input`, `weather.Info` and `weather.Status` heaps. These structures are used to write the heap data.

```json
{"wind": 6, "humidity": 97, "temperature": 18, "precipitation": 15} 
{"daytime": "Monday 8:00 AM"}
{"comment": "Mostly cloudy"}
```

In eliona these heaps are handled as `weather_location` asset type with appropriate attributes created during the [initialization](database/init.sql).


