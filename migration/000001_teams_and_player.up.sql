CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Player position enum
CREATE type player_position AS ENUM (
  'S', -- striker
  'M', -- mid
  'B', -- back
  'K' -- keeper
);

-- TEAMS
CREATE TABLE IF NOT EXISTS teams (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name            VARCHAR(100) NOT NULL,
  logo            TEXT DEFAULT NULL,
  founded_year    VARCHAR(4),
  address         TEXT NOT NULL,
  city            VARCHAR(100) NOT NULL,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  is_deleted      BOOL DEFAULT false,
  deleted_at      TIMESTAMPTZ DEFAULT NULL
);


-- PLAYERS
CREATE TABLE IF NOT EXISTS players (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id         UUID NOT NULL REFERENCES teams(id),
    name            VARCHAR(100) NOT NULL,
    height_cm       NUMERIC(5,2) NOT NULL CHECK (height_cm > 0),
    weight_kg       NUMERIC(5,2) NOT NULL CHECK (weight_kg > 0),
    position        player_position NOT NULL,
    jersey_number   SMALLINT NOT NULL CHECK (jersey_number BETWEEN 1 AND 99),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    is_deleted      BOOL DEFAULT false,
    deleted_at      TIMESTAMPTZ DEFAULT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_players_team_jersey
ON players (team_id, jersey_number)
WHERE is_deleted = false;
