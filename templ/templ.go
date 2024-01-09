package templ

import "embed"

//go:embed mail.html
var EmailHTML embed.FS

//go:embed companyXYZLogo.png
var CompanyXYZLogo embed.FS
