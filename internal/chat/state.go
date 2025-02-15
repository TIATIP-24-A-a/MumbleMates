package chat

type Message struct {
}

type State struct {
	Name string
}

func NewState() *State {
	return &State{}
}
