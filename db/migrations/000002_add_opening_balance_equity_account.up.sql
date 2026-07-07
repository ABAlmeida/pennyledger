ALTER TABLE accounts
ADD COLUMN account_type TEXT NOT NULL DEFAULT 'customer';

ALTER TABLE accounts
ADD CONSTRAINT accounts_account_type_valid
CHECK (account_type IN ('customer', 'internal'));

ALTER TABLE accounts
DROP CONSTRAINT accounts_balance_non_negative;

ALTER TABLE accounts
ADD CONSTRAINT accounts_customer_balance_non_negative
CHECK (account_type = 'internal' OR balance_pence >= 0);

INSERT INTO accounts (id, owner_name, account_type, balance_pence, status)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'Opening Balance Equity',
    'internal',
    0,
    'active'
);