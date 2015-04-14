package main

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateService(t *testing.T) {
	var svc Service
	Convey("Preparing the test", t, func() {
		setup("spaces")
		svc = Service{
			Name:    "svc1",
			Service: "myService",
			Plan:    "myPlan",
		}
	})

	Convey("Create Services", t, func() {
		err := svc.create()
		So(err, ShouldBeNil)
		So(test_fakeOut, ShouldResemble, []string{
			"create-service myService myPlan svc1",
		})
	})
	Convey("Error Creating Services", t, func() {
		test_cliConn.CliCommandReturns([]string{}, errors.New("Error Creating Services"))
		test_cliConn.CliCommandStub = nil
		err := svc.create()
		So(err, ShouldNotBeNil)
	})
}

func TestDeleteService(t *testing.T) {
	var svc Service
	Convey("Preparing the test", t, func() {
		setup("spaces")
		svc = Service{
			Name:    "svc1",
			Service: "myService",
			Plan:    "myPlan",
		}
	})

	Convey("Delete Services", t, func() {
		err := svc.delete()
		So(err, ShouldBeNil)
		So(test_fakeOut, ShouldResemble, []string{
			"delete-service svc1 -f",
		})
	})
	Convey("Error Deleting Services", t, func() {
		test_cliConn.CliCommandReturns([]string{}, errors.New("Error Deleting Services"))
		test_cliConn.CliCommandStub = nil
		err := svc.delete()
		So(err, ShouldNotBeNil)
	})
}

func TestServiceAccessEnable(t *testing.T) {
	var svc Service
	Convey("Preparing the test", t, func() {
		setup("spaces")
		svc = Service{
			Name:    "svc1",
			Service: "myService",
		}
	})

	Convey("Enable service access", t, func() {
		err := svc.enableAccess()
		So(err, ShouldBeNil)
		So(test_fakeOut, ShouldResemble, []string{
			"enable-service-access myService",
		})
	})

	Convey("Enable service access with plan and org", t, func() {
		setup("spaces")
		svc = Service{Name: "testBroker", Service: "myService", Plan: "MyPlan", Org: "MyOrg"}
		err := svc.enableAccess()
		So(err, ShouldBeNil)
		So(test_fakeOut, ShouldResemble, []string{
			"enable-service-access myService -p MyPlan -o MyOrg",
		})
	})
	Convey("Error enabling service access", t, func() {
		test_cliConn.CliCommandReturns([]string{}, errors.New("Error enabling Service access"))
		test_cliConn.CliCommandStub = nil
		err := svc.enableAccess()
		So(err, ShouldNotBeNil)
	})
}

func TestServiceDisableAccess(t *testing.T) {
	var svc Service
	Convey("Preparing the test", t, func() {
		setup("spaces")
		svc = Service{
			Name:    "svc1",
			Service: "myService",
		}
	})

	Convey("Disable service access", t, func() {
		err := svc.disableAccess()
		So(err, ShouldBeNil)
		So(test_fakeOut, ShouldResemble, []string{
			"disable-service-access myService",
		})
	})

	Convey("Disable service access with plan and org", t, func() {
		setup("spaces")
		svc = Service{Name: "testBroker", Service: "myService", Plan: "MyPlan", Org: "MyOrg"}
		err := svc.disableAccess()
		So(err, ShouldBeNil)
		So(test_fakeOut, ShouldResemble, []string{
			"disable-service-access myService -p MyPlan -o MyOrg",
		})
	})
	Convey("Error disabling service access", t, func() {
		test_cliConn.CliCommandReturns([]string{}, errors.New("Error disabling Service access"))
		test_cliConn.CliCommandStub = nil
		err := svc.disableAccess()
		So(err, ShouldNotBeNil)
	})
}
