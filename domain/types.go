package domain

// Types that are reusable across the app

type Coords struct {
	Latitude  float32
	Longitude float32
}

type Temperature string

const (
	TempUnknown Temperature = "unknown"
	TempHot     Temperature = "hot"
	TempCold    Temperature = "cold"
	TempMod     Temperature = "moderate"
)

type Weather struct {
	Coords      Coords
	States      []string
	Temperature Temperature
}

// RepoWeather purely existing so that WeatherService.CurrentIn actually does something.
// In a normal case we would convert the repo data into domain data.  AKA join states, and convert the temperature.
type RepoWeather struct {
	Coords      Coords
	States      []string
	Temperature float32
}
