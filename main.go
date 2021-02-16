package main

import (
	"context"
	"fmt"
	"net/http"

	"encoding/json"

	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
)

var ctx = context.Background()

type Author struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	e := echo.New()
	e.GET("/name", func(c echo.Context) error {
		name := c.QueryParam("name")
		val, err := get(client, name)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"status":  "FAIL",
				"message": err.Error(),
			})
		}
		data := Author{}
		json.Unmarshal([]byte(val), &data)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "OK",
			"data":   data,
		})
	})

	e.POST("/name", func(c echo.Context) error {
		m := Author{}
		if err := c.Bind(&m); err != nil {
			return err
		}
		name := m.Name
		age := m.Age
		err := set(client, name, Author{Name: name, Age: age})
		if err != nil {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"status": "FAIL",
				"name":   name,
				"age":    age,
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "OK",
			"name":   name,
			"age":    age,
		})
	})
	e.Logger.Fatal(e.Start(":1323"))
}

func get(c *redis.Client, key string) (string, error) {
	val, err := c.Get(key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func set(c *redis.Client, key string, value interface{}) error {
	json, err := json.Marshal(value)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = c.Set(key, json, 0).Err()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
