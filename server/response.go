package server

import "strconv"

type errorResponse struct {
	Status  int    `json:"status"`
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// This specifies our output float precision
type preciseFloat32 float32

func (pf *preciseFloat32) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatFloat(float64(*pf), 'f', 6, 32)), nil
}

type getCurrentByCoordsResponse struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attribtues currentAttributes
	// Links      apiLinks
	// Relationships apiRelationships
	// meta apiMeta
}

type currentAttributes struct {
	Latitude    preciseFloat32 `json:"latitude"`
	Longitude   preciseFloat32 `json:"longitude"`
	Temperature string         `json:"temperature"`
	Condition   string         `json:"condition"`
}
