package timestampvalue

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/parsers/value"
)

type TimestampValueParser struct {
	MetricName  string
	DataType    string
	Delimiter   string
	TimeLayout  string
	DefaultTags map[string]string
}

func New(name, dataType, delimiter, layout string) TimestampValueParser {
	delim := delimiter
	if delimiter == "" {
		delim = " "
	}

	return TimestampValueParser{
		MetricName: name,
		DataType:   dataType,
		Delimiter:  delim,
		TimeLayout: selectTimeLayout(layout),
	}
}

func (v *TimestampValueParser) Parse(buf []byte) ([]telegraf.Metric, error) {
	parts := bytes.Split(buf, []byte(v.Delimiter))
	if len(parts) < 2 {
		return nil, fmt.Errorf("timestamp-values must both a timestamp and a value")
	}

	var timestamp time.Time
	var err error
	ts := string(parts[0])
	if v.TimeLayout == "unix" {
		i, err := strconv.ParseInt(ts, 10, 64)
		if err != nil {
			return nil, err
		}

		timestamp = time.Unix(0, i)

	} else {
		timestamp, err = time.Parse(v.TimeLayout, ts)
		if err != nil {
			return nil, err
		}
	}

	value, err := value.Parse(v.DataType, parts[1])
	if err != nil {
		return nil, err
	}

	fields := map[string]interface{}{"value": value}
	metric, err := telegraf.NewMetric(v.MetricName, v.DefaultTags,
		fields, timestamp)
	if err != nil {
		return nil, err
	}

	return []telegraf.Metric{metric}, nil
}

func (v *TimestampValueParser) ParseLine(line string) (telegraf.Metric, error) {
	metrics, err := v.Parse([]byte(line))

	if err != nil {
		return nil, err
	}

	if len(metrics) < 1 {
		return nil, fmt.Errorf("Can not parse the line: %s, for data format: value", line)
	}

	return metrics[0], nil
}

func (v *TimestampValueParser) SetDefaultTags(tags map[string]string) {
	v.DefaultTags = tags
}

func selectTimeLayout(layout string) string {
	layout = strings.ToLower(layout)
	switch layout {
	case "ansic":
		return time.ANSIC
	case "unixdate":
		return time.UnixDate
	case "rubydate":
		return time.RubyDate
	case "rfc822":
		return time.RFC822
	case "rfc822z":
		return time.RFC822Z
	case "rfc850":
		return time.RFC850
	case "rfc1123":
		return time.RFC1123
	case "rfc1123z":
		return time.RFC1123Z
	case "rfc3339":
		return time.RFC3339
	case "rfc3339nano":
		return time.RFC3339Nano
	case "kitchen":
		return time.Kitchen
	case "stamp":
		return time.Stamp
	case "stampmilli":
		return time.StampMilli
	case "stampmicro":
		return time.StampMicro
	case "stampnano":
		return time.StampNano
	case "":
		return "unix"
	default:
		return layout
	}
}
