package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/codegangsta/cli"
	"gopkg.in/yaml.v2"
)

var conn plugin.CliConnection

type SeederManifest struct {
	Organizations map[string]*Organization
}

type Organization struct {
	Name   string
	Spaces map[string]*Space
}

type Deplist struct {
	Apps     []string `yaml:",omitempty"`
	Services []string `yaml:",omitempty"`
}

//SeedPlugin empty struct for plugin
type SeedPlugin struct{}

//SeedRepo of cli
type SeedRepo struct {
	fileName string
	Manifest SeederManifest
}

func fatalIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stdout, "error:", err)
		os.Exit(1)
	}
}

func main() {
	plugin.Start(&SeedPlugin{})
}

//Run of seeder plugin
func (plugin SeedPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	app := cli.NewApp()
	app.Name = "seed"
	app.Version = VERSION
	app.Author = "Long Nguyen"
	app.Email = "long.nguyen11288@gmail.com"
	app.Usage = "Seeds Cloud Foundry and setups apps/orgs/services on a given Cloud Foundry setup"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "f",
			Value: "",
			Usage: "seed manifest for seeding Cloud Foundry",
		},
		cli.BoolFlag{
			Name:  "c",
			Usage: "cleanup all things created by the manifest",
		},
	}
	app.Action = func(c *cli.Context) {
		if !c.IsSet("f") {
			cli.ShowAppHelp(c)
			os.Exit(1)
		}
		fileName := c.String("f")
		conn = cliConnection
		seedRepo := NewSeedRepo(fileName)

		err := seedRepo.readManifest()
		fatalIf(err)

		if c.Bool("c") {
			err = seedRepo.cleanup()
		} else {
			err = seedRepo.deploy()
		}
		fatalIf(err)
	}

	app.Run(args)
}

//GetMetadata of plugin
func (SeedPlugin) GetMetadata() plugin.PluginMetadata {
	versionParts := strings.Split(string(VERSION), ".")
	major, _ := strconv.Atoi(versionParts[0])
	minor, _ := strconv.Atoi(versionParts[1])
	patch, _ := strconv.Atoi(strings.TrimSpace(versionParts[2]))

	return plugin.PluginMetadata{
		Name: "cf-plugin-seed",
		Version: plugin.VersionType{
			Major: major,
			Minor: minor,
			Build: patch,
		},
		Commands: []plugin.Command{
			{
				Name:     "seed",
				HelpText: "Seeds Cloud Foundry and setups apps/orgs/services on new Cloud Foundry setup",
			},
		},
	}
}

func NewSeedRepo(fileName string) *SeedRepo {
	return &SeedRepo{
		fileName: fileName,
	}
}

func (self *SeedRepo) readManifest() error {
	file, err := ioutil.ReadFile(self.fileName)
	if err != nil {
		return err
	}
	self.Manifest = SeederManifest{}

	err = yaml.Unmarshal(file, &self.Manifest)
	if err != nil {
		return err
	}

	for o, org := range self.Manifest.Organizations {
		org.Name = o
		for s, space := range org.Spaces {
			space.Name = s
			space.org = org
			for a, app := range space.Apps {
				app.Name = a
				app.space = space
			}
			for s, svc := range space.Services {
				svc.Name = s
				svc.space = space
			}
		}
	}

	return nil
}

func (self *SeedRepo) deploy() error {
	for _, org := range self.Manifest.Organizations {
		_, err := conn.CliCommand("create-org", org.Name)
		if err != nil {
			return err
		}
		conn.CliCommand("target", "-o", org.Name)
		for _, space := range org.Spaces {
			err := space.create()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (self *SeedRepo) cleanup() error {
	for _, org := range self.Manifest.Organizations {
		conn.CliCommand("target", "-o", org.Name)
		for _, space := range org.Spaces {
			err := space.delete()
			if err != nil {
				return err
			}
		}
		_, err := conn.CliCommand("delete-org", org.Name, "-f")
		if err != nil {
			return err
		}
	}
	return nil
}
