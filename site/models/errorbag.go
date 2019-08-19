package models

type ErrorBag map[string]string

func (e *ErrorBag) Merge(b ErrorBag) {
	for k, v := range b {
		(*e)[k] = v
	}
}

func (e ErrorBag) Passes() bool {
	return len(e) == 0
}

func (e ErrorBag) Fails() bool {
	return !e.Passes()
}

func (e ErrorBag) HasErrors(field string) bool {
	_, ok := e[field]
	return ok
}

func (e ErrorBag) AsMap() map[string]string {
	return map[string]string(e)
}
