package cmd

import (
	"context"
	"path"
	"time"

	"github.com/jacobweinstock/vvalidator/pkg/vsphere"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	sizeCmd = &cobra.Command{
		Use:   "size",
		Short: "Gets the total size of all disk(s) attached to a vmName",
		Long:  "Gets the total size of all disk(s) attached to a vmName",
		Run: func(cmd *cobra.Command, args []string) {
			var size sizeResponse
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

func (c *sizeResponse) run() error {
	var err error
	tout := time.Duration(timeout) * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), tout)
	defer cancel()
	c.VMName = vmName
	client, err := vsphere.NewClient(ctx, url, user, password)
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

func (c *sizeResponse) response(err error) {
	r := c.ToLogrusFields()
	r["responseFile"] = path.Join(responseFileDirectory, responseFileName)
	if err != nil {
		r["errorMsg"] = err.Error()
		log.WithFields(r).Fatal()
	}
	log.WithFields(r).Info()
}
