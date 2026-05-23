package strategy

func BuildStrategy(name string) Strategy {

	switch name {

	case "round_robin":
		return NewRoundRobin()

	case "weighted_round_robin":
		return NewWeightedRoundRobin()

	case "least_connections":
		return NewLeastConnections()

	default:
		return NewRoundRobin()
	}
}
