package logic

import (
	"context"
	"log"
	"terminal/ent"
	"terminal/ent/folders"
	"terminal/pkg/syncmapx"
	"terminal/termx"
)

// Logic struct
type Logic struct {
	Ctx    context.Context
	db     *ent.Client
	ptyMap *syncmapx.Map[string, termx.PtyX]
}

// NewApp creates a new App application struct
func NewApp() *Logic {
	l := &Logic{
		ptyMap: syncmapx.New[string, termx.PtyX](),
	}
	client, err := ent.Open("sqlite3", "terminal.db?cache=shared&mode=rwc&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	//initFold(client)
	l.db = client
	return l
}

func initFold(client *ent.Client) {
	count, err := client.Folders.Query().Where(folders.LabelEQ("root")).Count(context.Background())
	if err != nil {
		return
	}
	if count == 0 {
		err := client.Folders.Create().SetLabel("root").Exec(context.Background())
		if err != nil {
			log.Fatalf("failed creating folders resources: %v\n", err)
		}
	}
}
