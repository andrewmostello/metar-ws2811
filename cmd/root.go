package cmd

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/andrewmostello/metar-ws281x/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "metar-ws281x",
	Short: "METAR WS281x LED controller",
	Long:  `METAR WS281x LED controller`,
}

func recoverAndLog() {
	if err := recover(); err != nil {
		filenm := fmt.Sprintf("panic.%d.log", time.Now().Unix())

		var w io.Writer = os.Stdout

		if out, err := os.Create(filenm); err == nil {
			w = out
			defer out.Close()
		}

		fmt.Fprintln(w, err)
		w.Write(debug.Stack())

		fmt.Fprintf(os.Stderr, "\nAn error has occurred. See %s for details.\n", filenm)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	defer recoverAndLog()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(Initialize)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", `Config file path. Default is to look for a file called "config" in /etc/metar-ws281x with an accepted extension.
Can be a json, yaml, or toml file; use the appropriate extension on the config file.
e.g. /etc/metar-ws281x/config.json`)

	// Logging options
	config.AddLogFlags(rootCmd)
}

func Initialize() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath("/etc/metar-ws281x/")
		viper.SetConfigName("config")
	}

	// Config available as env vars, with _'s instead of .'s in a key path
	viper.SetEnvPrefix("metar_ws281x")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	usingCfgFile := err == nil

	if usingCfgFile {
		slog.Info("loaded config file", "configFile", viper.ConfigFileUsed())
	}
}
