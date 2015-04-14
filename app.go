package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/cloudfoundry-community/cftype"
	"github.com/cloudfoundry/cli/cf/api/resources"
	"github.com/cloudfoundry/cli/cf/configuration/config_helpers"
	"github.com/cloudfoundry/cli/cf/configuration/core_config"
)

type App struct {
	Name          string        `yaml:",omitempty"`
	Repo          string        `yaml:",omitempty"`
	Path          string        `yaml:",omitempty"`
	Disk          string        `yaml:",omitempty"`
	Memory        string        `yaml:",omitempty"`
	Instances     string        `yaml:",omitempty"`
	Hostname      string        `yaml:",omitempty"`
	Domain        string        `yaml:",omitempty"`
	Buildpack     string        `yaml:",omitempty"`
	Manifest      string        `yaml:",omitempty"`
	ServiceBroker ServiceBroker `yaml:"service_broker,omitempty"`
	ServiceAccess []Service     `yaml:"service_access,omitempty"`
	Requires      Deplist       `yaml:",omitempty"`
	space         *Space
}

func (self *App) create() error {
	err := self.deploy()
	if err != nil {
		return err
	}
	emptyServiceBroker := ServiceBroker{}
	if self.ServiceBroker != emptyServiceBroker {
		fmt.Println("setting app as service")
		err := self.setAsService()
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *App) delete() error {
	emptyServiceBroker := ServiceBroker{}
	if self.ServiceBroker != emptyServiceBroker {
		err := self.unsetAsService()
		if err != nil {
			return err
		}
	}
	_, err := conn.CliCommand("delete", self.Name, "-f", "-r")
	if err != nil {
		return err
	}
	return nil
}

//deployApp deploys a single app
func (self *App) deploy() error {
	args := []string{"push", self.Name}
	if self.Repo != "" {
		wd, _ := os.Getwd()
		appPath := wd + "/apps/" + self.Name
		os.MkdirAll(appPath, 0777)

		files, _ := ioutil.ReadDir(appPath)

		if len(files) == 0 {
			gitPath, err := exec.LookPath("git")
			if err != nil {
				return err
			}
			err = exec.Command(gitPath, "clone", self.Repo, appPath).Run()
			if err != nil {
				return nil
			}
		}
		args = append(args, "-p", appPath)

	} else if self.Path != "" {
		args = append(args, "-p", self.Path)
	} else {
		errMsg := fmt.Sprintf("%s needs either a 'repo' or a 'path' set", self.Name)
		return errors.New(errMsg)
	}

	if self.Disk != "" {
		args = append(args, "-k", self.Disk)
	}
	if self.Memory != "" {
		args = append(args, "-m", self.Memory)
	}
	if self.Instances != "" {
		args = append(args, "-i", self.Instances)
	}
	if self.Hostname != "" {
		args = append(args, "-n", self.Hostname)
	}
	if self.Domain != "" {
		args = append(args, "-d", self.Domain)
	}
	if self.Buildpack != "" {
		args = append(args, "-b", self.Buildpack)
	}
	if self.Manifest != "" {
		args = append(args, "-f", self.Manifest)
	}

	conn.CliCommand(args...)

	return nil
}

func (self *App) setAsService() error {
	appInfo := self.getInfo()
	appRoute, err := firstAppRoute(appInfo)
	if err != nil {
		return err
	}
	self.ServiceBroker.Url = "https://" + appRoute
	err = self.ServiceBroker.create()
	if err != nil {
		return err
	}
	for _, service := range self.ServiceAccess {
		err := service.enableAccess()
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *App) unsetAsService() error {
	appInfo := self.getInfo()
	appRoute, err := firstAppRoute(appInfo)
	if err != nil {
		return err
	}
	self.ServiceBroker.Url = "https://" + appRoute
	for _, service := range self.ServiceAccess {
		err := service.disableAccess()
		if err != nil {
			return err
		}
	}
	err = self.ServiceBroker.delete()
	if err != nil {
		return err
	}
	return nil
}

func (self *App) getInfo() *cftype.RetrieveAParticularApp {
	confRepo := core_config.NewRepositoryFromFilepath(config_helpers.DefaultFilePath(), fatalIf)
	spaceGUID := confRepo.SpaceFields().Guid

	appGUID := findAppGUID(spaceGUID, self.Name)

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
