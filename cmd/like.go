/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Feedl struct {
	Feeds []Feed `json:"feeds`
}

type Feed struct {
	Id       int       `json:"id"`
	Origin   string    `json:"origin"`
	Subjects []Subject `json:"subjects`
}

type Subject struct {
	Id      string `json:"subjectId"`
	Pairing bool   `json:"pairing"`
	Token   string `json:"ratingToken"`
}

type Identity struct {
	IdentityId string  `json:"identityId`
	Prof       Profile `json:"profile"`
}

type Profile struct {
	Firstname string `json:"firstName"`
}

var Auth string
var Sess string
var feeds Feedl
var Identities []Identity

// likeCmd represents the like command
var likeCmd = &cobra.Command{
	Use:   "like",
	Short: "Like [n] subjects",
	Long:  `Likes a number of subject, passed as an argument.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return err
		}
		if _, err := strconv.Atoi(args[0]); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		viper.SetConfigName("hinger")
		viper.SetConfigType("json")
		viper.AddConfigPath("$HOME/.config/")
		viper.AddConfigPath(".")
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				// Config file not found; create it
				var openaikey string
				fmt.Println("Config not found! Enter openai api key")
				fmt.Scanln(&openaikey)
				viper.Set("openai", openaikey)
				err := viper.WriteConfigAs("hinger.json")
				if err != nil {
					panic(err)
				}
			} else {
				panic(err)
			}
		}
		count, err := strconv.Atoi(args[0])
		if err != nil {
			panic(err)
		}
		httpclient := &http.Client{}
		client := hingeClient{client: *httpclient, auth: Auth}
		fmt.Println(count)
		fmt.Println(client.getRecs())
	},
}

// seekCmd represents the seek command
var seekCmd = &cobra.Command{
	Use:   "seek",
	Short: "Seek a certain person",
	Long:  `Want to find a special person? seek them`,
	Run: func(cmd *cobra.Command, args []string) {
		m := make(map[string]string)
		httpclient := &http.Client{}
		client := hingeClient{client: *httpclient, auth: Auth}
		found := false
		for !found {
			recs, err := client.getRecs()
			//print(recs)
			if err != nil {
				panic(err)
			}
			err = json.Unmarshal([]byte(recs), &feeds)
			if err != nil {
				panic("Couldn't unmarshal json")
			}
			ids := ""
			//fmt.Println(feeds.Feeds)
			for _, feed := range feeds.Feeds {
				if feed.Id != 0 {
					//print(feed.Id)
					continue
				}
				for _, subject := range feed.Subjects {
					//fmt.Println(subject.Token)
					m[subject.Id] = subject.Token
					ids = ids + "," + subject.Id
				}
				profiles, err := client.getUsers(ids)
				if err != nil {
					panic(err)
				}
				//fmt.Println(profiles)
				err = json.Unmarshal([]byte(profiles), &Identities)
				if err != nil {
					panic(err)
				}
				for _, Identity := range Identities {
					if strings.Contains(strings.ToLower(Identity.Prof.Firstname), strings.ToLower(args[0])) {
						found = true
						fmt.Println("Found possible match!!!!")
						fmt.Println(Identity.IdentityId)
						os.Exit(0)
					}
					//fmt.Println(m[Identity.IdentityId], Identity.IdentityId)
					client.doRating("skip", m[Identity.IdentityId], Sess, Identity.IdentityId)
					fmt.Printf(" on %s \n", Identity.Prof.Firstname)
					time.Sleep(500 * time.Millisecond)
				}
			}
			time.Sleep(5 * time.Second)
		}
	},
}

func init() {
	rootCmd.AddCommand(likeCmd)
	rootCmd.AddCommand(seekCmd)
	likeCmd.Flags().StringVarP(&Auth, "auth", "a", "", "Authentication header (required)")
	likeCmd.MarkFlagRequired("auth")
	seekCmd.Flags().StringVarP(&Auth, "auth", "a", "", "Authentication header (required)")
	seekCmd.MarkFlagRequired("auth")
	seekCmd.Flags().StringVarP(&Sess, "sess", "s", "", "Session id")
	seekCmd.MarkFlagRequired("sess")
}
