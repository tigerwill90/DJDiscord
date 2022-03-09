package config

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerate(t *testing.T) {
	cases := []struct {
		name string
		opts *TemplateOption
		want string
	}{
		{
			name: "generate config with default option",
			opts: &TemplateOption{
				Token:  "supersecrettokenyoushouldnotcommitongithub",
				UserID: 123456789,
			},
			want: `token = supersecrettokenyoushouldnotcommitongithub
owner = 123456789
prefix = "@mention"
game = "DEFAULT"
status = ONLINE
songinstatus=false
altprefix = "NONE"
success = "🎶"
warning = "💡"
error = "🚫"
loading = "⌚"
searching = "🔎"
help = help
npimages = false
stayinchannel = false
maxtime = 0
alonetimeuntilstop = 0
playlistsfolder = "Playlists"
updatealerts=true
lyrics.default = "A-Z Lyrics"
aliases {
  settings = [ status ]
  lyrics = []
  nowplaying = [ np, current ]
  play = []
  playlists = [ pls ]
  queue = [ list ]
  remove = [ delete ]
  scsearch = []
  search = [ ytsearch ]
  shuffle = []
  skip = [ voteskip ]
  prefix = [ setprefix ]
  setdj = []
  settc = []
  setvc = []
  forceremove = [ forcedelete, modremove, moddelete ]
  forceskip = [ modskip ]
  movetrack = [ move ]
  pause = []
  playnext = []
  repeat = []
  skipto = [ jumpto ]
  stop = []
  volume = [ vol ]
}
eval=false`,
		},
	}

	buf := bytes.NewBuffer(nil)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer buf.Reset()
			require.NoError(t, Generate(buf, tc.opts))
			assert.Equal(t, tc.want, buf.String())
		})
	}
}
