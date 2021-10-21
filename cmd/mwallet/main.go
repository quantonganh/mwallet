package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log"
	"github.com/spf13/viper"

	"github.com/quantonganh/mwallet"
	"github.com/quantonganh/mwallet/account"
	"github.com/quantonganh/mwallet/payment"
	"github.com/quantonganh/mwallet/postgresql"
)

const (
	defaultPort = "8080"
)

func main() {
	var (
		addr = envString("PORT", defaultPort)
		httpAddr          = flag.String("http.addr", ":"+addr, "HTTP listen address")
	)

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		_ = logger.Log("msg", "failed to load the config file", "error", err)
		os.Exit(1)
	}

	var config *mwallet.Config
	if err := viper.Unmarshal(&config); err != nil {
		_ = logger.Log("msg", "failed to unmarshal the config", "error", err)
		os.Exit(1)
	}

	psqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.DB.Host, config.DB.Port, config.DB.User, config.DB.Password, config.DB.Name)

	db, err := sql.Open("postgres", psqlConn)
	if err != nil {
		_ = logger.Log("msg", "failed to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	var (
		accountRepo = postgresql.NewAccountRepository(db)
		paymentRepo = postgresql.NewPaymentRepository(db)
	)

	as := account.NewService(accountRepo)
	ps := payment.NewService(accountRepo, paymentRepo)

	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()

	mux.Handle("/opening/", account.MakeHandler(as, httpLogger))
	mux.Handle("/transferring/", payment.MakeHandler(ps, httpLogger))

	http.Handle("/", accessControl(mux))

	errs := make(chan error, 2)
	go func() {
		_ = logger.Log("transport", "http", "address", *httpAddr, "msg", "listening")
		errs <- http.ListenAndServe(*httpAddr, nil)
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()
	_ = logger.Log("terminated", <-errs)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == http.MethodOptions {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}