package main

import (
	"os"

	"github.com/Bowbaq/scala-imports"
	"github.com/codegangsta/cli"
)

var (
	Version string
	config  = scalaimports.Config{
		Internal: []string{"ai", "common", "dataImport", "emailService", "workflowEngine", "mailgunWebhookService"},
		Lang:     []string{"scala", "java", "javax"},
		Rewrites: map[string]string{
			"Tap._":          "util.Tap._",
			"MustMatchers._": "org.scalatest.MustMatchers._",
			"concurrent.":    "scala.concurrent.",
			"collection.":    "scala.collection.",
			"Keys._":         "sbt.Keys._",
		},
		Ignore: []string{
			"scala.collection.JavaConversions",
			"scala.collection.JavaConverters",
			"scala.concurrent.ExecutionContext.Implicits",
			"scala.language.implicitConversions",
			"scala.language.higherKinds",
			"scala.sys.process",
			"ai.somatix.data.csv.CanBuildFromCsv",
		},
		Remove: []string{
			"import scala.Some",
		},

		MaxLineLength: 110,
	}
)

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
