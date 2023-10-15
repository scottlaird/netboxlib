package netbox

import (
	"net/netip"

	"github.com/netbox-community/go-netbox/v3/netbox/client"
	"github.com/netbox-community/go-netbox/v3/netbox/client/ipam"
	"github.com/netbox-community/go-netbox/v3/netbox/models"
)

type IPPrefix struct {
	ID                 int64
	Prefix            netip.Prefix
	Display            string
	Family             string
	Status             string
	Description        string
	Tags               map[string]bool // Tags.Name -> true
	VRF                string
}

type IPPrefixes []*IPPrefix

func ipamPrefixToIPPrefix(i *models.Prefix) (*IPPrefix, error) {
	ip := &IPPrefix{
		Description:        i.Description,
		Display:            i.Display,
		ID:                 i.ID,
		Tags:               make(map[string]bool),
	}

	if i.Prefix != nil {
		prefix, err := netip.ParsePrefix(String(i.Prefix))
		if err != nil {
			return nil, err
		}
		ip.Prefix = prefix
	}
	if i.Family != nil {
		ip.Family = String(i.Family.Label)
	}
	if i.Status != nil {
		ip.Status = String(i.Status.Value)
	}
	if i.Vrf != nil {
		ip.VRF = String(i.Vrf.Name)
	}
	for _, t := range i.Tags {
		ip.Tags[*t.Name] = true
	}

	return ip, nil
}

func ListIPPrefixes(c *client.NetBoxAPI) (IPPrefixes, error) {
	var limit int64
	limit = 0

	r := ipam.NewIpamPrefixesListParams()
	r.Limit = &limit

	rs, err := c.Ipam.IpamPrefixesList(r, nil)
	if err != nil {
		return nil, err
	}

	ips := make(IPPrefixes, len(rs.Payload.Results))
	for i, result := range rs.Payload.Results {
		ip, err := ipamPrefixToIPPrefix(result)
		if err != nil {
			return nil, err
		}

		ips[i] = ip
	}

	return ips, nil
}

func GetIPPrefix(c *client.NetBoxAPI, id int64) (*IPPrefix, error) {
	var limit int64
	limit = 0
	idStr := string(id)

	r := ipam.NewIpamPrefixesListParams()
	r.Limit = &limit
	r.ID = &idStr

	rs, err := c.Ipam.IpamPrefixesList(r, nil)
	if err != nil {
		return nil, err
	}

	return ipamPrefixToIPPrefix(rs.Payload.Results[0])
}
