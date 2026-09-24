CREATE TABLE IF NOT EXISTS goals (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_id        UUID NOT NULL REFERENCES matches(id),
    player_id       UUID NOT NULL REFERENCES players(id),
    goal_minute     VARCHAR(25),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    is_deleted      BOOLEAN DEFAULT false,
    deleted_at      TIMESTAMPTZ DEFAULT NULL
);

CREATE INDEX idx_goals_match_id ON goals (match_id);
CREATE INDEX idx_goals_player_id ON goals (player_id);