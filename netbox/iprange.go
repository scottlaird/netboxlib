package netbox

import (
	"fmt"
	"net/netip"
	"reflect"

	"github.com/netbox-community/go-netbox/v3/netbox/client"
	"github.com/netbox-community/go-netbox/v3/netbox/client/ipam"
	"github.com/netbox-community/go-netbox/v3/netbox/models"
)

type IPRange struct {
	CustomFields       map[string]reflect.Value
	StartAddress            netip.Prefix
	EndAddress netip.Prefix
	Description        string
	Display            string
	Family             string
	ID                 int64
	Status             string
	Tags               map[string]bool // Tags.Name -> true
	VRF                string
}

type IPRanges []*IPRange


func ipamRangeToIPRange(i *models.IPRange) (*IPRange, error) {
	r := &IPRange{
		CustomFields:       make(map[string]reflect.Value),
		Description:        i.Description,
		Display:            i.Display,
		ID:                 i.ID,
		Tags:               make(map[string]bool),
	}

	if i.StartAddress != nil {
		prefix, err := netip.ParsePrefix(String(i.StartAddress))
		if err != nil {
			return nil, err
		}
		r.StartAddress = prefix
	}
	if i.EndAddress != nil {
		prefix, err := netip.ParsePrefix(String(i.EndAddress))
		if err != nil {
			return nil, err
		}
		r.EndAddress = prefix
	}
	if i.Family != nil {
		r.Family = String(i.Family.Label)
	}
	if i.Status != nil {
		r.Status = String(i.Status.Value)
	}
	if i.Vrf != nil {
		r.VRF = String(i.Vrf.Name)
	}
	for _, t := range i.Tags {
		r.Tags[*t.Name] = true
	}

	v := reflect.ValueOf(i.CustomFields)
	if v.Kind() == reflect.Map {
		for _, key := range v.MapKeys() {
			r.CustomFields[key.String()] = v.MapIndex(key)
		}
	}

	return r, nil
}

func ListIPRanges(c *client.NetBoxAPI) (IPRanges, error) {
	limit := int64(0)

	r := ipam.NewIpamIPRangesListParams()
	r.Limit = &limit

	rs, err := c.Ipam.IpamIPRangesList(r, nil)
	if err != nil {
		return nil, err
	}

	ranges := make(IPRanges, len(rs.Payload.Results))
	for i, result := range rs.Payload.Results {
		r, err := ipamRangeToIPRange(result)
		if err != nil {
			return nil, err
		}

		ranges[i] = r
	}

	return ranges, nil
}

func GetIPRange(c *client.NetBoxAPI, id int64) (*IPRange, error) {
	limit := int64(0)
	idStr := fmt.Sprint(id)

	r := ipam.NewIpamIPRangesListParams()
	r.Limit = &limit
	r.ID = &idStr

	rs, err := c.Ipam.IpamIPRangesList(r, nil)
	if err != nil {
		return nil, err
	}

	return ipamRangeToIPRange(rs.Payload.Results[0])
}
