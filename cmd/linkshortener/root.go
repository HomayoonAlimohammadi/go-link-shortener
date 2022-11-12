package main

import (
	"log"
	"os"

	"github.com/homayoonalimohammadi/go-link-shortener/linkshortener/internal/app/database"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	Author   string
	UseViper bool
	License  string
	Postgres database.PostgresConfig
	Redis    database.RedisConfig
}

var (
	// Used for flags.
	cfgFile     string
	userLicense string
	config      Config
	rootCmd     = &cobra.Command{
		Short: "Shortens any given link to a digestable token.",
		Long: `Link Shortener is a web app to shorten any given link
to a unique phrase so you can use that any where.`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add commands
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(versionCmd)

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is ./.cobra.yaml)")
	rootCmd.PersistentFlags().StringP("author", "a", "Homayoon Alimohammadi", "author name for copyright attribution")
	rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	rootCmd.PersistentFlags().BoolP("viper", "v", true, "use Viper for configuration")
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("author", "Homayoon Alimohammadi <homayoonalimohammadi@gmail.com>")
	viper.SetDefault("license", "MIT")

}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		currentDir, err := os.Getwd()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(currentDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cobra")
		log.Println("reading config file from:", currentDir)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		log.Println("could not unmarshal config")
	}
	log.Printf("Config: %+v", config)
}
