--
-- PostgreSQL database dump
--

-- Dumped from database version 17.2
-- Dumped by pg_dump version 17.2

-- Started on 2025-05-06 14:20:12

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 218 (class 1259 OID 16676)
-- Name: deposits; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.deposits (
    id integer NOT NULL,
    commitment character varying,
    depositor character varying,
    leaf_index integer,
    "timestamp" timestamp without time zone,
    tx_hash character varying
);


ALTER TABLE public.deposits OWNER TO postgres;

--
-- TOC entry 217 (class 1259 OID 16675)
-- Name: deposits_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.deposits_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.deposits_id_seq OWNER TO postgres;

--
-- TOC entry 4867 (class 0 OID 0)
-- Dependencies: 217
-- Name: deposits_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.deposits_id_seq OWNED BY public.deposits.id;


--
-- TOC entry 221 (class 1259 OID 16693)
-- Name: kyc; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.kyc (
    citizen_id character varying NOT NULL,
    wallet_address character varying,
    full_name character varying,
    phone_number character varying,
    date_of_birth date,
    nationality character varying,
    kyc_verified_at timestamp without time zone,
    verifier character varying,
    is_active boolean,
    wallet_signature character varying
);


ALTER TABLE public.kyc OWNER TO postgres;

--
-- TOC entry 220 (class 1259 OID 16685)
-- Name: withdrawals; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.withdrawals (
    id integer NOT NULL,
    recipient character varying,
    nullifier_hash character varying,
    relayer character varying,
    fee numeric,
    "timestamp" timestamp without time zone,
    tx_hash character varying
);


ALTER TABLE public.withdrawals OWNER TO postgres;

--
-- TOC entry 219 (class 1259 OID 16684)
-- Name: withdrawals_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.withdrawals_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.withdrawals_id_seq OWNER TO postgres;

--
-- TOC entry 4868 (class 0 OID 0)
-- Dependencies: 219
-- Name: withdrawals_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.withdrawals_id_seq OWNED BY public.withdrawals.id;


--
-- TOC entry 4704 (class 2604 OID 16679)
-- Name: deposits id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.deposits ALTER COLUMN id SET DEFAULT nextval('public.deposits_id_seq'::regclass);


--
-- TOC entry 4705 (class 2604 OID 16688)
-- Name: withdrawals id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.withdrawals ALTER COLUMN id SET DEFAULT nextval('public.withdrawals_id_seq'::regclass);


--
-- TOC entry 4858 (class 0 OID 16676)
-- Dependencies: 218
-- Data for Name: deposits; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.deposits (id, commitment, depositor, leaf_index, "timestamp", tx_hash) FROM stdin;
1	0xabc123	0x1a2b3c	0	2025-05-01 10:00:00	0x123txhash
2	0xdef456	0x4d5e6f	1	2025-05-02 11:00:00	0x456txhash
3	0xghi789	0x7g8h9i	2	2025-05-03 12:00:00	0x789txhash
\.


--
-- TOC entry 4861 (class 0 OID 16693)
-- Dependencies: 221
-- Data for Name: kyc; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.kyc (citizen_id, wallet_address, full_name, phone_number, date_of_birth, nationality, kyc_verified_at, verifier, is_active, wallet_signature) FROM stdin;
ID001	0x1a2b3c	Nguyen Van A	0123456789	1990-01-01	Vietnam	2025-04-30 09:00:00	Authority1	t	signature001
ID002	0x4d5e6f	Tran Thi B	0987654321	1985-02-02	Vietnam	2025-04-29 09:30:00	Authority2	t	signature002
ID003	0x7g8h9i	Le Van C	0912345678	1995-03-03	Vietnam	2025-04-28 10:00:00	Authority3	f	signature003
\.


--
-- TOC entry 4860 (class 0 OID 16685)
-- Dependencies: 220
-- Data for Name: withdrawals; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.withdrawals (id, recipient, nullifier_hash, relayer, fee, "timestamp", tx_hash) FROM stdin;
1	0x1a2b3c	0xnull123	0xrelay001	0.01	2025-05-04 10:30:00	0xwith123tx
2	0x4d5e6f	0xnull456	0xrelay002	0.02	2025-05-05 11:30:00	0xwith456tx
3	0x7g8h9i	0xnull789	0xrelay003	0.03	2025-05-06 12:30:00	0xwith789tx
\.


--
-- TOC entry 4869 (class 0 OID 0)
-- Dependencies: 217
-- Name: deposits_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.deposits_id_seq', 3, true);


--
-- TOC entry 4870 (class 0 OID 0)
-- Dependencies: 219
-- Name: withdrawals_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.withdrawals_id_seq', 3, true);


--
-- TOC entry 4707 (class 2606 OID 16683)
-- Name: deposits deposits_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.deposits
    ADD CONSTRAINT deposits_pkey PRIMARY KEY (id);


--
-- TOC entry 4711 (class 2606 OID 16699)
-- Name: kyc kyc_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.kyc
    ADD CONSTRAINT kyc_pkey PRIMARY KEY (citizen_id);


--
-- TOC entry 4709 (class 2606 OID 16692)
-- Name: withdrawals withdrawals_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.withdrawals
    ADD CONSTRAINT withdrawals_pkey PRIMARY KEY (id);


-- Completed on 2025-05-06 14:20:12

--
-- PostgreSQL database dump complete
--

