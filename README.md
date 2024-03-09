# JustQuiz
A quiz game for many players, on diferent topics, and difficulties

## Endpoints
```POST /v1/players``` - Create new user. Requires only `name`

```GET /v1/players/{id}``` - Get player by `{id}`

```PUT /v1 players/{id}``` - Update player name and score

```DELETE /v1/players/{id}``` - Delete player by `{id}`

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
