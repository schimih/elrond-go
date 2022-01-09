package result

type Port struct {
	ID       uint
	Number   int
	Protocol string
	State    string
	Owner    string
}

func NewPort(id uint,
	number int,
	protocol string,
	state string,
	owner string) Port {

	return Port{
		ID:       id,
		Number:   number,
		Protocol: protocol,
		State:    state,
		Owner:    owner,
	}
}
