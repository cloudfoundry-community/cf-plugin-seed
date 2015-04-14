package main

type Space struct {
	Name     string
	Apps     map[string]*App
	Services map[string]*Service
	org      *Organization
	seen     []string
}

func (self *Space) create() error {
	_, err := conn.CliCommand("create-space", self.Name)
	if err != nil {
		return err
	}
	conn.CliCommand("target", "-o", self.org.Name, "-s", self.Name)
	for _, app := range self.Apps {
		err = app.create()
		if err != nil {
			return err
		}
	}
	for _, svc := range self.Services {
		err = svc.create()
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *Space) delete() error {
	conn.CliCommand("target", "-o", self.org.Name, "-s", self.Name)
	for _, svc := range self.Services {
		err := svc.delete()
		if err != nil {
			return err
		}
	}
	for _, app := range self.Apps {
		err := app.delete()
		if err != nil {
			return err
		}
	}
	_, err := conn.CliCommand("delete-space", self.Name, "-f")
	if err != nil {
		return err
	}
	return nil
}
