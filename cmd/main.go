package main

import (
	"TraveLite/config"
	"TraveLite/database/postgres"
	_ "TraveLite/docs" // you need to update github.com/rizalgowandy/go-swag-sample with your own project path
	"TraveLite/internal/app"
	"TraveLite/pkg/logger"
	"flag"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"log"
)

// @title Echo Swagger Example API
// @version 1.0
// @description This is a sample server server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http

func main() {
	configFile := flag.String("config", "config/config.yaml", "set -c [config_path] flag")
	logFile := flag.String("log", "stdout", "set -l [log_file] flag")
	flag.Parse()

	c, err := config.ParseConfig(*configFile)
	if err != nil {
		log.Fatalf("Can't parse config with err %v\n", err)
	}

	l, err := logger.InitLogger(c.LogLevel, *logFile)
	if err != nil {
		log.Fatalf("Can't init logger with err %v", err)
	}

	dsn := postgres.NewDSN(c.Postgres.User, c.Postgres.DBName, c.Postgres.Password, c.Postgres.Host, c.Postgres.Port)
	p, err := postgres.NewPostgres(dsn)
	defer p.Close()
	if err != nil {
		log.Fatalf("Can't init database with err %v", err)
	}

	e := app.HandlersInit(p.GetPostgres())

	e.Logger = l

	listerAddr := c.Host + c.Port

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(listerAddr))
}
