# JustQuiz
### By Kalykul Margulan
A quiz game for many players, on diferent topics, and difficulties

## Running the app
Run this command to run the app
```
$ env POSTGRES_PASSWORD="postgres" APP_DSN="postgres://postgres:postgres@db:5432/postgres?sslmode=disable" docker-compose --env-file .env.example up --build
```

## Endpoints
* For players
```POST /v1/players``` - Create new player. Requires only `name`

```GET /v1/players/{id}``` - Get player by `{id}`

```PUT /v1/players/{id}``` - Update player name and score

```DELETE /v1/players/{id}``` - Delete player by `{id}`. Requires `menus:write` permission.

```GET /v1/healthcheck``` - For healthcheck

```GET /v1/players``` - Get a list of all players

* For users

	```POST /v1/users``` - Register new user

	```PUT /v1/users/activated``` - Activate user

	```POST /v1/users/login``` - Login user


## DB Structure
```
Table players {
  id bigserial [primary key]
  name text
  joined timestamp
  last_update timestamp
  score integer
}

Table quizes {
  id integer [primary key]
  created_at timestamp
	updated_at timestamp
	category  text
	dificulty integer
	q1 text
	a1 text
	q2 text
	a2 text
}

Table games {
  id bigserial [primary key]
  time timestamp
  place text
  quiz bigserial
  player bigserial
}

Ref: games.quiz < quizes.id

Ref: games.player < players.id
```
