package models

import (
	"fmt"
	"time"
)

//swagger:model StateObject
type StateObject struct {
	// State code
	StateCode string `json:"statecode"`
	// Total number of cases
	TotalNo float64 `json:"totalno"`
	// Last update time
	LastUpdated time.Time `json:"lastupdated"`
	//State name
	StateName string `json:"statename"`
}

func (this StateObject) ToString() string {
	result := fmt.Sprintf("\nstateObject stateCode: %s", this.StateCode)
	result = result + fmt.Sprintf("\nstate info : %f", this.TotalNo)
	result = result + fmt.Sprintf("\nlast updated : %s", this.LastUpdated)
	result = result + fmt.Sprintf("\nstatename: %s", this.StateName)
	return result
}
