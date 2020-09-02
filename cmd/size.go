package cmd

import (
	"github.com/jacobweinstock/vvalidator/pkg/vsphere"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"path"
)

var (
	sizeCmd = &cobra.Command{
		Use:   "size",
		Short: "Gets the total size of all disk(s) attached to a vmName",
		Long:  "Gets the total size of all disk(s) attached to a vmName",
		Run: func(cmd *cobra.Command, args []string) {
			var size SizeResponse
			err := size.run()
			size.response(err)
		},
	}
)

func init() {

	sizeCmd.PersistentFlags().StringVarP(&vmName, "vmName", "m", "", "name of an existing VM")
	_ = sizeCmd.MarkPersistentFlagRequired("vmName")
	rootCmd.AddCommand(sizeCmd)
}

func (c *SizeResponse) run() error {
	var err error
	c.VMName = vmName
	client, err := vsphere.NewClient(url, user, password)
	if err != nil {
		return err
	}
	client.Datacenter, err = client.GetDatacenter(datacenter)
	if err != nil {
		return err
	}
	c.TotalDiskSize, err = client.GetVMTotalStorageSize(c.VMName)
	if err != nil {
		return err
	}
	c.Success = true
	return err
}

func (c *SizeResponse) response(err error) {
	r := c.ToLogrusFields()
	r["responseFile"] = path.Join(responseFileDirectory, responseFileName)
	if err != nil {
		r["errorMsg"] = err.Error()
		log.WithFields(r).Fatal()
	}
	log.WithFields(r).Info()
}
