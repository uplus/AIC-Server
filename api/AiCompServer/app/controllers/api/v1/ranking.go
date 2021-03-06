package controllers

import (
	"AiCompServer/app/db"
	"AiCompServer/app/models"
	"errors"
	"github.com/revel/revel"
	"sort"
)

type Rank struct {
	Rank     int     `json:"rank"`
	Username string  `json:"username"`
	Score    float64 `json:"score"`
}

type Ranks []Rank

type ResponseRanking struct {
	Ranking Ranks `json:"ranking"`
}

func RankScore() (Ranks, error) {
	users := []models.User{}
	if err := db.DB.Order("id desc").Find(&users).Error; err != nil {
		return nil, errors.New("Record Find Failure")
	}
	var rank Ranks
	if PreRanking() {
		answer := []models.PreAnswer{}
		for _, user := range users {
			if user.Role == "admin" {
				continue
			}
			score := 0.0
			if err := db.DB.Find(&answer, "user_id = ?", user.ID).Error; err != nil {
				return nil, err
			}
			for _, ans := range answer {
				score = score + ans.Score
			}
			rank = append(rank, Rank{Rank: 0, Username: user.Username, Score: score})
		}
		sort.Slice(rank, func(i, j int) bool { return rank[i].Score > rank[j].Score })
		return rank, nil
	} else {
		answer := []models.Answer{}
		for _, user := range users {
			if user.Role == "admin" {
				continue
			}
			score := 0.0
			if err := db.DB.Find(&answer, "user_id = ?", user.ID).Error; err != nil {
				return nil, err
			}
			for _, ans := range answer {
				score = score + ans.Score
			}
			rank = append(rank, Rank{Rank: 0, Username: user.Username, Score: score})
		}
		sort.Slice(rank, func(i, j int) bool { return rank[i].Score > rank[j].Score })
		return rank, nil
	}
}

func (c ApiChallenge) Ranking() revel.Result {
	if err := CheckToken(c.ApiV1Controller); err != nil {
		return err
	}
	rank, err := RankScore()
	if err != nil {
		return c.HandleNotFoundError(err.Error())
	}
	Rindex := 1
	TmpScore := 0.0
	for Tindex, _ := range rank {
		// スコアが前の人と一緒の時は前の人の順位と同じにする
		if rank[Tindex].Score == TmpScore {
			rank[Tindex].Rank = Rindex
			TmpScore = rank[Tindex].Score
		} else {
			rank[Tindex].Rank = Tindex + 1
			Rindex = Tindex + 1
			TmpScore = rank[Tindex].Score
		}
		rank[Tindex].Score = float64(int((rank[Tindex].Score*1000)+0.5)) / 1000
	}
	r := Response{ResponseRanking{rank}}
	return c.RenderJSON(r)
}

func ChallengeRankScore(id int) (Ranks, error) {
	users := []models.User{}
	if err := db.DB.Order("id desc").Find(&users).Error; err != nil {
		return nil, errors.New("Record Find Failure")
	}
	var rank Ranks
	if PreRanking() {
		answer := models.PreAnswer{}
		for _, user := range users {
			if user.Role == "admin" {
				continue
			}
			if err := db.DB.Where("user_id = ? AND challenge_id = ?", user.ID, id).First(&answer).Error; err != nil {
				continue
			}
			rank = append(rank, Rank{Rank: 0, Username: user.Username, Score: answer.Score})
		}
		sort.Slice(rank, func(i, j int) bool { return rank[i].Score > rank[j].Score })
		return rank, nil
	} else {
		answer := models.Answer{}
		for _, user := range users {
			if user.Role == "admin" {
				continue
			}
			if err := db.DB.Where("user_id = ? AND challenge_id = ?", user.ID, id).First(&answer).Error; err != nil {
				continue
			}
			rank = append(rank, Rank{Rank: 0, Username: user.Username, Score: answer.Score})
		}
		sort.Slice(rank, func(i, j int) bool { return rank[i].Score > rank[j].Score })
		return rank, nil
	}
}

func (c ApiChallenge) RankingChallenge(id int) revel.Result {
	if err := CheckToken(c.ApiV1Controller); err != nil {
		return err
	}
	rank, err := ChallengeRankScore(id)
	if err != nil {
		return c.HandleNotFoundError(err.Error())
	}
	Rindex := 1
	TmpScore := 0.0
	for Tindex, _ := range rank {
		// スコアが前の人と一緒の時は前の人の順位と同じにする
		if rank[Tindex].Score == TmpScore {
			rank[Tindex].Rank = Rindex
			TmpScore = rank[Tindex].Score
		} else {
			rank[Tindex].Rank = Tindex + 1
			Rindex = Tindex + 1
			TmpScore = rank[Tindex].Score
		}
		rank[Tindex].Score = float64(int((rank[Tindex].Score*1000)+0.5)) / 1000
	}
	r := Response{ResponseRanking{rank}}
	return c.RenderJSON(r)
}
