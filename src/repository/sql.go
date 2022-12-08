package repository

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/dietzy1/discordbot/src/bot/emotes"
)

type Sqlrepository interface {
	IncrementEmote(ctx context.Context, emote *emotes.Emote) error
	GetUserEmotes(ctx context.Context, emote *emotes.Emote) ([]emotes.Emote, error)
	GetServerEmote(ctx context.Context, emote *emotes.Emote) ([]emotes.Emote, error)
}

type Sql struct {
	sql *sql.DB
}

func NewSql() (*Sql, error) {
	db, err := sql.Open("postgres", os.Getenv("DB"))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	a := &Sql{sql: db}

	return a, nil
}

func (db *Sql) IncrementEmote(ctx context.Context, emote *emotes.Emote) error {
	_, err := db.sql.Exec("INSERT INTO emotes (emote, user, guild, count) VALUES ($1, $2, $3, $4) ON CONFLICT (emote, user, guild) DO UPDATE SET count = emotes.count + 1", emote.Emote, emote.User, emote.Guild, 1)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (db *Sql) GetUserEmotes(ctx context.Context, emote *emotes.Emote) ([]emotes.Emote, error) {
	rows, err := db.sql.Query("SELECT emote, user, guild, count FROM emotes WHERE user = $1 AND guild = $2 ORDER BY count DESC", emote.User, emote.Guild)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	result := []emotes.Emote{}
	for rows.Next() {
		emote := emotes.Emote{}
		err := rows.Scan(&emote.Emote, &emote.User, &emote.Guild, &emote.Count)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		result = append(result, emote)
	}
	return result, nil
}

func (db *Sql) GetServerEmote(ctx context.Context, emote *emotes.Emote) ([]emotes.Emote, error) {
	rows, err := db.sql.Query("SELECT emote, user, guild, count FROM emotes WHERE emote = $1 AND guild = $2 ORDER BY count DESC", emote.Emote, emote.Guild)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	result := []emotes.Emote{}
	for rows.Next() {
		emote := emotes.Emote{}
		err := rows.Scan(&emote.Emote, &emote.User, &emote.Guild, &emote.Count)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		result = append(result, emote)
	}
	return result, nil
}

// Path: src/repository/mongo.go
