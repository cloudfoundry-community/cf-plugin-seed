package main

import (
	"errors"
	"fmt"
	"strings"
)

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
	order, err := self.resolveOrder()
	if err != nil {
		return err
	}
	for i := 0; i < len(order); i++ {
		vals := strings.SplitN(order[i], ".", 2)
		objType := vals[0]
		key := vals[1]
		if objType == "app" {
			err = self.Apps[key].create()
		} else if objType == "svc" {
			err = self.Services[key].create()
		} else {
			err = errors.New(fmt.Sprintf("Unsupported object type %q", objType))
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *Space) delete() error {
	conn.CliCommand("target", "-o", self.org.Name, "-s", self.Name)
	order, err := self.resolveOrder()
	if err != nil {
		return err
	}
	for i := len(order) - 1; i >= 0; i-- {
		vals := strings.SplitN(order[i], ".", 2)
		objType := vals[0]
		key := vals[1]
		if objType == "app" {
			err = self.Apps[key].delete()
		} else if objType == "svc" {
			err = self.Services[key].delete()
		} else {
			err = errors.New(fmt.Sprintf("Unsupported object type %q", objType))
		}
		if err != nil {
			return err
		}
	}
	_, err = conn.CliCommand("delete-space", self.Name, "-f")
	if err != nil {
		return err
	}
	return nil
}

func (self *Space) resolveOrder() ([]string, error) {
	order := []string{}
	seen := map[string]bool{}
	processed := map[string]bool{}
	for _, svc := range self.Services {
		err := self.resolveDeps("svc", svc.Name, svc.Requires, &order, seen, processed)
		if err != nil {
			return nil, err
		}
	}
	for _, app := range self.Apps {
		err := self.resolveDeps("app", app.Name, app.Requires, &order, seen, processed)
		if err != nil {
			return nil, err
		}
	}
	return order, nil
}

func (self *Space) resolveDeps(kind string, name string, deps Deplist, order *[]string, seen map[string]bool, processed map[string]bool) error {
	k := kind + "." + name
	if _, ok := seen[k]; ok {
		return errors.New(fmt.Sprintf("Cicrular dependency detected on %q", k))
	}
	seen[k] = true
	if _, ok := processed[k]; ok {
		return nil
	}

	if len(deps.Services) > 0 {
		for _, svcName := range deps.Services {
			svc, ok := self.Services[svcName]
			if !ok {
				fmt.Printf("Failed to find %q", svcName)
				return errors.New(fmt.Sprintf("Unable to find dependent service %q", svcName))
			}
			err := self.resolveDeps("svc", svcName, svc.Requires, order, seen, processed)
			if err != nil {
				return err
			}
		}
	}
	if len(deps.Apps) > 0 {
		for _, appName := range deps.Apps {
			app, ok := self.Apps[appName]
			if !ok {
				return errors.New(fmt.Sprintf("Unable to find dependent app %q", appName))
			}
			err := self.resolveDeps("app", appName, app.Requires, order, seen, processed)
			if err != nil {
				return err
			}
		}
	}

	*order = append(*order, k)
	delete(seen, k)
	processed[k] = true
	return nil
}
