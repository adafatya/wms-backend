package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/adafatya/wms-backend/api"
	db "github.com/adafatya/wms-backend/db/sqlc"
	"github.com/adafatya/wms-backend/util"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	err = util.InitLogger()
	if err != nil {
		log.Fatal("cannot initialize logger:", err)
	}
	defer util.Logger.Sync()

	dbDriver := os.Getenv("DB_DRIVER")
	dbSource := os.Getenv("DB_SOURCE")
	serverAddress := ":" + os.Getenv("PORT")

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		util.Logger.Fatal("cannot connect to db", zap.Error(err))
	}

	err = conn.Ping()
	if err != nil {
		util.Logger.Fatal("cannot ping db", zap.Error(err))
	}

	store := db.NewStore(conn)
	server := api.NewServer(store, util.Logger)

	util.Logger.Info("starting server", zap.String("address", serverAddress))
	err = server.Start(serverAddress)
	if err != nil {
		util.Logger.Fatal("cannot start server", zap.Error(err))
	}
}
