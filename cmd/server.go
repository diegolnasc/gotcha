package cmd

import (
	"github.com/diegolnasc/gotcha/pkg/handler"
	c "github.com/diegolnasc/gotcha/pkg/model"

	g "github.com/diegolnasc/gotcha/pkg/provider/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Provider string

const (
	GitHub Provider = "github"
)

// config is the server command configuration.
type config struct {
	providerWorker handler.ProviderWorker
	// providerName   Provider
	config c.Settings
}

// Init returns the command cli.
func Init() *cobra.Command {
	c := &config{}
	return &cobra.Command{
		Use:   "server",
		Short: "Show the current version",
		PreRun: func(cmd *cobra.Command, args []string) {
			c.preRun()
		},
		Run: func(cmd *cobra.Command, args []string) {
			c.providerWorker = g.New(&c.config)
			h := &handler.ServerHandler{
				Config:         c.config,
				ProviderWorker: c.providerWorker,
			}
			h.NewServer()
		},
	}
}

//preRun loads the version control provider configuration.
func (c *config) preRun() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./build")
	if err := viper.ReadInConfig(); err != nil {
		panic("Unable to read the configuration file (config.yaml).")
	}
	if err := viper.Unmarshal(&c.config); err != nil {
		panic("Unable to unmarshal the configuration.")
	}

}
