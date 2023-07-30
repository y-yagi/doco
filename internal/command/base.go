package command

type Command interface {
	Run() error
}
