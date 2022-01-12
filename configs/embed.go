package configs

import "embed"

//go:embed config.toml
var emb embed.FS
