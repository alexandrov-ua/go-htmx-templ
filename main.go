package main

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"htmx-templ/templates"

	"github.com/gin-gonic/gin"
)

var todoState = templates.TodoContext{Todos: []templates.TodoModel{templates.TodoModel{Title: "qwe", Description: "desc"}}}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		templates.Index("/hello?name=london").Render(c.Request.Context(), c.Writer)
	})

	r.GET("/hello", func(c *gin.Context) {

		if _, ok := c.Request.Header["Hx-Request"]; !ok {
			templates.Index("/hello?name="+c.Query("name")).Render(c.Request.Context(), c.Writer)
		} else {
			templates.Hello(c.Query("name")).Render(c.Request.Context(), c.Writer)
		}
	})

	number := 0
	r.GET("/counter", func(c *gin.Context) {

		if _, ok := c.Request.Header["Hx-Request"]; !ok {
			templates.Index("/counter").Render(c.Request.Context(), c.Writer)
		} else {
			templates.Counter(number).Render(c.Request.Context(), c.Writer)
		}
	})

	r.GET("/counter/click", func(c *gin.Context) {
		number += 1
		templates.Counter(number).Render(c.Request.Context(), c.Writer)
	})
	r.GET("/todo", func(c *gin.Context) {
		if _, ok := c.Request.Header["Hx-Request"]; !ok {
			templates.Index("/todo").Render(c.Request.Context(), c.Writer)
		} else {
			templates.Todo(todoState).Render(c.Request.Context(), c.Writer)
		}
	})
	r.POST("/todo", func(c *gin.Context) {
		todo := templates.TodoModel{}
		c.Bind(&todo)

		if slices.ContainsFunc(todoState.Todos, func(m templates.TodoModel) bool {
			return m.Title == todo.Title
		}) {
			c.Status(422)
			templates.TodoForm(&todo, "Item with same title already exists").Render(c.Request.Context(), c.Writer)
			return
		}

		todoState.Todos = append(todoState.Todos, todo)
		templates.TodoItem(todo).Render(c.Request.Context(), c.Writer)
		templates.TodoForm(nil, "").Render(c.Request.Context(), c.Writer)
	})

	r.DELETE("/todo/:title", func(c *gin.Context) {
		title := c.Param("title")
		fmt.Println(title)
		if i := slices.IndexFunc(todoState.Todos, func(m templates.TodoModel) bool {
			return m.Title == title
		}); i > 0 {
			todoState.Todos = append(todoState.Todos[:i], todoState.Todos[i+1:]...)
			c.Status(200)
			return
		}
		c.Status(404)
	})

	r.GET("/scroll", func(c *gin.Context) {
		if _, ok := c.Request.Header["Hx-Request"]; !ok {
			templates.Index("/scroll").Render(c.Request.Context(), c.Writer)
		} else {
			if page, er := strconv.Atoi(c.Query("page")); er != nil {
				items := make([]string, 10)
				for i := 0; i < 10; i++ {
					items[i] = strconv.Itoa(i)
				}
				templates.Scroll(items, page).Render(c.Request.Context(), c.Writer)
			} else {
				items := make([]string, 10)
				for i := 0; i < 10; i++ {
					n := page*10 + i
					items[i] = strconv.Itoa(n)
				}
				templates.ScrollItems(items, page).Render(c.Request.Context(), c.Writer)
			}

		}
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run("0.0.0.0:8080")
}
