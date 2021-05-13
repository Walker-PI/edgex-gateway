package lb

type LoadBalanceFactory struct{}

func (f *LoadBalanceFactory) NewLoadBalance(_type string) (LoadBalance, error) {

	switch _type {
	case "rand":
		return newRandBalance(), nil
	}
	return nil, nil
}
