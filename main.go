package main

import "github.com/dghubble/oauth1"
import "github.com/dghubble/go-twitter/twitter"
import "io/ioutil"
import "log"
import "time"
import "rand"

func main() {
	ck, err := ioutil.ReadFile("consumer_key")
	if err != nil {
		log.Fatalf("Couldn't read consumer key: %v", err)
	}

	cs, err := ioutil.ReadFile("consumer_secret")
	if err != nil {
		log.Fatalf("Couldn't read consumer secret: %v", err)
	}

	at, err := ioutil.ReadFile("access_token")
	if err != nil {
		log.Fatalf("Couldn't read access token: %v", err)
	}

	as, err := ioutil.ReadFile("as")
	if err != nil {
		log.Fatalf("Couldn't read access secret: %v", err)
	}

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	// Loading followers is mandatory on startup
	followers := getFollowers(cl)
	if followers == nil {
		log.Fatal("Failed to load followers at startup")
	}

	for {
		go loop(cl, followers)
		time.Sleep(time.Minute * 15)
	}
}

// On startup we load our follower list, and once a day we reload it
func getFollowers(cl twitter.Client) []string {
	params := &twitter.FollowerListParams{ScreenName: "DeleteEveryAcct"}
	followers, _, err := cl.Followers.List(params)
	if err != nil {
		log.Errorf("Failed to get followers: %v")
		return nil
	}

	var fol []string
	for _, user := range followers {
		fol = append(fol, user.ScreenName)
	}
	cursor := followers.NextCursor

	for ; ; cursor != 0 {
		params.Cursor = cursor
		followers, res, err := cl.Followers.List(params)
		if res == 429 {
			// sleep 15 minutes, then rerun
			time.Sleep(time.Minute * 15)
			continue
		}

		// if we succeded, update cursor and keep going
		cursor = followers.cursor
		for _, user := range followers {
			fol = append(fol, user.ScreenName)
		}
	}

	return fol
}

func loop(cl twitter.Client, followers []string) {
	// pick a random follower!
	i := rand.Intn(len(followers))
	ts := fmt.Sprintf("@%s, you should delete your account! You'd be free of this website!", followers[i])
	if twete, _, err := cl.Statuses.Update(ts, nil); err != nil {
		log.Errorf("Failed to post tweet %v: %v", twete, err)
	}
}
