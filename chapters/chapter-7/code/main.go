package main

import (
	"fmt"
	"io"
	"os"
	"slices"
)

type Team struct {
	Name    string
	Players []string
}

type League struct {
	Teams []Team
	Wins  map[string]int
}

func NewLeague(teams []Team) League {
	league := League{
		Teams: teams,
		Wins:  make(map[string]int, len(teams)),
	}
	for _, v := range teams {
		league.Wins[v.Name] = 0
	}
	return league
}

func (l League) ValidateTeam(teamName string) {
	_, ok := l.Wins[teamName]
	if !ok {
		l.Wins[teamName] = 0
	}
}

func (l League) MatchResult(team1 string, score1 int, team2 string, score2 int) {
	l.ValidateTeam(team1)
	l.ValidateTeam(team2)

	switch {
	case score1 > score2:
		l.Wins[team1] += 1
	case score2 > score1:
		l.Wins[team2] += 1
	}
}

// Descending order
func (l League) Ranking() []string {
	// Build list of team names
	teams := make([]string, 0, len(l.Teams))
	for _, t := range l.Teams {
		teams = append(teams, t.Name)
	}

	slices.SortFunc(teams, func(t1, t2 string) int {
		return l.Wins[t2] - l.Wins[t1]
	})

	return teams
}

type Ranker interface {
	Ranking() []string
}

func RankPrinter(r Ranker, w io.Writer) {
	rankings := r.Ranking()

	for i, v := range rankings {
		line := fmt.Sprintf("%v. %v\n", i+1, v)
		w.Write([]byte(line))
	}
}

func main() {
	teams := []Team{
		Team{
			Name: "Markham",
			Players: []string{
				"Marcus",
				"Emma",
			},
		},
		Team{
			Name: "Sauga",
			Players: []string{
				"Olly",
				"Bony",
			},
		},
		Team{
			Name: "Foo",
			Players: []string{
				"a",
				"b",
			},
		},
	}

	league := NewLeague(teams)

	league.MatchResult("Markham", 10, "Foo", 12)
	league.MatchResult("Sauga", 10, "Foo", 12)
	league.MatchResult("Sauga", 10, "Markham", 12)

	rankedTeams := league.Ranking()

	fmt.Println(league)
	fmt.Println(rankedTeams)

	RankPrinter(league, os.Stdout)
}
