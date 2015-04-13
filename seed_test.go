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

func TestreadManifest(t *testing.T) {
	setup()
	Convey("Read good manifest", t, func() {
		err := repo.readManifest()
		So(err, ShouldBeNil)
		So(repo.Manifest.Organizations[0].Name, ShouldEqual, "org1")
	})

	Convey("No Manifest file", t, func() {
		repo = NewSeedRepo("fake")
		err := repo.readManifest()
		So(err, ShouldNotBeNil)
	})

	Convey("Bad Manifest file", t, func() {
		repo = NewSeedRepo("bad.yml")
		err := repo.readManifest()
		So(err, ShouldNotBeNil)
	})
}

func TestOrganizations(t *testing.T) {
	setup()
	repo.readManifest()
	Convey("Create Organizations", t, func() {
		err := repo.createOrganizations()
		So(err, ShouldBeNil)
		So(cliConn.CliCommandCallCount(), ShouldEqual, 2)
	})
	Convey("Delete Organizations", t, func() {
		err := repo.deleteOrganizations()
		So(err, ShouldBeNil)
		So(cliConn.CliCommandCallCount(), ShouldEqual, 4)
	})
	Convey("Error Creating Organizations", t, func() {
		cliConn.CliCommandReturns([]string{}, errors.New("Error Creating Org"))
		err := repo.createOrganizations()
		So(err, ShouldNotBeNil)
	})
	Convey("Error Deleting Organizations", t, func() {
		cliConn.CliCommandReturns([]string{}, errors.New("Error Deleting Org"))
		err := repo.deleteOrganizations()
		So(err, ShouldNotBeNil)
	})
}

func TestSpaces(t *testing.T) {
	setup()
	repo.readManifest()
	Convey("Create Spaces", t, func() {
		err := repo.createSpaces()
		So(err, ShouldBeNil)
		So(cliConn.CliCommandCallCount(), ShouldEqual, 5)
	})
	Convey("Delete Spaces", t, func() {
		err := repo.deleteSpaces()
		So(err, ShouldBeNil)
		So(cliConn.CliCommandCallCount(), ShouldEqual, 10)
	})
	Convey("Error Creating Spaces", t, func() {
		cliConn.CliCommandReturns([]string{}, errors.New("Error Creating Space"))
		err := repo.createSpaces()
		So(err, ShouldNotBeNil)
	})
	Convey("Error Deleting Spaces", t, func() {
		cliConn.CliCommandReturns([]string{}, errors.New("Error Deleting Space"))
		err := repo.deleteSpaces()
		So(err, ShouldNotBeNil)
	})
}

func TestCreateApps(t *testing.T) {
	setup()
	repo.readManifest()
	Convey("Create Apps", t, func() {
		err := repo.createApps()
		So(err, ShouldBeNil)
		So(cliConn.CliCommandCallCount(), ShouldEqual, 7)
	})
	Convey("Delete Apps", t, func() {
		err := repo.deleteApps()
		So(err, ShouldBeNil)
		So(cliConn.CliCommandCallCount(), ShouldEqual, 14)
	})
	Convey("Error Creating Apps", t, func() {
		cliConn.CliCommandReturns([]string{}, errors.New("Error Creating Apps"))
		repo.Manifest.Organizations[0].Spaces[0].Apps =
			append(repo.Manifest.Organizations[0].Spaces[0].Apps, App{Name: "foo"})
		err := repo.createApps()
		So(err, ShouldNotBeNil)
	})
	Convey("Error Delete Apps", t, func() {
		cliConn.CliCommandReturns([]string{}, errors.New("Error Deleting Apps"))
		err := repo.deleteApps()
		So(err, ShouldNotBeNil)
	})
}

func TestCreateServices(t *testing.T) {
	setup()
	repo.readManifest()
	Convey("Create Services", t, func() {
		err := repo.createServices()
		So(err, ShouldBeNil)
		So(cliConn.CliCommandCallCount(), ShouldEqual, 5)
	})
	Convey("Delete Services", t, func() {
		err := repo.deleteServices()
		So(err, ShouldBeNil)
		So(cliConn.CliCommandCallCount(), ShouldEqual, 10)
	})
	Convey("Error Creating Services", t, func() {
		cliConn.CliCommandReturns([]string{}, errors.New("Error Creating Services"))
		err := repo.createServices()
		So(err, ShouldNotBeNil)
	})
	Convey("Error Deleting Services", t, func() {
		cliConn.CliCommandReturns([]string{}, errors.New("Error Deleting Services"))
		err := repo.deleteServices()
		So(err, ShouldNotBeNil)
	})
}

func TestGetAppInfo(t *testing.T) {
	// setup()
	// repo.readManifest()
	// Convey("Get App Info", t, func() {
	// 	app := App{Name: "foo"}
	// 	cliConn.CliCommandWithoutTerminalOutputReturns([]string{summaryPayload}, nil)
	// 	repo.GetAppInfo(app)
	// 	So(cliConn.CliCommandWithoutTerminalOutputCallCount(), ShouldEqual, 1)
	// })
}

func TestDeployApps(t *testing.T) {
	setup()
	tempDir := os.TempDir()
	os.Chdir(tempDir)
	Convey("Deploy App with repo not cloned", t, func() {
		app := App{Name: "testApp", Repo: "https://github.com/cloudfoundry-community/cf-env"}
		err := repo.deployApp(app)
		So(err, ShouldBeNil)
	})
	Convey("Deploy App with repo cloned", t, func() {
		app := App{Name: "testApp", Repo: "https://github.com/cloudfoundry-community/cf-env"}
		err := repo.deployApp(app)
		So(err, ShouldBeNil)
	})
	Convey("Deploy App with path", t, func() {
		app := App{Name: "testApp", Path: "test/path"}
		err := repo.deployApp(app)
		args := cliConn.CliCommandArgsForCall(2)
		So(err, ShouldBeNil)
		So(args, ShouldResemble, []string{"push", "testApp", "-p", "test/path"})
	})
	Convey("Deploy App with disk, memory, instances, domain, hostname, and buildpack", t, func() {
		app := App{Name: "testApp",
			Path:      "test/path",
			Disk:      "1g",
			Memory:    "256m",
			Instances: "3",
			Domain:    "xip.io",
			Hostname:  "cf-env",
			Buildpack: "my.awesome.buildpack"}
		err := repo.deployApp(app)
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
		app := App{Name: "testApp"}
		err := repo.deployApp(app)
		So(err, ShouldNotBeNil)
	})
}

func TestServiceBroker(t *testing.T) {
	setup()
	repo.readManifest()
	Convey("Create service broker", t, func() {
		broker := ServiceBroker{Name: "testBroker"}
		err := repo.createServiceBroker(broker)
		So(err, ShouldBeNil)
		So(cliConn.CliCommandCallCount(), ShouldEqual, 1)
	})
	Convey("Delete service broker", t, func() {
		broker := ServiceBroker{Name: "testBroker"}
		err := repo.deleteServiceBroker(broker)
		So(err, ShouldBeNil)
		So(cliConn.CliCommandCallCount(), ShouldEqual, 2)
	})
}

func TestServiceAccess(t *testing.T) {
	setup()
	repo.readManifest()
	Convey("Enable service access", t, func() {
		serviceAccess := Service{Name: "testBroker"}
		err := repo.enableServiceAccess(serviceAccess)
		So(err, ShouldBeNil)
		So(cliConn.CliCommandCallCount(), ShouldEqual, 1)
	})

	Convey("Enable service access with plan and org", t, func() {
		serviceAccess := Service{Name: "testBroker", Service: "myService", Plan: "MyPlan", Org: "MyOrg"}
		err := repo.enableServiceAccess(serviceAccess)
		args := cliConn.CliCommandArgsForCall(1)
		So(err, ShouldBeNil)
		So(cliConn.CliCommandCallCount(), ShouldEqual, 2)
		So(args, ShouldResemble, []string{"enable-service-access", "myService",
			"-p", "MyPlan",
			"-o", "MyOrg"})
	})

	Convey("Disable service access", t, func() {
		serviceAccess := Service{Name: "testBroker"}
		err := repo.disableServiceAccess(serviceAccess)
		So(err, ShouldBeNil)
		So(cliConn.CliCommandCallCount(), ShouldEqual, 3)
	})

	Convey("Disable service access with plan and org", t, func() {
		serviceAccess := Service{Name: "testBroker", Service: "myService", Plan: "MyPlan", Org: "MyOrg"}
		err := repo.disableServiceAccess(serviceAccess)
		args := cliConn.CliCommandArgsForCall(3)
		So(err, ShouldBeNil)
		So(cliConn.CliCommandCallCount(), ShouldEqual, 4)
		So(args, ShouldResemble, []string{"disable-service-access", "myService",
			"-p", "MyPlan",
			"-o", "MyOrg"})
	})
}

func setup() {
	cliConn = &fakes.FakeCliConnection{}
	conn = cliConn
	repo = NewSeedRepo("example.yml")
}

const summaryPayload = `{"guid":"ec2d33f6-fd1c-49a5-9a90-031454d1f1ac","name":"cf-env-test","routes":[{"guid":"9c2f2820-6bd6-43e2-a0a2-c6cd986d67fb","host":"cf-env-test","domain":{"guid":"64796f18-3412-4d8b-82b9-82f7303d58ea","name":"gotapaas.com"}}],"running_instances":1,"services":[{"guid":"d991be9a-fb4a-4b11-9f40-d3d6204479e3","name":"ls-test","bound_app_count":1,"last_operation":null,"dashboard_url":"http://cf-containers-broker.gotapaas.com/manage/instances/e72113ed-fde7-45fd-8758-aff41e6c5507/5218782d-7fab-4534-92b8-434204d88c7b/d991be9a-fb4a-4b11-9f40-d3d6204479e3","service_plan":{"guid":"fff955de-321c-4c7e-bdee-b9622ddce0ca","name":"free","service":{"guid":"c9bbb615-294b-41eb-b6f2-2e77575fa1cc","label":"logstash14","provider":null,"version":null}}}],"available_domains":[{"guid":"64796f18-3412-4d8b-82b9-82f7303d58ea","name":"gotapaas.com"},{"guid":"49b90200-04de-47da-b82b-55f6f62ddeb2","name":"gotapaas.internal"}],"name":"cf-env-test","production":false,"space_guid":"6b6017dd-7333-426d-ab41-8ebd1783ec06","stack_guid":"2c531037-68a2-4e2c-a9e0-71f9d0abf0d4","buildpack":null,"detected_buildpack":"Ruby","environment_json":{},"memory":256,"instances":1,"disk_quota":1024,"state":"STARTED","version":"c4009b76-a64d-49cf-a2fe-2784f4d8cb27","command":null,"console":false,"debug":null,"staging_task_id":"1c751e6695ff4403b45ae93f9b0c76a1","package_state":"STAGED","health_check_type":"port","health_check_timeout":null,"staging_failed_reason":null,"diego":false,"docker_image":null,"package_updated_at":"2015-02-28T01:56:43Z","detected_start_command":"rackup -p $PORT"}`
