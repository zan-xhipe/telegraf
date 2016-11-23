package timestampvalue

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseValidValues(t *testing.T) {
	parser := TimestampValueParser{
		MetricName: "timestamp_value_test",
		DataType:   "integer",
		Delimiter:  " ",
		TimeLayout: "unix",
	}
	metrics, err := parser.Parse([]byte("154516231 55"))
	assert.NoError(t, err)
	assert.Len(t, metrics, 1)
	assert.Equal(t, "timestamp_value_test", metrics[0].Name())
	assert.Equal(t, map[string]interface{}{
		"value": int64(55),
	}, metrics[0].Fields())
	assert.Equal(t, map[string]string{}, metrics[0].Tags())
	assert.Equal(t, time.Unix(0, 154516231), metrics[0].Time())

	// since the value parsing is the same as the normal Value parser it is not
	// nessecary to test all the parsable types as the Value parser already tests that
}

func TestParseMultipleValues(t *testing.T) {
	parser := TimestampValueParser{
		MetricName: "value_test",
		DataType:   "integer",
		Delimiter:  " ",
		TimeLayout: "unix",
	}
	metrics, err := parser.Parse([]byte(`564984312 55
45
223
12
999
`))
	assert.NoError(t, err)
	assert.Len(t, metrics, 1)
	assert.Equal(t, "value_test", metrics[0].Name())
	assert.Equal(t, map[string]interface{}{
		"value": int64(999),
	}, metrics[0].Fields())
	assert.Equal(t, map[string]string{}, metrics[0].Tags())
	assert.Equal(t, time.Unix(0, 564984312), metrics[0].Time())
}

func TestParseLineValues(t *testing.T) {
	parser := TimestampValueParser{
		MetricName: "timestamp_value_test",
		DataType:   "integer",
		Delimiter:  " ",
		TimeLayout: "unix",
	}
	metric, err := parser.ParseLine("154516231 55")
	assert.NoError(t, err)
	assert.Equal(t, "timestamp_value_test", metric.Name())
	assert.Equal(t, map[string]interface{}{
		"value": int64(55),
	}, metric.Fields())
	assert.Equal(t, map[string]string{}, metric.Tags())
	assert.Equal(t, time.Unix(0, 154516231), metric.Time())

	// since the value parsing is the same as the normal Value parser it is not
	// nessecary to test all the parsable types as the Value parser already tests that
}

func TestParseInvalidTimestamps(t *testing.T) {
	parser := TimestampValueParser{
		MetricName: "timestamp_value_test",
		DataType:   "integer",
		Delimiter:  " ",
		TimeLayout: "unix",
	}
	metrics, err := parser.Parse([]byte("55"))
	assert.Error(t, err)
	assert.Len(t, metrics, 0)

	metrics, err = parser.Parse([]byte("1556189 55.0"))
	assert.Error(t, err)
	assert.Len(t, metrics, 0)

	metrics, err = parser.Parse([]byte("2006-01-02T15:04:05Z07:00 55"))
	assert.Error(t, err)
	assert.Len(t, metrics, 0)

	metrics, err = parser.Parse([]byte("514320136+55"))
	assert.Error(t, err)
	assert.Len(t, metrics, 0)
}

func TestParseLineInvalidTimestamps(t *testing.T) {
	parser := TimestampValueParser{
		MetricName: "timestamp_value_test",
		DataType:   "integer",
		Delimiter:  " ",
		TimeLayout: "unix",
	}
	_, err := parser.ParseLine("55")
	assert.Error(t, err)

	_, err = parser.ParseLine("1556189 55.0")
	assert.Error(t, err)

	_, err = parser.ParseLine("2006-01-02T15:04:05Z07:00 55")
	assert.Error(t, err)

	_, err = parser.ParseLine("514320136+55")
	assert.Error(t, err)
}

func TestParseNonDefaultTimeLayouts(t *testing.T) {
	t1, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05+07:00")
	assert.NoError(t, err)

	parser := TimestampValueParser{
		MetricName: "timestamp_value_test",
		DataType:   "integer",
		Delimiter:  " ",
		TimeLayout: time.RFC3339,
	}
	metrics, err := parser.Parse([]byte("2006-01-02T15:04:05+07:00 55"))
	assert.NoError(t, err)
	assert.Len(t, metrics, 1)
	assert.Equal(t, "timestamp_value_test", metrics[0].Name())
	assert.Equal(t, map[string]interface{}{
		"value": int64(55),
	}, metrics[0].Fields())
	assert.Equal(t, map[string]string{}, metrics[0].Tags())
	assert.Equal(t, t1, metrics[0].Time())

	t2, err := time.Parse("2006/01/02,15:04:05Z07:00", "2006/01/02,15:04:05-07:00")
	assert.NoError(t, err)

	parser = TimestampValueParser{
		MetricName: "timestamp_value_test",
		DataType:   "integer",
		Delimiter:  " ",
		TimeLayout: "2006/01/02,15:04:05-07:00",
	}
	metrics, err = parser.Parse([]byte("2006/01/02,15:04:05-07:00 55"))
	assert.NoError(t, err)
	assert.Len(t, metrics, 1)
	assert.Equal(t, "timestamp_value_test", metrics[0].Name())
	assert.Equal(t, map[string]interface{}{
		"value": int64(55),
	}, metrics[0].Fields())
	assert.Equal(t, map[string]string{}, metrics[0].Tags())
	assert.Equal(t, t2, metrics[0].Time())
}

func TestNewParser(t *testing.T) {
	t1, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05+07:00")
	assert.NoError(t, err)

	parser := New(
		"timestamp_value_test",
		"integer",
		"",
		"rfc3339",
	)
	metrics, err := parser.Parse([]byte("2006-01-02T15:04:05+07:00 55"))
	assert.NoError(t, err)
	assert.Len(t, metrics, 1)
	assert.Equal(t, "timestamp_value_test", metrics[0].Name())
	assert.Equal(t, map[string]interface{}{
		"value": int64(55),
	}, metrics[0].Fields())
	assert.Equal(t, map[string]string{}, metrics[0].Tags())
	assert.Equal(t, t1, metrics[0].Time())

	t2, err := time.Parse("2006/01/02,15:04:05Z07:00", "2006/01/02,15:04:05-07:00")
	assert.NoError(t, err)

	parser = New(
		"timestamp_value_test",
		"integer",
		" ",
		"2006/01/02,15:04:05-07:00",
	)
	metrics, err = parser.Parse([]byte("2006/01/02,15:04:05-07:00 55"))
	assert.NoError(t, err)
	assert.Len(t, metrics, 1)
	assert.Equal(t, "timestamp_value_test", metrics[0].Name())
	assert.Equal(t, map[string]interface{}{
		"value": int64(55),
	}, metrics[0].Fields())
	assert.Equal(t, map[string]string{}, metrics[0].Tags())
	assert.Equal(t, t2, metrics[0].Time())

	parser = New(
		"timestamp_value_test",
		"integer",
		"",
		"",
	)
	metrics, err = parser.Parse([]byte("165498161 55"))
	assert.NoError(t, err)
	assert.Len(t, metrics, 1)
	assert.Equal(t, "timestamp_value_test", metrics[0].Name())
	assert.Equal(t, map[string]interface{}{
		"value": int64(55),
	}, metrics[0].Fields())
	assert.Equal(t, map[string]string{}, metrics[0].Tags())
	assert.Equal(t, time.Unix(0, 165498161), metrics[0].Time())
}
