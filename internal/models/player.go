package models

const (
	DirPath   = "./Downloads/"
	m4aExt    = "m4a"
	mp3Ext    = "mp3"
	mp4Ext    = "mp4"
	audioType = "audio/mpeg"
	videoType = "video/mp4"
)

type Player struct {
	ID      string
	Title   string
	IsAudio bool
	Path    string
	Type    string
}

func (p *Player) FillByName(name string) {
	if p == nil || name == "" || len(name) < 3 {
		return
	}

	ext := name[len(name)-3:]
	switch ext {
	case m4aExt, mp3Ext:
		p.IsAudio = true
		p.Type = audioType
	case mp4Ext:
		p.IsAudio = false
		p.Type = videoType
	default:
	}

	p.Title = name
	p.Path = "/resource/" + name
}
