package command

// Command represents a Redis command with its name and arguments
type Command struct {
	Cmd  string
	Args []string
}
