package services

import (
	"strings"

	"qmaru-api/configs"
	"qmaru-api/utils"
)

func tweetStatusID(url string) (i string) {
	firstSplit := strings.Split(url, "/")
	idRaw := firstSplit[len(firstSplit)-1]
	lastSplit := strings.Split(idRaw, "?")
	i = lastSplit[0]
	return
}

func tweetToken(k, s string) (t string) {
	// https://developer.twitter.com/en/docs/basics/authentication/api-reference/token
	oauthAPI := "https://api.twitter.com/oauth2/token"

	data := utils.MiniFormData{
		"grant_type": "client_credentials",
	}

	auth := utils.MiniAuth{k, s}

	res := utils.Minireq.Post(oauthAPI, data, auth)
	bodyJSON := res.RawJSON().(map[string]interface{})

	t = bodyJSON["access_token"].(string)
	return
}

func tweetData(statusID, token string) (vurl string) {
	showAPI := "https://api.twitter.com/1.1/statuses/show.json"
	headers := utils.MiniHeaders{
		"Authorization": "Bearer " + token,
	}

	params := utils.MiniParams{
		"id":         statusID,
		"tweet_mode": "extended",
	}

	res := utils.Minireq.Get(showAPI, headers, params)
	resJSON := res.RawJSON().(map[string]interface{})

	if _, ok := resJSON["extended_entities"]; ok {
		tweetExtendedEntities := resJSON["extended_entities"].(map[string]interface{})
		tweetMediaList := tweetExtendedEntities["media"].([]interface{})
		tweetMedia := tweetMediaList[0].(map[string]interface{})
		tweetVideoInfo := tweetMedia["video_info"].(map[string]interface{})
		tweetVariants := tweetVideoInfo["variants"].([]interface{})
		tweetBitrate := 0.0
		for _, v := range tweetVariants {
			tweetValue := v.(map[string]interface{})
			if _, ok := tweetValue["bitrate"]; ok {
				tweetBitrateM := tweetValue["bitrate"].(float64)
				if tweetBitrateM > tweetBitrate {
					tweetBitrate = tweetBitrateM
					vurl = tweetValue["url"].(string)
				}
			}
		}
	} else {
		vurl = ""
	}
	return
}

// TweetVideo 获取 Tweet 视频
func TweetVideo(url string) (vurl []interface{}) {
	cfg := configs.TweetCfg()
	token := cfg["token"].(string)
	if token == "" {
		key := cfg["twitter_key"].(string)
		secret := cfg["twitter_secret"].(string)
		token = tweetToken(key, secret)
	}

	statusID := tweetStatusID(url)
	turl := tweetData(statusID, token)
	vurl = []interface{}{
		turl,
	}
	return
}
