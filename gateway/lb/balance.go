package lb

type LoadBalance interface {
	Select() string
}
