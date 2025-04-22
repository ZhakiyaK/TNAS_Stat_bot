package entities

import "fmt"

type Stats struct {
	Status     string
	MemoryLeft int
	Avail      string
}

func (s *Stats) String() string {
	return fmt.Sprintf("Статус: %s\nОсталось места: %s", s.Status, s.Avail)
}
