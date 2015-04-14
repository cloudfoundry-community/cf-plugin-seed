package main

type Service struct {
	Name     string
	Service  string  `yaml:",omitempty"`
	Plan     string  `yaml:",omitempty"`
	Org      string  `yaml:",omitempty"`
	Requires Deplist `yaml:",omitempty"`
	space    *Space
}

func (self *Service) create() error {
	_, err := conn.CliCommand("create-service", self.Service, self.Plan, self.Name)
	if err != nil {
		return err
	}
	return nil
}

func (self *Service) delete() error {
	_, err := conn.CliCommand("delete-service", self.Name, "-f")
	if err != nil {
		return err
	}
	return nil
}

func (self *Service) enableAccess() error {
	args := []string{"enable-service-access", self.Service}
	if self.Plan != "" {
		args = append(args, "-p", self.Plan)
	}
	if self.Org != "" {
		args = append(args, "-o", self.Org)
	}
	_, err := conn.CliCommand(args...)
	return err
}

func (self *Service) disableAccess() error {
	args := []string{"disable-service-access", self.Service}
	if self.Plan != "" {
		args = append(args, "-p", self.Plan)
	}
	if self.Org != "" {
		args = append(args, "-o", self.Org)
	}
	_, err := conn.CliCommand(args...)
	return err
}
