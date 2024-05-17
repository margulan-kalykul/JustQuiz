#!/bin/sh

# Run migrations
./migrate -path=./pkg/quiz/migrations -database=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable up

# Start the application
exec ./justquiz