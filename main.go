package main

import (
	"GO_web/repository"
	"GO_web/router"
	"GO_web/service"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("加载 .env 文件失败:", err)
	}
	dsn := os.Getenv("DSN")
	if dsn == "" {
		log.Fatalln("DSN 变量未设置")
	}

	db, err := repository.OpenMySQL(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("MySQL 连接成功")

	todoRepo := repository.NewMySQLTodoRepository(db)
	todoService := service.NewTodoService(todoRepo)
	sessions := repository.NewSessionStore()
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

// what can I say
