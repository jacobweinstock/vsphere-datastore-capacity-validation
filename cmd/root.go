package cmd

import (
	"fmt"
	"github.com/spf13/pflag"
	"io"
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile                       string
	url                           string
	user                          string
	password                      string
	datacenter                    string
	timeout                       int
	responseFileDirectory         = "./"
	responseFileName              = "response.json"
	responseFileDirectoryFallback = "./"

	rootCmd = &cobra.Command{
		Use:   "vvalidator",
		Short: "generic vsphere validations cli",
		Long:  `vvalidator is a CLI library that does generic vsphere validations.`,
	}
)

type Validation interface {
	run()
	response(error)
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig, initLogging)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vvalidator.yaml)")
	rootCmd.PersistentFlags().StringVar(&responseFileDirectory, "dir", "./", "directory to write response file")
	rootCmd.PersistentFlags().StringVarP(&url, "url", "u", "", "vCenter url")
	rootCmd.PersistentFlags().StringVarP(&user, "user", "n", "", "vCenter username")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "vCenter password")
	rootCmd.PersistentFlags().StringVarP(&datacenter, "datacenter", "c", "", "vCenter datacenter name")
	rootCmd.PersistentFlags().IntVarP(&timeout, "timeout", "t", 5, "timeout")
	_ = rootCmd.MarkPersistentFlagRequired("url")
	_ = rootCmd.MarkPersistentFlagRequired("user")
	_ = rootCmd.MarkPersistentFlagRequired("password")
	_ = rootCmd.MarkPersistentFlagRequired("datacenter")
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
	postInitCommands(rootCmd.Commands())
}

func initLogging() {
	log.SetFormatter(&log.JSONFormatter{})
	if _, err := os.Stat(responseFileDirectory); os.IsNotExist(err) {
		err := os.Mkdir(responseFileDirectory, 0755)
		if err != nil {
			responseFileDirectory = responseFileDirectoryFallback
		}
	}
	respFile, err := os.Create(path.Join(responseFileDirectory, responseFileName))
	if err != nil {
		responseFileDirectory = responseFileDirectoryFallback
		respFile, err = os.Create(responseFileDirectory + responseFileName)
		if err != nil {
			log.SetOutput(os.Stdout)
			log.WithFields(log.Fields{
				"responseFile": path.Join(responseFileDirectory, responseFileName),
			}).Fatal("could not create response file")
		}

	}
	mw := io.MultiWriter(os.Stdout, respFile)
	log.SetOutput(mw)
}

func postInitCommands(commands []*cobra.Command) {
	for _, cmd := range commands {
		presetRequiredFlags(cmd)
		if cmd.HasSubCommands() {
			postInitCommands(cmd.Commands())
		}
	}
}

func presetRequiredFlags(cmd *cobra.Command) {
	viper.BindPFlags(cmd.Flags())
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if viper.IsSet(f.Name) && viper.GetString(f.Name) != "" {
			cmd.Flags().Set(f.Name, viper.GetString(f.Name))
		}
	})
}
