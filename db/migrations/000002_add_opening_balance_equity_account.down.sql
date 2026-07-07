DELETE FROM accounts
WHERE id = '00000000-0000-0000-0000-000000000001';

ALTER TABLE accounts
DROP CONSTRAINT accounts_customer_balance_non_negative;

ALTER TABLE accounts
ADD CONSTRAINT accounts_balance_non_negative
CHECK (balance_pence >= 0);

ALTER TABLE accounts
DROP CONSTRAINT accounts_account_type_valid;

ALTER TABLE accounts
DROP COLUMN account_type;