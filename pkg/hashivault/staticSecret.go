package hashivault

type staticSecret struct {
	sec *secret
}

func (s *staticSecret) get() map[string]any {
	return s.sec.GetData()
}
