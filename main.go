package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Todo struct {
	Id        int    `json:"id"`
	Completed bool   `json:"completed"`
	Content   string `json:"content"`
}

func main() {
	app := fiber.New()

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	todosMap := make(map[int]Todo)

	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := &Todo{}
		if err := c.BodyParser(todo); err != nil {
			return err
		}

		if todo.Content == "" {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Todo content is empty"})
		}

		size := len(todosMap) + 1
		todo.Id = size
		todosMap[size] = *todo

		return c.Status(http.StatusCreated).JSON(fiber.Map{"todos": listTodos(todosMap)})
	})

	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		intId, _ := strconv.Atoi(id)
		delete(todosMap, intId)
		return c.Status(http.StatusOK).JSON(fiber.Map{"todos": listTodos(todosMap)})
	})

	app.Get("/api/todos", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(fiber.Map{"todos": listTodos(todosMap)})
	})

	app.Get("/api/todos/:id", func(c *fiber.Ctx) error {
		if id, err := strconv.Atoi(c.Params("id")); err == nil {
			if todo, exists := todosMap[id]; exists {
				return c.Status(http.StatusOK).JSON(fiber.Map{"todo": todo})
			}
		} else {
			return err
		}
		return c.Status(http.StatusOK).JSON(fiber.Map{"todo": nil})
	})

	app.Patch("/api/todos", func(c *fiber.Ctx) error {
		todo := &Todo{}

		if err := c.BodyParser(todo); err != nil {
			return err
		}

		if todo.Id == 0 {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Todo Id is empty"})
		}

		if _, exists := todosMap[todo.Id]; !exists {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Todo doesn't exist"})
		}

		todosMap[todo.Id] = *todo

		return c.Status(http.StatusOK).JSON(fiber.Map{"todos": listTodos(todosMap)})
	})

	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	log.Fatal(app.Listen(":" + port))
}

func listTodos(todosMap map[int]Todo) []Todo {
	todos := make([]Todo, 0, len(todosMap))

	for _, v := range todosMap {
		todos = append(todos, v)
	}
	return todos
}
