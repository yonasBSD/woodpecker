// Copyright 2023 Woodpecker Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cron

import (
	"context"
	"html/template"
	"os"

	"github.com/urfave/cli/v3"

	"go.woodpecker-ci.org/woodpecker/v3/cli/common"
	"go.woodpecker-ci.org/woodpecker/v3/cli/internal"
	"go.woodpecker-ci.org/woodpecker/v3/woodpecker-go/woodpecker"
)

var cronCreateCmd = &cli.Command{
	Name:      "add",
	Usage:     "add a cron job",
	ArgsUsage: "[repo-id|repo-full-name]",
	Action:    cronCreate,
	Flags: []cli.Flag{
		common.RepoFlag,
		&cli.StringFlag{
			Name:     "name",
			Usage:    "cron name",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "branch",
			Usage: "cron branch",
		},
		&cli.StringFlag{
			Name:     "schedule",
			Usage:    "cron schedule",
			Required: true,
		},
		common.FormatFlag(tmplCronList, true),
	},
}

func cronCreate(ctx context.Context, c *cli.Command) error {
	var (
		cronName         = c.String("name")
		branch           = c.String("branch")
		schedule         = c.String("schedule")
		repoIDOrFullName = c.String("repository")
		format           = c.String("format") + "\n"
	)
	if repoIDOrFullName == "" {
		repoIDOrFullName = c.Args().First()
	}

	client, err := internal.NewClient(ctx, c)
	if err != nil {
		return err
	}

	repoID, err := internal.ParseRepo(client, repoIDOrFullName)
	if err != nil {
		return err
	}

	cron := &woodpecker.Cron{
		Name:     cronName,
		Branch:   branch,
		Schedule: schedule,
	}
	cron, err = client.CronCreate(repoID, cron)
	if err != nil {
		return err
	}
	tmpl, err := template.New("_").Parse(format)
	if err != nil {
		return err
	}
	return tmpl.Execute(os.Stdout, cron)
}
