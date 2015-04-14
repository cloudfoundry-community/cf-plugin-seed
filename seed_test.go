package main

import (
	"errors"
	"strings"
	"testing"

	"github.com/cloudfoundry/cli/plugin/fakes"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	test_cliConn *fakes.FakeCliConnection
	test_repo    *SeedRepo
	test_fakeOut []string
)

func setup(file string) {
	test_cliConn = &fakes.FakeCliConnection{
		CliCommandStub: func(args ...string) ([]string, error) {
			output := strings.Join(args, " ")
			test_fakeOut = append(test_fakeOut, output)
			return []string{output}, nil
		},
	}
	test_fakeOut = []string{}
	conn = test_cliConn
	test_repo = NewSeedRepo("fixtures/" + file + ".yml")
}

func TestReadManifest(t *testing.T) {
	setup("basic_tree")
	Convey("Read good manifest", t, func() {
		err := test_repo.readManifest()
		So(err, ShouldBeNil)
		org := test_repo.Manifest.Organizations["testOrg"]
		space := org.Spaces["testSpace"]
		app := space.Apps["testApp"]
		svc := space.Services["testSvc"]
		So(org.Name, ShouldEqual, "testOrg")
		So(space.Name, ShouldEqual, "testSpace")
		So(space.org, ShouldEqual, org)
		So(app.Name, ShouldEqual, "testApp")
		So(app.space, ShouldEqual, space)
		So(svc.Name, ShouldEqual, "testSvc")
		So(svc.space, ShouldEqual, space)
	})

	setup("not-a-file")
	Convey("No Manifest file", t, func() {
		err := test_repo.readManifest()
		So(err, ShouldNotBeNil)
	})

	setup("bad")
	Convey("Bad Manifest file", t, func() {
		err := test_repo.readManifest()
		So(err, ShouldNotBeNil)
	})
}

func TestDeployment(t *testing.T) {
	Convey("Preparing the test", t, func() {
		setup("orgs")
		err := test_repo.readManifest()
		So(err, ShouldBeNil)
	})

	Convey("Deploy", t, func() {
		err := test_repo.deploy()
		So(err, ShouldBeNil)
		So(test_fakeOut, ShouldResemble, []string{
			"create-org org1",
			"target -o org1",
			"create-space space1",
			"target -o org1 -s space1",
		})
	})
	Convey("Error Deploying", t, func() {
		test_cliConn.CliCommandReturns([]string{}, errors.New("Error Creating Org"))
		test_cliConn.CliCommandStub = nil
		err := test_repo.deploy()
		So(err, ShouldNotBeNil)
	})
}

func TestCleanup(t *testing.T) {
	Convey("Preparing the test", t, func() {
		setup("orgs")
		err := test_repo.readManifest()
		So(err, ShouldBeNil)
	})

	Convey("Clean Up", t, func() {
		err := test_repo.cleanup()
		So(err, ShouldBeNil)
		So(test_fakeOut, ShouldResemble, []string{
			"target -o org1",
			"target -o org1 -s space1",
			"delete-space space1 -f",
			"delete-org org1 -f",
		})
	})
	Convey("Error Cleaning Up", t, func() {
		test_cliConn.CliCommandReturns([]string{}, errors.New("Error Deleting Org"))
		test_cliConn.CliCommandStub = nil
		err := test_repo.cleanup()
		So(err, ShouldNotBeNil)
	})
}
