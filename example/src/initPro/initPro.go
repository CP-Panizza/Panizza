package initPro

import (
	
	."github.com/CP-Panizza/Panizza"
	."github.com/CP-Panizza/Panizza/example/src/hello/controller"
	."github.com/CP-Panizza/Panizza/example/src/hello/service"
)

var	App = New()

type Components struct {
	HelloController `@Controller`
	HelloService `@Service`
}

func init() {
	RegisterComponents(new(Components))
}