package command

type Registry struct {
	commands []Command
}

func NewRegistry() *Registry {
	return &Registry{
		commands: make([]Command, 0),
	}
}

func (r *Registry) Regiseter(cmd Command) {
	r.commands = append(r.commands, cmd)
}

func (r *Registry) Find(prefix, name string) *Command {
	for i := range r.commands {
		cmd := &r.commands[i]

		if cmd.Prefix == prefix && cmd.Name == name {
			return cmd
		}
	}

	return nil
}

func (r *Registry) PublicComands() []Command {
	var public []Command

	for _, cmd := range r.commands {
		if !cmd.AdminOnly {
			public = append(public, cmd)
		}
	}

	return public
}
