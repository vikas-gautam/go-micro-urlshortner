--
-- Name: url; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.url_mapping (
	id serial NOT NULL,
    url character varying(255),
    generated_id character varying(255),
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);

ALTER TABLE public.url_mapping OWNER TO postgres;
