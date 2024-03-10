package resources

import (
	"fmt"
	"log"

	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/network"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func UpdateDnsRecord(ctx *pulumi.Context, zoneName string, aRecord string, ipAddress pulumi.StringPtrOutput) *network.PrivateRecordSet {

	name := fmt.Sprintf("%s.%s", aRecord, zoneName)
	myDns, err := network.NewPrivateRecordSet(ctx, name, &network.PrivateRecordSetArgs{
		ARecords: network.ARecordArray{
			network.ARecordArgs{
				Ipv4Address: ipAddress,
			},
		},
		Metadata: pulumi.StringMap{
			"key1": pulumi.String("value1"),
		},
		RecordType:            pulumi.String("A"),
		RelativeRecordSetName: pulumi.String(aRecord),
		ResourceGroupName:     pulumi.String(WorkerDnsRG),
		Ttl:                   pulumi.Float64(1),
		PrivateZoneName:       pulumi.String(DnsZone),
	})

	if err != nil {
		log.Fatalf(err.Error())
	}

	//ctx.Export("workerDns", myDns)

	return myDns
}
