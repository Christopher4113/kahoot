package main

import (
    "log"
	"os"

	"github.com/joho/godotenv"

    "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/contrib/websocket"
	"context"
    "time"

	"go.mongodb.org/mongo-driver/v2/bson"
    "go.mongodb.org/mongo-driver/v2/mongo"
    "go.mongodb.org/mongo-driver/v2/mongo/options"
)
var quizCollection *mongo.Collection

func main() {
    app := fiber.New()
	app.Use(cors.New(cors.Config{
        AllowOrigins: "http://localhost:5173",
        AllowMethods: "GET,POST,HEAD,PUT,DELETE,OPTIONS",
        AllowHeaders: "Origin, Content-Type, Accept",
    }))
	setUpMongoDB()
	

    app.Get("/", func (c *fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })
	app.Get("/api/quizzes", getQuizzes)

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", msg)

			if err = c.WriteMessage(mt, msg); err != nil {
				log.Println("write:", err)
				break
			}
		}
	}))

    log.Fatal(app.Listen(":3000"))
}

func setUpMongoDB() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    MONGO_URL := os.Getenv("MONGO_URL")

    clientOpts := options.Client().ApplyURI(MONGO_URL)
    client, err := mongo.Connect(clientOpts)
    if err != nil {
        panic(err)
    }
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    if err := client.Ping(ctx, nil); err != nil {
        panic(err)
    }

    quizCollection = client.Database("kahoot").Collection("quizzes")
}

func getQuizzes(c *fiber.Ctx) error {
	cursor, err := quizCollection.Find(context.Background(), bson.M{})
	if err != nil {
		panic(err)
	}
	quizzes := []map[string]any{}
	err = cursor.All(context.Background(), &quizzes)
	if err != nil {
		panic(err)
	}	
	return c.JSON(quizzes)
}