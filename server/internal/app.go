package internal

import (
    "log"
	"os"

	"github.com/joho/godotenv"

    "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/contrib/websocket"
	"context"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
	"quiz.com/quiz/internal/service"
	"quiz.com/quiz/internal/collection"
	"quiz.com/quiz/internal/controller"
)


type App struct{
	httpServer *fiber.App
	database  *mongo.Database

	quizService  *service.QuizService
	netService  *service.NetService
}

func (a *App) Init() {
	a.setUpMongoDB()
	a.setUpServices()
	a.setUpHttp()

	log.Fatal(a.httpServer.Listen(":3000"))
}

func (a *App) setUpHttp() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	quizController := controller.Quiz(a.quizService)
	app.Get("/api/quizzes", quizController.GetQuizzes)

	wsController := controller.Ws(a.netService)
	app.Get("/ws", websocket.New(wsController.Ws))
	a.httpServer = app
}

func (a *App) setUpServices() {
	a.quizService = service.Quiz(collection.Quiz(a.database.Collection("quizzes")))
	a.netService = service.Net(a.quizService)
}


func (a *App) setUpMongoDB() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    MONGO_URL := os.Getenv("MONGO_URL")

	clientOpts := options.Client().ApplyURI(MONGO_URL)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		panic(err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		panic(err)
	}

	a.database = client.Database("kahoot")
}

