CREATE TABLE IF NOT EXISTS stocks (
  id BIGINT PRIMARY KEY DEFAULT unique_rowid(),
  ticker STRING NOT NULL,
  target_from STRING,
  target_to STRING,
  company STRING,
  action STRING,
  brokerage STRING,
  rating_from STRING,
  rating_to STRING,
  time TIMESTAMPTZ,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_stocks_ticker ON stocks (ticker);
CREATE INDEX IF NOT EXISTS idx_stocks_time ON stocks (time);
CREATE UNIQUE INDEX IF NOT EXISTS ux_stocks_ticker_time ON stocks (ticker, time);

