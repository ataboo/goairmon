package session

import (
	"testing"
)

func _assertStack(stack *IdStack, count int, hasId string, t *testing.T) {
	if stack.Count() != count {
		t.Errorf("unnexpected stack count %d, %d", count, stack.Count())
	}

	if stack.IndexOf(hasId) < 0 {
		t.Errorf("failed to find id in stack")
	}
}

func TestIdStack(t *testing.T) {
	stack := NewIdStack()

	_, err := stack.Pop()
	if err == nil {
		t.Error("expected fail to pop empty")
	}

	err = stack.Remove("some id")
	if err == nil {
		t.Error("expected error when removing from empty stack")
	}

	id := stack.Peak()
	if id != "" {
		t.Error("expected empty id")
	}

	stack.PushBack("first_id")

	_assertStack(stack, 1, "first_id", t)

	stack.PushBack("second_id")
	stack.PushBack("third_id")

	idx := stack.IndexOf("not_in_stack")
	if idx >= 0 {
		t.Error("expected index not found")
	}

	id, err = stack.Pop()
	if err != nil || id != "first_id" {
		t.Errorf("unexpected result from pop %s, %s", id, err)
	}

	_assertStack(stack, 2, "third_id", t)
}
