CREATE TABLE IF NOT EXISTS account (
    id VARCHAR(50) PRIMARY KEY,
    balance NUMERIC(13, 2),
    currency VARCHAR(3)
);

CREATE TABLE IF NOT EXISTS payment (
    id SERIAL,
    account VARCHAR(50) REFERENCES account (id) ON DELETE CASCADE,
    amount NUMERIC(13, 2),
    from_account VARCHAR(50),
    to_account VARCHAR(50),
    direction VARCHAR(8)
);
