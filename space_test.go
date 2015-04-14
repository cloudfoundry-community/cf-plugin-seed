package main

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateSpaces(t *testing.T) {
	Convey("Preparing the test", t, func() {
		setup("spaces")
		err := test_repo.readManifest()
		So(err, ShouldBeNil)
	})

	Convey("Create Spaces", t, func() {
		err := test_repo.deploy()
		So(err, ShouldBeNil)
		So(test_fakeOut, ShouldResemble, []string{
			"create-org org1",
			"target -o org1",
			"create-space space1",
			"target -o org1 -s space1",
			"create-service myservice myplan svc1",
			"push app1 -p tmp/app1",
		})
	})
	Convey("Error Creating Spaces", t, func() {
		test_cliConn.CliCommandReturns([]string{}, errors.New("Error Creating Space"))
		test_cliConn.CliCommandStub = nil
		err := test_repo.deploy()
		So(err, ShouldNotBeNil)
	})
}

func TestDeleteSpaces(t *testing.T) {
	Convey("Preparing the test", t, func() {
		setup("spaces")
		err := test_repo.readManifest()
		So(err, ShouldBeNil)
	})

	Convey("Delete Spaces", t, func() {
		err := test_repo.cleanup()
		So(err, ShouldBeNil)
		So(test_fakeOut, ShouldResemble, []string{
			"target -o org1",
			"target -o org1 -s space1",
			"delete app1 -f -r",
			"delete-service svc1 -f",
			"delete-space space1 -f",
			"delete-org org1 -f",
		})
	})
	Convey("Error Deleting Spaces", t, func() {
		test_cliConn.CliCommandReturns([]string{}, errors.New("Error Deleting Space"))
		test_cliConn.CliCommandStub = nil
		err := test_repo.cleanup()
		So(err, ShouldNotBeNil)
	})
}
