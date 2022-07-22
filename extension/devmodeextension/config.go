// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package devmodeextension // import "github.com/open-telemetry/opentelemetry-collector-contrib/extension/devmodeextension"

import (
	"errors"

	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/confignet"
)

// Config stores the configuration for OAuth2 Client Credentials (2-legged OAuth2 flow) setup.
type Config struct {
	config.ExtensionSettings `mapstructure:",squash"`
	TCPAddr                  confignet.TCPAddr `mapstructure:",squash"`
}

var _ config.Extension = (*Config)(nil)

// Validate checks if the extension configuration is valid
func (cfg *Config) Validate() error {
	if cfg.TCPAddr.Endpoint == "" {
		return errors.New("\"endpoint\" is required when using the \"devmode\" extension")
	}
	return nil
}
