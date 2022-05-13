--  This file is part of the eliona project.
--  Copyright © 2022 LEICOM iTEC AG. All Rights Reserved.
--  Authors: Adam Lange, et al.
--  ______ _ _
-- |  ____| (_)
-- | |__  | |_  ___  _ __   __ _
-- |  __| | | |/ _ \| '_ \ / _` |
-- | |____| | | (_) | | | | (_| |
-- |______|_|_|\___/|_| |_|\__,_|
--
--  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
--  BUT NOT LIMITED  TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
--  NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
--  DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
--  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

-- This script (database/init.sql) is called once when the app starts first time within an eliona environment.
-- Put there all initialisations need for this app. For changing installed apps use patch scripts (database/patches/*.sql).

-- Creates a schema named like the app within the eliona database.
-- All persistent and configurable data needed by the app should be located in this schema.
create schema if not exists weather;

-- Create a table for global configuration like endpoints, secrets and so on.
-- This table should be made editable by eliona frontend.
create table if not exists weather.configuration
(
    name text primary key,
    value text not null
);

-- Create a table to map real weather locations with eliona assets.
-- This table should be made editable by eliona frontend.
create table if not exists weather.locations
(
    location text not null,
    asset_id integer unique references public.asset(asset_id) on delete cascade primary key
);

-- Create asset type for weather locations
insert into public.asset_type (asset_type, vendor, translation, urldoc, icon) values
    ('weather_location', 'weatherDB by Dron Bhattacharya & Rituraj Datta', '{"de": "Wetterstandort", "en": "Weather location"}', 'https://weatherdbi.herokuapp.com/documentation/v1', 'weather')
    on conflict(asset_type) do update set vendor = excluded.vendor, translation = excluded.translation, urldoc = excluded.urldoc, icon = excluded.icon;

-- Create attributes to structuring data stored by weather locations
insert into public.attribute_schema (asset_type, attribute_type, attribute, subtype, enable, translation, unit, scale, pipeline_mode, pipeline_raster, viewer, ar) values
    ('weather_location', 'humidity', 'humidity', '', true, '{"de": "Luftfeuchte", "en": "Humidity"}', '%', 0, 'avg', '{M15,H1,DAY}', true, true),
    ('weather_location', 'weather', 'precipitation', '', true, '{"de": "Niederschlag", "en": "Precipitation"}', '%', 0, 'avg', '{M15,H1,DAY}', true, true),
    ('weather_location', 'weather', 'wind', '', true, '{"de": "Wind", "en": "Wind"}', 'km/h', 0, 'avg', '{M15,H1,DAY}', true, true),
    ('weather_location', 'temperature', 'temperature', '', true, '{"de": "Temperatur", "en": "Temperature"}', '°C', 0, 'avg', '{M15,H1,DAY}', true, true),
    ('weather_location', 'weather', 'comment', 'status', true, '{"de": "Kommentar", "en": "Comment"}', null, null, null, null, true, true),
    ('weather_location', null, 'daytime', 'info', true, '{"de": "Tageszeit", "en": "Daytime"}', null, null, null, null, true, true)
    on conflict do nothing;

-- A good practice is to initialize the app configuration with default values and create visual elements for the eliona
-- frontend. So the user can see how this app works, what can be done and what needs to be configured.
\ir defaults.sql
