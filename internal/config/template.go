package config

type TemplateOption struct {
	Token  string
	UserID int
	Prefix string
	Game   string
	Status string
}

var tmpl = `token = {{.Token}}
owner = {{.UserID}}
prefix = "{{if .Prefix}}{{.Prefix}}{{else}}@mention{{end}}"
game = "{{if .Game}}{{.Game}}{{else}}DEFAULT{{end}}"
status = {{if .Status}}{{.Status}}{{else}}ONLINE{{end}}
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
eval=false`
