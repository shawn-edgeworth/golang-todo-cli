package todo

type Todo struct {
	Value string
	Done  bool
}

func NewTodo(value string) *Todo {
	return &Todo{
		Value: value,
		Done:  false,
	}
}

func (t *Todo) Complete() {
	t.Done = true
}

func (t *Todo) Undone() {
	t.Done = false
}
