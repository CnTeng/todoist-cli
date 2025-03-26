package cli

type cli struct {
	icons *icons
}

func NewCLI(iconType iconType) *cli {
	return &cli{icons: newIcons(iconType).withColor()}
}
