package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	BotToken string
	TBAToken string
	client   = &http.Client{}
)

func init() {
	flag.StringVar(&BotToken, "bot", "", "Bot token")
	flag.StringVar(&TBAToken, "tba", "", "TBA token")
	flag.Parse()
}

func main() {
	dg, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatal("Error creating Discord session:", err)
		return
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		log.Fatal("Error opening connection:", err)
		return
	}

	log.Println("Bot now running. Press CTRL+C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "~!") {
		in := strings.Split(strings.TrimPrefix(m.Content, "~!"), " ")
		cmd := in[0]
		var params []string
		if len(in) > 1 {
			params = in[1:]
		}

		var emb *discordgo.MessageEmbed

		switch cmd {
		case "team":
			emb = formatTeamInfo(getTeamInfo(params[0]))
		case "awards":
			emb = formatTeamAwards(params[0], getTeamAwards(params[0]))
		}

		s.ChannelMessageSendEmbed(m.ChannelID, emb)
	}
}

func formatTeamAwards(teamkey string, info TBATeamAwards) *discordgo.MessageEmbed {

	list := make(map[int][]string)

	for _, award := range info {
		list[award.Year] = append(list[award.Year], award.Name)
	}

	var formatted string

	for k, v := range list {
		formatted += "\n" + strconv.Itoa(k) + ":\n"

		for _, a := range v {
			formatted += "\t" + a + "\n"
		}
	}

	teaminfo := getTeamInfo(teamkey)

	emb := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name: "Awards for " + teaminfo.Name + " (" + strconv.Itoa(teaminfo.TeamNumber) + ")",
			URL:  teaminfo.Website,
		},
		Description: formatted,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Info provided by The Blue Alliance",
		},
	}

	return emb
}

func getTeamAwards(teamkey string) TBATeamAwards {
	var info TBATeamAwards

	req, err := http.NewRequest("GET", "http://www.thebluealliance.com/api/v3/team/frc"+teamkey+"/awards", nil)
	if err != nil {
		log.Println("Failed team info request:", err)
		return TBATeamAwards{}
	}

	req.Header.Set("X-TBA-Auth-Key", TBAToken)

	res, err := client.Do(req)
	if err != nil {
		log.Println("Failed team info request:", err)
		return TBATeamAwards{}
	}

	err = json.NewDecoder(res.Body).Decode(&info)
	if err != nil {
		log.Println("Failed team info request:", err)
		return TBATeamAwards{}
	}

	return info
}

func formatTeamInfo(info TBATeam) *discordgo.MessageEmbed {
	motto := info.Motto
	if motto == "" {
		motto = "No motto found"
	}

	embfields := []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "Nickname",
			Value:  info.Nickname,
			Inline: false,
		},
		&discordgo.MessageEmbedField{
			Name:   "Motto",
			Value:  motto,
			Inline: false,
		},
		&discordgo.MessageEmbedField{
			Name:   "Rookie Year",
			Value:  strconv.Itoa(info.RookieYear),
			Inline: false,
		},
		&discordgo.MessageEmbedField{
			Name:   "Location",
			Value:  info.City + ", " + info.SateProv,
			Inline: false,
		},
	}

	emb := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name: info.Name + " (" + strconv.Itoa(info.TeamNumber) + ")",
			URL:  info.Website,
		},
		Description: "\nInfo for team " + info.Key + ":\n",
		Fields:      embfields,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Info provided by The Blue Alliance",
		},
	}

	return emb
}

func getTeamInfo(teamkey string) TBATeam {
	var info TBATeam

	req, err := http.NewRequest("GET", "http://www.thebluealliance.com/api/v3/team/frc"+teamkey, nil)
	if err != nil {
		log.Println("Failed team info request:", err)
		return TBATeam{}
	}

	req.Header.Set("X-TBA-Auth-Key", TBAToken)

	res, err := client.Do(req)
	if err != nil {
		log.Println("Failed team info request:", err)
		return TBATeam{}
	}

	err = json.NewDecoder(res.Body).Decode(&info)
	if err != nil {
		log.Println("Failed team info request:", err)
		return TBATeam{}
	}

	return info
}
