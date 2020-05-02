package main

import (
	"fmt"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/go-routeros/routeros"
	"github.com/go-routeros/routeros/proto"
)

func main() {
	reply := routeros.Reply{
		Re: []*proto.Sentence{
			{
				Word: "!re",
				List: []proto.Pair{
					{
						Key:   "name",
						Value: "test script",
					},
				},
			},
		},
	}
	var script client.Script
	client.Unmarshal(reply, &script)

	fmt.Printf("%#v", script)
}
