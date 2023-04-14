package util

type Set map[string]struct{}

func NewSet() *Set {
	s := make(Set)
	return &s
}

func (s Set) Add(key ...string) {
	for _, v := range key {
		t := v
		(s)[t] = struct{}{}
	}
}

func (s Set) Remove(key string) {
	delete(s, key)
}

func (s Set) Has(key string) bool {
	_, ok := (s)[key]
	return ok
}

func (s Set) Len() int {
	return len(s)
}
