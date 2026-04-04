package app

import (
	"log"

	"github.com/adafatya/wms-backend/internal/config"
	"github.com/adafatya/wms-backend/internal/db"
	"github.com/adafatya/wms-backend/internal/db/sqlc"
	"github.com/adafatya/wms-backend/internal/server"
	"github.com/adafatya/wms-backend/pkg/utils"
	"go.uber.org/zap"
)

func Run() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	err = utils.InitLogger()
	if err != nil {
		log.Fatal("cannot initialize logger:", err)
	}
	defer utils.Logger.Sync()

	conn, err := db.ConnectDB(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		utils.Logger.Fatal("cannot connect to db", zap.Error(err))
	}

	store := sqlc.NewStore(conn)
	srv := server.NewServer(store, utils.Logger)

	utils.Logger.Info("starting server", zap.String("address", cfg.ServerAddress))
	err = srv.Start(cfg.ServerAddress)
	if err != nil {
		utils.Logger.Fatal("cannot start server", zap.Error(err))
	}
}
