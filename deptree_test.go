package main

import (
	//	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDeptreeResolution(t *testing.T) {
	Convey("Preparing the test", t, func() {
		setup("deptree")
		err := test_repo.readManifest()
		So(err, ShouldBeNil)
	})

	Convey("Deploy + rollback", t, func() {
		err := test_repo.deploy()
		So(err, ShouldBeNil)
		err = test_repo.cleanup()
		So(test_fakeOut, ShouldResemble, []string{
			"create-org testOrg",
			"target -o testOrg",
			"create-space testSpace",
			"target -o testOrg -s testSpace",
			"create-service postgresql free pg",
			"push pgApp -p bar",
			"create-service app shared appAsSvc",
			"push app1 -p foo",
			// Now for rollback
			"target -o testOrg",
			"target -o testOrg -s testSpace",
			"delete app1 -f -r",
			"delete-service appAsSvc -f",
			"delete pgApp -f -r",
			"delete-service pg -f",
			"delete-space testSpace -f",
			"delete-org testOrg -f",
		})
	})
}

func TestDeptreeFailures(t *testing.T) {
	Convey("Circular dependencies", t, func() {
		setup("circular-deps")
		err := test_repo.readManifest()
		So(err, ShouldBeNil)

		err = test_repo.deploy()
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "Cicrular dependency detected on \"app.myApp\"")
	})
	Convey("Unsolvable dependencies", t, func() {
		setup("unsolved-deps")
		err := test_repo.readManifest()
		So(err, ShouldBeNil)

		err = test_repo.deploy()
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "Unable to find dependent service \"postgresql\"")
	})
}
