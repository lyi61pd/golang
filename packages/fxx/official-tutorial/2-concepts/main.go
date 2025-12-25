package fx

type App
  func New(opts ...Option) *App
  func (app *App) Run()

type Option
  func Provide(constructors ...interface{}) Option
  func Invoke(funcs ...interface{}) Option