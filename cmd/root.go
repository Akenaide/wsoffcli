// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	serieNumber string
	neo         string
	allRarity   bool
)

// Card info to export
type Card struct {
	Set               string   `json:"set"`
	SetName           string   `json:"setName"`
	Side              string   `json:"side"`
	Release           string   `json:"release"`
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	JpName            string   `json:"jpName"`
	CardType          string   `json:"cardType"`
	Colour            string   `json:"colour"`
	Level             string   `json:"level"`
	Cost              string   `json:"cost"`
	Power             string   `json:"power"`
	Soul              string   `json:"soul"`
	Rarity            string   `json:"rarity"`
	BreakDeckbuilding bool     `json:"breakDeckbuilding"`
	ENEquivalent      bool     `json:"EN_Equivalent"`
	FlavourText       string   `json:"flavourText"`
	Trigger           []string `json:"trigger"`
	Ability           []string `json:"ability"`
	SpecialAttrib     []string `json:"specialAttrib"`
	Version           string   `json:"version"`
	Cardcode          string   `json:"cardcode"`
	ImageURL          string   `json:"imageURL"`
	Tags              []string `json:"tags"`
}

// CardModelVersion : Card format version
const CardModelVersion = "4"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wsoffcli",
	Short: "Collect data from https://ws-tcg.com/",
	Long: `Collect data from https://ws-tcg.com/.

Create a json file for each card with most information.

Example:
'wsoffcli fetch -n IMC' will fetch all cards with a code starting with 'IMC'

If you want more than one use '##' as seperator like 'wsoffcli fetch -n BD##IM'

'--serie' use a hidden number in the official site, this number is increment for each new set (e.g Kadokawa is number 259, Goblin 260 ...).

To use environ variable, use the prefix 'WSOFF'.
	 `,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {
	// },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVarP(&serieNumber, "serie", "s", "", "serie number")
	rootCmd.PersistentFlags().StringVarP(&neo, "neo", "n", "", "Neo standar by set")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".wsoffcli" (without extension).
		viper.SetEnvPrefix("wsoff")
		viper.AddConfigPath(home)
		viper.SetConfigName(".wsoffcli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
