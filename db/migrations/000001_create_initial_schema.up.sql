CREATE TABLE accounts (
    id UUID PRIMARY KEY,
    owner_name TEXT NOT NULL,
    currency CHAR(3) NOT NULL DEFAULT 'GBP',
    balance_pence BIGINT NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT accounts_currency_gbp CHECK (currency = 'GBP'),
    CONSTRAINT accounts_balance_non_negative CHECK (balance_pence >= 0),
    CONSTRAINT accounts_status_valid CHECK (status IN ('active', 'closed'))
);

CREATE TABLE ledger_transactions (
    id UUID PRIMARY KEY,
    kind TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT ledger_transactions_kind_valid CHECK (kind IN ('opening_balance', 'transfer'))
);

CREATE TABLE ledger_entries (
    id UUID PRIMARY KEY,
    ledger_transaction_id UUID NOT NULL REFERENCES ledger_transactions(id),
    account_id UUID NOT NULL REFERENCES accounts(id),
    amount_pence BIGINT NOT NULL,
    balance_after_pence BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX ledger_entries_account_id_created_at_idx
    ON ledger_entries(account_id, created_at DESC);

CREATE INDEX ledger_entries_ledger_transaction_id_idx
    ON ledger_entries(ledger_transaction_id);