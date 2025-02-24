package config

import (
	"reflect"

	"github.com/mitchellh/hashstructure"
	"github.com/signalfx/signalfx-agent/internal/core/dpfilters"
	"github.com/signalfx/signalfx-agent/internal/monitors/types"
	log "github.com/sirupsen/logrus"
)

// MonitorConfig is used to configure monitor instances.  One instance of
// MonitorConfig may be used to configure multiple monitor instances.  If a
// monitor's discovery rule does not match any discovered services, the monitor
// will not run.
type MonitorConfig struct {
	// The type of the monitor
	Type string `yaml:"type" json:"type"`
	// The rule used to match up this configuration with a discovered endpoint.
	// If blank, the configuration will be run immediately when the agent is
	// started.  If multiple endpoints match this rule, multiple instances of
	// the monitor type will be created with the same configuration (except
	// different host/port).
	DiscoveryRule string `yaml:"discoveryRule" json:"discoveryRule"`
	// If true, a warning will be emitted if a discovery rule contains
	// variables that will never possibly match a rule.  If using multiple
	// observers, it is convenient to set this to false to suppress spurious
	// errors.  The top-level setting `validateDiscoveryRules` acts as a
	// default if this isn't set.
	ValidateDiscoveryRule *bool `yaml:"validateDiscoveryRule"`
	// A set of extra dimensions (key:value pairs) to include on datapoints emitted by the
	// monitor(s) created from this configuration. To specify metrics from this
	// monitor should be high-resolution, add the dimension `sf_hires: 1`
	ExtraDimensions map[string]string `yaml:"extraDimensions" json:"extraDimensions"`
	// A mapping of extra dimension names to a [discovery rule
	// expression](https://docs.signalfx.com/en/latest/integrations/agent/auto-discovery.html)
	// that is used to derive the value of the dimension.  For example, to use
	// a certain container label as a dimension, you could use something like this
	// in your monitor config block: `extraDimensionsFromEndpoint: {env: 'Get(container_labels, "myapp.com/environment")'}`
	ExtraDimensionsFromEndpoint map[string]string `yaml:"extraDimensionsFromEndpoint" json:"extraDimensionsFromEndpoint"`
	// A set of mappings from a configuration option on this monitor to
	// attributes of a discovered endpoint.  The keys are the config option on
	// this monitor and the value can be any valid expression used in discovery
	// rules.
	ConfigEndpointMappings map[string]string `yaml:"configEndpointMappings" json:"configEndpointMappings"`
	// The interval (in seconds) at which to emit datapoints from the
	// monitor(s) created by this configuration.  If not set (or set to 0), the
	// global agent intervalSeconds config option will be used instead.
	IntervalSeconds int `yaml:"intervalSeconds" json:"intervalSeconds"`
	// If one or more configurations have this set to true, only those
	// configurations will be considered. This setting can be useful for testing.
	Solo bool `yaml:"solo" json:"solo"`
	// DEPRECATED in favor of the `datapointsToExclude` option.  That option
	// handles negation of filter items differently.
	MetricsToExclude []MetricFilter `yaml:"metricsToExclude" json:"metricsToExclude" default:"[]"`
	// A list of datapoint filters.  These filters allow you to comprehensively
	// define which datapoints to exclude by metric name or dimension set, as
	// well as the ability to define overrides to re-include metrics excluded
	// by previous patterns within the same filter item.  See [monitor
	// filtering](https://github.com/signalfx/signalfx-agent/tree/master/docs/filtering.md#monitor-level-filtering)
	// for examples and more information.
	DatapointsToExclude []MetricFilter `yaml:"datapointsToExclude" json:"datapointsToExclude" default:"[]"`
	// Some monitors pull metrics from services not running on the same host
	// and should not get the host-specific dimensions set on them (e.g.
	// `host`, `AWSUniqueId`, etc).  Setting this to `true` causes those
	// dimensions to be omitted.  You can disable this globally with the
	// `disableHostDimensions` option on the top level of the config.
	DisableHostDimensions bool `yaml:"disableHostDimensions" json:"disableHostDimensions" default:"false"`
	// This can be set to true if you don't want to include the dimensions that
	// are specific to the endpoint that was discovered by an observer.  This
	// is useful when you have an endpoint whose identity is not particularly
	// important since it acts largely as a proxy or adapter for other metrics.
	DisableEndpointDimensions bool `yaml:"disableEndpointDimensions" json:"disableEndpointDimensions"`
	// A map from dimension names emitted by the monitor to the desired
	// dimension name that will be emitted in the datapoint that goes to
	// SignalFx.  This can be useful if you have custom metrics from your
	// applications and want to make the dimensions from a monitor match those.
	// Also can be useful when scraping free-form metrics, say with the
	// `prometheus-exporter` monitor.  Right now, only static key/value
	// transformations are supported.  Note that filtering by dimensions will
	// be done on the *original* dimension name and not the new name. Note that
	// it is possible to remove unwanted dimensions via this configuration, by
	// making the desired dimension name an empty string.
	DimensionTransformations map[string]string `yaml:"dimensionTransformations" json:"dimensionTransformations"`
	// Extra metrics to enable besides the default included ones.  This is an
	// [overridable filter](https://docs.signalfx.com/en/latest/integrations/agent/filtering.html#overridable-filtering).
	ExtraMetrics []string `yaml:"extraMetrics" json:"extraMetrics"`
	// Extra metric groups to enable in addition to the metrics that are
	// emitted by default.  A metric group is simply a collection of metrics,
	// and they are defined in each monitor's documentation.
	ExtraGroups []string `yaml:"extraGroups" json:"extraGroups"`
	// OtherConfig is everything else that is custom to a particular monitor
	OtherConfig map[string]interface{} `yaml:",inline" neverLog:"omit"`
	// ValidationError is where a message concerning validation issues can go
	// so that diagnostics can output it.
	Hostname        string          `yaml:"-" json:"-"`
	BundleDir       string          `yaml:"-" json:"-"`
	ValidationError string          `yaml:"-" json:"-" hash:"ignore"`
	MonitorID       types.MonitorID `yaml:"-" hash:"ignore"`
}

var _ CustomConfigurable = &MonitorConfig{}

// Validate ensures the config is correct beyond what basic YAML parsing
// ensures
func (mc *MonitorConfig) Validate() error {
	var err error
	if _, err = mc.OldFilterSet(); err != nil {
		return err
	}
	if _, err = mc.NewFilterSet(); err != nil {
		return err
	}
	return nil
}

// OldFilterSet makes a filter set using the old filter style
func (mc *MonitorConfig) OldFilterSet() (*dpfilters.FilterSet, error) {
	return makeOldFilterSet(mc.MetricsToExclude, nil)
}

// NewFilterSet makes a filter set using the new filter style
func (mc *MonitorConfig) NewFilterSet() (*dpfilters.FilterSet, error) {
	return makeNewFilterSet(mc.DatapointsToExclude)
}

// Equals tests if two monitor configs are sufficiently equal to each other.
// Two monitors should only be equal if it doesn't make sense for two
// configurations to be active at the same time.
func (mc *MonitorConfig) Equals(other *MonitorConfig) bool {
	return mc.Type == other.Type && mc.DiscoveryRule == other.DiscoveryRule &&
		reflect.DeepEqual(mc.OtherConfig, other.OtherConfig)
}

// ExtraConfig returns generic config as a map
func (mc *MonitorConfig) ExtraConfig() (map[string]interface{}, error) {
	return mc.OtherConfig, nil
}

// HasAutoDiscovery returns whether the monitor is static (i.e. doesn't rely on
// autodiscovered services and is manually configured) or dynamic.
func (mc *MonitorConfig) HasAutoDiscovery() bool {
	return mc.DiscoveryRule != ""
}

// ShouldValidateDiscoveryRule return ValidateDiscoveryRule or false if that is
// nil.
func (mc *MonitorConfig) ShouldValidateDiscoveryRule() bool {
	if mc.ValidateDiscoveryRule == nil || !*mc.ValidateDiscoveryRule {
		return false
	}
	return true
}

// MonitorConfigCore provides a way of getting the MonitorConfig when embedded
// in a struct that is referenced through a more generic interface.
func (mc *MonitorConfig) MonitorConfigCore() *MonitorConfig {
	return mc
}

// Hash calculates a unique hash value for this config struct
func (mc *MonitorConfig) Hash() uint64 {
	hash, err := hashstructure.Hash(mc, nil)
	if err != nil {
		log.WithError(err).Error("Could not get hash of MonitorConfig struct")
		return 0
	}
	return hash
}

// MonitorCustomConfig represents monitor-specific configuration that doesn't
// appear in the MonitorConfig struct.
type MonitorCustomConfig interface {
	MonitorConfigCore() *MonitorConfig
}

// ExtraMetrics interface for monitors that support generating additional metrics to allow through.
type ExtraMetrics interface {
	GetExtraMetrics() []string
}
