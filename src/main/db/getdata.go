package db

import (
	"fmt"
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

	//sqlDB, errDB := db.DB()
	//if errDB != nil{
	//	fmt.Println(errDB)
	////}//else{
	////	defer sqlDB.Close()
	////}

	if err != nil {
		panic("failed to connect database")
	}

	initialMigration(db)

	return db
}

func Test () func(echo.Context) error{
	return func(c echo.Context) error {

		return c.JSON(http.StatusOK, "succes !!!")
	}
}

func AllTodos (db *gorm.DB) func(echo.Context) error{
	return func(c echo.Context) error {
		var todos []Todo
		db.Find(&todos)
		fmt.Println("select todos -->  ", todos)

		return c.JSON(http.StatusOK, todos)
	}
}

func NewTodo(db *gorm.DB) func(echo.Context) error {
	return func(c echo.Context) error {
		msg := c.QueryParam("msg") //c.Param("msg")
		fmt.Println("done param = ",c.QueryParam("done"))
		fmt.Println("msg param= ", c.QueryParam("msg"))
		done,err := strconv.ParseBool(c.QueryParam("done"))

		if err!= nil {
			fmt.Println(err)
		}else{

			var nrOfRows int64
			statement := db.Where(
				"msg = ?", msg,
			).Find(&Todo{}).Count(&nrOfRows);

			fmt.Println("Stmt: ",statement)
			fmt.Println("NrOfRows: ",nrOfRows)

			msg = strings.Replace(msg, "%20", " ", -1)

			if nrOfRows == 0 {
				db.Create(&Todo{Msg: msg, Done: done})

				var todos []Todo
				db.Find(&todos)
				fmt.Println("select todos -->  ", todos)

				return c.JSON(http.StatusOK, todos)
			}
		}

		return nil
	}
}

func DeleteTodo(db *gorm.DB) func(echo.Context) error {
	return func(c echo.Context) error {
		msg := c.QueryParam("msg")

		msg = strings.Replace(msg, "%20", " ", -1)

		var todo Todo
		db.Where("msg=?", msg).Find(&todo)
		db.Delete(&todo)

		var todos []Todo
		db.Find(&todos)
		fmt.Println("select todos -->  ", todos)

		return c.JSON(http.StatusOK, todos)
	}
}

func UpdateTodo(db *gorm.DB) func(echo.Context) error {
	return func(c echo.Context) error {
		msg := c.QueryParam("msg")
		done, err := strconv.ParseBool(c.Param("done"))

		msg = strings.Replace(msg, "%20", " ", -1)

		if err != nil{
			fmt.Println(err)
		}else{
			var todo Todo
			db.Where("msg=?", msg).Find(&todo)
			todo.Done = done
			db.Save(&todo)

			var todos []Todo
			db.Find(&todos)
			fmt.Println("select todos -->  ", todos)

			return c.JSON(http.StatusOK, todos)
		}
		return nil
	}
}


