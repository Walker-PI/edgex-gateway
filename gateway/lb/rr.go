package lb

type RRBalance struct {
}

func newRRBalance() LoadBalance {
	return &RRBalance{}
}

func (rr *RRBalance) Select() string {

	return ""
}
