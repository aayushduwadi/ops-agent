// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/GoogleCloudPlatform/ops-agent/confgenerator"
)

var (
	service  = flag.String("service", "", "service to generate config for")
	outDir   = flag.String("out", os.Getenv("RUNTIME_DIRECTORY"), "directory to write configuration files to")
	input    = flag.String("in", "/etc/google-cloud-ops-agent/config.yaml", "path to the user specified agent config")
	logsDir  = flag.String("logs", "/var/log/google-cloud-ops-agent", "path to store agent logs")
	stateDir = flag.String("state", "/var/lib/google-cloud-ops-agent", "path to store agent state like buffers")
)

func main() {
	flag.Parse()
	if err := run(); err != nil {
		log.Fatalf("The agent config file is not valid. Detailed error: %s", err)
	}
}
func run() error {
	// TODO(lingshi) Move this to a shared place across Linux and Windows.
	confDebugFolder := filepath.Join(os.Getenv("RUNTIME_DIRECTORY"), "conf", "debug")
	if err := confgenerator.MergeConfFiles(*input, confDebugFolder, "linux"); err != nil {
		return err
	}
	if err := logConfigFiles(filepath.Join(confDebugFolder, "built-in-config.yaml"), filepath.Join(confDebugFolder, "merged-config.yaml")); err != nil {
		return err
	}
	return confgenerator.GenerateFiles(filepath.Join(confDebugFolder, "merged-config.yaml"), *service, *logsDir, *stateDir, *outDir)
}

// logConfigFiles logs the built-in and merged config files to STDOUT. These are then written by journald to var/log/syslog and so to
// Cloud Logging once the ops-agent is running.
func logConfigFiles(builtInConfigFile, mergedConfigFile string) error {
	builtInConfig, err := ioutil.ReadFile(builtInConfigFile)
	if err != nil {
		return err
	}

	mergedConfig, err := ioutil.ReadFile(mergedConfigFile)
	if err != nil {
		return err
	}

	log.Printf("Built-in config: %s", builtInConfig)
	log.Printf("Merged config: %s", mergedConfig)

	return nil
}