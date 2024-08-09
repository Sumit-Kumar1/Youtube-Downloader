package handler

import (
	"os"
	"strconv"

	"ytdl_http/internal/models"
)

func getDownloadedFiles() ([]models.Player, error) {
	entries, err := os.ReadDir("Downloads")
	if err != nil {
		return nil, err
	}

	var (
		pl []models.Player
		p  models.Player
	)

	for i := range entries {
		info, err := entries[i].Info()
		if err != nil {
			return nil, err
		}

		p.FillByName(info.Name())

		p.ID = strconv.Itoa(i)
		pl = append(pl, p)
	}

	return pl, nil
}
