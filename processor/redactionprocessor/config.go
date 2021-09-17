// Copyright  OpenTelemetry Authors
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

package redactionprocessor

import (
	"go.opentelemetry.io/collector/config"
)

type Config struct {
	config.ProcessorSettings `mapstructure:",squash"`

	AllowedKeys   []string `mapstructure:"allowed_keys"`
	BlockedValues []string `mapstructure:"blocked_values"`
	DryRun        bool     `mapstructure:"dry_run"`
	Limits        Limits   `mapstructure:"limits"`
	MetricTags    []string `mapstructure:"metric_tags"`
	Summary       string   `mapstructure:"summary"`
}

type Limits struct {
	MaxValueLength  int      `mapstructure:"max_value_length"`
	LimitExceptions []string `mapstructure:"limit_exceptions"`
}
