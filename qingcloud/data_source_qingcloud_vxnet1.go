package qingcloud

import (
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

const (
	dataSourceVxnetType   = "vxnet_type"
	dataSourceZone        = "zone"
	dataSourceVerbose     = "verbose"
	dataSourceSearchWorld = "search_word"
	dataSourceVxNetID     = "vxnet_id"
)

func dataSourceQingcloudVxNet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVxNetRead,

		Schema: map[string]*schema.Schema{
			dataSourceVxnetType: {
				Type:         schema.TypeInt,
				Required:     true,
				Default:      1,
				ValidateFunc: withinArrayInt(0, 1, 2),
			},
			dataSourceZone: {
				Type:     schema.TypeString,
				Optional: true,
			},
			dataSourceVerbose: {
				Type:         schema.TypeInt,
				Computed:     true,
				Default:      0,
				ValidateFunc: withinArrayInt(0, 1),
			},
			dataSourceVxNetID: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVxNetRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).vxnet
	input := new(qc.DescribeVxNetsInput)
	input.VxNetType = qc.Int(d.Get(resourceVxnetType).(int))
	input.SearchWord = getSetStringPointer(d, dataSourceSearchWorld)
	vxnets := getSetStringPointer(d, resourceInstanceManagedVxnetID)
	var vxnetsStringPointer []*string
	for _, s := range strings.Split(*vxnets, ",") {
		vxnetsStringPointer = append(vxnetsStringPointer, &s)
	}
	input.VxNets = vxnetsStringPointer

	var output *qc.DescribeVxNetsOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.DescribeVxNets(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	return nil
}
