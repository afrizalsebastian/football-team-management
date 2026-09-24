CREATE TABLE IF NOT EXISTS match_results (
    id              UUID PRIMARY KEY REFERENCES matches(id) ,
    home_score      SMALLINT NOT NULL CHECK (home_score >= 0),
    away_score      SMALLINT NOT NULL CHECK (away_score >= 0),
    status          SMALLINT NOT NULL CHECK (status IN (-1, 0, 1)),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
