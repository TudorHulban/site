package domain

type Browser uint8

const (
	Chrome Browser = iota + 1
	Safari
	Edge
	Firefox
	Opera
	Brave
	SamsungInternet
)

func (b Browser) String() string {
	switch b {
	case Chrome:
		return "Chrome"
	case Safari:
		return "Safari"
	case Edge:
		return "Edge"
	case Firefox:
		return "Firefox"
	case Opera:
		return "Opera"
	case Brave:
		return "Brave"
	case SamsungInternet:
		return "Samsung Internet"

	default:
		return "Unknown"
	}
}
