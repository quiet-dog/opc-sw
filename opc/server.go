package opc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/id"
	"github.com/gopcua/opcua/ua"
)

type OpcClient struct {
	Endpoint string
	Duration time.Duration
	client   *opcua.Client
	sub      *opcua.Subscription
	ctx      context.Context
	gateway  chan Data
}

type NodeId struct {
	ID   uint64
	Node string
}

type Data struct {
	ID         uint64
	DataType   string
	Value      interface{}
	SourceTime time.Time
	Param      string
}

func (o *OpcClient) connect() {
	ctx, cancel := context.WithTimeout(context.Background(), o.Duration)
	defer cancel()

	endpoints, err := opcua.GetEndpoints(ctx, o.Endpoint)
	if err != nil {
		panic(err)
	}
	ep, err := opcua.SelectEndpoint(endpoints, "", ua.MessageSecurityModeFromString(""))
	if err != nil {
		log.Fatal(err)
	}
	ep.EndpointURL = o.Endpoint

	opts := []opcua.Option{
		opcua.SecurityPolicy(""),
		opcua.SecurityModeString(""),
		opcua.CertificateFile(""),
		opcua.PrivateKeyFile(""),
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

	o.client = c
	notifyCh := make(chan *opcua.PublishNotificationData)

	sub, err := c.Subscribe(ctx, &opcua.SubscriptionParameters{
		Interval: opcua.DefaultSubscriptionInterval,
	}, notifyCh)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Cancel(ctx)
	// log.Printf("Created subscription with id %v", sub.SubscriptionID)

	// id, err := ua.ParseNodeID(*nodeID)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// var miCreateRequest *ua.MonitoredItemCreateRequest
	// var eventFieldNames []string

	// miCreateRequest = valueRequest(id)
	// res, err := sub.Monitor(ctx, ua.TimestampsToReturnBoth, miCreateRequest)
	// if err != nil || res.Results[0].StatusCode != ua.StatusOK {
	// 	log.Fatal(err)
	// }
	o.sub = sub

	for {
		select {
		case <-ctx.Done():
			fmt.Println("opc client connect done")
			return
		case res := <-notifyCh:
			if res.Error != nil {
				log.Print(res.Error)
				continue
			}

			switch x := res.Value.(type) {
			case *ua.DataChangeNotification:
				for _, item := range x.MonitoredItems {
					data := item.Value.Value.Value()
					log.Printf("MonitoredItem with client handle %v = %v", item.ClientHandle, data)
					if item.Value != nil {
						data := Data{
							ID:         uint64(item.ClientHandle),
							DataType:   item.Value.Value.Type().String(),
							Value:      item.Value.Value.Value(),
							SourceTime: item.Value.SourceTimestamp,
						}
						// 判断gateway是否关闭
						o.gateway <- data
					}
				}

			case *ua.EventNotificationList:
				// for _, item := range x.Events {
				// 	log.Printf("Event for client handle: %v\n", item.ClientHandle)
				// 	for i, field := range item.EventFields {
				// 		log.Printf("%v: %v of Type: %T", eventFieldNames[i], field.Value(), field.Value())
				// 	}
				// 	log.Println()
				// }

			default:
				log.Printf("what's this publish result? %T", res.Value)
			}
		}
	}
}

func (o *OpcClient) AddNodeID(n NodeId) error {
	id, err := ua.ParseNodeID(n.Node)
	if err != nil {
		return err
	}
	var miCreateRequest *ua.MonitoredItemCreateRequest
	// var eventFieldNames []string
	// if isEvent {
	// 	miCreateRequest, eventFieldNames = eventRequest(id)
	// } else {
	miCreateRequest = o.valueRequest(id, uint32(n.ID))

	// }
	res, err := o.sub.Monitor(o.ctx, ua.TimestampsToReturnBoth, miCreateRequest)
	if err != nil || res.Results[0].StatusCode != ua.StatusOK {
		return err
	}
	log.Printf("Added new monitored item for NodeID: %s", n.Node)
	return nil
}

func (o *OpcClient) valueRequest(nodeID *ua.NodeID, handle uint32) *ua.MonitoredItemCreateRequest {
	// handle := uint32(42)
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
