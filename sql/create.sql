--
-- PostgreSQL database dump
--

-- Dumped from database version 10.11 (Ubuntu 10.11-1.pgdg16.04+1)
-- Dumped by pg_dump version 10.11 (Ubuntu 10.11-1.pgdg16.04+1)

-- Started on 2019-12-01 23:48:38 MSK

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

--
-- TOC entry 3002 (class 1262 OID 33514)
-- Name: roadmaps; Type: DATABASE; Schema: -; Owner: -
--

CREATE DATABASE roadmaps WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'en_US.UTF-8' LC_CTYPE = 'en_US.UTF-8';


connect roadmaps

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

--
-- TOC entry 1 (class 3079 OID 13051)
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- TOC entry 3004 (class 0 OID 0)
-- Dependencies: 1
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET default_with_oids = false;

--
-- TOC entry 200 (class 1259 OID 33548)
-- Name: plans; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.plans (
    id integer NOT NULL,
    title character varying(256) NOT NULL,
    topic character varying(128) NOT NULL,
    owner character varying(36) NOT NULL,
    points integer DEFAULT 0 NOT NULL
);


--
-- TOC entry 199 (class 1259 OID 33546)
-- Name: plans_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.plans_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 3005 (class 0 OID 0)
-- Dependencies: 199
-- Name: plans_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.plans_id_seq OWNED BY public.plans.id;


--
-- TOC entry 204 (class 1259 OID 33583)
-- Name: sources; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sources (
    id bigint NOT NULL,
    title character varying(256) NOT NULL,
    identifier character varying(256) NOT NULL,
    normalizedidentifier character varying(256) NOT NULL,
    type character varying(24) NOT NULL,
    properties character varying(1024),
    img character varying(256),
    description character varying(4096)
);


--
-- TOC entry 203 (class 1259 OID 33581)
-- Name: sources_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.sources_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 3006 (class 0 OID 0)
-- Dependencies: 203
-- Name: sources_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.sources_id_seq OWNED BY public.sources.id;


--
-- TOC entry 202 (class 1259 OID 33569)
-- Name: steps; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.steps (
    id bigint NOT NULL,
    planid integer NOT NULL,
    referenceid bigint NOT NULL,
    referencetype character varying(24) NOT NULL,
    "position" integer NOT NULL
);


--
-- TOC entry 201 (class 1259 OID 33567)
-- Name: steps_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.steps_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 3007 (class 0 OID 0)
-- Dependencies: 201
-- Name: steps_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.steps_id_seq OWNED BY public.steps.id;


--
-- TOC entry 205 (class 1259 OID 33596)
-- Name: steps_sources; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.steps_sources (
    stepid bigint NOT NULL,
    sourceid bigint NOT NULL
);


--
-- TOC entry 197 (class 1259 OID 33517)
-- Name: topics; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.topics (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    title character varying(100) NOT NULL,
    description character varying(1024),
    creator character varying(64) NOT NULL
);


--
-- TOC entry 196 (class 1259 OID 33515)
-- Name: topics_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.topics_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 3008 (class 0 OID 0)
-- Dependencies: 196
-- Name: topics_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.topics_id_seq OWNED BY public.topics.id;


--
-- TOC entry 198 (class 1259 OID 33526)
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id character varying(36) NOT NULL,
    name character varying(64) NOT NULL,
    normalizedname character varying(100) NOT NULL,
    email character varying(128) NOT NULL,
    emailconfirmed boolean NOT NULL,
    emailconfirmation character varying(64),
    img character varying(64),
    tokens character varying(4096),
    rights integer NOT NULL,
    password bytea NOT NULL,
    salt bytea NOT NULL
);


--
-- TOC entry 206 (class 1259 OID 33613)
-- Name: usersplans; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.usersplans (
    userid character varying(36) NOT NULL,
    topic character varying(100) NOT NULL,
    planid integer NOT NULL
);


--
-- TOC entry 2831 (class 2604 OID 33551)
-- Name: plans id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.plans ALTER COLUMN id SET DEFAULT nextval('public.plans_id_seq'::regclass);


--
-- TOC entry 2834 (class 2604 OID 33586)
-- Name: sources id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sources ALTER COLUMN id SET DEFAULT nextval('public.sources_id_seq'::regclass);


--
-- TOC entry 2833 (class 2604 OID 33572)
-- Name: steps id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.steps ALTER COLUMN id SET DEFAULT nextval('public.steps_id_seq'::regclass);


--
-- TOC entry 2830 (class 2604 OID 33520)
-- Name: topics id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.topics ALTER COLUMN id SET DEFAULT nextval('public.topics_id_seq'::regclass);


--
-- TOC entry 2849 (class 2606 OID 33554)
-- Name: plans plans_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.plans
    ADD CONSTRAINT plans_pkey PRIMARY KEY (id);


--
-- TOC entry 2854 (class 2606 OID 33591)
-- Name: sources sources_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sources
    ADD CONSTRAINT sources_pkey PRIMARY KEY (id);


--
-- TOC entry 2852 (class 2606 OID 33574)
-- Name: steps steps_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.steps
    ADD CONSTRAINT steps_pkey PRIMARY KEY (id);


--
-- TOC entry 2862 (class 2606 OID 33600)
-- Name: steps_sources steps_sources_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.steps_sources
    ADD CONSTRAINT steps_sources_pkey PRIMARY KEY (stepid, sourceid);


--
-- TOC entry 2837 (class 2606 OID 33525)
-- Name: topics topics_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.topics
    ADD CONSTRAINT topics_pkey PRIMARY KEY (id);


--
-- TOC entry 2856 (class 2606 OID 33595)
-- Name: sources u_sources_identifier; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sources
    ADD CONSTRAINT u_sources_identifier UNIQUE (identifier);


--
-- TOC entry 2858 (class 2606 OID 33593)
-- Name: sources u_sources_normalizedidentifier; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sources
    ADD CONSTRAINT u_sources_normalizedidentifier UNIQUE (normalizedidentifier);


--
-- TOC entry 2839 (class 2606 OID 33545)
-- Name: topics u_topics_name; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.topics
    ADD CONSTRAINT u_topics_name UNIQUE (name);


--
-- TOC entry 2841 (class 2606 OID 33535)
-- Name: users u_users_email; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT u_users_email UNIQUE (email);


--
-- TOC entry 2843 (class 2606 OID 33537)
-- Name: users u_users_name; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT u_users_name UNIQUE (normalizedname);


--
-- TOC entry 2845 (class 2606 OID 33533)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- TOC entry 2846 (class 1259 OID 33566)
-- Name: fki_fk_plans_owner; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX fki_fk_plans_owner ON public.plans USING btree (owner);


--
-- TOC entry 2847 (class 1259 OID 33560)
-- Name: fki_fk_plans_topic; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX fki_fk_plans_topic ON public.plans USING btree (topic);


--
-- TOC entry 2850 (class 1259 OID 33580)
-- Name: fki_fk_steps_planid; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX fki_fk_steps_planid ON public.steps USING btree (planid);


--
-- TOC entry 2859 (class 1259 OID 33612)
-- Name: fki_fk_steps_sources_sourceid; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX fki_fk_steps_sources_sourceid ON public.steps_sources USING btree (sourceid);


--
-- TOC entry 2860 (class 1259 OID 33606)
-- Name: fki_fk_steps_sources_stepid; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX fki_fk_steps_sources_stepid ON public.steps_sources USING btree (stepid);


--
-- TOC entry 2835 (class 1259 OID 33543)
-- Name: fki_fk_topics_creator; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX fki_fk_topics_creator ON public.topics USING btree (creator);


--
-- TOC entry 2863 (class 1259 OID 33634)
-- Name: fki_fk_userplans_planid; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX fki_fk_userplans_planid ON public.usersplans USING btree (planid);


--
-- TOC entry 2864 (class 1259 OID 33628)
-- Name: fki_fk_userplans_topic; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX fki_fk_userplans_topic ON public.usersplans USING btree (topic);


--
-- TOC entry 2865 (class 1259 OID 33622)
-- Name: fki_fk_usersplans_userid; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX fki_fk_usersplans_userid ON public.usersplans USING btree (userid);


--
-- TOC entry 2866 (class 1259 OID 33616)
-- Name: ix_usersplans_userid_topic; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX ix_usersplans_userid_topic ON public.usersplans USING btree (userid, topic);


--
-- TOC entry 2869 (class 2606 OID 33561)
-- Name: plans fk_plans_owner; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.plans
    ADD CONSTRAINT fk_plans_owner FOREIGN KEY (owner) REFERENCES public.users(id) ON UPDATE CASCADE NOT VALID;


--
-- TOC entry 2868 (class 2606 OID 33555)
-- Name: plans fk_plans_topic; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.plans
    ADD CONSTRAINT fk_plans_topic FOREIGN KEY (topic) REFERENCES public.topics(name) ON UPDATE CASCADE NOT VALID;


--
-- TOC entry 2870 (class 2606 OID 33575)
-- Name: steps fk_steps_planid; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.steps
    ADD CONSTRAINT fk_steps_planid FOREIGN KEY (planid) REFERENCES public.plans(id) ON UPDATE CASCADE NOT VALID;


--
-- TOC entry 2872 (class 2606 OID 33607)
-- Name: steps_sources fk_steps_sources_sourceid; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.steps_sources
    ADD CONSTRAINT fk_steps_sources_sourceid FOREIGN KEY (sourceid) REFERENCES public.sources(id) ON UPDATE CASCADE NOT VALID;


--
-- TOC entry 2871 (class 2606 OID 33601)
-- Name: steps_sources fk_steps_sources_stepid; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.steps_sources
    ADD CONSTRAINT fk_steps_sources_stepid FOREIGN KEY (stepid) REFERENCES public.steps(id) ON UPDATE CASCADE NOT VALID;


--
-- TOC entry 2867 (class 2606 OID 33538)
-- Name: topics fk_topics_creator; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.topics
    ADD CONSTRAINT fk_topics_creator FOREIGN KEY (creator) REFERENCES public.users(id) ON UPDATE CASCADE NOT VALID;


--
-- TOC entry 2875 (class 2606 OID 33629)
-- Name: usersplans fk_userplans_planid; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.usersplans
    ADD CONSTRAINT fk_userplans_planid FOREIGN KEY (planid) REFERENCES public.plans(id) ON UPDATE CASCADE NOT VALID;


--
-- TOC entry 2874 (class 2606 OID 33623)
-- Name: usersplans fk_userplans_topic; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.usersplans
    ADD CONSTRAINT fk_userplans_topic FOREIGN KEY (topic) REFERENCES public.topics(name) ON UPDATE CASCADE NOT VALID;


--
-- TOC entry 2873 (class 2606 OID 33617)
-- Name: usersplans fk_usersplans_userid; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.usersplans
    ADD CONSTRAINT fk_usersplans_userid FOREIGN KEY (userid) REFERENCES public.users(id) NOT VALID;


-- Completed on 2019-12-01 23:48:38 MSK

--
-- PostgreSQL database dump complete
--

