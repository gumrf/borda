package main

import (
	borda "borda/internal/app"
)

const configPath string = "configs"

func main() {
	borda.Run(configPath)
}
