package ui

import (
	_ "embed"
	"path/filepath"
)

const (
	defaultFontSize    = 48
	defaultLineSpacing = 10
)

var (
	// mainFont     = assets.Fonts["JSE_AmigaAmos"]
	mainMenuPath = filepath.Join("assets", "menus", "game", "mainmenu.yaml")
)
