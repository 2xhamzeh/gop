package template

import "embed"

//go:embed all:templates
var templateFS embed.FS

var EmptyTemplate = templateConfig{
	fs:            templateFS,
	basePath:      "templates/empty",
	oldModuleName: "example.com/app",
}

var AppTemplate = templateConfig{
	fs:            templateFS,
	basePath:      "templates/app",
	oldModuleName: "example.com/app",
}
