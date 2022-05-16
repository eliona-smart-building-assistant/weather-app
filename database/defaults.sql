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

-- This script is called by apps.Init() before the app is launched for the first time. After that,
-- the execution of this script is skipped.

-- A good practice is to initialize the app configuration with default values and create visual elements for the eliona
-- frontend. So the user can see how this app works and what can be done and what needs to be configured.

-- Default configuration values.
-- This table should be made editable by eliona frontend.
insert into weather.configuration (name, value) values
    ('endpoint', 'https://weatherdbi.herokuapp.com/data/weather/'), -- where is the weatherDB located
    ('polling_interval', '10') -- with interval in seconds is used to poll the weatherDB
    on conflict (name) do update set value = excluded.value;

-- Example assets for weather locations in Switzerland.
-- This table should be made editable by eliona frontend.
insert into public.asset (gai, name, proj_id, asset_type, tags) values
    ('47.5056400, 8.7241300', 'Winterthur', (select min(proj_id) from public.eliona_project), 'weather_location', '{"hs:dis": "Winterthur, Schweiz"}'),
    ('47.3666700, 8.5500000', 'Zürich', (select min(proj_id) from public.eliona_project), 'weather_location', '{"hs:dis": "Zürich, Schweiz"}'),
    ('46.9480900, 7.4474400', 'Bern', (select min(proj_id) from public.eliona_project), 'weather_location', '{"hs:dis": "Bern, Schweiz"}')
    on conflict (gai, proj_id) do update set name = excluded.name, tags = excluded.tags;

-- Map the example assets above to real locations.
-- This table should be made editable by eliona frontend.
insert into weather.locations (location, asset_id) values
    ('winterthur', (select asset_id from public.asset where gai = '47.5056400, 8.7241300' and proj_id = (select min(proj_id) from public.eliona_project))),
    ('zurich', (select asset_id from public.asset where gai = '47.3666700, 8.5500000' and proj_id = (select min(proj_id) from public.eliona_project))),
    ('bern', (select asset_id from public.asset where gai = '46.9480900, 7.4474400' and proj_id = (select min(proj_id) from public.eliona_project)))
    on conflict do nothing;
