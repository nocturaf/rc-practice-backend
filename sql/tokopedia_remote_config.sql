CREATE TABLE users (
    id serial NOT NULL,
    first_name character varying(255) NOT NULL,
    last_name character varying(255),
    email character varying(255) NOT NULL,
    password text NOT NULL,
    role int NOT NULL DEFAULT 0,

    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT email_unique UNIQUE (email)
);