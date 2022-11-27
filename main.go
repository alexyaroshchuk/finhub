package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"finHubPipeline/calucaltor"
	"finHubPipeline/structs"

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

	symbols := []string{"BINANCE:BTCUSDT", "BINANCE:ETHUSDT", "BINANCE:ADAUSDT"}
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

	maBTC := calucaltor.New(intWindow, db)
	maETH := calucaltor.New(intWindow, db)
	maADA := calucaltor.New(intWindow, db)

	for {
		err := w.ReadJSON(&msg)
		if err != nil {
			log.Fatalf("Error readJSON, %v", err)
		}
		for i := 0; i < len(msg.Data); i++ {
			if msg.Data[i].S == "BINANCE:BTCUSDT" {
				maBTC.CalculateData(msg.Data[i])
				continue
			}
			if msg.Data[i].S == "BINANCE:ETHUSDT" {
				maETH.CalculateData(msg.Data[i])
				continue
			}
			if msg.Data[i].S == "BINANCE:ADAUSDT" {
				maADA.CalculateData(msg.Data[i])
				continue
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
