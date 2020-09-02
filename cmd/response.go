package cmd

import (
	"github.com/sirupsen/logrus"
)

type BaseResponse struct {
	Success  bool   `json:"success"`
	ErrorMsg string `json:"errorMsg"`
}

type capacityResponse struct {
	RequestedSpaceInGBs float64 `json:"requestedSpaceInGBs"`
	FreeSpaceInGBs      float64 `json:"freeSpaceInGBs"`
	SpaceAvailable      bool    `json:"spaceAvailable"`
	BaseResponse        `json:",inline"`
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

type SizeResponse struct {
	VMName        string  `json:"vmName"`
	TotalDiskSize float64 `json:"totalDiskSize"`
	BaseResponse  `json:",inline"`
}

// ToLogrusFields is a helper for the logrus library
func (c SizeResponse) ToLogrusFields() logrus.Fields {
	return logrus.Fields{
		"vmName":        c.VMName,
		"totalDiskSize": c.TotalDiskSize,
		"success":       c.Success,
		"errorMsg":      c.ErrorMsg,
	}
}
