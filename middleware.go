package otto

// Middleware defines the inteface for a otto Middleware
type Middleware func(HandlerFunc) HandlerFunc

type middlewareStack []Middleware

func (m *middlewareStack) Add(mm ...Middleware) {
	stack := append(middlewareStack{}, mm...)
	*m = append(stack, *m...)
}

func (m middlewareStack) Handle(r Route) HandlerFunc {
	h := r.HandlerFunc

	h = func(_ HandlerFunc) HandlerFunc {
		return h
	}(h)

	for _, mf := range m {
		h = mf(h)
	}

	return h
}

func (m middlewareStack) Copy() middlewareStack {
	c := middlewareStack{}
	return append(c, m...)
}
