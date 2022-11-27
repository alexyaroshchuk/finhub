package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"finnhubPipeline/calculator"
	"finnhubPipeline/structs"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	}
	apiKey := os.Getenv("API_KEY")
	w, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("wss://ws.finnhub.io?token=%s", apiKey), nil)
	if err != nil {
		log.Fatalf("Error dial connection, %v", err)
	}
	defer w.Close()

	symbols := strings.Split(os.Getenv("SYMBOLS"), ",")
	for _, s := range symbols {
		msg, _ := json.Marshal(map[string]interface{}{"type": "subscribe", "symbol": s})

		err = w.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Fatalf("Error write message, %v", err)
		}
	}

	var msg structs.Msg
	db := connectToPgSQL()
	intWindow, err := strconv.Atoi(os.Getenv("WINDOW"))
	if err != nil {
		log.Fatalf("Error get window, %v", err)
	}

	coinsMap := map[string]*calculator.MovingAverage{}
	for _, v := range symbols {
		coinsMap[v] = calculator.New(intWindow, db)
	}

	for {
		err := w.ReadJSON(&msg)
		if err != nil {
			log.Fatalf("Error readJSON, %v", err)
		}
		for i := 0; i < len(msg.Data); i++ {
			if calc, ok := coinsMap[msg.Data[i].S]; ok {
				calc.CalculateData(msg.Data[i])
			}
		}

		log.Printf("Message from server: %+v\n", msg)
	}
}

func connectToPgSQL() *sql.DB {
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Fatalf("Error create pgsql connection, %v", err)
	}

	return db
}
