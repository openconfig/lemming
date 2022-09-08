package main

import (
	"context"
	"flag"
	"fmt"
	"strconv"
	"time"

	log "github.com/golang/glog"
	"github.com/openconfig/lemming/gnmi/gnmit"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/ygnmi/ygnmi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/credentials/local"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
)

var (
	host   = flag.String("host", "localhost", "name of the host hosting the gNMI service")
	port   = flag.Int("port", 6030, "port of the gNMI service on the host.")
	target = flag.String("target", "dut", "name of the fake target")
)

func main() {
	flag.Parse()

	datastoreConn, err := grpc.DialContext(context.Background(), fmt.Sprintf("unix:%s", gnmit.DatastoreAddress), grpc.WithBlock(), grpc.WithTransportCredentials(local.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial %s: %v", gnmit.DatastoreAddress, err)
	}
	datastoreClient := gpb.NewGNMIClient(datastoreConn)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	serverAddr := fmt.Sprintf("%s:%d", *host, *port)
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial %s: %v", serverAddr, err)
	}
	client := gpb.NewGNMIClient(conn)

	yclient, err := ygnmi.NewClient(client, ygnmi.WithTarget(*target))
	if err != nil {
		log.Fatalf("Error while creating ygnmi client: %v", err)
	}

	// A background task that updates the interface description once every second.
	go func() {
		for count := 0; ; count++ {
			time.Sleep(1 * time.Second)
			p, _, err := ygnmi.ResolvePath(ocpath.Root().Interface("test-eth0").Description().Config().PathStruct())
			if err != nil {
				log.Fatal(err)
			}
			if _, err = datastoreClient.Set(context.Background(), &gpb.SetRequest{
				Prefix: &gpb.Path{Origin: "openconfig", Target: *target},
				Update: []*gpb.Update{{
					Path: p,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "testagent" + strconv.Itoa(count)}},
				}},
			}); err != nil {
				log.Fatal(err)
			}
		}
	}()

	// rootWatcher := ygnmi.Watch(
	// 	context.Background(),
	// 	yclient,
	// 	ocpath.Root().State(),
	// 	func(root *ygnmi.Value[*oc.Root]) error {
	// 		log.Infof("%+v\n", root)
	// 		return ygnmi.Continue
	// 	},
	// )
	interfaceWatcher := ygnmi.Watch(
		context.Background(),
		yclient,
		ocpath.Root().Interface("test-eth0").State(),
		func(root *ygnmi.Value[*oc.Interface]) error {
			intf, ok := root.Val()
			var desc string
			if ok && intf.Description != nil {
				desc = *intf.Description
			}
			log.Infof("%v, %+v\n", desc, root)
			return ygnmi.Continue
		},
	)
	if _, err := interfaceWatcher.Await(); err != nil {
		log.Fatal(err)
	}
}
