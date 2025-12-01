package http

// applyMiddlewares applies a chain of middlewares to a handler
func applyMiddlewares[T any](handler func(T) error, middlewares []any) func(T) error {
	result := handler
	// Apply middlewares in reverse order so they execute in the defined order
	for i := len(middlewares) - 1; i >= 0; i-- {
		if mw, ok := middlewares[i].(func(func(T) error) func(T) error); ok {
			result = mw(result)
		}
	}
	return result
}
