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
	sizeCmd.PersistentFlags().StringVarP(&url, "url", "u", "", "vCenter url")
	sizeCmd.PersistentFlags().StringVarP(&user, "user", "n", "", "vCenter username")
	sizeCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "vCenter password")
	sizeCmd.PersistentFlags().StringVarP(&datacenter, "datacenter", "c", "", "vCenter datacenter name")
	sizeCmd.PersistentFlags().StringVarP(&vmName, "vmName", "m", "", "name of an existing vmName (takes precedence over size)")
	_ = sizeCmd.MarkPersistentFlagRequired("url")
	_ = sizeCmd.MarkPersistentFlagRequired("user")
	_ = sizeCmd.MarkPersistentFlagRequired("password")
	_ = sizeCmd.MarkPersistentFlagRequired("datacenter")
	_ = sizeCmd.MarkPersistentFlagRequired("vmName")
	rootCmd.AddCommand(sizeCmd)
}

func (s *SizeResponse) run() error {
	var err error
	client, err := vsphere.NewClient(url, user, password)
	if err != nil {
		return err
	}
	client.Datacenter, err = client.GetDatacenter(datacenter)
	if err != nil {
		return err
	}

	totalSize, err := client.GetVMTotalStorageSize(vmName)
	if err != nil {
		return err
	}
	s.TotalDiskSize = totalSize
	s.VMName = vmName
	s.Success = true
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
