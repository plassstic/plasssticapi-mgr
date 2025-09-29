package utils

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
	. "plassstic.tech/gopkg/golang-manager/internal/logic/api"
)

func escapeMarkdownV2(text string) string {
	specialChars := []string{
		"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+",
		"-", "=", "|", "{", "}", ".", "!",
	}

	for _, char := range specialChars {
		text = strings.ReplaceAll(text, char, "\\"+char)
	}

	return text
}

func clamp(s string, to ...interface{}) string {
	bound := 60
	if len(to) != 0 {
		if v, ok := to[0].(int); ok {
			bound = v
		}
	}
	if len(s) > bound {
		return s[:(bound-3)] + "\\.\\.\\."
	}
	return s
}

func GetFormattedPlayer(player *PlayerSI) string {
	trackName := strings.ToLower(escapeMarkdownV2(player.Item.Name))

	artists := strings.Join(lo.Map(player.Item.Artists, func(item ArtistSI, index int) string { return strings.ToLower(escapeMarkdownV2(item.Name)) }), ", ")

	albumName := strings.ToLower(escapeMarkdownV2(player.Item.Album.Name))
	albYear := strings.Split(player.Item.Album.ReleaseDate, "-")[0]
	minutesPlayed := player.ProgressMs / 60000
	secondsPlayed := (player.ProgressMs % 60000) / 1000
	minutesTotal := player.Item.DurationMs / 60000
	secondsTotal := (player.Item.DurationMs % 60000) / 1000

	getDevEmoji := func(devType string) string {
		switch devType {
		case "Smartphone":
			return "📱"
		case "Computer":
			return "💻"
		default:
			return "✨"
		}
	}

	getVolEmoji := func(volPercent int) string {
		if volPercent < 20 {
			return "🔈"
		}
		if volPercent < 40 {
			return "🔉"
		}
		if volPercent <= 100 {
			return "🔊"
		}
		return "🥚"
	}

	var playlist string
	if player.Context != nil {
		playlist = escapeMarkdownV2(player.Context.ExternalURLs["spotify"])
		if strings.Contains(playlist, "/collection/tracks") {
			playlist = ""
		}
	}

	var link string
	if len(player.Item.ExternalURLs) > 0 {
		if spotifyURL, ok := player.Item.ExternalURLs["spotify"]; ok {
			link = escapeMarkdownV2(spotifyURL)
		}
	}

	artStr := []string{">", "•", "🎭", fmt.Sprintf("_%s_", clamp(artists))}
	trackStr := []string{">", "•", "🎶", fmt.Sprintf("_*%s*_", clamp(trackName))}
	albumStr := []string{">", "•", "💽", fmt.Sprintf("_*%s*", clamp(albumName, 49)), fmt.Sprintf("\\(%s\\)_", albYear)}

	var statusStr []string
	if player.IsPlaying {
		statusStr = []string{
			">➜",
			getDevEmoji(player.Device.Type),
			"_streaming on:",
			fmt.Sprintf("*%s*_", strings.ToLower(player.Device.Name)),
		}
	} else {
		statusStr = []string{
			">➜",
			getDevEmoji(player.Device.Type),
			"streaming paused on:",
			fmt.Sprintf("*%s*", strings.ToLower(player.Device.Name)),
		}
	}

	volStr := []string{
		">", "•",
		getVolEmoji(player.Device.VolumePercent),
		"_volume:",
		fmt.Sprintf("*%d%%*_", player.Device.VolumePercent),
	}

	calcLength := func(arr []string) int {
		total := 0
		for _, s := range arr {
			total += len(s)
		}
		return total + (len(arr)-1)*2
	}

	largest := max(calcLength(trackStr), calcLength(albumStr)-10, calcLength(statusStr))

	a := int(float64(largest) * 0.75)
	borderStr := fmt.Sprintf(">✦\t\t%s%s%s%s%s",
		strings.Repeat("┅", a/7),
		strings.Repeat("┉", a/9),
		strings.Repeat("┉", a/9),
		strings.Repeat("┄", a/9),
		strings.Repeat("┈", a/11),
	)
	border := []string{borderStr}

	progressBarSpace := int(float64(a) * 0.4)
	plRatio := float64(player.ProgressMs) / float64(player.Item.DurationMs)
	npRatio := 1.0 - plRatio

	playChar := "■"
	if player.IsPlaying {
		playChar = "◉"
	}

	progBar := "`" +
		strings.Repeat("▬", int(plRatio*float64(progressBarSpace))) +
		playChar +
		strings.Repeat("─", int(npRatio*float64(progressBarSpace))) +
		"`"

	progressStr := []string{
		">", "•", "🎧",
		fmt.Sprintf("_%02d:%02d_", minutesPlayed, secondsPlayed),
		progBar,
		fmt.Sprintf("_%02d:%02d_", minutesTotal, secondsTotal),
	}

	hasDynamic := link != "" || playlist != ""

	message := [][]string{statusStr, volStr, border, artStr, trackStr, albumStr, border, progressStr}

	if hasDynamic {
		message = append(message, border)
		if link != "" {
			message = append(message, []string{">", "•", "🖇", fmt.Sprintf("[*›\\___track__*](%s)", link)})
		}
		if playlist != "" {
			message = append(message, []string{">", "•", "📃", fmt.Sprintf("[*›\\___playlist__*](%s)", playlist)})
		}
		message[len(message)-1] = append(message[len(message)-1], "\n")
	}

	var lines []string
	for _, part := range message {
		lines = append(lines, strings.Join(part, "\t\t"))
	}

	return strings.Join(lines, "\n")
}
