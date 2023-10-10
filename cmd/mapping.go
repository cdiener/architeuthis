/*
Copyright © 2023 Christian Diener <mail(a)cdiener.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"github.com/spf13/cobra"
)

// mappingCmd represents the mapping command
var mappingCmd = &cobra.Command{
	Use:   "mapping",
	Short: "Analyze read and k-mer level mapping.",
	Long: `The mapping command helps with analyzing the read- and k-mer level
taxonomic assignments of your Kraken output.`,
}

func init() {
	rootCmd.AddCommand(mappingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mappingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mappingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
