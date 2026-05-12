//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yong/doc-generation-mcp-server/examples/eino/shared"
)

func main() {
	ctx := context.Background()
	client, err := shared.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	tools, err := client.ListTools(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Connected to %s\n\n", shared.MCPServerURL())
	for _, item := range tools {
		fmt.Printf("- %s: %s\n", item.Name, item.Description)
	}
}
