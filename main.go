package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/debug"
	"github.com/gopcua/opcua/id"
	"github.com/gopcua/opcua/ua"
)

func main() {
	var (
		endpoint = flag.String("endpoint", "opc.tcp://localhost:4840", "OPC UA Endpoint URL")
		policy   = flag.String("policy", "", "Security policy: None, Basic128Rsa15, Basic256, Basic256Sha256. Default: auto")
		mode     = flag.String("mode", "", "Security mode: None, Sign, SignAndEncrypt. Default: auto")
		certFile = flag.String("cert", "", "Path to cert.pem. Required for security mode/policy != None")
		keyFile  = flag.String("key", "", "Path to private key.pem. Required for security mode/policy != None")
		nodeID   = flag.String("node", "", "node id to subscribe to")
		event    = flag.Bool("event", false, "subscribe to node event changes (Default: node value changes)")
		interval = flag.Duration("interval", opcua.DefaultSubscriptionInterval, "subscription interval")
	)
	flag.BoolVar(&debug.Enable, "debug", false, "enable debug logging")
	flag.Parse()
	log.SetFlags(0)
	// *nodeID = "ns=2;s=Sauter.EMS.第三方集成.隔离器系统.B区隔离器2021-075lg.北生108隔离器_2021-075Ig_Device_1_B1_PDT02_LL"

	// add an arbitrary timeout to demonstrate how to stop a subscription
	// with a context.
	d := 60 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
	log.Printf("Subscription will stop after %s for demonstration purposes", d)

	endpoints, err := opcua.GetEndpoints(ctx, *endpoint)
	if err != nil {
		fmt.Println("错误1")
		log.Fatal(err)
	}
	ep, err := opcua.SelectEndpoint(endpoints, *policy, ua.MessageSecurityModeFromString(*mode))
	if err != nil {
		fmt.Println("错误2")
		log.Fatal(err)
	}
	ep.EndpointURL = *endpoint

	fmt.Println("*", ep.SecurityPolicyURI, ep.SecurityMode)

	opts := []opcua.Option{
		opcua.SecurityPolicy(*policy),
		opcua.SecurityModeString(*mode),
		opcua.CertificateFile(*certFile),
		opcua.PrivateKeyFile(*keyFile),
		opcua.AuthAnonymous(),
		opcua.SecurityFromEndpoint(ep, ua.UserTokenTypeAnonymous),
	}

	c, err := opcua.NewClient(ep.EndpointURL, opts...)
	if err != nil {
		log.Fatal(err)
	}
	if err := c.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	defer c.Close(ctx)

	notifyCh := make(chan *opcua.PublishNotificationData)

	sub, err := c.Subscribe(ctx, &opcua.SubscriptionParameters{
		Interval: *interval,
	}, notifyCh)
	if err != nil {
		fmt.Println("错误3")
		log.Fatal(err)
	}
	defer sub.Cancel(ctx)
	log.Printf("Created subscription with id %v", sub.SubscriptionID)

	id, err := ua.ParseNodeID(*nodeID)
	if err != nil {
		fmt.Println("错误4")
		log.Fatal(err)
	}

	var miCreateRequest *ua.MonitoredItemCreateRequest
	var eventFieldNames []string
	if *event {
		miCreateRequest, eventFieldNames = eventRequest(id)
	} else {
		miCreateRequest = valueRequest(id)
	}
	res, err := sub.Monitor(ctx, ua.TimestampsToReturnBoth, miCreateRequest)
	if err != nil || res.Results[0].StatusCode != ua.StatusOK {
		fmt.Println("错误5")
		log.Fatal(err)
	}

	// read from subscription's notification channel until ctx is cancelled
	for {
		select {
		case <-ctx.Done():
			return
		case res := <-notifyCh:
			if res.Error != nil {
				fmt.Println("错误6")
				log.Print(res.Error)
				continue
			}

			switch x := res.Value.(type) {
			case *ua.DataChangeNotification:
				for _, item := range x.MonitoredItems {
					// // 打印值的类型
					data := item.Value.Value.Value()
					switch v := data.(type) {
					case bool:
						log.Printf("Type: Boolean, Value: %v", v)
					case int32:
						log.Printf("Type: Int32, Value: %v", v)
					case float64:
						log.Printf("Type: Float64, Value: %v", v)
					case string:
						log.Printf("Type: String, Value: %v", v)
					default:
						log.Printf("Unknown type: %T, Value: %v", v, v)
					}
					log.Printf("MonitoredItem with client handle %v = %v", item.ClientHandle, data)
				}

			case *ua.EventNotificationList:
				for _, item := range x.Events {
					log.Printf("Event for client handle: %v\n", item.ClientHandle)
					for i, field := range item.EventFields {
						log.Printf("%v: %v of Type: %T", eventFieldNames[i], field.Value(), field.Value())
					}
					log.Println()
				}

			default:
				log.Printf("what's this publish result? %T", res.Value)
			}
		}
	}
}

func valueRequest(nodeID *ua.NodeID) *ua.MonitoredItemCreateRequest {
	handle := uint32(42)
	return opcua.NewMonitoredItemCreateRequestWithDefaults(nodeID, ua.AttributeIDValue, handle)
}

func eventRequest(nodeID *ua.NodeID) (*ua.MonitoredItemCreateRequest, []string) {
	fieldNames := []string{"EventId", "EventType", "Severity", "Time", "Message"}
	selects := make([]*ua.SimpleAttributeOperand, len(fieldNames))

	for i, name := range fieldNames {
		selects[i] = &ua.SimpleAttributeOperand{
			TypeDefinitionID: ua.NewNumericNodeID(0, id.BaseEventType),
			BrowsePath:       []*ua.QualifiedName{{NamespaceIndex: 0, Name: name}},
			AttributeID:      ua.AttributeIDValue,
		}
	}

	wheres := &ua.ContentFilter{
		Elements: []*ua.ContentFilterElement{
			{
				FilterOperator: ua.FilterOperatorGreaterThanOrEqual,
				FilterOperands: []*ua.ExtensionObject{
					{
						EncodingMask: 1,
						TypeID: &ua.ExpandedNodeID{
							NodeID: ua.NewNumericNodeID(0, id.SimpleAttributeOperand_Encoding_DefaultBinary),
						},
						Value: ua.SimpleAttributeOperand{
							TypeDefinitionID: ua.NewNumericNodeID(0, id.BaseEventType),
							BrowsePath:       []*ua.QualifiedName{{NamespaceIndex: 0, Name: "Severity"}},
							AttributeID:      ua.AttributeIDValue,
						},
					},
					{
						EncodingMask: 1,
						TypeID: &ua.ExpandedNodeID{
							NodeID: ua.NewNumericNodeID(0, id.LiteralOperand_Encoding_DefaultBinary),
						},
						Value: ua.LiteralOperand{
							Value: ua.MustVariant(uint16(0)),
						},
					},
				},
			},
		},
	}

	filter := ua.EventFilter{
		SelectClauses: selects,
		WhereClause:   wheres,
	}

	filterExtObj := ua.ExtensionObject{
		EncodingMask: ua.ExtensionObjectBinary,
		TypeID: &ua.ExpandedNodeID{
			NodeID: ua.NewNumericNodeID(0, id.EventFilter_Encoding_DefaultBinary),
		},
		Value: filter,
	}

	handle := uint32(42)
	req := &ua.MonitoredItemCreateRequest{
		ItemToMonitor: &ua.ReadValueID{
			NodeID:       nodeID,
			AttributeID:  ua.AttributeIDEventNotifier,
			DataEncoding: &ua.QualifiedName{},
		},
		MonitoringMode: ua.MonitoringModeReporting,
		RequestedParameters: &ua.MonitoringParameters{
			ClientHandle:     handle,
			DiscardOldest:    true,
			Filter:           &filterExtObj,
			QueueSize:        10,
			SamplingInterval: 1.0,
		},
	}

	return req, fieldNames
}
