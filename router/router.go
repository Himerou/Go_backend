package router

import (
	"net/http"
	"path/filepath"

	"GO_web/handler"
	"GO_web/service"
)

func NewRouter(
	todoService *service.TodoService,
	authService *service.AuthService,
	userService *service.UserService,
	sessions service.SessionRepository,
) http.Handler {
	mux := http.NewServeMux()

	// #公开路由
	mux.HandleFunc(
		"POST /api/register",
		handler.RegisterTodo(authService),
	)
	mux.HandleFunc(
		"POST /api/login",
		handler.LoginTodo(authService),
	)
	mux.HandleFunc(
		"GET /api/users/{username}/profile",
		handler.GetPublicProfile(userService),
	)
	// 静态资源
	mux.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.Dir("./web")),
		),
	)
	// 页面路由
	mux.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("web", "login.html"))
	})
	mux.HandleFunc("GET /register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("web", "register.html"))
	})
	//

	// #需要登录的路由
	protectedMux := http.NewServeMux()
	// 页面路由
	protectedMux.HandleFunc("GET /profile", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("web", "profile.html"))
	})
	//

	protectedMux.HandleFunc(
		"GET /api/profile",
		handler.GetProfile(userService),
	)
	protectedMux.HandleFunc(
		"PUT /api/profile",
		handler.UpdateProfile(userService),
	)
	protectedMux.HandleFunc(
		"PUT /api/password",
		handler.UpdatePassword(userService),
	)
	protectedMux.HandleFunc(
		"PUT /api/username",
		handler.UpdateUsername(userService),
	)

	protectedMux.HandleFunc(
		"PUT /api/todos/{id}",
		handler.UpdateTodo(todoService),
	)
	protectedMux.HandleFunc(
		"GET /api/todos",
		handler.ListTodos(todoService),
	)
	protectedMux.HandleFunc(
		"POST /api/todos",
		handler.CreateTodo(todoService),
	)
	protectedMux.HandleFunc(
		"DELETE /api/todos/{id}",
		handler.DeleteTodo(todoService),
	)
	protectedMux.HandleFunc(
		"GET /api/todos/{id}",
		handler.GetTodo(todoService),
	)
	protectedMux.HandleFunc(
		"POST /api/logout",
		handler.LogoutTodo(authService),
	)
	protectedMux.HandleFunc(
		"DELETE /api/user",
		handler.DeleteUser(userService),
	)
	protectedHandler := handler.Middleware(sessions)(protectedMux)

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, filepath.Join("web", "index.html"))
			return
		}

		protectedHandler.ServeHTTP(w, r)
	}))
	return mux
}
