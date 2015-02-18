package main

import (
	"errors"
	"os"
	"testing"

	"github.com/cloudfoundry/cli/plugin/fakes"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	cliConn *fakes.FakeCliConnection
	repo    *SeedRepo
)

func TestReadManifest(t *testing.T) {
	setup()
	Convey("Read good manifest", t, func() {
		err := repo.ReadManifest()
		So(err, ShouldBeNil)
		So(repo.Manifest.Organizations[0].Name, ShouldEqual, "org1")
	})

	Convey("No Manifest file", t, func() {
		repo = NewSeedRepo(cliConn, "fake")
		err := repo.ReadManifest()
		So(err, ShouldNotBeNil)
	})

	Convey("Bad Manifest file", t, func() {
		repo = NewSeedRepo(cliConn, "bad.yml")
		err := repo.ReadManifest()
		So(err, ShouldNotBeNil)
	})
}

func TestOrganizations(t *testing.T) {
	setup()
	repo.ReadManifest()
	Convey("Create Organizations", t, func() {
		err := repo.CreateOrganizations()
		So(err, ShouldBeNil)
		So(cliConn.CliCommandCallCount(), ShouldEqual, 2)
	})
	Convey("Delete Organizations", t, func() {
		err := repo.DeleteOrganizations()
		So(err, ShouldBeNil)
		So(cliConn.CliCommandCallCount(), ShouldEqual, 4)
	})
	Convey("Error Creating Organizations", t, func() {
		cliConn.CliCommandReturns([]string{}, errors.New("Error Creating Org"))
		err := repo.CreateOrganizations()
		So(err, ShouldNotBeNil)
	})
	Convey("Error Deleting Organizations", t, func() {
		cliConn.CliCommandReturns([]string{}, errors.New("Error Deleting Org"))
		err := repo.DeleteOrganizations()
		So(err, ShouldNotBeNil)
	})
}

func TestSpaces(t *testing.T) {
	setup()
	repo.ReadManifest()
	Convey("Create Spaces", t, func() {
		err := repo.CreateSpaces()
		So(err, ShouldBeNil)
		So(cliConn.CliCommandCallCount(), ShouldEqual, 5)
	})
	Convey("Delete Spaces", t, func() {
		err := repo.DeleteSpaces()
		So(err, ShouldBeNil)
		So(cliConn.CliCommandCallCount(), ShouldEqual, 10)
	})
	Convey("Error Creating Spaces", t, func() {
		cliConn.CliCommandReturns([]string{}, errors.New("Error Creating Space"))
		err := repo.CreateSpaces()
		So(err, ShouldNotBeNil)
	})
	Convey("Error Deleting Spaces", t, func() {
		cliConn.CliCommandReturns([]string{}, errors.New("Error Deleting Space"))
		err := repo.DeleteSpaces()
		So(err, ShouldNotBeNil)
	})
}

func TestCreateApps(t *testing.T) {
	setup()
	repo.ReadManifest()
	Convey("Create Apps", t, func() {
		err := repo.CreateApps()
		So(err, ShouldBeNil)
		So(cliConn.CliCommandCallCount(), ShouldEqual, 7)
	})
	Convey("Delete Apps", t, func() {
		err := repo.DeleteApps()
		So(err, ShouldBeNil)
		So(cliConn.CliCommandCallCount(), ShouldEqual, 14)
	})
	Convey("Error Creating Apps", t, func() {
		cliConn.CliCommandReturns([]string{}, errors.New("Error Creating Apps"))
		repo.Manifest.Organizations[0].Spaces[0].Apps =
			append(repo.Manifest.Organizations[0].Spaces[0].Apps, DeployApp{Name: "foo"})
		err := repo.CreateApps()
		So(err, ShouldNotBeNil)
	})
	Convey("Error Delete Apps", t, func() {
		cliConn.CliCommandReturns([]string{}, errors.New("Error Deleting Apps"))
		err := repo.DeleteApps()
		So(err, ShouldNotBeNil)
	})
}

func TestCreateServices(t *testing.T) {
	setup()
	repo.ReadManifest()
	Convey("Create Services", t, func() {
		err := repo.CreateServices()
		So(err, ShouldBeNil)
		So(cliConn.CliCommandCallCount(), ShouldEqual, 5)
	})
	Convey("Delete Services", t, func() {
		err := repo.DeleteServices()
		So(err, ShouldBeNil)
		So(cliConn.CliCommandCallCount(), ShouldEqual, 10)
	})
	Convey("Error Creating Services", t, func() {
		cliConn.CliCommandReturns([]string{}, errors.New("Error Creating Services"))
		err := repo.CreateServices()
		So(err, ShouldNotBeNil)
	})
	Convey("Error Deleting Services", t, func() {
		cliConn.CliCommandReturns([]string{}, errors.New("Error Deleting Services"))
		err := repo.DeleteServices()
		So(err, ShouldNotBeNil)
	})
}

func TestSetAppAsService(t *testing.T) {
	setup()
	repo.ReadManifest()
	Convey("Create Services", t, func() {
		app := DeployApp{Name: "foo"}
		err := repo.SetAppAsService(app)
		So(err, ShouldBeNil)
		So(cliConn.CliCommandCallCount(), ShouldEqual, 0)
	})
}

func TestGetAppInfo(t *testing.T) {
	setup()
	repo.ReadManifest()
	Convey("Get App Info", t, func() {
		app := DeployApp{Name: "foo"}
		err := repo.GetAppInfo(app)
		So(err, ShouldBeNil)
		So(cliConn.CliCommandCallCount(), ShouldEqual, 1)
	})
}

func TestDeployApps(t *testing.T) {
	setup()
	tempDir := os.TempDir()
	os.Chdir(tempDir)
	Convey("Deploy App with repo not cloned", t, func() {
		app := DeployApp{Name: "testApp", Repo: "https://github.com/cloudfoundry-community/cf-env"}
		err := repo.DeployApp(app)
		So(err, ShouldBeNil)
	})
	Convey("Deploy App with repo cloned", t, func() {
		app := DeployApp{Name: "testApp", Repo: "https://github.com/cloudfoundry-community/cf-env"}
		err := repo.DeployApp(app)
		So(err, ShouldBeNil)
	})
	Convey("Deploy App with path", t, func() {
		app := DeployApp{Name: "testApp", Path: "test/path"}
		err := repo.DeployApp(app)
		args := cliConn.CliCommandArgsForCall(2)
		So(err, ShouldBeNil)
		So(args, ShouldResemble, []string{"push", "testApp", "-p", "test/path"})
	})
	Convey("Deploy App with disk, memory, instances, domain, hostname, and buildpack", t, func() {
		app := DeployApp{Name: "testApp",
			Path:      "test/path",
			Disk:      "1g",
			Memory:    "256m",
			Instances: "3",
			Domain:    "xip.io",
			Hostname:  "cf-env",
			Buildpack: "my.awesome.buildpack"}
		err := repo.DeployApp(app)
		args := cliConn.CliCommandArgsForCall(3)
		So(err, ShouldBeNil)
		So(args, ShouldResemble, []string{"push", "testApp",
			"-p", "test/path",
			"-k", "1g",
			"-m", "256m",
			"-i", "3",
			"-n", "cf-env",
			"-d", "xip.io",
			"-b", "my.awesome.buildpack"})
	})
	Convey("Deploy App with no repo or path", t, func() {
		app := DeployApp{Name: "testApp"}
		err := repo.DeployApp(app)
		So(err, ShouldNotBeNil)
	})
}

func setup() {
	cliConn = &fakes.FakeCliConnection{}
	repo = NewSeedRepo(cliConn, "example.yml")
}
