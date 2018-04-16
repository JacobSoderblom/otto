package otto

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Bind_DefaultBinder_JSON(t *testing.T) {

	b, err := json.Marshal(map[string]interface{}{
		"a": "b",
		"c": 2,
	})

	assert.NoError(t, err, "should not cast error on json marshal")

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Set(HeaderContentType, MIMEApplicationJSON)

	c := &context{
		req:      req,
		bindFunc: DefaultBinder,
	}

	var body struct {
		A string `json:"a"`
		C int    `json:"c"`
	}

	err = c.Bind(&body)
	assert.NoError(t, err, "should not cast error on Bind json")

	assert.Equal(t, "b", body.A)
	assert.Equal(t, 2, body.C)
}

func Test_Bind_DefaultBinder_No_Support(t *testing.T) {

	b, err := json.Marshal(map[string]interface{}{
		"a": "b",
		"c": 2,
	})

	assert.NoError(t, err, "should not cast error on json marshal")

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Set(HeaderContentType, "some content type")

	c := &context{
		req:      req,
		bindFunc: DefaultBinder,
	}

	var body struct {
		A string `json:"a"`
		C int    `json:"c"`
	}

	err = c.Bind(&body)
	assert.Contains(t, err.Error(), "No support for content type")
}

func Test_Bind_DefaultBinder_GET_DELETE(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)

	c := &context{
		req:      req,
		bindFunc: DefaultBinder,
	}

	err := c.Bind(nil)
	assert.Contains(t, err.Error(), "Bind is not supported for GET method")

	req = httptest.NewRequest("DELETE", "/", nil)

	c = &context{
		req:      req,
		bindFunc: DefaultBinder,
	}

	err = c.Bind(nil)
	assert.Contains(t, err.Error(), "Bind is not supported for DELETE method")
}

func Test_Bind_DefaultBinder_Content_length_Error(t *testing.T) {
	req := httptest.NewRequest("POST", "/", nil)

	c := &context{
		req:      req,
		bindFunc: DefaultBinder,
	}

	err := c.Bind(nil)
	assert.Contains(t, err.Error(), "Request body cannot be empty")
}

func Test_Bind_DefaultBinder_JSON_UnmarshalTypeError(t *testing.T) {

	b, err := json.Marshal(map[string]interface{}{
		"a": "b",
		"c": 2,
	})

	assert.NoError(t, err, "should not cast error on json marshal")

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Set(HeaderContentType, MIMEApplicationJSON)

	c := &context{
		req:      req,
		bindFunc: DefaultBinder,
	}

	var body struct {
		A int `json:"a"`
		C int `json:"c"`
	}

	err = c.Bind(&body)
	assert.Error(t, err, "should cast unmarshal type error")
	assert.Contains(t, err.Error(), "Unmarshal type error:")
}

func Test_Bind_DefaultBinder_JSON_SyntaxError(t *testing.T) {

	data := "{ a: b}"

	req := httptest.NewRequest("POST", "/", bytes.NewReader([]byte(data)))
	req.Header.Set(HeaderContentType, MIMEApplicationJSON)

	c := &context{
		req:      req,
		bindFunc: DefaultBinder,
	}

	var body struct {
		A int `json:"a"`
		C int `json:"c"`
	}

	err := c.Bind(&body)
	assert.Error(t, err, "should cast syntax error")
	assert.Contains(t, err.Error(), "Syntax error:")
}

func Test_Bind_DefaultBinder_JSON_Error(t *testing.T) {

	data := "{"

	req := httptest.NewRequest("POST", "/", bytes.NewReader([]byte(data)))
	req.Header.Set(HeaderContentType, MIMEApplicationJSON)

	c := &context{
		req:      req,
		bindFunc: DefaultBinder,
	}

	var body struct {
		A int `json:"a"`
		C int `json:"c"`
	}

	err := c.Bind(&body)
	assert.Error(t, err, "should cast error")
	assert.Contains(t, err.Error(), "Could not decode json")
}
