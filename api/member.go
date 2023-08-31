package api

type Membership interface {
	Members() []Endpoint
	IsLeader(e Endpoint) (bool, error)
}

func GetLeader(m Membership) (Endpoint, error) {
	for _, e := range m.Members() {
		isLeader, err := m.IsLeader(e)
		if err != nil {
			return nil, err
		}
		if isLeader {
			return e, nil
		}
	}
	return nil, nil
}
