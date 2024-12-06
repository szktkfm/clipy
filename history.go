package clipy

import (
	"fmt"
	"log"
	"sync"

	"github.com/ktr0731/go-fuzzyfinder"
)

func update(c *[]Clipboard, db *Repository, offset, limit int) {
	page, err := db.read(offset, limit)
	if err != nil {
		log.Fatal(err)
	}

	if len(page) == 0 {
		return
	}

	*c = append(*c, page...)
}

func ListHistories(dbPath string) error {
	db, err := NewRepository(dbPath)
	if err != nil {
		return err
	}

	offset, limit := 0, 20
	total, err := db.count()
	if err != nil {
		return err
	}

	history, err := db.read(offset, limit)
	if err != nil {
		return err
	}

	// Channel to signal when more items are needed
	loadMore := make(chan struct{}, 1)
	go func() {
		for range loadMore {
			if offset < int(total) {
				offset += limit
				update(&history, db, offset, limit)
			}
		}
	}()

	var mu sync.Mutex
	// fuzzy-finder
	idx, err := fuzzyfinder.FindMulti(
		&history,
		func(i int) string {
			return history[i].ClipText
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				mu.Lock()
				defer mu.Unlock()
				loadMore <- struct{}{}
				return ""
			}
			return fmt.Sprint(history[i].ClipText)
		}),
		fuzzyfinder.WithHotReloadLock(&mu),
	)
	if err != nil {
		return err
	}

	// Print selected item
	if len(idx) > 0 {
		fmt.Print(string(history[idx[0]].ClipText))
	}

	return nil
}
