package utils

import (
	"fmt"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/signalfx/golib/datapoint"
)

func sortedDimensionString(dims map[string]string) string {
	var keys []string
	for k := range dims {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	table.SetHeader(keys)
	var vals []string
	for _, k := range keys {
		vals = append(vals, dims[k])
	}
	table.Append(vals)
	table.SetAutoFormatHeaders(false)

	table.Render()
	return tableString.String()
}

func dpTypeToString(t datapoint.MetricType) string {
	switch t {
	case datapoint.Gauge:
		return "gauge"
	case datapoint.Count:
		return "counter"
	case datapoint.Counter:
		return "cumulative counter"
	default:
		return fmt.Sprintf("unsupported type %d", t)
	}
}

// DatapointToString pretty prints a datapoint in a consistent manner for
// logging purposes.  The most important thing here is to sort the dimension
// dict so it is consistent so that it is easier to visually scan a large list
// of datapoints.
func DatapointToString(dp *datapoint.Datapoint) string {
	var tsStr string
	if !dp.Timestamp.IsZero() {
		tsStr = dp.Timestamp.String()
	}
	return fmt.Sprintf("%s: %s (%s) %s\n%s\n", dp.Metric, dp.Value, strings.ToUpper(dpTypeToString(dp.MetricType)), tsStr, sortedDimensionString(dp.Dimensions))
}

// BoolToInt returns 1 if b is true and 0 otherwise.  It is useful for
// datapoints which track a binary value since we don't support boolean
// datapoint directly.
func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// TruncateDimensionValuesInPlace restricts the value of dimensions to 256
// characters.  If the dim exceeds this limit, the value is truncated.
func TruncateDimensionValuesInPlace(dims map[string]string) {
	for k, v := range dims {
		// Not sure if our backend enforces character length or byte length.
		// If values include multi-byte unicode chars, this might not work.
		if len(v) > 256 {
			dims[k] = v[:256]
		}
	}
}
