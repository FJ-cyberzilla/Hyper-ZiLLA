package trace

type Span struct {
	name string
	tags map[string]interface{}
}

func StartSpan(name string) *Span {
	return &Span{
		name: name,
		tags: make(map[string]interface{}),
	}
}

func (s *Span) End() {}

func (s *Span) SetTag(key string, value interface{}) {
	s.tags[key] = value
}
