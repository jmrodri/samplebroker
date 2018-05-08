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

package main

import (
	"net/url"
	"os"

	"github.com/automationbroker/bundle-lib/registries"
	bladapters "github.com/automationbroker/bundle-lib/registries/adapters"
	flags "github.com/jessevdk/go-flags"
	"github.com/jmrodri/samplebroker/pkg/registries/adapters"
	"github.com/openshift/ansible-service-broker/pkg/app"
	log "github.com/sirupsen/logrus"
)

func main() {

	var args app.Args
	var err error

	// Writing directly to stderr because log has not been bootstrapped
	if args, err = app.CreateArgs(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	// To add your custom registries, define an entry in this array.
	regs := []registries.Registry{}

	// File adapter
	c := registries.Config{
		URL:        "",
		User:       "",
		Pass:       "",
		Org:        "jmrodri",
		Tag:        "latest",
		Type:       "file", // bundle-lib registry.go needs to change for this
		Name:       "foo",
		Images:     []string{"hello-world-db-apb"},
		Namespaces: []string{"openshift"},
		Fail:       false,
		WhiteList:  []string{".*-apb$"},
		BlackList:  []string{},
		AuthType:   "",
		AuthName:   "",
		Runner:     "",
	}

	fadapter := adapters.FileAdapter{Name: "foobar"}
	log.Info("Creating a new custom registry")
	reg, err := registries.NewCustomRegistry(c, fadapter, "openshift")
	if err != nil {
		log.Errorf(
			"Failed to initialize foo Registry err - %v \n", err)
		os.Exit(1)
	}

	regs = append(regs, reg)

	// Dockerhub adapter
	c = registries.Config{
		URL:        "https://registry.hub.docker.com",
		User:       "",
		Pass:       "",
		Org:        "jmrodri",
		Tag:        "latest",
		Type:       "dockerhub",
		Name:       "cdh",
		Images:     []string{},
		Namespaces: []string{},
		Fail:       false,
		WhiteList:  []string{".*-apb$"},
		BlackList:  []string{},
		AuthType:   "",
		AuthName:   "",
		Runner:     "",
	}

	theurl, _ := url.Parse(c.URL)

	// would be better to use a conversion function
	adapterconfig := bladapters.Configuration{
		URL:        theurl,
		User:       c.User,
		Pass:       c.Pass,
		Org:        c.Org,
		Runner:     c.Runner,
		Images:     c.Images,
		Namespaces: c.Namespaces,
		Tag:        c.Tag,
	}

	dhadapter := bladapters.DockerHubAdapter{Config: adapterconfig}
	log.Info("Creating a custom dockerhub registry")
	reg, err = registries.NewCustomRegistry(c, dhadapter, "openshift")
	if err != nil {
		log.Errorf(
			"Failed to initialize dh Registry err - %v \n", err)
		os.Exit(1)
	}

	regs = append(regs, reg)

	// CreateApp passing in the args and registries
	app := app.CreateApp(args, regs)
	app.Start()
}
