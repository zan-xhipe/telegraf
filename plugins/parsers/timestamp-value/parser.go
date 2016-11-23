package timestampvalue

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/influxdata/telegraf"
)

type TimestampValueParser struct {
	MetricName  string
	DataType    string
	Delimiter   string
	TimeLayout  string
	DefaultTags map[string]string
}

func (v *TimestampValueParser) Parse(buf []byte) ([]telegraf.Metric, error) {
	vStr := string(bytes.TrimSpace(bytes.Trim(buf, "\x00")))
	parts := strings.Split(vStr, v.Delimiter)

	timestamp, err := time.Parse(v.TimeLayout, parts[0])
	if err != nil {
		return nil, err
	}

	vStr = parts[1]
	// unless it's a string, separate out any fields in the buffer,
	// ignore anything but the last.
	if v.DataType != "string" {
		values := strings.Fields(vStr)
		if len(values) < 1 {
			return []telegraf.Metric{}, nil
		}
		vStr = string(values[len(values)-1])
	}

	var value interface{}
	switch v.DataType {
	case "", "int", "integer":
		value, err = strconv.Atoi(vStr)
	case "float", "long":
		value, err = strconv.ParseFloat(vStr, 64)
	case "str", "string":
		value = vStr
	case "bool", "boolean":
		value, err = strconv.ParseBool(vStr)
	}
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
