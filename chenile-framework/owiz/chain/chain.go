package chain

import "context"

type Command[T any] interface {
	Execute(context.Context, *T) error
}

type CommandFunc[T any] func(context.Context, *T) error

func (f CommandFunc[T]) Execute(ctx context.Context, value *T) error {
	return f(ctx, value)
}

type Chain[T any] struct {
	commands []Command[T]
}

func New[T any](commands ...Command[T]) *Chain[T] {
	return &Chain[T]{commands: append([]Command[T]{}, commands...)}
}

func (c *Chain[T]) Add(command Command[T]) {
	c.commands = append(c.commands, command)
}

func (c *Chain[T]) Execute(ctx context.Context, value *T) error {
	for _, command := range c.commands {
		if err := command.Execute(ctx, value); err != nil {
			return err
		}
	}
	return nil
}
