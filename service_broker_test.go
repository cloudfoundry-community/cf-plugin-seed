package main

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateServiceBroker(t *testing.T) {
	var broker ServiceBroker
	Convey("Preparing the test", t, func() {
		setup("spaces")
		broker = ServiceBroker{
			Name:     "broker1",
			Username: "myuser",
			Password: "mypass",
			Url:      "http://example.com/url",
		}
	})

	Convey("Create ServiceBroker", t, func() {
		err := broker.create()
		So(err, ShouldBeNil)
		So(test_fakeOut, ShouldResemble, []string{
			"create-service-broker broker1 myuser mypass http://example.com/url",
		})
	})
	Convey("Error Creating ServiceBroker", t, func() {
		test_cliConn.CliCommandReturns([]string{}, errors.New("Error Creating ServiceBroker"))
		test_cliConn.CliCommandStub = nil
		err := broker.create()
		So(err, ShouldNotBeNil)
	})
}

func TestDeleteServiceBroker(t *testing.T) {
	var broker ServiceBroker
	Convey("Preparing the test", t, func() {
		setup("spaces")
		broker = ServiceBroker{
			Name:     "broker1",
			Username: "myuser",
			Password: "mypass",
			Url:      "http://example.com/url",
		}
	})

	Convey("Delete ServiceBroker", t, func() {
		err := broker.delete()
		So(err, ShouldBeNil)
		So(test_fakeOut, ShouldResemble, []string{
			"delete-service-broker broker1 -f",
		})
	})
	Convey("Error Deleting ServiceBroker", t, func() {
		test_cliConn.CliCommandReturns([]string{}, errors.New("Error Deleting ServiceBroker"))
		test_cliConn.CliCommandStub = nil
		err := broker.create()
		So(err, ShouldNotBeNil)
	})
}
