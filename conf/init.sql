--  This file is part of the eliona project.
--  Copyright © 2022 LEICOM iTEC AG. All Rights Reserved.
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

-- Creates a schema named like the app within the eliona init.
-- All persistent and configurable data needed by the app should be located in this schema.
create schema if not exists weather;

-- Committing schema creation because this cannot be wrapped inside transactions
commit;

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
