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
	"regexp"

	"go.opencensus.io/tag"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/model/pdata"
	"go.uber.org/zap"

	"context"
)

var _ component.TracesProcessor = (*redaction)(nil)

type redaction struct {
	// Attribute keys allowed in a span
	allowList map[string]string
	// Attribute values blocked in a span
	blockRegexList map[string]*regexp.Regexp
	// Redaction processor configuration
	config *Config
	// Attributes which would never be truncated
	limitExceptionList map[string]string
	// Logger
	logger *zap.Logger
	// Attribute keys persisted as metric dimensions or tags
	metricTagKeys map[string]tag.Key
	// Next trace consumer in line
	next consumer.Traces
}

const (
	ErrBlockListProcessing = "error_block_list_processing"
	ErrTraceProcessing     = "error_trace_processing"
	ErrTraceConsume        = "error_trace_consume"
	ErrMetricsRecording    = "error_metrics_recording"
	ErrRegexCompilation    = "error_regex_compilation"
)

// newRedaction creates a new instance of the redaction processor
func newRedaction(ctx context.Context, config *Config, logger *zap.Logger, next consumer.Traces) (*redaction, error) {
	if config.Summary == "" {
		config.Summary = Info
	}
	if config.DryRun {
		logger.Info("Redaction processor is configured for dry run mode.")
	}

	allowList := makeAllowList(config)
	blockRegexList, err := makeBlockRegexList(ctx, config, logger)
	if err != nil {
		logger.Error("Error processing block list", zap.String("error", err.Error()))
		return nil, err
	}
	limitExceptionList := makeTruncateExceptionsMap(config)
	metricTagKeys := map[string]tag.Key{}

	return &redaction{
		allowList:          allowList,
		blockRegexList:     blockRegexList,
		config:             config,
		limitExceptionList: limitExceptionList,
		logger:             logger,
		metricTagKeys:      metricTagKeys,
		next:               next,
	}, nil
}

// ProcessTraces is a helper function that processes the incoming data and returns the data to be sent to the next component.
// If error is returned then returned data are ignored. It MUST not call the next component.
func (s *redaction) ProcessTraces(_ context.Context, batch pdata.Traces) (pdata.Traces, error) {
	// TODO: Implementation to follow in the next PR
	return batch, nil
}

// ConsumeTraces implements the SpanProcessor interface
func (s *redaction) ConsumeTraces(_ context.Context, _ pdata.Traces) error {
	// TODO: Implementation to follow in the next PR
	return nil
}

const (
	Debug               = "debug"
	Info                = "info"
	Silent              = "silent"
	RedactedKeys        = "redacted_keys"
	RedactedKeyCount    = "redacted_key_count"
	MaskedValues        = "masked_values"
	MaskedValueCount    = "masked_value_count"
	TruncatedValues     = "truncated_values"
	TruncatedValueCount = "truncated_value_count"
)

// makeAllowList sets up a lookup table of allowed span attribute keys
func makeAllowList(c *Config) map[string]string {
	redactionKeys := []string{RedactedKeys, RedactedKeyCount, MaskedValues, MaskedValueCount, TruncatedValues, TruncatedValueCount}
	allowList := make(map[string]string, len(c.AllowedKeys)+len(redactionKeys))
	for _, key := range c.AllowedKeys {
		allowList[key] = key
	}
	for _, key := range redactionKeys {
		allowList[key] = key
	}
	return allowList
}

// makeBlockRegexList precompiles all the blocked regex patterns
func makeBlockRegexList(_ context.Context, config *Config, logger *zap.Logger) (map[string]*regexp.Regexp, error) {
	blockRegexList := make(map[string]*regexp.Regexp, len(config.BlockedValues))
	for _, pattern := range config.BlockedValues {
		re, err := regexp.Compile(pattern)
		if err != nil {
			logger.Error("Error compiling regex in block list", zap.String("error", err.Error()))
			return nil, err
		}
		blockRegexList[pattern] = re
	}
	return blockRegexList, nil
}

// makeTruncateExceptionsMap creates a lookup table of attributes to not truncate
func makeTruncateExceptionsMap(config *Config) map[string]string {
	truncateExceptions := make(map[string]string, len(config.Limits.LimitExceptions))
	for _, ex := range config.Limits.LimitExceptions {
		truncateExceptions[ex] = ex
	}
	return truncateExceptions
}

// Capabilities specifies what this processor does, such as whether it mutates data
func (s *redaction) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}

// Start the redaction processor
func (s *redaction) Start(_ context.Context, _ component.Host) error {
	s.logger.Info("Starting redaction processor")
	return nil
}

// Shutdown the redaction processor
func (s *redaction) Shutdown(context.Context) error {
	s.logger.Info("Shutting down redaction processor")
	return nil
}
