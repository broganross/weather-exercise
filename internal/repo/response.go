package repo

// TODO: extend to parse unix time stamps correctly, instead of using int64
type currentWeatherResponse struct {
	Coord struct {
		Lat float32 `json:"lat"`
		Lon float32 `json:"lon"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp        float32 `json:"temp"`
		FeelsLike   float32 `json:"feels_like"`
		Min         float32 `json:"temp_min"`
		Max         float32 `json:"temp_max"`
		Pressure    int     `json:"pressure"`
		Humidity    int     `json:"humidity"`
		SeaLevel    int     `json:"sea_level"`
		GroundLevel int     `json:"grnd_level"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed   float32 `json:"speed"`
		Degrees int     `json:"deg"`
		Gust    float32 `json:"gust"`
	} `json:"wind"`
	Rain struct {
		OneHour   float32 `json:"1h"`
		ThreeHour float32 `json:"3h"`
	} `json:"rain"`
	Clouds struct {
		All float32 `json:"all"`
	} `json:"clouds"`
	Snow struct {
		OneHour   float32 `json:"1h"`
		ThreeHour float32 `json:"3h"`
	} `json:"snow"`
	DateTime int64 `json:"dt"`
	Sys      struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int64  `json:"sunrise"`
		Sunset  int64  `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
}
