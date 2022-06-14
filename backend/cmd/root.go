package cmd

import (
	"context"
	"net"
	"net/http"
	"os"

	"github.com/ryanmccarthy031/view-list/pkg/routes"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-envconfig"
	"github.com/spf13/cobra"
)

type Config struct {
	Hostname string `env:"HTTP_HOST"`
	Port     string `env:"HTTP_PORT,default=8080"`
}

var c Config

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal().Err(err).Msg("failed to process environment variables")
	}

	return os.Getenv(key)
}
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("failed to execute root command")
	}
}

var rootCmd = &cobra.Command{
	Use: "view-list-api",
	RunE: func(cmd *cobra.Command, args []string) error {
		r := routes.Setup()
		listenAddr := net.JoinHostPort(goDotEnvVariable("HTTP_HOST"), goDotEnvVariable("HTTP_PORT"))
		log.Info().Msgf("Listen: %s", listenAddr)
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
