package main

import (
	"GO_web/repository"
	"GO_web/router"
	"GO_web/service"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := "todo:Fz0826327826@tcp(127.0.0.1:3306)/todo_app?parseTime=true&charset=utf8mb4"

	db, err := repository.OpenMySQL(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("MySQL 连接成功")

	repo := repository.NewMySQLTodoRepository(db)
	todoService := service.NewTodoService(repo)

	sessions := service.NewSessionStore()

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, sessions)
	userService := service.NewUserService(userRepo)

	mux := router.NewRouter(
		todoService,
		authService,
		userService,
		sessions,
	)

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println(err)
	}
}
// waht can i say
