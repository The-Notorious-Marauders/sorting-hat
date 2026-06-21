package main

import (
	"github.com/The-Notorious-Marauders/sorting-hat/rest/bootstrap"
	"go.uber.org/fx"
)

// @title           Sorting Hat Auth Service
// @version         1.0
// @description     Authentication and authorization service.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	fx.New(bootstrap.All()).Run()
}
