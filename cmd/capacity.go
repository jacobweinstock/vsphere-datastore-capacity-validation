package cmd

import (
	"context"
	"github.com/jacobweinstock/vvalidator/pkg/vsphere"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path"
	"time"
)

var (
	datastore string
	vmName    string
	vmSize    float64

	capacityCmd = &cobra.Command{
		Use:   "capacity",
		Short: "Validate that a vmName's disk(s) or a size, given in GB, is available on a specified datastore",
		Long:  "Validate that a vmName's disk(s) or a size, given in GB, is available on a specified datastore",
		Args: func(cmd *cobra.Command, args []string) error {
			if vmName == "" && vmSize == 0 || viper.Get("vmName") == "" && viper.Get("vmSize") == 0 {
				return errors.New("please set a vm name or a size. use either a flag (--vmName, --vmSize) or an env var (VVALIDATOR_VMNAME, VVALIDATOR_VMSIZE)")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			var capacity capacityResponse
			err := capacity.run()
			capacity.response(err)
		},
	}
)

func init() {
	capacityCmd.PersistentFlags().StringVarP(&datastore, "datastore", "s", "", "vCenter datastore name")
	capacityCmd.PersistentFlags().StringVarP(&vmName, "vmName", "m", "", "name of an existing VM (takes precedence over vmSize)")
	capacityCmd.PersistentFlags().Float64VarP(&vmSize, "vmSize", "z", 0, "vm disk size in GBs")
	_ = capacityCmd.MarkPersistentFlagRequired("datastore")
	rootCmd.AddCommand(capacityCmd)
}

func (c *capacityResponse) run() error {
	var err error
	tout := time.Duration(timeout) * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), tout)
	defer cancel()
	client, err := vsphere.NewClient(ctx, url, user, password)
	if err != nil {
		return err
	}
	client.Datacenter, err = client.GetDatacenter(datacenter)
	if err != nil {
		return err
	}
	client.Datastore, err = client.GetDatastore(datastore)
	if err != nil {
		return err
	}

	if vmName != "" {
		c.RequestedSpaceInGBs, err = client.GetVMTotalStorageSize(vmName)
		if err != nil {
			return err
		}
	} else {
		c.RequestedSpaceInGBs = vmSize
	}

	_, c.FreeSpaceInGBs, err = client.DatastoreCapacity()
	if err != nil {
		c.SpaceAvailable = false
		c.Success = false
		return err
	}
	if c.RequestedSpaceInGBs <= c.FreeSpaceInGBs {
		c.SpaceAvailable = true
	} else {
		c.SpaceAvailable = false
	}
	c.Success = true
	return err
}

func (c *capacityResponse) response(err error) {
	r := c.ToLogrusFields()
	r["responseFile"] = path.Join(responseFileDirectory, responseFileName)
	if err != nil {
		r["errorMsg"] = err.Error()
		log.WithFields(r).Fatal()
	}
	log.WithFields(r).Info()
}
