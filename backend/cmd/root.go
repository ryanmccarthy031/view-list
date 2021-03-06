package cmd

import (
	"context"
	"net"
	"net/http"

	"github.com/ryanmccarthy031/view-list/pkg/routes"

	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-envconfig"
	"github.com/spf13/cobra"
)

type Config struct {
	Hostname string `env:"HTTP_HOST"`
	Port     string `env:"PORT,default=8080"`
}

var c Config

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("failed to execute root command")
	}
}

var rootCmd = &cobra.Command{
	Use: "view-list-api",
	RunE: func(cmd *cobra.Command, args []string) error {
		r := routes.Setup()
		listenAddr := net.JoinHostPort(routes.GoDotEnvVariable("HTTP_HOST"), routes.GoDotEnvVariable("PORT"))
		log.Fatal().Err(http.ListenAndServe(listenAddr, r)).Msg("failed to start server")
		return nil
	},
}

func init() {
	if err := envconfig.Process(context.Background(), &c); err != nil {
		log.Info().Msg("Processing vars")
		log.Fatal().Err(err).Msg("failed to process environment variables")
	}
}
