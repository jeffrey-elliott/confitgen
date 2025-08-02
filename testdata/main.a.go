package main

import (
	"fmt"

	"hhgttg/internal/confit/guidebook"
)

type AppConfig struct {
	Guidebook guidebook.Guidebook
}

func LoadConfig() (AppConfig, error) {
	g, err := guidebook.Load("guidebook.confit.values.json")
	if err != nil {
		return AppConfig{}, err
	}
	return AppConfig{Guidebook: g}, nil
}

func main() {
	cf, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", cf.Guidebook)
	fmt.Println("THEME:", cf.Guidebook.Display.Theme)
	fmt.Println("BRITISHNESS:", cf.Guidebook.Narration.Britishness)
	fmt.Println("INSTALLED MODULES:", cf.Guidebook.InstalledModules)
}
