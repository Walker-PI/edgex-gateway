package lb

type URLHashBalance struct {
}

func newURLHashBalance() LoadBalance {
	return &URLHashBalance{}
}

func (urlHash *URLHashBalance) Select() string {
	return ""
}
