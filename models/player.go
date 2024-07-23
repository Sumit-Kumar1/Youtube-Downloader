package models

const (
	m4aExt    = "m4a"
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
	case m4aExt:
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
