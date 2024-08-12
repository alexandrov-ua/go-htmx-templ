package templates

type TodoContext struct {
	Todos []TodoModel
}

type TodoModel struct {
	Title, Description string
}
