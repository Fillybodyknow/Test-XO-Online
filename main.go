package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"

	"tic-tac-toe-game/internal/hanlder"
	"tic-tac-toe-game/internal/router"
	"tic-tac-toe-game/internal/service"
	"tic-tac-toe-game/internal/socket"
)

func main() {
	r := gin.Default()

	// เปิด CORS
	r.Use(cors.Default())

	gameService := service.GameService{}
	gameHandler := hanlder.NewGameHandler(gameService)

	router := router.NewGameRouter(gameHandler)
	socket.InitSocketServer(gameHandler)

	server := socket.GetSocketServer()

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("connected:", s.ID())
		return nil
	})

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer server.Close()

	r.GET("/socket.io/*any", gin.WrapH(server))
	r.POST("/socket.io/*any", gin.WrapH(server))

	API := r.Group("/api")
	{
		router.InitRouter(API)
	}

	if err := r.Run(":3000"); err != nil {
		log.Fatal("failed run app: ", err)
	}
}
