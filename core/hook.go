package core

type Hook struct {
	Trigger     string
	Branch      string
	Event       string
	Timestamp   int64
	Title       string
	Message     string
	Before      string
	After       string
	Ref         string
	Author      string
	AuthorEmail string
	Changes     []string
}
