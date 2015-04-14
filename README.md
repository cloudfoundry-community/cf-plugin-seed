# Overview

Seed creates orgs, spaces, services, and pushes apps, to initialize a Cloud Foundry environment from a manifest file.  This is useful if you want to keep the consistency between different Cloud Foundry Environments, or automate the deployment of your Cloud Foundry Environment. Eventually, Future support will include basically anything in the cf cli commands.

## installation

```
$ go get github.com/cloudfoundry-community/cf-plugin-seed
$ cf install-plugin $GOPATH/bin/cf-plugin-seed
```

## usage

```
$ cf seed -f example.yml
```

## manifest

Example of seed manifest.yml

```yaml
---
organizations:
  org1:
    spaces:
      space1:
        apps:
          app1:
            repo: https://github.com/cloudfoundry-community/cf-env
          app2:
            path: apps/app1
            memory: 256m
            disk: 1g
            instances: 2
      space2: {}
  org2:
    spaces:
      space3:
        apps:
          app3:
            path: "blah"
          app4:
            path: "foo"

```

Example of a seed manifest.yml with services and dependencies:

```yaml
---
organizations:
  org1:
    spaces:
      space1: {}
      space2: {}
  org2:
    spaces:
      space3:
        apps:
          app5:
            repo: https://github.com/cloudfoundry-community/worlds-simplest-service-broker
            manifest: manifests/wssb.yml
            service_broker:
              name: wssb
              username: admin
              password: admin
            service_access:
            - name: wssb
              service: wssb
            requires:
              services:
              - postgresql
        services:
          wssb-service:
            service: wssb
            plan: shared
            requires:
              apps:
              - app5
          postgresql:
            service: postgresql
```

## upgrading seed

```bash
cf uninstall-plugin cf-plugin-seed
go get -u github.com/cloudfoundry-community/cf-plugin-seed
cf install-plugin cf-plugin-seed
```

*Note* if you are running <0.0.1 of seed please run first.

```
cf uninstall-plugin seed
```
