CREATE TABLE IF NOT EXISTS matches (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  match_date      DATE NOT NULL,
  match_time      TIME NOT NULL,
  home_team_id    UUID NOT NULL REFERENCES teams(id),
  away_team_id    UUID NOT NULL REFERENCES teams(id),
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  is_deleted      BOOL DEFAULT false,
  deleted_at      TIMESTAMPTZ DEFAULT NULL,
  CONSTRAINT chk_matches_different_teams CHECK (home_team_id <> away_team_id)
);

CREATE INDEX idx_matches_date ON matches (match_date);
CREATE INDEX idx_matches_home_team ON matches (home_team_id);
CREATE INDEX idx_matches_away_team ON matches (away_team_id);