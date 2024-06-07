package main

import (
	"context"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"terminal/ent"
	"testing"
)

func TestName(t *testing.T) {
	ctx := context.Background()
	client, err := ent.Open("sqlite3", "terminal.db?cache=shared&mode=rwc&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	err = client.Folders.Create().SetLabel("root").Exec(ctx)
	if err != nil {
		log.Fatalf("failed creating folder: %v", err)
	}
	all, err := client.Folders.Query().All(ctx)
	if err != nil {
		log.Fatalf("failed querying folders: %v", err)
	}
	for _, v := range all {
		fmt.Println(v)
	}
}
