package config

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestGenerateOption(t *testing.T) {
	cases := []struct {
		name string
		opts *TemplateOption
		want []string
	}{
		{
			name: "all options are set",
			opts: &TemplateOption{
				Token:  "supersecrettoken",
				UserID: 123456789,
				Prefix: "!",
				Game:   "Minecraft",
				Status: "DND",
			},
			want: []string{
				"token = supersecrettoken",
				"owner = 123456789",
				`prefix = "!"`,
				`game = "Minecraft"`,
				"status = DND",
			},
		},
		{
			name: "values keep their special characters",
			opts: &TemplateOption{
				Token:  "a&b",
				UserID: 1,
				Prefix: ">",
				Game:   "Rock'n'Roll <live>",
			},
			want: []string{
				"token = a&b",
				`prefix = ">"`,
				`game = "Rock'n'Roll <live>"`,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			require.NoError(t, Generate(buf, tc.opts))
			for _, want := range tc.want {
				assert.Contains(t, buf.String(), want)
			}
		})
	}
}
