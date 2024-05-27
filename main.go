package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/supertokens/supertokens-golang/recipe/dashboard"
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/recipe/passwordless"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type application struct {
	logger *slog.Logger
	config *configuration
}

func newApp() *application {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	config, err := loadConfig()
	if err != nil {
		logger.Error("failed to load configuration", slog.Any("error", err))
		os.Exit(1)
	}

	return &application{
		logger: logger,
		config: config,
	}
}

func main() {
	app := newApp()

	app.initSupertokens(app.config.supertokens)

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      setupRouter(app.config.supertokens.websiteDomain),
		ErrorLog:     slog.NewLogLogger(slog.NewTextHandler(os.Stderr, nil), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fullAddr := fmt.Sprintf("http://localhost%s", app.config.addr)
	app.logger.Info("starting server", slog.String("addr", fullAddr))

	err := srv.ListenAndServe()
	app.logger.Error("server error", slog.Any("inner_error", err))
	os.Exit(1)

}

func (a *application) initSupertokens(config *supertokensConfig) {
	authBasePath := "/auth"
	err := supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: config.connectionURI,
			APIKey:        config.apiKey,
		},
		AppInfo: supertokens.AppInfo{
			AppName:         config.appName,
			APIDomain:       config.apiDomain,
			WebsiteDomain:   config.websiteDomain,
			APIBasePath:     &authBasePath,
			WebsiteBasePath: &authBasePath,
		},
		RecipeList: []supertokens.Recipe{
			passwordless.Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE",
				ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
					Enabled: true,
				},
			}),
			session.Init(nil),
			dashboard.Init(&dashboardmodels.TypeInput{
				Admins: &config.adminEmails,
			}),
		},
	})

	if err != nil {
		a.logger.Error("failed to initialize supertokens", slog.Any("error", err))
		os.Exit(1)
	}
}
