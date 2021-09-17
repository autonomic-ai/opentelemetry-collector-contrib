# Redaction processor

Supported pipeline types: traces

This processor deletes span attributes that don't match an allowlist, masks
span attribute values that match a blocked value regex list, as well as
truncates span attribute values that exceed a configured maximum length limit.

Typical use-cases:

* Prevent sensitive fields from accidentally leaking into traces
* Guarantee compliance with legal, privacy, or security requirements

Please refer to [config.go](./config.go) for the config spec.

Examples:

```yaml
processors:
  redaction:
    # Allowlist for span attribute keys
    allowed_keys:
      - description
      - group
      - id
      - name
    # Blocklist for span attribute values
    blocked_values:
      - "4[0-9]{12}(?:[0-9]{3})?" ## Visa credit card number for testing purpose
      - "(5[1-5][0-9]{14})"
    # Summarize the redactions but do not do the redactions
    dry_run: true
    limits:
      # Max length for span attribute values
      max_value_length: 8
      # Don't apply the length limit to these attributes
      limit_exceptions:
        - exception_statement
    # Span attribute keys to use as metric dimensions
    metric_tags:
      - group
      - name
    # Verbosity: debug vs info vs silent
    summary: debug
```

## Configuration

Refer to [config.yaml](./testdata/config.yaml) for detailed examples on using the processor.

Leaving the `allowed_keys` property blank will remove all span attributes
unless the processor is in dry run mode.

## Metrics

The processor records the following metrics:

* `num_redacted_keys` represents the number of span attribute keys redacted by processor
* `num_masked_values` represents the number of span attribute values masked by processor
* `num_truncated_values` represents the number of span attribute values truncated by processor
