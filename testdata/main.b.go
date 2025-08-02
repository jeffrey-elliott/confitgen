package main

import (
	"fmt"

	"hhgttg/internal/confit/guidebook"
	"hhgttg/internal/confit/starship"
)

type AppConfig struct {
	Guidebook guidebook.Guidebook
	Starship  starship.Starship
}

func LoadConfig() (AppConfig, error) {
	g, err := guidebook.Load("guidebook.confit.values.json")
	if err != nil {
		return AppConfig{}, err
	}
	s, err := starship.Load("starship.confit.values.json")
	if err != nil {
		return AppConfig{}, err
	}
	return AppConfig{Guidebook: g, Starship: s}, nil
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

	fmt.Printf("%v", cf.Starship)
	fmt.Println("INFINITE IMPROBABILITY:", cf.Starship.InfiniteImprobabilityEnabled)
	fmt.Println("DESTINATION:", cf.Starship.Destination)
	fmt.Println("CREW:", cf.Starship.Crew)
}
