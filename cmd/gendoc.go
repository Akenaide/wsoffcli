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

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// gendocCmd represents the gendoc command
var gendocCmd = &cobra.Command{
	Use:   "gendoc",
	Short: "Generate doc with Cobra",
	Long:  `Generate documentation with Cobra`,
	Run: func(cmd *cobra.Command, args []string) {
		var docPath string = "doc"
		fmt.Println("gendoc called")
		os.Mkdir(docPath, 0744)
		linkHandler := func(filename string) string {
			if filename == "wsoffcli.md" {
				return "../README.md"
			}
			return fmt.Sprintf("%v/%v", docPath, filename)
		}

		filePrepender := func(name string) string {
			return ""
		}
		err := doc.GenMarkdownTreeCustom(rootCmd, docPath, filePrepender, linkHandler)
		if err != nil {
			fmt.Println(err)
		}
		os.Rename(fmt.Sprintf("%v/wsoffcli.md", docPath), fmt.Sprintf("%v/../README.md", docPath))
	},
}

func init() {
	rootCmd.AddCommand(gendocCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gendocCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gendocCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
