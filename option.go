package goweb

type Option func(*App)

func WithValidator(validate func(any) error) Option {
	return func(a *App) {
		a.validate = validate
	}
}

func WithErrorHandler(handle func(*Context, error) Response) Option {
	return func(a *App) {
		a.handleError = handle
	}
}
