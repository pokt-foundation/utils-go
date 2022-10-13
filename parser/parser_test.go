package parser

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegerFromJSONPayload(t *testing.T) {
	c := require.New(t)

	number, err := IntegerFromJSONPayload(bytes.NewReader([]byte(`{"ajua":12}`)), "ajua")
	c.NoError(err)
	c.Equal(int64(12), number)

	number, err = IntegerFromJSONPayload(bytes.NewReader([]byte(`{"ajua":"12"}`)), "ajua")
	c.NoError(err)
	c.Equal(int64(12), number)

	number, err = IntegerFromJSONPayload(bytes.NewReader([]byte(`{"ajua":["12"]}`)), "ajua")
	c.EqualError(err, "error parsing field ajua: invalid type for payload: []interface {}")
	c.Empty(number)
}

func TestIntegerJSONString(t *testing.T) {
	c := require.New(t)

	number, err := IntegerJSONString(`{"ajua":12}`, "ajua")
	c.NoError(err)
	c.Equal(int64(12), number)

	number, err = IntegerJSONString(`{"ajua":"12"}`, "ajua")
	c.NoError(err)
	c.Equal(int64(12), number)

	number, err = IntegerJSONString(`{"ajua":["12"]}`, "ajua")
	c.EqualError(err, "error parsing field ajua: invalid type for payload: []interface {}")
	c.Empty(number)
}

func TestParseNestedIntegerJSONString(t *testing.T) {
	c := require.New(t)

	number, err := IntegerJSONString(`{"ajua": {"papolo": 12}}`, "ajua.papolo")
	c.NoError(err)
	c.Equal(int64(12), number)

	number, err = IntegerJSONString(`{"ajua": {"papolo": "12"}}`, "ajua.papolo")
	c.NoError(err)
	c.Equal(int64(12), number)

	number, err = IntegerJSONString(`{"ajua": {"papolo": [12]}}`, "ajua.papolo")
	c.EqualError(err, "error parsing field ajua.papolo: invalid type for payload: []interface {}")
	c.Empty(number)
}
