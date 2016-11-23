package timestampvalue

import (
	"bytes"
	"fmt"
	"strconv"
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

func (v *TimestampValueParser) Parse(buf []byte) ([]telegraf.Metric, error) {
	delim := []byte(v.Delimiter)
	// set default delimiter
	if v.Delimiter == "" {
		delim = []byte{' '}
	}
	parts := bytes.Split(buf, delim)
	if len(parts) < 2 {
		return nil, fmt.Errorf("timestamp-values must both a timestamp and a value")
	}

	var timestamp time.Time
	var err error
	ts := string(parts[0])
	if v.TimeLayout == "" {
		// parse nanosecond unix timestamp
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
