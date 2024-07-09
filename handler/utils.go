package handler

import (
	"os"
	"strconv"
	"ytdl_http/models"
)

func getDownloadedFiles() ([]models.Player, error) {
	entries, err := os.ReadDir("Downloads")
	if err != nil {
		return nil, err
	}

	var pl []models.Player
	var p models.Player
	for i := range entries {
		info, err := entries[i].Info()
		if err != nil {
			return nil, err
		}

		name := info.Name()

		switch name[len(name)-3:] {
		case "m4a":
			p.IsAudio = true
			p.Type = "audio/mpeg"
		case "mp4":
			p.IsAudio = false
			p.Type = "video/mp4"
		default:
			continue
		}

		p.Title = name
		p.ID = strconv.Itoa(i)
		p.Path = "/resource/" + name

		pl = append(pl, p)
	}

	return pl, nil
}
