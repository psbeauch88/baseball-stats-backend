package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

type AtBatResult string

const (
	Single            AtBatResult = "single"
	Double            AtBatResult = "double"
	Triple            AtBatResult = "triple"
	HomeRun           AtBatResult = "home_run"
	StrikeoutSwinging AtBatResult = "strikeout_swinging"
	StrikeoutLooking  AtBatResult = "strikeout_looking"
	Walk              AtBatResult = "walk"
	GroundOut         AtBatResult = "ground_out"
	FlyOut            AtBatResult = "fly_out"
	LineOut           AtBatResult = "line_out"
	FielderChoice     AtBatResult = "fielder_choice"
	SacrificeFly      AtBatResult = "sacrifice_fly"
	SacrificeBunt     AtBatResult = "sacrifice_bunt"
	HitByPitch        AtBatResult = "hit_by_pitch"
)

type AtBat struct {
	ID         int         `json:"id"`
	PlayerID   int         `json:"player_id"`
	Result     AtBatResult `json:"result"`
	PitchCount int         `json:"pitch_count"`
	Date       string      `json:"date"`
}

type Player struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Team        string  `json:"team"`
	Position    string  `json:"position"`
	GamesPlayed int     `json:"games_played"`
	RunsScored  int     `json:"runs_scored"`
	AtBats      []AtBat `json:"at_bats"` // Slice to store at-bats
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Setup Gin router
	router := gin.Default()

	// Add CORS middleware with specific settings
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // React.js dev server
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Database connection
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	defer conn.Close(context.Background())

	log.Print("DB URL: ", os.Getenv("DATABASE_URL"))

	// Test endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/players", func(c *gin.Context) {
		var players []Player
		rows, _ := conn.Query(context.Background(), "SELECT * FROM players")
		for rows.Next() {
			var player Player
			rows.Scan(&player.ID, &player.Name, &player.Team, &player.Position, &player.GamesPlayed, &player.RunsScored, &player.Hits, &player.HomeRuns)
			players = append(players, player)
		}
		c.JSON(200, players)
	})

	// Run server
	router.Run(":8080")
}
