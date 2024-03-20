CREATE TABLE public.categories
(
    id          bigserial    NOT NULL,
    "name"      varchar(255) NOT NULL,
    description varchar(255) NULL,
    status      int2         NOT NULL DEFAULT '1':: smallint,
    is_default  bool         NULL DEFAULT false,
    created_at  timestamp(0) NULL,
    updated_at  timestamp(0) NULL,
    CONSTRAINT categories_pkey PRIMARY KEY (id)
);