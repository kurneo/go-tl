CREATE TABLE public.admin_access_tokens
(
    id         bigserial NOT NULL,
    token      text,
    created_at timestamp(0) NULL,
    expired_at timestamp(0) NULL,
    admin_id    int8      NOT NULL,
    CONSTRAINT admin_access_tokens_pkey PRIMARY KEY (id)
);

ALTER TABLE public.admin_access_tokens
    ADD CONSTRAINT admin_access_tokens_admin_id_foreign FOREIGN KEY (admin_id) REFERENCES public.admins (id);

CREATE INDEX admin_access_tokens_token_index ON admin_access_tokens (token);