CREATE TABLE public.users (
    id serial NOT NULL,
    first_name character varying(255),
    last_name character varying(255),
    email character varying(255),
    password text,
    UNIQUE(email)
);