package models

import (
	"fmt"
)

//swagger:model LatLong
type LatLong struct {
	// Latitude of user
	Lat string `json:"lat" query:"lat"`
	// Longititude of user
	Long string `json:"long" query:"long"`
}

func (this LatLong) ToString() string {
	result := fmt.Sprintf("\nstateObject stateCode: %s", this.Lat)
	result = result + fmt.Sprintf("\nlast updated : %s", this.Long)
	return result
}
