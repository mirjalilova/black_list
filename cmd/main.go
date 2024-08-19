package main

import (
	"github.com/mirjalilova/black_list/internal/config"
	"github.com/mirjalilova/black_list/pkg/app"
)

func main() {
	cnf := config.Load()

	app.Run(&cnf)
}
