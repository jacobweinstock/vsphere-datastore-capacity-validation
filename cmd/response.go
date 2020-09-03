package cmd

import (
	"github.com/sirupsen/logrus"
)

type baseResponse struct {
	Success  bool   `json:"success"`
	ErrorMsg string `json:"errorMsg"`
}

type capacityResponse struct {
	RequestedSpaceInGBs float64 `json:"requestedSpaceInGBs"`
	FreeSpaceInGBs      float64 `json:"freeSpaceInGBs"`
	SpaceAvailable      bool    `json:"spaceAvailable"`
	baseResponse        `json:",inline"`
}

// ToLogrusFields is a helper for the logrus library
func (s capacityResponse) ToLogrusFields() logrus.Fields {
	return logrus.Fields{
		"requestedSpaceInGBs": s.RequestedSpaceInGBs,
		"freeSpaceInGBs":      s.FreeSpaceInGBs,
		"spaceAvailable":      s.SpaceAvailable,
		"success":             s.Success,
		"errorMsg":            s.ErrorMsg,
	}
}

type sizeResponse struct {
	VMName        string  `json:"vmName"`
	TotalDiskSize float64 `json:"totalDiskSize"`
	baseResponse  `json:",inline"`
}

// ToLogrusFields is a helper for the logrus library
func (c sizeResponse) ToLogrusFields() logrus.Fields {
	return logrus.Fields{
		"vmName":        c.VMName,
		"totalDiskSize": c.TotalDiskSize,
		"success":       c.Success,
		"errorMsg":      c.ErrorMsg,
	}
}
