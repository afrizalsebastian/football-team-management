DROP INDEX IF EXISTS idx_matches_away_team;
DROP INDEX IF EXISTS idx_matches_home_team;
DROP INDEX IF EXISTS idx_matches_away_team;
DROP CONSTRAINT IF EXISTS chk_matches_different_teams;
DROP TABLE IF EXISTS matches;