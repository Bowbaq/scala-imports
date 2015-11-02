package main

import (
	"fmt"
	"os"

	"github.com/Bowbaq/scala-imports"
	"github.com/codegangsta/cli"
	"github.com/spf13/viper"
)

var (
	Version string
	config  scalaimports.Config
)

func init() {
	viper.SetConfigName(".fix-imports")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file: %s\n", err)
		os.Exit(-1)
	}

	viper.Unmarshal(&config)
}

func main() {
	app := cli.NewApp()
	app.Name = "fix-imports"
	app.Usage = "organize imports in a scala project"
	app.Version = Version
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "enable debug output",
		},
	}
	app.Action = func(c *cli.Context) {
		if c.Bool("verbose") {
			config.Verbose = true
		}

		scalaimports.SetConfig(config)

		if len(c.Args()) > 0 {
			for _, path := range c.Args() {
				scalaimports.Format(path)
			}
		} else {
			scalaimports.Format(".")
		}
	}

	app.Run(os.Args)
}
