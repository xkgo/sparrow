package sparrow

type Option func(app *Application)

type Runner func(app *Application) (err error)

/**
设置 APP 应用名称，默认会读取环境中的 sparrow.application.name
*/
func WithName(appName string) Option {
	return func(app *Application) {
		app.Name = appName
	}
}

/**
当运行环境、IOC注入都已经好了之后
*/
func WithRunner(runner Runner) Option {
	return func(app *Application) {
		app.runner = runner
	}
}

/**
程序退出时候执行销毁处理
*/
func WithDestroyer(destroyer Runner) Option {
	return func(app *Application) {
		app.destroyer = destroyer
	}
}

func WithBeforeInitHandler(handlers ...Runner) Option {
	return func(app *Application) {
		if len(handlers) > 0 {
			if nil == app.beforeInitHandlers {
				app.beforeInitHandlers = handlers
			}
			app.beforeInitHandlers = append(app.beforeInitHandlers, handlers...)
		}
	}
}
