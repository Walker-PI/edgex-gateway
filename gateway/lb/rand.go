package lb

type RandBalance struct {
}

func newRandBalance() LoadBalance {
	return &RandBalance{}
}

func (rb *RandBalance) Select() string {

	return ""
}
