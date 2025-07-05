# auth_service


## Project Structure

```plaintext
myapp/
├── cmd/
│   └── app/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── db/
│   │   └── db.go
│   ├── repository/
│   │   ├── user_repository.go
│   │   └── user_repository_pg.go
│   ├── services/
│   │   ├── user_service.go
│   │   └── user_service_impl.go
│   ├── controllers/
│   │   └── user_controller.go
│   └── routes/
│       └── routes.go
├── migrations/               # SQL или миграции с помощью инструмента (golang-migrate)
├── go.mod
└── go.sum
```

---

// cmd/app/main.go
```go
package main

import (
    "log"
    "net/http"

    "myapp/internal/config"
    "myapp/internal/db"
    "myapp/internal/repository"
    "myapp/internal/routes"
    "myapp/internal/services"
)

func main() {
    // Загрузка конфигурации (DB_URL, PORT и т.д.)
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("config load: %v", err)
    }

    // Инициализация DB
    sqlDB, err := db.InitDB(cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("db init: %v", err)
    }
    defer sqlDB.Close()

    // Репозиторий и сервис
    userRepo := repository.NewUserRepoPg(sqlDB)
    userSvc := services.NewUserService(userRepo)

    // Настройка роутов
    router := routes.NewRouter(userSvc)
    log.Printf("Starting server on %s", cfg.Port)
    if err := http.ListenAndServe(cfg.Port, router); err != nil {
        log.Fatal(err)
    }
}
```

---

// internal/config/config.go
```go
package config

import (
    "os"
    "fmt"
)

// Config хранит настройки приложения
type Config struct {
    DatabaseURL string
    Port        string
}

// Load загружает настройки из окружения
func Load() (*Config, error) {
    dbURL := os.Getenv("DATABASE_URL")
    port := os.Getenv("PORT")
    if dbURL == "" || port == "" {
        return nil, fmt.Errorf("DATABASE_URL and PORT must be set")
    }
    return &Config{DatabaseURL: dbURL, Port: port}, nil
}
```

---

// internal/db/db.go
```go
package db

import (
    "database/sql"
    _ "github.com/lib/pq"
)

// InitDB открывает и проверяет соединение
func InitDB(dsn string) (*sql.DB, error) {
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, err
    }
    if err := db.Ping(); err != nil {
        db.Close()
        return nil, err
    }
    return db, nil
}
```

---

// internal/repository/user_repository.go
```go
package repository

import (
    "context"
    "myapp/internal/models"
)

// UserRepository описывает операции хранения пользователей
type UserRepository interface {
    Create(ctx context.Context, u models.User) (int, error)
    GetByID(ctx context.Context, id int) (models.User, error)
    Update(ctx context.Context, u models.User) error
    Delete(ctx context.Context, id int) error
}
```

---

// internal/repository/user_repository_pg.go
```go
package repository

import (
    "context"
    "database/sql"
    "myapp/internal/models"
)

// userRepoPg — PostgreSQL-реализация UserRepository
type userRepoPg struct { db *sql.DB }

func NewUserRepoPg(db *sql.DB) UserRepository {
    return &userRepoPg{db: db}
}

func (r *userRepoPg) Create(ctx context.Context, u models.User) (int, error) {
    var id int
    query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`
    err := r.db.QueryRowContext(ctx, query, u.Name, u.Email).Scan(&id)
    return id, err
}

func (r *userRepoPg) GetByID(ctx context.Context, id int) (models.User, error) {
    var u models.User
    query := `SELECT id, name, email FROM users WHERE id=$1`
    err := r.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Name, &u.Email)
    return u, err
}

func (r *userRepoPg) Update(ctx context.Context, u models.User) error {
    query := `UPDATE users SET name=$1, email=$2 WHERE id=$3`
    res, err := r.db.ExecContext(ctx, query, u.Name, u.Email, u.ID)
    if err != nil {
        return err
    }
    if cnt, _ := res.RowsAffected(); cnt == 0 {
        return sql.ErrNoRows
    }
    return nil
}

func (r *userRepoPg) Delete(ctx context.Context, id int) error {
    query := `DELETE FROM users WHERE id=$1`
    res, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return err
    }
    if cnt, _ := res.RowsAffected(); cnt == 0 {
        return sql.ErrNoRows
    }
    return nil
}
```

---

// internal/services/user_service.go
```go
package services

import (
    "context"
    "myapp/internal/models"
    "myapp/internal/repository"
)

// UserService содержит бизнес-логику пользователей
type UserService interface {
    CreateUser(ctx context.Context, u models.User) (models.User, error)
    GetUser(ctx context.Context, id int) (models.User, error)
    UpdateUser(ctx context.Context, u models.User) (models.User, error)
    DeleteUser(ctx context.Context, id int) error
}
```

---

// internal/services/user_service_impl.go
```go
package services

import (
    "context"
    "database/sql"
    "myapp/internal/models"
    "myapp/internal/repository"
)

// userServiceImpl реализует UserService через репозиторий
type userServiceImpl struct {
    repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
    return &userServiceImpl{repo: repo}
}

func (s *userServiceImpl) CreateUser(ctx context.Context, u models.User) (models.User, error) {
    id, err := s.repo.Create(ctx, u)
    if err != nil {
        return models.User{}, err
    }
    u.ID = id
    return u, nil
}

func (s *userServiceImpl) GetUser(ctx context.Context, id int) (models.User, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *userServiceImpl) UpdateUser(ctx context.Context, u models.User) (models.User, error) {
    if err := s.repo.Update(ctx, u); err != nil {
        return models.User{}, err
    }
    return u, nil
}

func (s *userServiceImpl) DeleteUser(ctx context.Context, id int) error {
    return s.repo.Delete(ctx, id)
}
```

---

// internal/controllers/user_controller.go
```go
package controllers

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/go-chi/chi/v5"
    "myapp/internal/models"
    "myapp/internal/services"
)

// UserController оборачивает UserService
type UserController struct { svc services.UserService }

func NewUserController(svc services.UserService) *UserController {
    return &UserController{svc: svc}
}

func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
    var u models.User
    if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    res, err := c.svc.CreateUser(r.Context(), u)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(res)
}

func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        http.Error(w, "invalid id", http.StatusBadRequest)
        return
    }
    u, err := c.svc.GetUser(r.Context(), id)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "not found", http.StatusNotFound)
        } else {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }
    json.NewEncoder(w).Encode(u)
}

func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        http.Error(w, "invalid id", http.StatusBadRequest)
        return
    }
    var u models.User
    if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    u.ID = id
    updated, err := c.svc.UpdateUser(r.Context(), u)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(updated)
}

func (c *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        http.Error(w, "invalid id", http.StatusBadRequest)
        return
    }
    if err := c.svc.DeleteUser(r.Context(), id); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}
```

---

// internal/routes/routes.go
```go
package routes

import (
    "net/http"

    "github.com/go-chi/chi/v5"
    "myapp/internal/controllers"
    "myapp/internal/services"
)

// NewRouter настраивает маршруты для UserController
func NewRouter(userSvc services.UserService) http.Handler {
    r := chi.NewRouter()
    userCtrl := controllers.NewUserController(userSvc)

    r.Route("/users", func(r chi.Router) {
        r.Post("/", userCtrl.CreateUser)
        r.Get("/{id}", userCtrl.GetUser)
        r.Put("/{id}", userCtrl.UpdateUser)
        r.Delete("/{id}", userCtrl.DeleteUser)
    })
    return r
}
```

---

// internal/models/user.go
```go
package models

// User соответствует таблице users
// Пример миграции:
//  CREATE TABLE users (
//    id SERIAL PRIMARY KEY,
//    name TEXT NOT NULL,
//    email TEXT UNIQUE NOT NULL
//  );

type User struct {
    ID    int    `json:"id" db:"id"`
    Name  string `json:"name" db:"name"`
    Email string `json:"email" db:"email"`
}
```
