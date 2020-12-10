package db

import (
	"fmt"
	log2 "github.com/labstack/gommon/log"
	"log"
	"github.com/labstack/echo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)


type Todo struct {
	gorm.Model `json:"model"`
	Msg       string `json:"msg"`
	Done      bool `json:"done"`
}

func initialMigration (db *gorm.DB){
	db.AutoMigrate(&Todo{})
}

func ConnectToDB() *gorm.DB{
	dsn := "root:carpdaniela@tcp(localhost:3306)/tema?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	initialMigration(db)

	return db
}

func AllTodos (db *gorm.DB) func(echo.Context) error{
	return func(c echo.Context) error {

		log.Println("Get all Todos Request --> ",c.QueryParams())

		todos := getTodos(db)

		return c.JSON(http.StatusOK, todos)
	}
}

func NewTodo(db *gorm.DB) func(echo.Context) error {
	return func(c echo.Context) error {
		log.Println("Add Request --> ",c.QueryParams())
		msg := c.QueryParam("msg") //c.Param("msg")
		done,err := strconv.ParseBool(c.QueryParam("done"))

		if err!= nil {
			log2.Error(err)
		}else{

			var nrOfRows int64
			db.Where(
				"msg = ?", msg,
			).Find(&Todo{}).Count(&nrOfRows);

			msg = strings.Replace(msg, "%20", " ", -1)

			if nrOfRows == 0 {
				db.Create(&Todo{Msg: msg, Done: done})

				todos := getTodos(db)

				return c.JSON(http.StatusOK, todos)
			}
		}

		return nil
	}
}

func DeleteTodo(db *gorm.DB) func(echo.Context) error {
	return func(c echo.Context) error {
		log.Println("Delete Request --> ",c.QueryParams())
		msg := c.QueryParam("msg")

		msg = strings.Replace(msg, "%20", " ", -1)

		var todo Todo
		db.Where("msg=?", msg).Find(&todo)
		db.Delete(&todo)

		todos := getTodos(db)

		return c.JSON(http.StatusOK, todos)
	}
}

func UpdateTodo(db *gorm.DB) func(echo.Context) error {
	return func(c echo.Context) error {

		log.Println("Update Request --> ",c.QueryParams())
		msg := c.QueryParam("msg")
		done, err := strconv.ParseBool(c.QueryParam("done"))

		msg = strings.Replace(msg, "%20", " ", -1)

		if err != nil{
			fmt.Println(err)
		}else{
			var todo Todo
			db.Where("msg=?", msg).Find(&todo)
			todo.Done = done
			db.Save(&todo)

			todos := getTodos(db)

			return c.JSON(http.StatusOK, todos)
		}
		return nil
	}
}

func getTodos(db *gorm.DB) []Todo{
	var todos []Todo
	db.Find(&todos)
	return todos
}


