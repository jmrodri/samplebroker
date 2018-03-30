//
// Copyright (c) 2018 Red Hat, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package adapters

import (
	"github.com/automationbroker/bundle-lib/apb"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v1"
)

const apbyaml = `version: 1.0
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

// FileAdapter - Docker Hub Adapter
type FileAdapter struct {
	name string
}

// RegistryName - Retrieve the registry name
func (r FileAdapter) RegistryName() string {
	return r.name
}

// GetImageNames - retrieve the images
func (r FileAdapter) GetImageNames() ([]string, error) {
	var apbData []string
	apbData = append(apbData, "hello-world-db-apb")
	return apbData, nil
}

// FetchSpecs - retrieve the spec for the image names.
func (r FileAdapter) FetchSpecs(imageNames []string) ([]*apb.Spec, error) {
	specs := []*apb.Spec{}
	for _, imageName := range imageNames {
		spec, err := r.loadSpec(imageName)
		if err != nil {
			log.Errorf("Failed to retrieve spec data for image %s - %v", imageName, err)
		}
		if spec != nil {
			specs = append(specs, spec)
		}
	}
	return specs, nil
}

func (r FileAdapter) loadSpec(imageName string) (*apb.Spec, error) {
	var spec apb.Spec

	err := yaml.Unmarshal([]byte(apbyaml), &spec)
	if err != nil {
		return nil, err
	}

	return &spec, nil
}
