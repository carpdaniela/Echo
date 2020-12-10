package main

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/todo", dana)
	e.POST("/todo", func(c echo.Context) error {
		u := User{
			Name:  "Dana POST",
			Email: "iulianadana97@gmail.com",
		}
		resp, _ := ToJSON(u)
		return c.String(http.StatusOK, resp)
	})
	e.Logger.Fatal(e.Start(":1323"))
}

func dana(c echo.Context) error {
	u := User{
		Name:  "Dana",
		Email: "iulianadana97@gmail.com",
	}
	resp, _ := ToJSON(u)
	return c.String(http.StatusOK, resp)
}

type User struct {
	Name  string `json:"name" xml:"name" form:"name" query:"name"`
	Email string `json:"email" xml:"email" form:"email" query:"email"`
}

func ToJSON(obj interface{}) (string, error) {
	res, err := json.Marshal(obj)
	if err != nil {
		res = []byte("")
	}
	return string(res), err
}
