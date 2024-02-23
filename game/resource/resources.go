package resource

type Resource uint

const (
	Food Resource = iota
	Wood
	Stones
	EndResources
)

type ResourceProps struct {
	Name  string
	Short string
}

var Properties = [EndResources]ResourceProps{
	{Name: "food", Short: "F"},
	{Name: "wood", Short: "W"},
	{Name: "stones", Short: "S"},
}
