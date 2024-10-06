package bananas

import "embed"

//go:embed templates/cmd/*
var BaseFS embed.FS

//go:embed templates/databases/*
var DbFS embed.FS

//go:embed templates/pkg/*
var PkgFS embed.FS

//go:embed templates/scripts/*
var ScriptsFS embed.FS

//go:embed templates/middlewares/*
var MiddlewareFS embed.FS
