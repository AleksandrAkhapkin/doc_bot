package main

import (
	"doc_bot/config"
	"doc_bot/infrastructure/events"
	"doc_bot/infrastructure/repositories"
	"doc_bot/infrastructure/restapi"
	"doc_bot/internal/service/telegram"
	"doc_bot/libs/liblog"
	"doc_bot/pkg/application/endpoints"
	"flag"

	"github.com/prometheus/common/log"

	migrate "github.com/rubenv/sql-migrate"

	"gorm.io/driver/postgres"

	"fmt"

	"github.com/friendsofgo/errors"
	b "gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
)

var (
	cfg    *config.Config
	logger liblog.Logger
	db     *gorm.DB
)

func main() {
	var err error
	cfg, err = config.LoadConfig()
	if err != nil {
		log.Fatal(errors.Wrap(err, "config.LoadConfig() error: "))
	}
	if err := cfg.Validate(); err != nil {
		log.Fatal(errors.Wrap(err, "cfg.Validate() error: "))
	}

	logger, err = liblog.NewLogger(cfg.Logger)
	if err != nil {
		log.Fatal(errors.Wrap(err, "liblog.NewLogger error: "))
	}

	db, err = gorm.Open(postgres.Open(
		fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s search_path=%s",
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.User,
			cfg.Database.Name,
			cfg.Database.Password,
			cfg.Database.SSLMode,
			cfg.Database.Schema,
		),
	), &gorm.Config{})
	if err != nil {
		logger.WithError(errors.Wrap(err, "gorm.Open error: ")).Fatal("init doc_bot")
	}

	sqlBD, err := db.DB()
	if err != nil {
		logger.WithError(errors.Wrap(err, "db.DB() error: ")).Fatal("init doc_bot")
	}

	makeMigrate := new(bool)
	flag.BoolVar(makeMigrate, "migrate", false, "path to yaml config")
	flag.Parse()
	if *makeMigrate {
		migrations := &migrate.FileMigrationSource{
			Dir: "migrations",
		}
		count, err := migrate.Exec(sqlBD, cfg.Database.Dialect, migrations, migrate.Up)
		if err != nil {
			logger.WithError(errors.Wrap(err, "migrate.Exec error: ")).Fatal("init doc_bot")
		}
		fmt.Printf("%d MIGRATE IS APPLIED", count)
		return
	}

	shopRepo, err := repositories.NewAnswerRepository(db)
	if err != nil {
		logger.Fatal(err)
	}

	bot, err := telegram.NewBot(cfg.Telegram)
	if err != nil {
		logger.WithError(errors.Wrap(err, "telegram.NewBot error: ")).Fatal("init doc_bot")
	}

	commandLineEvents := events.NewAnswerQueries(shopRepo, logger)
	commandLineEndpoints := endpoints.NewCommandLineEndpoints(commandLineEvents)
	commandLineHandlers := restapi.NewCommandLineHandlers(logger, commandLineEndpoints, bot)

	bot.Handle(b.OnText, commandLineHandlers.TextQuestion)
	bot.Start()
}
