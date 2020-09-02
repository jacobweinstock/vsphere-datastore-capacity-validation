package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	jsonLogging bool
	verbose     bool

	rootCmd = &cobra.Command{
		Use:   "vvalidator",
		Short: "generic vsphere validations cli",
		Long:  `vvalidator is a CLI library that does generic vsphere validations.`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig, initLogging)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vvalidator.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&jsonLogging, "json", "j", false, "Enable json logging")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose logging")
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			er(err)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".vvalidator")
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("vvalidator")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func initLogging() {
	log.SetOutput(os.Stdout)
	if verbose {
		log.SetReportCaller(true)
		log.SetLevel(log.DebugLevel)
	}
	if jsonLogging {
		log.SetFormatter(&log.JSONFormatter{})
	}
}
