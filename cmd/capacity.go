package cmd

import (
	"github.com/jacobweinstock/vvalidator/pkg/vsphere"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	vCenterURL        string
	vCenterUser       string
	vCenterPassword   string
	vCenterDatastore  string
	vCenterDatacenter string
	vCenterVM         string
	vmSize            float64

	capacityCmd = &cobra.Command{
		Use:   "capacity",
		Short: "Validate a VM's disks or a size in GB is available on a specified datastore",
		Long:  "Validate a VM's disks or a size in GB is available on a specified datastore",
		Args: func(cmd *cobra.Command, args []string) error {
			if vCenterVM == "" && vmSize == 0 {
				return errors.New("please set either --vm or --size")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			validateCapacity()
		},
	}
)

func init() {
	capacityCmd.PersistentFlags().StringVarP(&vCenterURL, "url", "u", "", "vCenter URL")
	capacityCmd.PersistentFlags().StringVarP(&vCenterUser, "user", "n", "", "vCenter username")
	capacityCmd.PersistentFlags().StringVarP(&vCenterPassword, "password", "p", "", "vCenter password")
	capacityCmd.PersistentFlags().StringVarP(&vCenterDatacenter, "datacenter", "c", "", "vCenter datacenter name")
	capacityCmd.PersistentFlags().StringVarP(&vCenterDatastore, "datastore", "s", "", "vCenter datastore name")
	capacityCmd.PersistentFlags().StringVarP(&vCenterVM, "vm", "m", "", "name of an existing VM (takes precedence over size)")
	capacityCmd.PersistentFlags().Float64VarP(&vmSize, "size", "z", 0, "vm disk size in GBs")
	_ = capacityCmd.MarkPersistentFlagRequired("url")
	_ = capacityCmd.MarkPersistentFlagRequired("user")
	_ = capacityCmd.MarkPersistentFlagRequired("password")
	_ = capacityCmd.MarkPersistentFlagRequired("datacenter")
	_ = capacityCmd.MarkPersistentFlagRequired("datastore")
	rootCmd.AddCommand(capacityCmd)
}

func validateCapacity() {
	client, err := vsphere.NewClient(vCenterURL, vCenterUser, vCenterPassword)
	if err != nil {
		log.Fatal(err)
	}
	client.Datacenter, err = client.GetDatacenter(vCenterDatacenter)
	if err != nil {
		log.Fatal(err)
	}
	client.Datastore, err = client.GetDatastore(vCenterDatastore)
	if err != nil {
		log.Fatal(err)
	}

	var requestedDiskSpace float64
	if vCenterVM != "" {
		totalSize, err := client.GetVMTotalStorageSize(vCenterVM)
		if err != nil {
			log.Fatal(err)
		}

		requestedDiskSpace = totalSize
	} else {
		requestedDiskSpace = vmSize
	}

	_, free, err := client.DatastoreCapacity()
	if err != nil {
		log.WithFields(log.Fields{
			"requestedSpaceInGBs": requestedDiskSpace,
			"freeSpaceInGBs":      free,
			"spaceAvailable":      false,
		}).Fatal(err)
	}
	if requestedDiskSpace <= free {
		log.WithFields(log.Fields{
			"requestedSpaceInGBs": requestedDiskSpace,
			"freeSpaceInGBs":      free,
			"spaceAvailable":      true,
		}).Info("requested datastore space SHOULD be available")
	} else {
		log.WithFields(log.Fields{
			"requestedSpaceInGBs": requestedDiskSpace,
			"freeSpaceInGBs":      free,
			"spaceAvailable":      false,
		}).Info("requested datastore space is NOT available")
	}
}
