package session

import "fmt"

type IdStack struct {
	ids []string
}

func NewIdStack() *IdStack {
	stack := IdStack{
		ids: []string{},
	}

	return &stack
}

func (s *IdStack) PushBack(id string) {
	s.ids = append(s.ids, id)
}

func (s *IdStack) IndexOf(id string) int {
	for i, v := range s.ids {
		if v == id {
			return i
		}
	}

	return -1
}

func (s *IdStack) Remove(id string) error {
	idx := s.IndexOf(id)
	if idx < 0 {
		return fmt.Errorf("failed to find id")
	}

	s.ids = append(s.ids[:idx], s.ids[idx+1:]...)

	return nil
}

func (s *IdStack) Peak() string {
	if len(s.ids) == 0 {
		return ""
	}

	return s.ids[0]
}

func (s *IdStack) Pop() (string, error) {
	if len(s.ids) == 0 {
		return "", fmt.Errorf("no ids in stack")
	}
	val := s.ids[0]
	s.ids = s.ids[1:]

	return val, nil
}

func (s *IdStack) Count() int {
	return len(s.ids)
}
