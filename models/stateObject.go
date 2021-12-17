package models

import (
	"fmt"
	"time"
)

type StateObject struct {
	StateCode   string    `json:"statecode"`
	TotalNo     float64   `json:"totalno"`
	LastUpdated time.Time `json:"lastupdated"`
	StateName   string    `json:"statename"`
}

func (this StateObject) ToString() string {
	result := fmt.Sprintf("\nstateObject stateCode: %s", this.StateCode)
	result = result + fmt.Sprintf("\nstate info : %f", this.TotalNo)
	result = result + fmt.Sprintf("\nlast updated : %s", this.LastUpdated)
	result = result + fmt.Sprintf("\nstatename: %s", this.StateName)
	return result
}
