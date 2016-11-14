package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:123123@tcp(127.0.0.1:3306)/gotest")
	if err != nil {
		fmt.Println(err.Error())
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println(err.Error())
	}

	type Todo struct {
		Id        int
		Name      string
		Completed bool
	}

	r := gin.Default()

	//ROUTES

	//GET SINGLE TODO
	r.GET("/todo/:id", func(c *gin.Context) {
		var todo Todo
		var result gin.H

		id := c.Param("id")

		row := db.QueryRow("select id, name, completed from todo where id = ?;", id)

		err = row.Scan(&todo.Id, &todo.Name)
		if err != nil {
			result = gin.H{
				"result": nil,
				"count":  0,
			}
		} else {
			result = gin.H{
				"result": todo,
				"count":  1,
			}
		}

		c.JSON(http.StatusOK, result)
	})

	//GET ALL TODOS
	r.GET("/todos", func(c *gin.Context) {
		var todo Todo
		var todos []Todo

		rows, err := db.Query("select id, name, completed from todo;")
		if err != nil {
			fmt.Println(err.Error())
		}

		for rows.Next() {
			err = rows.Scan(&todo.Id, &todo.Name, &todo.Completed)
			todos = append(todos, todo)
			if err != nil {
				fmt.Println(err.Error())
			}
		}

		defer rows.Close()

		c.JSON(http.StatusOK, gin.H{
			"result": todos,
			"count":  len(todos),
		})
	})

	//POST NEW TODO
	r.POST("/todos", func(c *gin.Context) {
		var buffer bytes.Buffer

		name := c.PostForm("name")

		stmt, err := db.Prepare("insert into todo (name) values(?);")

		if err != nil {
			fmt.Print(err.Error())
		}

		_, err = stmt.Exec(name)
		if err != nil {
			fmt.Print(err.Error())
		}

		buffer.WriteString(name)

		defer stmt.Close()

		b := buffer.String()

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf(" %s successfully created", b),
		})
	})

	r.PATCH("/todos/:id", func(c *gin.Context) {
		var nameBuffer bytes.Buffer
		var completedBuffer bytes.Buffer
		var idBuffer bytes.Buffer

		id := c.Param("id")
		name := c.PostForm("name")
		completed := c.PostForm("completed")

		stmt, err := db.Prepare("update todo set name=?, completed=? where id= ?;")
		if err != nil {
			fmt.Println(err.Error())
		}

		_, err = stmt.Exec(name, completed, id)
		if err != nil {
			fmt.Println(err.Error())
		}

		nameBuffer.WriteString(name)
		completedBuffer.WriteString(completed)
		idBuffer.WriteString(id)
		nb := nameBuffer.String()
		cb := completedBuffer.String()
		ib := idBuffer.String()

		defer stmt.Close()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully updated ID %[1]s with a name = %[2]s and completed = %[3]s", ib, nb, cb),
		})
	})

	r.DELETE("/todos/:id", func(c *gin.Context) {
		var buffer bytes.Buffer

		id := c.Param("id")

		stmt, err := db.Prepare("delete from todo where id= ?;")
		if err != nil {
			fmt.Println(err.Error())
		}

		_, err = stmt.Exec(id)
		if err != nil {
			fmt.Println(err.Error())
		}

		buffer.WriteString(id)

		defer stmt.Close()

		b := buffer.String()

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("deleted todo with id=%s", b),
		})
	})

	fmt.Println("serving goodies at localhost:3000")
	r.Run(":3000")
}
