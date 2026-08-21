package main

import (
	"cmp"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"strings"
	"syscall"

	"github.com/tigerwill90/djdiscord/internal/build"
	"github.com/tigerwill90/djdiscord/internal/config"
	"github.com/tigerwill90/djdiscord/internal/exec"
)

const configFile = "config.txt"

var validStatus = []string{"ONLINE", "IDLE", "DND", "INVISIBLE"}

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
	token = cmp.Or(token, os.Getenv("DJ_DISCORD_TOKEN"))
	if token == "" {
		exitErr(1, errors.New("a discord token is required"))
	}

	if owner <= 0 {
		userID := os.Getenv("DJ_DISCORD_USER_ID")
		if userID == "" {
			exitErr(1, errors.New("a discord user id is required"))
		}
		var err error
		owner, err = strconv.Atoi(userID)
		if err != nil {
			exitErr(1, fmt.Errorf("invalid discord user id: %w", err))
		}
		if owner <= 0 {
			exitErr(1, errors.New("a discord user id must be a positive number"))
		}
	}

	prefix = cmp.Or(prefix, os.Getenv("DJ_DISCORD_BOT_PREFIX"))
	game = cmp.Or(game, os.Getenv("DJ_DISCORD_BOT_GAME"))
	status = cmp.Or(status, os.Getenv("DJ_DISCORD_BOT_STATUS"))

	if status != "" && !slices.Contains(validStatus, status) {
		exitErr(1, fmt.Errorf("valid values for status are %s", strings.Join(validStatus, ", ")))
	}

	if err := writeConfig(); err != nil {
		exitErr(1, err)
	}

	cmd := exec.New(build.Version)
	if err := cmd.Start(); err != nil {
		exitErr(1, fmt.Errorf("unable to start bot: %w", err))
	}

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	select {
	case s := <-sig:
		// signal.Stop restores the default handler, so a second signal stops
		// this program at once.
		signal.Stop(sig)
		fmt.Fprintf(os.Stderr, "received %s, stopping JMusicBot\n", s)
		if err := cmd.Kill(toSignal(s)); err != nil {
			exitErr(1, err)
		}
	case err := <-cmd.Wait():
		if err != nil {
			exitErr(1, err)
		}
	}
}

// writeConfig generates the JMusicBot configuration file. The file holds the
// discord token, hence the 0600 mode.
func writeConfig() error {
	f, err := os.OpenFile(configFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	if err := config.Generate(f, &config.TemplateOption{
		Token:  token,
		UserID: owner,
		Prefix: prefix,
		Game:   game,
		Status: status,
	}); err != nil {
		f.Close()
		return err
	}

	return f.Close()
}

// toSignal unwraps s for Kill. signal.Notify always delivers a syscall.Signal
// on unix, so the fallback is defensive.
func toSignal(s os.Signal) syscall.Signal {
	if sig, ok := s.(syscall.Signal); ok {
		return sig
	}
	return syscall.SIGTERM
}

func exitErr(code int, err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(code)
}
