package goproxyscrape

import "fmt"

type Anonymity uint8

const (
	Transparent Anonymity = iota
	Anonymous
	Elite
	UnknownAnonymity
)

func (a Anonymity) String() string {
	if a < Transparent || a > UnknownAnonymity {
		return "Anonymity(" + string(a) + ")"
	}
	return [...]string{"Transparent", "Anonymous", "Elite", "Unknown"}[a]
}

func (a *Anonymity) Parse(str string) error {
	switch str {
	case "Transparent":
		*a = Transparent
	case "Anonymous":
		*a = Anonymous
	case "Elite":
		*a = Elite
	case "Unknown":
		*a = UnknownAnonymity
	default:
		return fmt.Errorf("invalid anonymity: %s", str)
	}
	return nil
}
