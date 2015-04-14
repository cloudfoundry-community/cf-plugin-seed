package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateApps(t *testing.T) {
	var app App
	Convey("Preparing the test", t, func() {
		setup("spaces")
		app = App{
			Name: "app1",
			Path: "tmp/app1",
		}
	})

	Convey("Create Apps", t, func() {
		err := app.create()
		So(err, ShouldBeNil)
		So(test_fakeOut, ShouldResemble, []string{
			"push app1 -p tmp/app1",
		})
	})
	Convey("Error Creating Apps", t, func() {
		test_cliConn.CliCommandReturns([]string{}, errors.New("Error Creating Apps"))
		test_cliConn.CliCommandStub = nil
		app = App{Name: "willFail"}
		err := app.create()
		So(err, ShouldNotBeNil)
	})
}

func TestDeleteApp(t *testing.T) {
	var app App
	Convey("Preparing the test", t, func() {
		setup("spaces")
		app = App{
			Name: "app1",
			Path: "tmp/app1",
		}
	})

	Convey("Delete Apps", t, func() {
		err := app.delete()
		So(err, ShouldBeNil)
		So(test_fakeOut, ShouldResemble, []string{
			"delete app1 -f -r",
		})
	})
	Convey("Error Delete Apps", t, func() {
		test_cliConn.CliCommandReturns([]string{}, errors.New("Error Deleting Apps"))
		test_cliConn.CliCommandStub = nil
		err := app.delete()
		So(err, ShouldNotBeNil)
	})
}

func TestDeployApps(t *testing.T) {
	var tempDir string
	var orig_dir string
	var app App
	var err error

	Convey("Preparing the test", t, func() {
		//Create a new temp directory for testing
		tempDir, err = ioutil.TempDir("", "cf-plugin-seed-test-")
		So(err, ShouldBeNil)
		//Save previous working directory so we don't mess up other tests
		orig_dir, err = os.Getwd()
		So(err, ShouldBeNil)
		// Move to temp dir
		err = os.Chdir(tempDir)
		So(err, ShouldBeNil)
		// Grab the real name of the temp dir (symlinks happen from time to time)
		tempDir, err = os.Getwd()
		So(err, ShouldBeNil)

		setup("spaces")
		app = App{
			Name: "app1",
			Path: "tmp/app1",
		}
	})

	Convey("Deploy App with repo not cloned", t, func() {
		setup("spaces")
		app := App{Name: "testApp", Repo: "https://github.com/cloudfoundry-community/cf-env"}
		err := app.deploy()
		So(err, ShouldBeNil)
		So(test_fakeOut, ShouldResemble, []string{
			fmt.Sprintf("push testApp -p %s/apps/testApp", tempDir),
		})
	})

	Convey("Deploy App with repo cloned", t, func() {
		setup("spaces")
		app := App{Name: "testApp", Repo: "https://github.com/cloudfoundry-community/cf-env"}
		err := app.deploy()
		So(err, ShouldBeNil)
		So(test_fakeOut, ShouldResemble, []string{
			fmt.Sprintf("push testApp -p %s/apps/testApp", tempDir),
		})
	})
	Convey("Deploy App with path", t, func() {
		setup("spaces")
		app := App{Name: "testApp", Path: "test/path"}
		err := app.deploy()
		So(err, ShouldBeNil)
		So(test_fakeOut, ShouldResemble, []string{
			"push testApp -p test/path",
		})
	})
	Convey("Deploy App with disk, memory, instances, domain, hostname, and buildpack", t, func() {
		setup("spaces")
		app := App{Name: "testApp",
			Path:      "test/path",
			Disk:      "1g",
			Memory:    "256m",
			Instances: "3",
			Domain:    "xip.io",
			Hostname:  "cf-env",
			Buildpack: "my.awesome.buildpack"}
		err := app.deploy()
		So(err, ShouldBeNil)
		So(test_fakeOut, ShouldResemble, []string{
			"push testApp -p test/path -k 1g -m 256m -i 3 -n cf-env -d xip.io -b my.awesome.buildpack",
		})
	})
	Convey("Deploy App with no repo or path", t, func() {
		app := App{Name: "testApp"}
		err := app.deploy()
		So(err, ShouldNotBeNil)
	})
	//Clean up temp dir
	os.Chdir(orig_dir)
	exec.Command("/bin/rm", "-rf", tempDir).Run()
}

func TestGetAppInfo(t *testing.T) {
	// setup()
	// repo.readManifest()
	// Convey("Get App Info", t, func() {
	// 	app := App{Name: "foo"}
	// 	test_cliConn.CliCommandWithoutTerminalOutputReturns([]string{summaryPayload}, nil)
	// 	repo.GetAppInfo(app)
	// 	So(test_cliConn.CliCommandWithoutTerminalOutputCallCount(), ShouldEqual, 1)
	// })
}

const summaryPayload = `{"guid":"ec2d33f6-fd1c-49a5-9a90-031454d1f1ac","name":"cf-env-test","routes":[{"guid":"9c2f2820-6bd6-43e2-a0a2-c6cd986d67fb","host":"cf-env-test","domain":{"guid":"64796f18-3412-4d8b-82b9-82f7303d58ea","name":"gotapaas.com"}}],"running_instances":1,"services":[{"guid":"d991be9a-fb4a-4b11-9f40-d3d6204479e3","name":"ls-test","bound_app_count":1,"last_operation":null,"dashboard_url":"http://cf-containers-broker.gotapaas.com/manage/instances/e72113ed-fde7-45fd-8758-aff41e6c5507/5218782d-7fab-4534-92b8-434204d88c7b/d991be9a-fb4a-4b11-9f40-d3d6204479e3","service_plan":{"guid":"fff955de-321c-4c7e-bdee-b9622ddce0ca","name":"free","service":{"guid":"c9bbb615-294b-41eb-b6f2-2e77575fa1cc","label":"logstash14","provider":null,"version":null}}}],"available_domains":[{"guid":"64796f18-3412-4d8b-82b9-82f7303d58ea","name":"gotapaas.com"},{"guid":"49b90200-04de-47da-b82b-55f6f62ddeb2","name":"gotapaas.internal"}],"name":"cf-env-test","production":false,"space_guid":"6b6017dd-7333-426d-ab41-8ebd1783ec06","stack_guid":"2c531037-68a2-4e2c-a9e0-71f9d0abf0d4","buildpack":null,"detected_buildpack":"Ruby","environment_json":{},"memory":256,"instances":1,"disk_quota":1024,"state":"STARTED","version":"c4009b76-a64d-49cf-a2fe-2784f4d8cb27","command":null,"console":false,"debug":null,"staging_task_id":"1c751e6695ff4403b45ae93f9b0c76a1","package_state":"STAGED","health_check_type":"port","health_check_timeout":null,"staging_failed_reason":null,"diego":false,"docker_image":null,"package_updated_at":"2015-02-28T01:56:43Z","detected_start_command":"rackup -p $PORT"}`
