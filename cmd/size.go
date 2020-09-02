package cmd

import (
	"github.com/jacobweinstock/vvalidator/pkg/vsphere"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	sizeCmd = &cobra.Command{
		Use:   "size",
		Short: "Gets the total size of all disk(s) attached to a VM",
		Long:  "Gets the total size of all disk(s) attached to a VM",
		Run: func(cmd *cobra.Command, args []string) {
			getVMSize()
		},
	}
)

func init() {
	sizeCmd.PersistentFlags().StringVarP(&vCenterURL, "url", "u", "", "vCenter URL")
	sizeCmd.PersistentFlags().StringVarP(&vCenterUser, "user", "n", "", "vCenter username")
	sizeCmd.PersistentFlags().StringVarP(&vCenterPassword, "password", "p", "", "vCenter password")
	sizeCmd.PersistentFlags().StringVarP(&vCenterDatacenter, "datacenter", "c", "", "vCenter datacenter name")
	sizeCmd.PersistentFlags().StringVarP(&vCenterVM, "vm", "m", "", "name of an existing VM (takes precedence over size)")
	_ = sizeCmd.MarkPersistentFlagRequired("url")
	_ = sizeCmd.MarkPersistentFlagRequired("user")
	_ = sizeCmd.MarkPersistentFlagRequired("password")
	_ = sizeCmd.MarkPersistentFlagRequired("datacenter")
	_ = sizeCmd.MarkPersistentFlagRequired("vm")
	rootCmd.AddCommand(sizeCmd)
}

func getVMSize() {
	log.Info("getting VM disk(s) size")
	client, err := vsphere.NewClient(vCenterURL, vCenterUser, vCenterPassword)
	if err != nil {
		log.Fatal(err)
	}
	client.Datacenter, err = client.GetDatacenter(vCenterDatacenter)
	if err != nil {
		log.Fatal(err)
	}

	totalSize, err := client.GetVMTotalStorageSize(vCenterVM)
	if err != nil {
		log.Fatal(err)
	}
	log.WithFields(log.Fields{
		"vmName": vCenterVM,
		"totalDiskSize": totalSize,
	}).Info("total size of all disk(s)")
}
