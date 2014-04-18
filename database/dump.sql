--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

SET search_path = public, pg_catalog;

DROP INDEX public.exhibition_substring_idx;
DROP INDEX public.date_range;
ALTER TABLE ONLY public.gallery DROP CONSTRAINT gallery_pkey;
ALTER TABLE ONLY public.exhibition DROP CONSTRAINT exhibition_pkey;
DROP TABLE public.gallery;
DROP TABLE public.exhibition;
DROP EXTENSION plpgsql;
DROP SCHEMA public;
--
-- Name: public; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA public;


--
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON SCHEMA public IS 'standard public schema';


--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_with_oids = false;

--
-- Name: exhibition; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE exhibition (
    id character varying(500) NOT NULL,
    _byteid bytea NOT NULL,
    gallery_id uuid NOT NULL,
    title character varying(500) NOT NULL,
    description character varying(5000) NOT NULL,
    date_range daterange NOT NULL,
    created timestamp with time zone DEFAULT ('now'::text)::date,
    updated timestamp with time zone
);


--
-- Name: gallery; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE gallery (
    id uuid NOT NULL,
    name character varying(100) NOT NULL,
    meta json NOT NULL,
    about character varying(2000) NOT NULL,
    created timestamp with time zone DEFAULT ('now'::text)::date,
    updated timestamp with time zone
);


--
-- Name: exhibition_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY exhibition
    ADD CONSTRAINT exhibition_pkey PRIMARY KEY (_byteid);


--
-- Name: gallery_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY gallery
    ADD CONSTRAINT gallery_pkey PRIMARY KEY (id);


--
-- Name: date_range; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX date_range ON exhibition USING gist (date_range);


--
-- Name: exhibition_substring_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX exhibition_substring_idx ON exhibition USING btree ("substring"(_byteid, 5));


--
-- Name: public; Type: ACL; Schema: -; Owner: -
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM tom;
GRANT ALL ON SCHEMA public TO tom;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

