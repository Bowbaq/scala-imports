package main

import (
	"os"

	"github.com/Bowbaq/scala-imports"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "fix-imports"
	app.Usage = "organize imports in a scala project"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "enable debug output",
		},
	}
	app.Action = func(c *cli.Context) {
		if c.Bool("verbose") {
			scalaimports.Verbose = true
		}

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
