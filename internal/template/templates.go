/*
Copyright © 2025 2xhamzeh
*/
package template

import "embed"

//go:embed all:templates
var templateFS embed.FS

var EmptyTemplate = templateConfig{
	fs:            templateFS,
	basePath:      "templates/empty",
	oldModuleName: "example.com/app",
}

var RestTemplate = templateConfig{
	fs:            templateFS,
	basePath:      "templates/rest",
	oldModuleName: "example.com/rest",
}
