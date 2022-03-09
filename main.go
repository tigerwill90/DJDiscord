package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/tigerwill90/djdiscord/internal/build"
	"github.com/tigerwill90/djdiscord/internal/config"
	"github.com/tigerwill90/djdiscord/internal/exec"
	"os"
	"os/signal"
	"strconv"
)

var (
	token  string
	owner  int
	prefix string
	game   string
	status string
)

func init() {
	flag.StringVar(&token, "token", "", "set discord secret token")
	flag.IntVar(&owner, "owner", -1, "set discord user id")
	flag.StringVar(&prefix, "prefix", "", "set the prefix for the bot")
	flag.StringVar(&game, "game", "", "modify the default game for the bot")
	flag.StringVar(&status, "status", "", "modify the default status for the bot")
	flag.Parse()
}

func main() {
	if token == "" {
		if token = os.Getenv("DJ_DISCORD_TOKEN"); token == "" {
			exitErr(1, errors.New("a discord token is required"))
		}
	}
	if owner <= 0 {
		var err error
		owner, err = strconv.Atoi(os.Getenv("DJ_DISCORD_USER_ID"))
		if err != nil {
			exitErr(1, fmt.Errorf("invalid discord user id: %s", err))
		}
		if owner <= 0 {
			exitErr(1, errors.New("a discord user id is required"))
		}
	}

	if prefix == "" {
		prefix = os.Getenv("DJ_DISCORD_BOT_PREFIX")
	}
	if game == "" {
		prefix = os.Getenv("DJ_DISCORD_BOT_GAME")
	}
	if status == "" {
		status = os.Getenv("DJ_DISCORD_BOT_STATUS")
	}
	if status != "" && !sliceContain(status, []string{"ONLINE", "IDLE", "DND", "INVISIBLE"}) {
		exitErr(1, errors.New("valid values for status are ONLINE, IDLE, DND, INVISIBLE"))
	}

	f, err := os.OpenFile("config.txt", os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		exitErr(1, err)
	}

	if err := config.Generate(f, &config.TemplateOption{
		Token:  token,
		UserID: owner,
		Prefix: prefix,
		Game:   game,
	}); err != nil {
		exitErr(1, err)
	}

	if err := f.Close(); err != nil {
		exitErr(1, err)
	}

	cmd := exec.New(build.Version)
	if err := cmd.Start(); err != nil {
		exitErr(1, fmt.Errorf("unable to start bot: %s", err))
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)

	select {
	case <-sig:
		if err := cmd.Kill(); err != nil {
			exitErr(1, err)
		}
	case err := <-cmd.Wait():
		if err != nil {
			exitErr(1, err)
		}
	}
}

func exitErr(code int, err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(code)
}

func sliceContain(s string, slice []string) bool {
	for _, e := range slice {
		if s == e {
			return true
		}
	}
	return false
}
