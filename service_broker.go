package main

type ServiceBroker struct {
	Name     string
	Username string `yaml:",omitempty"`
	Password string `yaml:",omitempty"`
	Url      string `yaml:",omitempty"`
}

func (self *ServiceBroker) create() error {
	args := []string{"create-service-broker", self.Name, self.Username, self.Password, self.Url}
	_, err := conn.CliCommand(args...)
	return err
}

func (self *ServiceBroker) delete() error {
	args := []string{"delete-service-broker", self.Name, "-f"}
	_, err := conn.CliCommand(args...)
	return err
}
