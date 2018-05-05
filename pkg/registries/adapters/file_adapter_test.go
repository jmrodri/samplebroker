package adapters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testapb = `version: 1.0
name: hello-world-db-apb
description: A sample APB which deploys Hello World Database
bindable: True
async: optional
metadata:
  displayName: Hello World Database (APB)
  dependencies: ['docker.io/centos/postgresql-94-centos7']
  providerDisplayName: "Red Hat, Inc."
plans:
  - name: default
    description: A sample APB which deploys Hello World Database
    free: True
    metadata:
      displayName: Default
      longDescription: This plan deploys a Postgres Database the Hello World application can connect to
      cost: $0.00
    parameters:
      - name: postgresql_database
        title: PostgreSQL Database Name
        type: string
        default: admin
      - name: postgresql_user
        title: PostgreSQL User
        type: string
        default: admin
      - name: postgresql_password
        title: PostgreSQL Password
        type: string
        default: admin`

func TestRegistryName(t *testing.T) {
	fa := FileAdapter{Name: "testadapter"}
	assert.Equal(t, "testadapter", fa.RegistryName())
}

func TestGetImageNames(t *testing.T) {
	fa := FileAdapter{Name: "testadapter"}
	data, err := fa.GetImageNames()
	if err != nil {
		t.Fatal(err.Error())
	}
	assert.Equal(t, 1, len(data))
}

func TestFetchSpecs(t *testing.T) {
	fa := FileAdapter{Name: "testadapter"}
	imagenames := []string{"hello-world-db-apb"}
	specs, err := fa.FetchSpecs(imagenames)
	if err != nil {
		t.Fatal(err.Error())
	}
	assert.Equal(t, 1, len(specs))
}
