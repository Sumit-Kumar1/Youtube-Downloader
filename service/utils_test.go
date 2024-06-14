package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isPlaylist(t *testing.T) {
	pl := "https://www.youtube.com/playlist?list=PLZHQObOWTQDN52m7Y21ePrTbvXkPaWVSg"
	vid := "https://www.youtube.com/watch?v=wI2sxYte6RA&list=RDCLAK5uy_n41V9iqmjro6caBDuDD1E4eWs5yTb5_OY&index=1"

	tests := []struct {
		name   string
		url    string
		expErr error
		want   bool
	}{
		{name: "playlist url", url: pl, want: true, expErr: nil},
		{name: "video url", url: vid, want: false, expErr: nil},
		{name: "invalid url", url: "www.lazy.com", want: false, expErr: errors.New("not a valid link")},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isPlayList, err := validateURL(tt.url)

			assert.Equalf(t, tt.expErr, err, "TEST[%d] Failed - %s", i, tt.name)
			assert.Equalf(t, tt.want, isPlayList, "TEST[%d] Failed - %s", i, tt.name)
		})
	}
}
