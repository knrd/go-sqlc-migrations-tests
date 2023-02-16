CREATE TABLE balance_logs (
  id BIGSERIAL PRIMARY KEY,
  balance_id INTEGER NOT NULL REFERENCES balances (id),
  balance_before_change INTEGER NOT NULL,
  change INTEGER NOT NULL,
  note VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL
);
