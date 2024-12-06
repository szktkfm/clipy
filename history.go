package clipy

import (
	"fmt"
	"log"
	"sync"

	"github.com/ktr0731/go-fuzzyfinder"
)

func updateClipboards(c *[]Clipboard, db *Repository, offset, limit int) {
	newCbPage := db.read(offset, limit)
	if len(newCbPage) == 0 {
		return
	}
	*c = append(*c, newCbPage...)
}

func ListHistories() {

	db := NewRepository("./test.db")

	var mu sync.Mutex

	offset := 0
	limit := 20

	count := db.count()
	clipboards := db.read(offset, limit)

	noitem := make(chan *byte, 1)
	go func(msg chan *byte) {
		for range msg {
			if offset < int(count) {
				offset += limit
				updateClipboards(&clipboards, db, offset, limit)
			}
		}
	}(noitem)

	// fuzzy-finder
	idx, err := fuzzyfinder.FindMulti(
		&clipboards,
		func(i int) string {
			return clipboards[i].Bytes
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				mu.Lock()
				defer mu.Unlock()
				noitem <- nil
				return ""
			}
			return fmt.Sprint(clipboards[i].Bytes)
		}),
		fuzzyfinder.WithHotReloadLock(&mu),
	)
	if err != nil {
		log.Fatal(err)
	}
	selectedItem := string(clipboards[idx[0]].Bytes)
	fmt.Print(selectedItem)
}
