package wikipedia

import "github.com/seppo0010/wikipedia-go"

func NewWikipedia() *wikipedia.Wikipedia {
    return wikipedia.NewWikipedia()
}

type Page interface {
    Title() (string, error)
    Content() (string, error)
}
