package jutils

var (
	BasicKeyring = MakeContextKeyring()
)

type ContextKey string
type ContextKeyring map[string]ContextKey

func MakeContextKeyring() ContextKeyring {
	ring := ContextKeyring{}
	ring.AddKeyToRing("ReqUniquePath")
	return ring
}

func makeContextKey(name string) ContextKey {
	return ContextKey(name)
}

func (r ContextKeyring) AddKeyToRing(key string) {
	r[key] = makeContextKey(key)
}

func (r ContextKeyring) UseKey(key string) ContextKey {
	return r[key]
}
