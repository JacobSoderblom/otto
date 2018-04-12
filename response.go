package otto

import "net/http"

// Response that holds some information about the response
type Response struct {
	http.ResponseWriter
	size int
	code int
}

func (r *Response) Write(b []byte) (int, error) {
	s, err := r.ResponseWriter.Write(b)
	r.size = s
	return s, err
}

// WriteHeader writes the code to response writer
// and stores the status code
func (r *Response) WriteHeader(code int) {
	r.code = code
	r.ResponseWriter.WriteHeader(code)
}

// Size returns the size of the content
func (r Response) Size() int {
	return r.size
}
