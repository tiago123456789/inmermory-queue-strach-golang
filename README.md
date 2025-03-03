## ABOUT

The project has the main goal implement a inmemory queue from scratch.

### CONCEPTS USED ON PROJECT

- Goroutines.
- LinkedList.
- How to create package to reuse in another projects.

### HOW TO USE THE PACKAGGE

```go
go get "github.com/gofiber/fiber/v2"
go get "github.com/tiago123456789/inmermory-queue-strach-golang/queue"
```

```go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/tiago123456789/inmermory-queue-strach-golang/queue"
)

func Notify(url string, data interface{}) {
	body, _ := json.Marshal(data)
	payload := bytes.NewBuffer(body)
	resp, err := http.Post(url, "application/json", payload)

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return
	}

	defer func() {
		resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		return
	}

	return
}

func notifySubscriptions(message interface{}) error {
	Notify("https://webhook.site/e78fd48c-fb03-4f38-babc-37968ab3e736", message)
	return nil
}

func main() {
	var queueInMemory queue.IQueue
	queueInMemory = queue.New()
	queueInMemory.AddHandler(notifySubscriptions)
	queueInMemory.Start()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		message := map[string]string{
			"message": c.Query("message"),
		}
		queueInMemory.Publish(message)
		return c.SendString("Hello, World!")
	})

	app.Listen(":5000")
}
```
