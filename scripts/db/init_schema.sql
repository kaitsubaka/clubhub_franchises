
CREATE DATABASE franchises_db WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'en_US.utf8' LC_CTYPE = 'en_US.utf8';

ALTER DATABASE franchises_db OWNER TO postgres;

\connect franchises_db

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_with_oids = false;

create table if not EXISTS public.pending_franchises (
  id uuid,
  url VARCHAR(128) not null,
  status VARCHAR(128) not null,
  error VARCHAR(256),
  created_at TIMESTAMPTZ default NOW(),
  updated_at TIMESTAMPTZ default NOW(),
  deleted_at TIMESTAMPTZ,
  primary key (id)
);

create table if not EXISTS public.countries (
 id SERIAL,
  name VARCHAR(128) not null,
  created_at TIMESTAMPTZ default NOW(),
  updated_at TIMESTAMPTZ default NOW(),
  deleted_at TIMESTAMPTZ,
  primary key (id)
);

create table if not EXISTS public.cities (
  id SERIAL,
  name VARCHAR(64) not null,
  country_id int not null,
  created_at TIMESTAMPTZ default NOW(),
  updated_at TIMESTAMPTZ default NOW(),
  deleted_at TIMESTAMPTZ,
  primary key (id),
  constraint fk_country
    foreign key(country_id) 
	    references public.countries(id)
);

CREATE TABLE IF NOT EXISTS public.locations (
  id SERIAL,
  address VARCHAR(128) NOT NULL,
  zip_code VARCHAR(64) NOT NULL,
  city_id int NOT NULL,
  created_at TIMESTAMPTZ default NOW(),
  updated_at TIMESTAMPTZ default NOW(),
  deleted_at TIMESTAMPTZ,
  PRIMARY KEY (id),
  CONSTRAINT fk_city
    FOREIGN KEY(city_id) 
	    references public.cities(id)
);

CREATE TABLE IF NOT EXISTS public.companies (
  id uuid,
  owner_id uuid NOT NULL,
  name VARCHAR(64) NOT NULL,
  tax_number VARCHAR(64) NOT NULL,
  location_id int NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,
  PRIMARY KEY (id),
  CONSTRAINT fk_location
    FOREIGN KEY(location_id) 
	    references public.locations(id)
);

CREATE TABLE IF NOT EXISTS public.franchises (
  id uuid,
  company_id uuid NOT NULL,
  title VARCHAR(128) NOT NULL,
  site_name VARCHAR(128) NOT NULL,
  description VARCHAR(512) NOT NULL,
  image VARCHAR(128) NOT NULL,
  url VARCHAR(128) NOT NULL,
  protocol VARCHAR(16) NOT NULL,
  domain_jumps SMALLINT NOT NULL,
  server_names TEXT[] NOT NULL,
  domain_creation_date TIMESTAMPTZ NOT NULL,
  domain_expiration_date TIMESTAMPTZ NOT NULL,
  registrant_name VARCHAR(128) NOT NULL,
  contact_email VARCHAR(128) NOT NULL,
  location_id int NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,
  PRIMARY KEY (id),
  CONSTRAINT fk_location
    FOREIGN KEY(location_id) 
	    references public.locations(id),
  CONSTRAINT fk_pending_franchise
    FOREIGN KEY(id) 
      references public.pending_franchises(id),
  CONSTRAINT fk_company
    FOREIGN KEY(company_id) 
      references public.companies(id)
);
