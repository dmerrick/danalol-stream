CREATE TABLE users (
  id             SERIAL PRIMARY KEY,
  username       VARCHAR(64) NOT NULL,
  miles          REAL,
  num_visits     INTEGER,
  has_donated    BOOLEAN,
  first_seen     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  last_seen      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  date_created   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
