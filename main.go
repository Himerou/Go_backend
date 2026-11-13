package main

import (
	"GO_web/repository"
	"GO_web/router"
	"GO_web/service"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("未加载 .env 文件，尝试使用系统环境变量")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("MYSQL_USER")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DATABASE")

	if dbHost == "" || dbPort == "" || dbUser == "" ||
		dbPassword == "" || dbName == "" {
		log.Fatalln("数据库配置不完整")
	}

	cfg := mysql.NewConfig()
	cfg.User = dbUser
	cfg.Passwd = dbPassword
	cfg.Net = "tcp"
	cfg.Addr = dbHost + ":" + dbPort
	cfg.DBName = dbName
	cfg.ParseTime = true
	cfg.Loc = time.Local

	dsn := cfg.FormatDSN()

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
