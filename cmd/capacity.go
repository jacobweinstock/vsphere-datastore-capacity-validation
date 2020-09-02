package cmd

import (
	"github.com/jacobweinstock/vvalidator/pkg/vsphere"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path"
)

var (
	url        string
	user       string
	password   string
	datastore  string
	datacenter string
	vmName     string
	vmSize     float64

	capacityCmd = &cobra.Command{
		Use:   "capacity",
		Short: "Validate a vmName's disks or a size in GB is available on a specified datastore",
		Long:  "Validate a vmName's disks or a size in GB is available on a specified datastore",
		Args: func(cmd *cobra.Command, args []string) error {
			if vmName == "" && vmSize == 0 || viper.Get("vmName") == "" && viper.Get("vmSize") == 0 {
				return errors.New("please set a vm name or a size. use either a flag (--vmName, --vmSize) or an env var (VVALIDATOR_VMNAME, VVALIDATOR_VMSIZE)")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			flagEnvSetter()
			var capacity capacityResponse
			err := capacity.run()
			capacity.response(err)
		},
	}
)

func init() {
	capacityCmd.PersistentFlags().StringVarP(&url, "url", "u", "", "vCenter url")
	capacityCmd.PersistentFlags().StringVarP(&user, "user", "n", "", "vCenter username")
	capacityCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "vCenter password")
	capacityCmd.PersistentFlags().StringVarP(&datacenter, "datacenter", "c", "", "vCenter datacenter name")
	capacityCmd.PersistentFlags().StringVarP(&datastore, "datastore", "s", "", "vCenter datastore name")
	capacityCmd.PersistentFlags().StringVarP(&vmName, "vmName", "m", "", "name of an existing vmName (takes precedence over size)")
	capacityCmd.PersistentFlags().Float64VarP(&vmSize, "vmSize", "z", 0, "vm disk size in GBs")
	_ = capacityCmd.MarkPersistentFlagRequired("url")
	_ = capacityCmd.MarkPersistentFlagRequired("user")
	_ = capacityCmd.MarkPersistentFlagRequired("password")
	_ = capacityCmd.MarkPersistentFlagRequired("datacenter")
	_ = capacityCmd.MarkPersistentFlagRequired("datastore")
	rootCmd.AddCommand(capacityCmd)
}

// cobra and viper play nice but dont allow for a common variable to get either the flag or the env var
// this function make a variable the single source of the flag and env value
func flagEnvSetter() {
	if url == "" {
		url = viper.GetString("url")
	}
	if user == "" {
		user = viper.GetString("user")
	}
	if password == "" {
		password = viper.GetString("password")
	}
	if datacenter == "" {
		datacenter = viper.GetString("datacenter")
	}
	if datastore == "" {
		datastore = viper.GetString("datastore")
	}
	if vmName == "" {
		vmName = viper.GetString("vmName")
	}
	if vmSize == 0 {
		vmSize = viper.GetFloat64("vmSize")
	}
}

func (c *capacityResponse) run() error {
	var err error
	client, err := vsphere.NewClient(url, user, password)
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

	var requestedDiskSpace float64
	if vmName != "" {
		totalSize, err := client.GetVMTotalStorageSize(vmName)
		if err != nil {
			return err
		}
		requestedDiskSpace = totalSize
	} else {
		requestedDiskSpace = vmSize
	}

	_, free, err := client.DatastoreCapacity()
	if err != nil {
		c.RequestedSpaceInGBs = requestedDiskSpace
		c.FreeSpaceInGBs = free
		c.SpaceAvailable = false
		c.Success = false
		return err
	}
	if requestedDiskSpace <= free {
		c.SpaceAvailable = true
	} else {
		c.SpaceAvailable = false
	}
	c.RequestedSpaceInGBs = requestedDiskSpace
	c.FreeSpaceInGBs = free
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
