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
	"time"

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

	ticker := time.NewTicker(time.Minute * 3)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				err = dg.UpdateStreamingStatus(0, "for Team 6657", "http://github.com/team6657/frchan")
				if err != nil {
					log.Println("Error setting bot status:", err)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	log.Println("Bot now running. Press CTRL+C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
	close(quit)
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

		var ok bool

		switch cmd {
		case "help":
			embfields := []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:   "help",
					Value:  "```md\n# Displays this page\n< usage：>\n\t~!help\n```",
					Inline: false,
				},
				&discordgo.MessageEmbedField{
					Name:   "team",
					Value:  "```md\n# Displays FRC team info\n< usage：>\n\t~!team <team #>\n```",
					Inline: false,
				},
				&discordgo.MessageEmbedField{
					Name:   "awards",
					Value:  "```md\n# Displays FRC team awards in competitions over the years\n< usage：>\n\t~!awards <team #>\n```",
					Inline: false,
				},
				&discordgo.MessageEmbedField{
					Name:   "events",
					Value:  "```md\n# Displays FRC team rankings in events\n< usage：>\n\t~!events <team #> [event #]\n```",
					Inline: false,
				},
			}

			emb = &discordgo.MessageEmbed{
				Author: &discordgo.MessageEmbedAuthor{
					Name: "FRChan Help Guide",
					URL:  "http://github.com/team6657/frchan",
				},
				Description: "List of commands for FRChan:\n",
				Fields:      embfields,
				Color:       14490723,
			}
		case "team":
			emb, ok = formatTeamInfo(getTeamInfo(params[0]))
		case "awards":
			emb, ok = formatTeamAwards(params[0], getTeamAwards(params[0]))
		case "events":
			emb, ok = formatTeamEventStatus(params)
		}

		if emb == nil || !ok {
			return
		}

		_, err := s.ChannelMessageSendEmbed(m.ChannelID, emb)
		if err != nil {
			log.Println("Error sending message:", err)
			return
		}
	}
}

func formatTeamEventStatus(args []string) (*discordgo.MessageEmbed, bool) {
	_, err := strconv.Atoi(args[0])
	if err != nil {
		return &discordgo.MessageEmbed{}, false
	}

	teamkey := args[0]

	team := getTeamInfo(teamkey)
	events := getTeamEventsSimple(teamkey)

	var formatted string
	var emb *discordgo.MessageEmbed

	authstr := team.Nickname + " (" + strconv.Itoa(team.TeamNumber) + ")"

	if len(authstr) > 40 {
		authstr = team.Name + " (" + strconv.Itoa(team.TeamNumber) + ")"
		if len(authstr) > 40 {
			authstr = "Team " + strconv.Itoa(team.TeamNumber)
		}
	}

	if len(args) > 1 {
		if en := args[1]; en != "" {

			n, err := strconv.Atoi(en)
			if err != nil {
				log.Println("Error converting string to int:", err)
				return &discordgo.MessageEmbed{}, false
			}

			if !(n >= len(events)) {
				e := events[n]
				tei := getTeamEventStatus(teamkey, e.Key)

				var playoffstr string

				if tei.Playoff == (TBAPlayoff{}) {
					playoffstr = "Team did not make it to playoffs"
				} else {
					playoffstr = "Level: " + tei.Playoff.Level +
						"\nW-L-T:\t" + strconv.Itoa(tei.Playoff.Record.Wins) + "-" + strconv.Itoa(tei.Playoff.Record.Losses) + "-" + strconv.Itoa(tei.Playoff.Record.Ties)
				}

				embfields := []*discordgo.MessageEmbedField{
					&discordgo.MessageEmbedField{
						Name:   "Event",
						Value:  e.Name + " (" + strconv.Itoa(e.Year) + ")",
						Inline: true,
					},
					&discordgo.MessageEmbedField{
						Name:   "Location",
						Value:  e.City + ", " + e.StateProv + " " + e.Country,
						Inline: true,
					},
					&discordgo.MessageEmbedField{
						Name:   "Date",
						Value:  e.StartDate + " - " + e.EndDate,
						Inline: true,
					},
					&discordgo.MessageEmbedField{
						Name: "Qualifiers",
						Value: strconv.Itoa(tei.Qual.Ranking.Rank) + "/" + strconv.Itoa(tei.Qual.NumTeams) +
							"\nW-L-T:\t" + strconv.Itoa(tei.Qual.Ranking.Record.Wins) + "-" + strconv.Itoa(tei.Qual.Ranking.Record.Losses) + "-" + strconv.Itoa(tei.Qual.Ranking.Record.Ties),
						Inline: false,
					},
					&discordgo.MessageEmbedField{
						Name:   "Playoffs",
						Value:  playoffstr,
						Inline: false,
					},
				}

				emb = &discordgo.MessageEmbed{
					Author: &discordgo.MessageEmbedAuthor{
						Name: authstr,
						URL:  team.Website,
					},
					Description: team.Name + "\nNick: " + team.Nickname + "\nTeam event info:\n",
					Fields:      embfields,
					Color:       14490723,
					Footer: &discordgo.MessageEmbedFooter{
						Text: "Info provided by The Blue Alliance",
					},
				}
			}
		}
	} else {
		for i, e := range events {
			formatted +=
				"\n```md\n" + strconv.Itoa(i) + ". <" + e.Name +
					" (" + strconv.Itoa(e.Year) + ")>\n\t" +
					e.City + ", " +
					e.StateProv + " " +
					e.Country +
					"```"
		}

		emb = &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				Name: authstr,
				URL:  team.Website,
			},
			Description: team.Name + "\nNick: " + team.Nickname + "\n" + formatted,
			Color:       14490723,
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Info provided by The Blue Alliance",
			},
		}
	}

	return emb, true
}

func getTeamEventStatus(teamkey, eventkey string) TBATeamEventStatus {
	_, err := strconv.Atoi(teamkey)
	if err != nil {
		return TBATeamEventStatus{}
	}

	var info TBATeamEventStatus

	req, err := http.NewRequest("GET", "http://www.thebluealliance.com/api/v3/team/frc"+teamkey+"/event/"+eventkey+"/status", nil)
	if err != nil {
		log.Println("Failed event info request:", err)
		return TBATeamEventStatus{}
	}

	req.Header.Set("X-TBA-Auth-Key", TBAToken)

	res, err := client.Do(req)
	if err != nil {
		log.Println("Failed event info request:", err)
		return TBATeamEventStatus{}
	}

	err = json.NewDecoder(res.Body).Decode(&info)
	if err != nil {
		log.Println("Failed event info request:", err)
		return TBATeamEventStatus{}
	}

	return info
}

func getEventRankings(eventkey string) TBAEventRankings {
	var info TBAEventRankings

	req, err := http.NewRequest("GET", "http://www.thebluealliance.com/api/v3/event/"+eventkey+"/rankings", nil)
	if err != nil {
		log.Println("Failed event info request:", err)
		return TBAEventRankings{}
	}

	req.Header.Set("X-TBA-Auth-Key", TBAToken)

	res, err := client.Do(req)
	if err != nil {
		log.Println("Failed event info request:", err)
		return TBAEventRankings{}
	}

	err = json.NewDecoder(res.Body).Decode(&info)
	if err != nil {
		log.Println("Failed event info request:", err)
		return TBAEventRankings{}
	}

	return info
}

func getTeamEventsSimple(teamkey string) TBATeamEventsSimple {
	_, err := strconv.Atoi(teamkey)
	if err != nil {
		return TBATeamEventsSimple{}
	}

	var info TBATeamEventsSimple

	req, err := http.NewRequest("GET", "http://www.thebluealliance.com/api/v3/team/frc"+teamkey+"/events/simple", nil)
	if err != nil {
		log.Println("Failed team info request:", err)
		return TBATeamEventsSimple{}
	}

	req.Header.Set("X-TBA-Auth-Key", TBAToken)

	res, err := client.Do(req)
	if err != nil {
		log.Println("Failed team info request:", err)
		return TBATeamEventsSimple{}
	}

	err = json.NewDecoder(res.Body).Decode(&info)
	if err != nil {
		log.Println("Failed team info request:", err)
		return TBATeamEventsSimple{}
	}

	return info
}

func formatTeamAwards(teamkey string, info TBATeamAwards) (*discordgo.MessageEmbed, bool) {

	if info == nil {
		return &discordgo.MessageEmbed{}, false
	}

	_, err := strconv.Atoi(teamkey)
	if err != nil {
		return &discordgo.MessageEmbed{}, false
	}

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

	authstr := teaminfo.Nickname + " (" + strconv.Itoa(teaminfo.TeamNumber) + ")"

	if len(authstr) > 40 {
		authstr = teaminfo.Name + " (" + strconv.Itoa(teaminfo.TeamNumber) + ")"
		if len(authstr) > 40 {
			authstr = "Team " + strconv.Itoa(teaminfo.TeamNumber)
		}
	}

	emb := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name: authstr,
			URL:  teaminfo.Website,
		},
		Description: formatted,
		Color:       14490723,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Info provided by The Blue Alliance",
		},
	}

	return emb, true
}

func getTeamAwards(teamkey string) TBATeamAwards {

	_, err := strconv.Atoi(teamkey)
	if err != nil {
		return TBATeamAwards{}
	}

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

func formatTeamInfo(info TBATeam) (*discordgo.MessageEmbed, bool) {
	if info == (TBATeam{}) {
		return &discordgo.MessageEmbed{}, false
	}

	motto := info.Motto
	if motto == "" {
		motto = "No motto found"
	}

	authstr := info.Nickname + " (" + strconv.Itoa(info.TeamNumber) + ")"

	if len(authstr) > 40 {
		authstr = info.Name + " (" + strconv.Itoa(info.TeamNumber) + ")"
		if len(authstr) > 40 {
			authstr = "Team " + strconv.Itoa(info.TeamNumber)
		}
	}

	embfields := []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "Name",
			Value:  info.Name,
			Inline: false,
		},
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
			Name: authstr,
			URL:  info.Website,
		},
		Description: "\nInfo for team " + info.Key + ":\n",
		Fields:      embfields,
		Color:       14490723,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Info provided by The Blue Alliance",
		},
	}

	return emb, true
}

func getTeamInfo(teamkey string) TBATeam {
	_, err := strconv.Atoi(teamkey)
	if err != nil {
		return TBATeam{}
	}

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
