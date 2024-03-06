CREATE TABLE public.users
(
    id            bigserial    NOT NULL,
    name          varchar(255) NOT NULL,
    email         varchar(255) NULL,
    password      varchar(255) NOT NULL,
    last_login_at timestamp(0) NULL,
    created_at    timestamp(0) NULL,
    updated_at    timestamp(0) NULL,
	CONSTRAINT users_pkey PRIMARY KEY (id),
	CONSTRAINT users_email_unique UNIQUE (email)
);

INSERT INTO "users" ("name", "email", "password", "last_login_at", "created_at", "updated_at")
VALUES ('Admin', 'admin@example.com', '$2a$14$H4g2bAIPI7SYNJHrgbZhTu9IoD9/SwMFbFC3aqI3LtEfZiYu5b4xS', '2021-11-10 18:02:53.769', '2021-11-10 18:02:53.769', '2021-11-10 18:02:53.769')