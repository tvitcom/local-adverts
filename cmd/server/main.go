package main

import (
	"os"
	"time"
	"context"
	"syscall"
	"net/http"
	"os/signal"
	"database/sql"
	// "golang.org/x/oauth2"
	// "golang.org/x/oauth2/google"
	"github.com/go-ozzo/ozzo-dbx"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"	
	_ "github.com/go-sql-driver/mysql"

	"github.com/tvitcom/local-adverts/internal/advert"
	"github.com/tvitcom/local-adverts/internal/config"
	"github.com/tvitcom/local-adverts/internal/healthcheck"
	"github.com/tvitcom/local-adverts/pkg/dbcontext"
	logz "github.com/tvitcom/local-adverts/pkg/log"
)

// Version indicates the current version of the application.
var (
	Version = "0.0.1"
	// cliCredential GoogleOauth2ClientCredentials
	// oauthConf     *oauth2.Config
	// oauthState    string
)

func init() {
	//Google oauth2 init
	// googlecredential, err := ioutil.ReadFile("./data/credentials.json")
	// if err != nil {
	// 	log.Printf("File error: %v\n", err)
	// 	os.Exit(1)
	// }
	// json.Unmarshal(googlecredential, &cliCredential)

	// oauthConf = &oauth2.Config{
	// 	ClientID:     cliCredential.Cid,
	// 	ClientSecret: cliCredential.Csecret,
	// 	RedirectURL:  conf.GOOGLE_REDIRURL,
	// 	Scopes: []string{
	// 		"https://www.googleapis.com/auth/userinfo.email",   // You have to select your own scope from here -> https://developers.google.com/identity/protocols/oauth2/scopes#oauth2
	// 		"https://www.googleapis.com/auth/userinfo.profile", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/oauth2/scopes#oauth2
	// 	},
	// 	Endpoint: google.Endpoint,
	// }
}
func main() {
	// create root logger tagged with server version
	logger := logz.New().With(nil, "version", Version)

	// load application configurations
	cfg, err := config.Load(logger)
	if err != nil {
		logger.Errorf("failed to load application configuration: %s", err)
		os.Exit(-1)
	}
	db, err := dbx.MustOpen(cfg.DBType, cfg.DSN)
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}
	db.QueryLogFunc = logDBQuery(logger)
	db.ExecLogFunc = logDBExec(logger)
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error(err)
		}
	}()

    // Initialize standard Go html template engine
    renderEngine := html.New("./web/"+cfg.AppThemeUI, ".html")
	router := fiber.New(fiber.Config{
		Network: "tcp4",
		BodyLimit: (700 * 1024 * 1024),// TODO decrease body limit. LAst for my videos to kuzen
		ReadTimeout: 30 * time.Second,
		WriteTimeout: 60 * time.Second,
		EnableTrustedProxyCheck: false,
		Views: renderEngine,
        Prefork:       (cfg.AppMode == "prod"),
        DisableStartupMessage: (cfg.AppMode == "prod"),
        DisableKeepalive: true,
        DisableDefaultDate: true,
        ReduceMemoryUsage: true,
        CaseSensitive: false,
        StrictRouting: true,
        ServerHeader:  cfg.WebservName,
        ErrorHandler: func(c *fiber.Ctx, err error) error {
			 code := http.StatusInternalServerError
			 if e, ok := err.(*fiber.Error); ok {
			   code = e.Code
			 }
			 c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
			 return c.Status(code).SendString(err.Error())
		},
    })
    router.Server().MaxConnsPerIP = 50
    router.Static("/media", config.PictureAdvertsPath, fiber.Static{
	  Compress:      false,
	  ByteRange:     true,
	  Browse:        true,
	  Index:         "index.html",
	  CacheDuration: 48 * 60 * time.Minute,
	  MaxAge:        31536000,
	})
    router.Static("/userpic", config.PictureUserPath)
	router.Static("/assets", "./web/assets/"+cfg.AppThemeUI, fiber.Static{
	  Compress:      true,
	  ByteRange:     true,
	  Browse:        true,
	  Index:         "index.html",
	  CacheDuration: 120 * time.Minute,
	  MaxAge:        31536000,
	})
	mountDinamicRouters(router, logger, dbcontext.New(db), cfg)

	go func() {
		if err := router.Listen(cfg.HttpEntrypoint); err != nil {
			logger.Error(err)
			os.Exit(-1)
		}
	}()
	c := make(chan os.Signal, 1)   // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel
	_ = <-c // This blocks the main thread until an interrupt is received

	_ = router.Shutdown()
	
	// Your cleanup tasks go here
	if err := db.Close(); err != nil {
		logger.Error(err)
	}
	logger.Infof("Server %v Ver: %v was successful shutdown.", cfg.HttpEntrypoint, Version)
}

// mountDinamicRouters sets up the HTTP routing and builds an HTTP handler.
// func mountDinamicRouters(app *fiber.App, logger logz.Logger, db *dbcontext.DB, cfg *config.Config) http.Handler {
func mountDinamicRouters(router *fiber.App, logger logz.Logger, db *dbcontext.DB, cfg *config.Config) {
	
	healthcheck.RegisterHandlers(router, Version)
	
	advert.RegisterHandlers(router, advert.NewAgregator(advert.NewRepository(db, logger), logger), logger)
}

// logDBQuery returns a logging function that can be used to log SQL queries.
func logDBQuery(logger logz.Logger) dbx.QueryLogFunc {
	return func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {
		if err == nil {
			logger.With(ctx, "duration", t.Milliseconds(), "sql", sql).Info("DB query successful")
		} else {
			logger.With(ctx, "sql", sql).Errorf("DB query error: %v", err)
		}
	}
}

// logDBExec returns a logging function that can be used to log SQL executions.
func logDBExec(logger logz.Logger) dbx.ExecLogFunc {
	return func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {
		if err == nil {
			logger.With(ctx, "duration", t.Milliseconds(), "sql", sql).Info("DB execution successful")
		} else {
			logger.With(ctx, "sql", sql).Errorf("DB execution error: %v", err)
		}
	}
}
