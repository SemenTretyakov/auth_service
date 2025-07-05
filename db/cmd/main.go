package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	db "github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	dbDSN string = "host=localhost port=54321 dbname=auth_service user=postgres password=postgres sslmode=disable"
)

func main() {
	ctx := context.Background()

	// Создаем пул соединений с базой данных
	pool, err := pgxpool.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Делаем запрос на вставку записи в таблицу users
	builderInsert := db.Insert("users").
		PlaceholderFormat(db.Dollar).
		Columns("fullname").
		Values(gofakeit.Name()).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var userID int
	err = pool.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		log.Fatalf("failed to insert users: %v", err)
	}

	log.Printf("inserted users with id: %d", userID)

	// Делаем запрос на выборку записей из таблицы users
	builderSelect := db.Select("id", "fullname", "createdAt", "updatedAt").
		From("users").
		PlaceholderFormat(db.Dollar).
		OrderBy("id ASC").
		Limit(10)

	query, args, err = builderSelect.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to select notes: %v", err)
	}

	var id int
	var fullname string
	var createdAt time.Time
	var updatedAt sql.NullTime

	for rows.Next() {
		err = rows.Scan(&id, &fullname, &createdAt, &updatedAt)
		if err != nil {
			log.Fatalf("failed to scan users: %v", err)
		}

		log.Printf("id: %d, fullname: %s, created_at: %v, updated_at: %v\n", id, fullname, createdAt, updatedAt)
	}

	// Делаем запрос на обновление записи в таблице users
	builderUpdate := db.Update("users").
		PlaceholderFormat(db.Dollar).
		Set("fullname", gofakeit.FirstName()).
		Set("updatedAt", time.Now().UTC()).
		Where(db.Eq{"id": userID})

	query, args, err = builderUpdate.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	res, err := pool.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to update users: %v", err)
	}

	log.Printf("updated %d rows", res.RowsAffected())

	// Делаем запрос на получение измененной записи из таблицы users
	builderSelectOne := db.Select("id", "fullname", "createdAt", "updatedAt").
		From("users").
		PlaceholderFormat(db.Dollar).
		Where(db.Eq{"id": userID}).
		Limit(1)

	query, args, err = builderSelectOne.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	err = pool.QueryRow(ctx, query, args...).Scan(&id, &fullname, &createdAt, &updatedAt)
	if err != nil {
		log.Fatalf("failed to select notes: %v", err)
	}

	log.Printf("id: %d, fullname: %s, created_at: %v, updated_at: %v\n", id, fullname, createdAt, updatedAt)
}
