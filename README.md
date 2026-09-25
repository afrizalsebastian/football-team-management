<h1>Football Team Management</h1>
<p>
  <strong>App that can manage teams, team player, team matches, and match goals</strong>
</p>

## Prerequisite
- PostgreSQL
- Makefile
- Go
- migrate-cli 
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```
- sqlc-cli
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```
- swaggo-cli
```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

## Project Structure

```
├── api
│   └── web_response.go
├── application
│   ├── controllers
│   │   ├── admin_controlelr.go
│   │   ├── goals_controller.go
│   │   ├── hello_controller.go
│   │   ├── match_controller.go
│   │   ├── players_controller.go
│   │   ├── superadmin_controller.go
│   │   └── teams_controller.go
│   ├── dto
│   │   ├── admin_account.go
│   │   ├── matches.go
│   │   ├── player.go
│   │   └── team.go
│   ├── helper
│   │   ├── jwt_manager.go
│   │   ├── matches_helper.go
│   │   ├── pointer_util.go
│   │   └── validator.go
│   └── services
│       ├── admin_service.go
│       ├── goals_service.go
│       ├── hello_services.go
│       ├── matches_service.go
│       ├── players_service.go
│       ├── superadmin_service.go
│       └── teams_service.go
├── bootstrap
│   └── football_management_app.go
├── config
│   └── config.go
├── constants
│   ├── errors.go
│   ├── match_status.go
│   ├── player_position.go
│   ├── response_code.go
│   └── validation_code.go
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── domain
│   ├── dao
│   │   ├── admin_account.go
│   │   ├── goals.go
│   │   ├── matches.go
│   │   ├── players.go
│   │   └── teams.go
│   └── repository
│       ├── queries
│       │   ├── db
│       │   │   ├── db.go
│       │   │   ├── models.go
│       │   │   ├── querier.go
│       │   │   └── queries.sql.go
│       │   └── queries.sql
│       ├── admin_account.go
│       ├── goals_repository.go
│       ├── matches_repository.go
│       ├── players_repository.go
│       ├── teams_repository.go
│       └── util.go
├── middleware
│   ├── auth_middleware.go
│   ├── default_middleware.go
│   └── superadmin_middleware.go
├── migration
│   ├── 000001_teams_and_player.down.sql
│   ├── 000001_teams_and_player.up.sql
│   ├── 000002_match_table.down.sql
│   ├── 000002_match_table.up.sql
│   ├── 000003_player_goals.down.sql
│   ├── 000003_player_goals.up.sql
│   ├── 000004_match_result.down.sql
│   ├── 000004_match_result.up.sql
│   ├── 000005_admin_account.down.sql
│   ├── 000005_admin_account.up.sql
│   ├── 000006_alter_column_admin_identity.down.sql
│   └── 000006_alter_column_admin_identity.up.sql
├── module
│   ├── database
│   │   └── db_connections.go
│   └── logger
│       └── logger.go
├── routes
│   ├── admin_routes.go
│   ├── goal_routes.go
│   ├── hello_routes.go
│   ├── matches_routes.go
│   ├── players_routes.go
│   ├── superadmin_routes.go
│   └── teams_routes.go
├── server
│   ├── di.go
│   └── server.go
├── .env.example
├── .gitignore
├── Makefile
├── README.md
├── go.mod
├── go.sum
├── main.go
└── sqlc.yaml
```

## How To Run
-- Make sure the postgre is running <br>
-- Copy your own .env file from .env.example and fill the env <br>
-- Run the migration folder <br>
```bash
make migrate-up
```
-- Running locally
```bash
make run
```
-- Running in build file
```bash
make build-run
```

## Test the app
Open this in web browser to test the app from Swagger UI
```
<YOUR_BASE_URL_APP>/swagger/index.html#/
```
Example
```bash
http://localhost:8080/swagger/index.html#/
```