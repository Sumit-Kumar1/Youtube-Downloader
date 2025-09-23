package handler

import (
	"os"
	"strconv"

	"ytdl_http/internal/models"
)

func getDownloadedFiles() ([]models.Player, error) {
	entries, err := os.ReadDir(models.DirDownloads)
	if err != nil {
		return nil, err
	}

	var (
		pl = make([]models.Player, 0)
		p  models.Player
	)

	for i := range entries {
		info, err := entries[i].Info()
		if err != nil {
			return nil, err
		}

		p.FillByName(info.Name())

		switch p.Type {
		case models.TypeAudio, models.TypeVideo:
			p.ID = strconv.Itoa(i)
			pl = append(pl, p)
		}
	}

	return pl, nil
}
