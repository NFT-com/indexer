package logs

type ActionType int

const (
	Addition ActionType = iota + 1
	OwnerChange
)

func (e ActionType) String() string {
	switch e {
	case Addition:
		return "addition"
	case OwnerChange:
		return "owner_change"
	default:
		return ""
	}
}
