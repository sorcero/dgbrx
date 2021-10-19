package main

import (
	"github.com/urfave/cli/v2"
	"github.com/withmandala/go-log"
	"gitlab.com/sorcero/community/dgbrx/dgraph"
	"gitlab.com/sorcero/community/dgbrx/ops"
	"os"
)

var logger = log.New(os.Stdout)

var dgraphParams = []cli.Flag{
	&cli.StringFlag{
		Name:     "url",
		Usage:    "Base URL to the dgraph instance, or the base URL to the admin dgraph instance",
		EnvVars:  []string{"DGRAPH_API_URL"},
		Required: true,
	},
	&cli.StringFlag{
		Name:     "api-key",
		Usage:    "API key required for authenticating against Dgraph cloud instance",
		EnvVars:  []string{"DGRAPH_API_KEY"},
		Required: true,
	},
}

func main() {
	app := &cli.App{
		Name:  "dgbrx",
		Usage: "Dgraph Backup and Restore",
		Commands: []*cli.Command{
			{
				Name:  "backup",
				Usage: "Back up a Dgraph Instance locally",
				Flags: dgraphParams,
				Action: func(c *cli.Context) error {
					dg := dgraph.FromCliContext(c)
					return ops.Backup(dg)
				},
			},
			{
				Name:  "restore",
				Usage: "Restore a Dgraph instance from local files",
				Flags: append(dgraphParams, []cli.Flag{
					&cli.StringFlag{Name: "json", Usage: "Path to json backup file", Required: true},
					&cli.StringFlag{Name: "schema", Usage: "Path to schema file", Required: true},
				}...),
				Action: func(c *cli.Context) error {
					dg := dgraph.FromCliContext(c)
					return ops.Restore(
						dg,
						ops.InputOptions{
							JsonPath:   c.String("json"),
							SchemaPath: c.String("schema")},
						ops.OutputOptions{
							JsonPath:   "",
							SchemaPath: "",
						},
					)
				},
			},
			{
				Name:  "clean",
				Usage: "Remove the schema and the entire database",
				Flags: dgraphParams,
				Action: func(c *cli.Context) error {
					dg := dgraph.FromCliContext(c)
					return ops.Clean(dg)
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Fatal(err)
	}
}
