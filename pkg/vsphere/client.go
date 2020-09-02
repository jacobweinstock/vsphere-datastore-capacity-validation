package vsphere

import (
	"context"
	"net/url"
	"strings"

	"github.com/pkg/errors"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"

	"github.com/vmware/govmomi"
)

type Session struct {
	Conn         *govmomi.Client
	Datacenter   *object.Datacenter
	Datastore    *object.Datastore
	Folder       *object.Folder
	ResourcePool *object.ResourcePool
	Network      object.NetworkReference
	Ctx          context.Context
}

// NewClient returns a new vsphere Session
func NewClient(server string, username string, password string, ctx context.Context) (*Session, error) {
	sm := new(Session)
	if !strings.HasPrefix(server, "https://") && !strings.HasPrefix(server, "http://") {
		server = "https://" + server
	}
	nonAuthURL, err := url.Parse(server)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse vCenter url %s", server)
	}
	if !strings.HasSuffix(nonAuthURL.Path, "sdk") {
		nonAuthURL.Path = nonAuthURL.Path + "sdk"
	}
	authenticatedURL, err := url.Parse(nonAuthURL.String())
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse vCenter url %s", nonAuthURL.String())
	}
	client, err := govmomi.NewClient(ctx, nonAuthURL, true)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create new vSphere client")
	}
	authenticatedURL.User = url.UserPassword(username, password)
	if err = client.Login(ctx, authenticatedURL.User); err != nil {
		return nil, errors.Wrap(err, "unable to login to vSphere")
	}
	sm.Conn = client
	sm.Ctx = ctx

	return sm, nil
}

// GetDatacenter returns the govmomi object for a datacenter
func (s *Session) GetDatacenter(name string) (*object.Datacenter, error) {
	finder := find.NewFinder(s.Conn.Client, true)
	datacenter, err := finder.Datacenter(s.Ctx, name)
	if err != nil {
		return nil, errors.Wrapf(err, "error finding datacenter %s", name)
	}
	return datacenter, err
}

// GetDatastore returns the govmomi object for a datastore
func (s *Session) GetDatastore(name string) (*object.Datastore, error) {
	finder := find.NewFinder(s.Conn.Client, true)
	finder.SetDatacenter(s.Datacenter)
	datastore, err := finder.Datastore(s.Ctx, name)
	if err != nil {
		return nil, errors.Wrapf(err, "error finding datastore %s", name)
	}
	return datastore, err
}

// GetNetwork returns the govmomi object for a network
func (s *Session) GetNetwork(name string) (object.NetworkReference, error) {
	finder := find.NewFinder(s.Conn.Client, true)
	finder.SetDatacenter(s.Datacenter)
	network, err := finder.Network(s.Ctx, name)
	if err != nil {
		return nil, errors.Wrapf(err, "error finding network %s", name)
	}
	return network, err
}

// GetResourcePool returns the govmomi object for a resource pool
func (s *Session) GetResourcePool(name string) (*object.ResourcePool, error) {
	finder := find.NewFinder(s.Conn.Client, true)
	finder.SetDatacenter(s.Datacenter)
	resourcePool, err := finder.ResourcePool(s.Ctx, name)
	if err != nil {
		return nil, errors.Wrapf(err, "error finding resource pool %s", name)
	}
	return resourcePool, err
}

// GetVM returns the govmomi object for a virtual machine
func (s *Session) GetVM(name string) (*object.VirtualMachine, error) {
	finder := find.NewFinder(s.Conn.Client, true)
	finder.SetDatacenter(s.Datacenter)
	vm, err := finder.VirtualMachine(s.Ctx, name)
	if err != nil {
		return nil, errors.Wrapf(err, "error finding VM %s", name)
	}
	return vm, err
}
