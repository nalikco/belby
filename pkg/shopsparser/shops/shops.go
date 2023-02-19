package shops

const (
	DefaultUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"
)

type Shop interface {
	GetTitle() string
	Find(query string) (Product, error)
}
