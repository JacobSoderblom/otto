package otto

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Test_Params_String(t *testing.T) {
	r := mux.NewRouter()

	r.HandleFunc("/{id}", func(_ http.ResponseWriter, req *http.Request) {
		p := Params(mux.Vars(req))

		assert.Equal(t, "1", p.String("id"))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/1", ts.URL), nil)
	assert.NoError(t, err, "should not throw any error")
	_, err = http.DefaultClient.Do(req)
	assert.NoError(t, err, "should not throw any error")
}

func Test_Params_Int(t *testing.T) {
	r := mux.NewRouter()

	r.HandleFunc("/{id}", func(_ http.ResponseWriter, req *http.Request) {
		p := Params(mux.Vars(req))

		id, err := p.Int("id")
		assert.NoError(t, err, "should not cast error")

		assert.Equal(t, 1, id)
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/1", ts.URL), nil)
	assert.NoError(t, err, "should not throw any error")
	_, err = http.DefaultClient.Do(req)
	assert.NoError(t, err, "should not throw any error")
}

func Test_Params_Int_Not_Int(t *testing.T) {
	r := mux.NewRouter()

	r.HandleFunc("/{id}", func(_ http.ResponseWriter, req *http.Request) {
		p := Params(mux.Vars(req))

		_, err := p.Int("id")
		assert.Error(t, err, "should cast error")

		assert.Contains(t, err.Error(), "failed to parse 'a' to int")
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/a", ts.URL), nil)
	assert.NoError(t, err, "should not throw any error")
	_, err = http.DefaultClient.Do(req)
	assert.NoError(t, err, "should not throw any error")
}

func Test_Params_Bool(t *testing.T) {
	r := mux.NewRouter()

	r.HandleFunc("/{id}", func(_ http.ResponseWriter, req *http.Request) {
		p := Params(mux.Vars(req))

		b, err := p.Bool("id")
		assert.NoError(t, err, "should not cast error")

		assert.Equal(t, true, b)
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/1", ts.URL), nil)
	assert.NoError(t, err, "should not throw any error")
	_, err = http.DefaultClient.Do(req)
	assert.NoError(t, err, "should not throw any error")
}

func Test_Params_Bool_Not_Bool(t *testing.T) {
	r := mux.NewRouter()

	r.HandleFunc("/{id}", func(_ http.ResponseWriter, req *http.Request) {
		p := Params(mux.Vars(req))

		_, err := p.Bool("id")
		assert.Error(t, err, "should cast error")

		assert.Contains(t, err.Error(), "failed to parse 'a' to bool")
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/a", ts.URL), nil)
	assert.NoError(t, err, "should not throw any error")
	_, err = http.DefaultClient.Do(req)
	assert.NoError(t, err, "should not throw any error")
}

func Test_ValueParams_String(t *testing.T) {
	p := ValueParams{
		vals: url.Values{
			"id": []string{"123"},
		},
	}

	assert.Equal(t, "123", p.String("id"))
}

func Test_ValueParams_Strings(t *testing.T) {
	p := ValueParams{
		vals: url.Values{
			"id": []string{"123", "321"},
		},
	}

	assert.Equal(t, []string{"123", "321"}, p.Strings("id"))
}

func Test_ValueParams_Int(t *testing.T) {
	p := ValueParams{
		vals: url.Values{
			"id": []string{"123"},
		},
	}

	id, err := p.Int("id")
	assert.NoError(t, err, "should not cast error")
	assert.Equal(t, 123, id)
}

func Test_ValueParams_Int_Not_Int(t *testing.T) {
	p := ValueParams{
		vals: url.Values{
			"id": []string{"a"},
		},
	}

	_, err := p.Int("id")
	assert.Error(t, err, "should cast error")
	assert.Contains(t, err.Error(), "failed to parse 'a' to int")
}

func Test_ValueParams_Ints(t *testing.T) {
	p := ValueParams{
		vals: url.Values{
			"id": []string{"123", "321"},
		},
	}

	IDs, err := p.Ints("id")
	assert.NoError(t, err, "should not cast error")
	assert.Equal(t, []int{123, 321}, IDs)
}

func Test_ValueParams_Ints_Not_Int(t *testing.T) {
	p := ValueParams{
		vals: url.Values{
			"id": []string{"a", "123"},
		},
	}

	_, err := p.Ints("id")
	assert.Error(t, err, "should cast error")
	assert.Contains(t, err.Error(), "failed to parse 'a' to int")
}

func Test_ValueParams_Bool(t *testing.T) {
	p := ValueParams{
		vals: url.Values{
			"id": []string{"true"},
		},
	}

	id, err := p.Bool("id")
	assert.NoError(t, err, "should not cast error")
	assert.Equal(t, true, id)
}

func Test_ValueParams_Bool_Not_Bool(t *testing.T) {
	p := ValueParams{
		vals: url.Values{
			"id": []string{"a"},
		},
	}

	_, err := p.Bool("id")
	assert.Error(t, err, "should cast error")
	assert.Contains(t, err.Error(), "failed to parse 'a' to bool")
}

func Test_ValueParams_Bools(t *testing.T) {
	p := ValueParams{
		vals: url.Values{
			"id": []string{"true", "false"},
		},
	}

	bools, err := p.Bools("id")
	assert.NoError(t, err, "should not cast error")
	assert.Equal(t, []bool{true, false}, bools)
}

func Test_ValueParams_Bools_Not_Bool(t *testing.T) {
	p := ValueParams{
		vals: url.Values{
			"id": []string{"a", "true"},
		},
	}

	_, err := p.Bools("id")
	assert.Error(t, err, "should cast error")
	assert.Contains(t, err.Error(), "failed to parse 'a' to bool")
}
