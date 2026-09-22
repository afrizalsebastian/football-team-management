-- name: CreateTeam :one
INSERT INTO teams (
  name, logo, founded_year, address, city
) VALUES (
  @name, @logo, @founded_year, @address, @city
) RETURNING id;

-- name: GetListTeams :many
SELECT id, name, logo, founded_year
FROM teams
ORDER BY founded_year DESC;

-- name: SoftDeleteTeams :exec
UPDATE teams
SET is_deleted = true,
  deleted_at = now()
WHERE id = @id;