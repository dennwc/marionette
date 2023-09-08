package marionette

type Context int

const (
	CHROME Context = 1 + iota
	CONTENT
)

var contexts = [...]string{
	"chrome",
	"content",
}

// String returns the value name of the context ("chrome", "content").
func (c Context) String() string {
	return contexts[c-1]
}
