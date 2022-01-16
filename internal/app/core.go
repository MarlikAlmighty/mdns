package app

// Core application
type Core struct {
	Config Config `config:"-"`
	Logger Logger `logger:"-"`
}

// New application core initialization
func New(c Config, l Logger) *Core {
	return &Core{
		Config: c,
		Logger: l,
	}
}
