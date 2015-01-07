# Overview

Seed creates orgs, spaces, and pushes app. This plugin was created for automation of creating certain orgs, spaces, and apps. This is useful if you want to keep the consistency between different Cloud Foundry Environments. Future support should includes creation of users, services and basically anything in cf cli commands.

## installation

```
$ go get github.com/cloudfoundry-community/cf-plugin-seed
$ cf install-plugin $GOPATH/bin/seed
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
  - name: org1
    spaces:
    - name: space1
      apps:
      - name: app1
        repo: https://github.com/cloudfoundry-community/cf-env
      - name: app2
      path: apps/app1
      memory: 256m
      disk: 1g
      instances: 2
    - name: space2
  - name: org2
    spaces:
    - name: space3
      apps:
      - name: app3
        path: "blah"
      - name: app4
        path: "foo"

```
