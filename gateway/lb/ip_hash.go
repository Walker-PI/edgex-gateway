package lb

type IPHashBalance struct {
}

func newIPHashBalance() LoadBalance {
	return &IPHashBalance{}
}

func (ipHash *IPHashBalance) Select() string {
	return ""
}
