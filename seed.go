package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/cloudfoundry-community/cftype"
	"github.com/cloudfoundry/cli/cf/api/resources"
	"github.com/cloudfoundry/cli/cf/configuration/config_helpers"
	"github.com/cloudfoundry/cli/cf/configuration/core_config"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/codegangsta/cli"
	"gopkg.in/yaml.v2"
)

var conn plugin.CliConnection

func fatalIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stdout, "error:", err)
		os.Exit(1)
	}
}

func main() {
	plugin.Start(&SeedPlugin{})
}

//SeedPlugin empty struct for plugin
type SeedPlugin struct{}

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

	//			err = seedRepo.deleteApps()
	//			fatalIf(err)

	//			err = seedRepo.deleteServices()
	//			fatalIf(err)

	//			err = seedRepo.deleteSpaces()
	//			fatalIf(err)

	//			err = seedRepo.deleteOrganizations()
	//			fatalIf(err)
	//		} else {
	//			err = seedRepo.createOrganizations()
	//			fatalIf(err)

	//			err = seedRepo.createSpaces()
	//			fatalIf(err)

	//			err = seedRepo.createApps()
	//			fatalIf(err)

	//			err = seedRepo.createServices()
	//			fatalIf(err)
	//		}
	//	}
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

//SeedRepo of cli
type SeedRepo struct {
	fileName string
	Manifest SeederManifest
}

func NewSeedRepo(fileName string) *SeedRepo {
	return &SeedRepo{
		fileName: fileName,
	}
}

func (repo *SeedRepo) readManifest() error {
	file, err := ioutil.ReadFile(repo.fileName)
	if err != nil {
		return err
	}
	repo.Manifest = SeederManifest{}

	err = yaml.Unmarshal(file, &repo.Manifest)
	if err != nil {
		return err
	}

	return nil
}

func (repo *SeedRepo) deploy() error {
	for o, org := range repo.Manifest.Organizations {
		org.Name = o
		_, err := conn.CliCommand("create-org", org.Name)
		if err != nil {
			return err
		}
		conn.CliCommand("target", "-o", org.Name)
		for s, space := range org.Spaces {
			space.Name = s
			space.org = org
			space.create()
		}
	}
	return nil
}

func (repo *SeedRepo) cleanup() error {
	for o, org := range repo.Manifest.Organizations {
		conn.CliCommand("target", "-o", org.Name)
		org.Name = o
		for s, space := range org.Spaces {
			space.Name = s
			space.org = org
			space.delete()
		}
		_, err := conn.CliCommand("delete-org", org.Name, "-f")
		if err != nil {
			return err
		}
	}
	return nil
}

func (space *Space) create() error {
	_, err := conn.CliCommand("create-space", space.Name)
	if err != nil {
		return err
	}
	conn.CliCommand("target", "-o", space.org.Name, "-s", space.Name)
	for a, app := range space.Apps {
		app.Name = a
		app.space = space
		err = app.create()
		if err != nil {
			return err
		}
	}
	for s, svc := range space.Services {
		svc.Name = s
		svc.space = space
		err = svc.create()
		if err != nil {
			return err
		}
	}
	return nil
}

func (space *Space) delete() error {
	conn.CliCommand("target", "-o", space.org.Name, "-s", space.Name)
	for s, svc := range space.Services {
		svc.Name = s
		svc.space = space
		err := svc.delete()
		if err != nil {
			return err
		}
	}
	for a, app := range space.Apps {
		app.Name = a
		app.space = space
		err := app.delete()
		if err != nil {
			return err
		}
	}
	conn.CliCommand("target", "-o", space.org.Name)
	_, err := conn.CliCommand("delete-space", space.Name, "-f")
	if err != nil {
		return err
	}
	return nil
}

func (service *Service) create() error {
	_, err := conn.CliCommand("create-service", service.Service, service.Plan, service.Name)
	if err != nil {
		return err
	}
	return nil
}

func (service *Service) delete() error {
	_, err := conn.CliCommand("delete-service", service.Name, "-f")
	if err != nil {
		return err
	}
	return nil
}

func (app *App) create() error {
	err := app.deployApp()
	if err != nil {
		return err
	}
	emptyServiceBroker := ServiceBroker{}
	if app.ServiceBroker != emptyServiceBroker {
		fmt.Println("setting app as service")
		err := app.setAsService()
		if err != nil {
			return err
		}
	}
	return nil
}

func (app *App) delete() error {
	emptyServiceBroker := ServiceBroker{}
	if app.ServiceBroker != emptyServiceBroker {
		err := app.unsetAsService()
		if err != nil {
			return err
		}
	}
	err := app.deleteApp()
	if err != nil {
		return err
	}
	return nil
}

//DeleteApp deletes a single app
func (app *App) deleteApp() error {
	_, err := conn.CliCommand("delete", app.Name, "-f", "-r")
	if err != nil {
		return err
	}

	return nil
}

//deployApp deploys a single app
func (app *App) deployApp() error {
	//	for _, dep := range app.Requires {
	//		objType, name = strings.SplitN(dep, '.', 2)
	//	}

	args := []string{"push", app.Name}
	if app.Repo != "" {
		wd, _ := os.Getwd()
		appPath := wd + "/apps/" + app.Name
		os.MkdirAll(appPath, 0777)

		files, _ := ioutil.ReadDir(appPath)

		if len(files) == 0 {
			gitPath, err := exec.LookPath("git")
			if err != nil {
				return err
			}
			err = exec.Command(gitPath, "clone", app.Repo, appPath).Run()
			if err != nil {
				return nil
			}
		}
		args = append(args, "-p", appPath)

	} else if app.Path != "" {
		args = append(args, "-p", app.Path)
	} else {
		errMsg := fmt.Sprintf("App need repo or path %s", app.Name)
		return errors.New(errMsg)
	}

	if app.Disk != "" {
		args = append(args, "-k", app.Disk)
	}
	if app.Memory != "" {
		args = append(args, "-m", app.Memory)
	}
	if app.Instances != "" {
		args = append(args, "-i", app.Instances)
	}
	if app.Hostname != "" {
		args = append(args, "-n", app.Hostname)
	}
	if app.Domain != "" {
		args = append(args, "-d", app.Domain)
	}
	if app.Buildpack != "" {
		args = append(args, "-b", app.Buildpack)
	}
	if app.Manifest != "" {
		args = append(args, "-f", app.Manifest)
	}

	conn.CliCommand(args...)

	return nil
}

func (app *App) setAsService() error {
	appInfo := getAppInfo(app)
	appRoute, err := firstAppRoute(appInfo)
	if err != nil {
		return err
	}
	app.ServiceBroker.Url = "https://" + appRoute
	err = createServiceBroker(app.ServiceBroker)
	if err != nil {
		return err
	}
	for _, service := range app.ServiceAccess {
		err := enableServiceAccess(service)
		if err != nil {
			return err
		}
	}
	return nil
}

func (app *App) unsetAsService() error {
	appInfo := getAppInfo(app)
	appRoute, err := firstAppRoute(appInfo)
	if err != nil {
		return err
	}
	app.ServiceBroker.Url = "https://" + appRoute
	for _, service := range app.ServiceAccess {
		err := disableServiceAccess(service)
		if err != nil {
			return err
		}
	}
	err = deleteServiceBroker(app.ServiceBroker)
	if err != nil {
		return err
	}
	return nil
}

func getAppInfo(app *App) *cftype.RetrieveAParticularApp {
	confRepo := core_config.NewRepositoryFromFilepath(config_helpers.DefaultFilePath(), fatalIf)
	spaceGUID := confRepo.SpaceFields().Guid

	appGUID := findAppGUID(spaceGUID, app.Name)

	appInfo := findApp(appGUID)
	return appInfo
}

func firstAppRoute(app *cftype.RetrieveAParticularApp) (fullRoute string, err error) {
	routes := &cftype.ListAllRoutesForTheApp{}
	cmd := []string{"curl", app.Entity.RoutesURL}
	output, _ := conn.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &routes)

	if routes.TotalResults == 0 {
		return "", fmt.Errorf("App '%s' has no routes", app.Entity.Name)
	}
	route := routes.Resources[0]

	domain := &cftype.RetrieveAParticularDomain{}
	cmd = []string{"curl", route.Entity.DomainURL}
	output, _ = conn.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &domain)

	if route.Entity.Host != "" {
		return fmt.Sprintf("%s.%s", route.Entity.Host, domain.Entity.Name), nil
	}
	return domain.Entity.Name, nil
}

func findApp(appGUID string) (app *cftype.RetrieveAParticularApp) {
	app = &cftype.RetrieveAParticularApp{}
	cmd := []string{"curl", fmt.Sprintf("/v2/apps/%s", appGUID)}
	output, _ := conn.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &app)
	return app
}

func findAppGUID(spaceGUID string, appName string) string {
	appQuery := fmt.Sprintf("/v2/spaces/%v/apps?q=name:%v&inline-relations-depth=1", spaceGUID, appName)
	cmd := []string{"curl", appQuery}

	output, _ := conn.CliCommandWithoutTerminalOutput(cmd...)
	res := &resources.PaginatedApplicationResources{}
	json.Unmarshal([]byte(strings.Join(output, "")), &res)

	return res.Resources[0].Resource.Metadata.Guid
}

func createServiceBroker(broker ServiceBroker) error {
	args := []string{"create-service-broker", broker.Name, broker.Username, broker.Password, broker.Url}
	_, err := conn.CliCommand(args...)
	return err
}

func deleteServiceBroker(broker ServiceBroker) error {
	args := []string{"delete-service-broker", broker.Name, "-f"}
	_, err := conn.CliCommand(args...)
	return err
}

func enableServiceAccess(service Service) error {
	args := []string{"enable-service-access", service.Service}
	if service.Plan != "" {
		args = append(args, "-p", service.Plan)
	}
	if service.Org != "" {
		args = append(args, "-o", service.Org)
	}
	_, err := conn.CliCommand(args...)
	return err
}

func disableServiceAccess(service Service) error {
	args := []string{"disable-service-access", service.Service}
	if service.Plan != "" {
		args = append(args, "-p", service.Plan)
	}
	if service.Org != "" {
		args = append(args, "-o", service.Org)
	}
	_, err := conn.CliCommand(args...)
	return err
}
