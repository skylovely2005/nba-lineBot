package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const DATE_TIME_LAYOUT = "01 月02 日 15:04"

func GetNBATodayData() (*AutoGenerated, error) {
	nbaAPIURL, found := _config.Source["nba"]
	if !found {
		log.Print("config url nil")
	}
	resp, err := http.Get(nbaAPIURL)
	if err != nil {
		log.Printf("error: get error %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("error: get fail %s", resp.Body)
		return nil, fmt.Errorf("status code error %v", resp.Body)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error: ReadAll error %v", err)
		return nil, err
	}
	data := AutoGenerated{}
	json.Unmarshal(body, &data)
	return &data, err
}

func (data *AutoGenerated) ParseToMessage() string {
	message := "     主隊 : 客隊\n"
	for index, val := range data.Payload.Date.Games {
		var gameInfo string
		homeScore := val.Boxscore.HomeScore
		awayScore := val.Boxscore.AwayScore
		homeTeamName := val.HomeTeam.Profile.Name
		awayTeamName := val.AwayTeam.Profile.Name
		status := val.Boxscore.Status
		switch status {
		case "1": // 未開賽
			gameTimeStr := UtcMillis2TimeString(val.Profile.UtcMillis, DATE_TIME_LAYOUT)
			gameInfo = fmt.Sprintf("未開賽 | %s ", gameTimeStr)
		default: //2: 比賽中, 3: 結束
			gameInfo = fmt.Sprintf(" %3d : %3d | %s %s", homeScore, awayScore, val.Boxscore.StatusDesc, val.Boxscore.PeriodClock)
			// case "3":
			// 	gameInfo = fmt.Sprintf(" %3d : %3d | 結束", homeScore, awayScore)
		}
		teamMessage := fmt.Sprintf("#%d %s vs %s\n      %s", index+1, homeTeamName, awayTeamName, gameInfo)
		message += teamMessage + "\n"
	}
	return message
}
