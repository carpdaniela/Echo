package main

import (
	myTodos "./db"
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gorm.io/gorm"
)

func main() {

	db := myTodos.ConnectToDB();
	handleRequest(db)

}

func handleRequest(db *gorm.DB) {
	e := echo.New()

	e.Use(CORSMiddlewareWrapper)

	e.GET("/todos", myTodos.AllTodos(db))
	e.POST("/todos/add/", myTodos.NewTodo(db))
	e.DELETE("/todos/remove", myTodos.DeleteTodo(db))
	e.PUT("/user/todos/update", myTodos.UpdateTodo(db))

	e.Logger.Fatal(e.Start(":3000"))
}

func CORSMiddlewareWrapper(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		req := ctx.Request()
		dynamicCORSConfig := middleware.CORSConfig{
			AllowOrigins: []string{req.Header.Get("Origin")},
			AllowHeaders: []string{"Accept", "Cache-Control", "Content-Type", "X-Requested-With"},
		}
		CORSMiddleware := middleware.CORSWithConfig(dynamicCORSConfig)
		CORSHandler := CORSMiddleware(next)
		return CORSHandler(ctx)
	}
}

func ToJSON(obj interface{}) (string, error) {
	res, err := json.Marshal(obj)
	if err != nil {
		res = []byte("")
	}
	return string(res), err
}
