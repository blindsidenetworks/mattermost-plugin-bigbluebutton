/*
Copyright 2018 Blindside Networks

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"github.com/pkg/errors"
	"strings"
)

type Configuration struct {
	BaseURL            string `json:"BASE_URL"`
	Secret             string `json:"SALT"`
	AdminOnly          bool   `json:"ADMINONLY"`
	ProcessRecordings  bool   `json:"PROCESS_RECORDINGS"`
	AllowRecordings    bool   `json:"ALLOW_RECORDINGS"`
	AllowExternalUsers bool   `json:"ALLOW_EXTERNAL_USERS"`
}

func (p *Plugin) OnConfigurationChange() error {
	var newConfig Configuration
	// loads configuration from our config ui page
	err := p.API.LoadPluginConfiguration(&newConfig)

	newConfig.BaseURL = strings.Trim(newConfig.BaseURL, "/")
	newConfig.BaseURL = strings.Trim(newConfig.BaseURL, " ")

	oldConfig := p.config()

	// close running job if process recording is turned off
	if oldConfig != nil && (oldConfig.ProcessRecordings && !newConfig.ProcessRecordings) && p.job != nil {
		p.job.Close()
	}

	// stores the config in an Atomic.Value place
	p.configuration.Store(&newConfig)
	return err
}

func (p *Plugin) config() *Configuration {
	// returns the config file we had stored in Atomic.Value
	config := p.configuration.Load()
	if config == nil {
		return nil
	}
	return config.(*Configuration)
}

func (c *Configuration) IsValid() error {
	if len(c.BaseURL) == 0 {
		return errors.New("BASE URL is not configured.")
	} else if len(c.Secret) == 0 {
		return errors.New("Secret is not configured.")
	}

	return nil
}
