package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
	_ "unsafe"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "betterwhy <dependency with version to explain>",
	Short: "A better version of 'go mod why' which explains build requirements not actually required for execution",
	Long: `betterwhy traverses go mod's dependency tree to diagnose why a given dependency is required.  
It is different from 'go mod why' in that it will look at dependencies that are required to be pulled in, but not required for the build.
This is useful for scenarios like:
 - a direct dependency's test dependencies (no longer a problem in 1.16+)
 - explain why an exact version of a dependency was pulled in, allowing you to explain every entry in go.sum
It does not handle indirect dependencies (e.g. any dependency that does not declare its own dependencies, but still requires them).
Additionally, betterwhy finds A path to the desired dependency.
If two dependencies ultimately pull in the desired version, the returned path isn't guaranteed to be the shortest one.'`,
 	Example: "betterwhy github.com/stretchr/testify@v1.3.0",
 	Args: cobra.ExactArgs(1),
 	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		requested := args[0]
		if !strings.Contains(requested, "@") {
			panic("input must be an exact version of a module, pinned with the '@' symbol such as github.com/stretchr/testify@v1.3.0")
		}

		out, err := exec.Command("go",  "mod", "graph").Output()
		if err != nil {
			panic(errors.Wrap(err, "unable to execute 'go mod graph'"))
		}

		// first entry is the requested value.
		// lines will be reversed at the end to start with root at the top, and requested at the bottom
		lineage := []string{requested}
		nextSearch := requested
		for {
			// if we've reached the root node, we're done
			if !strings.Contains(nextSearch, "@") {
				break
			}

			// for each node in the tree
			for _, line := range strings.Split(string(out), "\n") {
				// if the node has what we're looking for
				if strings.Contains(line, nextSearch) {

					// and it's on the right side (e.g. someone looked for it)
					split := strings.Split(line, " ")
					if strings.Contains(split[1], nextSearch) {
						// then add the requester to the lineage
						lineage = append(lineage, split[0])
						nextSearch = split[0]
						break
					}
				}
			}
			if nextSearch == requested {
				fmt.Printf("requested dependency '%s' not found in tree.  Tidy and verify it exists in go.sum\n", requested)
				os.Exit(1)
			}
		}
		fmt.Printf("# %s\n", requested)
		for i := len(lineage)-1; i >= 0; i-- {
			fmt.Println(lineage[i])
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
