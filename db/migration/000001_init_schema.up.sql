CREATE TABLE public.transfers (
	id bigserial NOT NULL,
	from_account_id int8 NOT NULL,
	to_account_id int8 NOT NULL,
	ammount int8 NOT NULL,
	created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT pk_transfers PRIMARY KEY (id)
);

CREATE TABLE public.entries (
	id bigserial NOT NULL,
	account_id int8 NOT NULL,
	amount int8 NOT NULL,
	created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT pk_entries PRIMARY KEY (id)
);

CREATE TABLE public.accounts (
	id bigserial NOT NULL,
	"owner" varchar NOT NULL,
	balance int8 NULL,
	currency varchar NOT NULL,
	created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT pk_accounts PRIMARY KEY (id)
);


CREATE INDEX idx_transfers ON public.transfers USING btree (from_account_id, to_account_id);
CREATE INDEX idx_entries ON public.entries USING btree (account_id);
CREATE INDEX idx_accounts_owner ON public.accounts USING btree (owner);

ALTER TABLE public.entries ADD CONSTRAINT fk_entries_accounts FOREIGN KEY (account_id) REFERENCES public.accounts(id);
ALTER TABLE public.transfers ADD CONSTRAINT fk2_transfers_accounts FOREIGN KEY (to_account_id) REFERENCES public.accounts(id);
ALTER TABLE public.transfers ADD CONSTRAINT fk_transfers_accounts FOREIGN KEY (from_account_id) REFERENCES public.accounts(id);