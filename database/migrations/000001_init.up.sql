CREATE TABLE balances (
  id SERIAL PRIMARY KEY,
  amount INTEGER NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL CONSTRAINT proper_email CHECK (email ~* '^[A-Za-z0-9][A-Za-z0-9._+%-]*@[A-Za-z0-9.-]+[.][A-Za-z]+$'),
  created_at TIMESTAMP NOT NULL
);

CREATE TABLE balance_logs (
  id BIGSERIAL PRIMARY KEY,
  balance_id INTEGER NOT NULL REFERENCES balances (id),
  change INTEGER NOT NULL,
  note VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL
);
