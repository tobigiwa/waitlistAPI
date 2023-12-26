package templ

import "embed"

//go:embed mail.html
var EmailHTML embed.FS

//go:embed BlockRideLogo.png
var BlockRideLogo embed.FS
