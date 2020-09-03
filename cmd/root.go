package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/spf13/pflag"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	buildTime                     string
	gitCommit                     string
	cfgFile                       string
	url                           string
	user                          string
	password                      string
	datacenter                    string
	timeout                       int
	responseFileDirectory         string
	responseFileName              = "response.json"
	responseFileDirectoryFallback = "./"

	rootCmd = &cobra.Command{
		Use:     appName,
		Short:   "generic vsphere validations cli",
		Long:    fmt.Sprintf("%v is a CLI library that does generic vsphere validations.", appName),
		Version: version,
	}
)

var (
	appInfo = versionResponse{
		Version:   version,
		Name:      appName,
		GitCommit: gitCommit,
		Built:     buildTime,
	}
)

type versionResponse struct {
	Version   string `json:"version"`
	Name      string `json:"name"`
	GitCommit string `json:"gitCommit"`
	Built     string `json:"built"`
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig, initLogging)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is $HOME/.%v.yaml)", appName))
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
	info, _ := json.Marshal(appInfo)
	rootCmd.SetVersionTemplate(string(info))
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
		viper.SetConfigName("." + appName)
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix(appName)

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	postInitCommands(rootCmd.Commands())
}

func initLogging() {
	log.SetFormatter(&log.JSONFormatter{})

	if responseFileDirectory == "./" {
		curDir, err := os.Getwd()
		if err == nil {
			responseFileDirectory = path.Join(curDir, responseFileDirectory)
		}
	}

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
	_ = viper.BindPFlags(cmd.Flags())
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if viper.IsSet(f.Name) && viper.GetString(f.Name) != "" {
			_ = cmd.Flags().Set(f.Name, viper.GetString(f.Name))
		}
	})
}
