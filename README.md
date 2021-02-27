# goschedule

## Database

### Create migration

`go run github.com/golang-migrate/migrate/v4/cmd/migrate create -ext sql -dir db/migrations -seq <name>`

### Run migration

`go run ./cmd/migrate/main.go -up`
