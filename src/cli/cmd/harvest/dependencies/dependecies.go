package dependecies

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cldcvr/terrarium/src/cli/internal/config"
	"github.com/cldcvr/terrarium/src/pkg/db"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var depIfaceDirectoryFlag string

var dependencyCmd = &cobra.Command{
	Use:   "dependencies",
	Short: "Harvests dependencies from the given directory",
	Long:  "Harvests dependencies from all the yaml/yml files in the directory provided and adds it to the database.",
	Run: func(cmd *cobra.Command, args []string) {
		main()
	},
}

func GetCmd() *cobra.Command {
	addFlags()
	return dependencyCmd
}

func addFlags() {
	dependencyCmd.Flags().StringVarP(&depIfaceDirectoryFlag, "dir", "d", "", "path to dependency directory")
}

func main() {
	g, err := config.DBConnect()
	if err != nil {
		panic(err)
	}

	err = filepath.Walk(depIfaceDirectoryFlag, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file is a YAML file (ends with .yaml or .yml)
		if !info.IsDir() && (strings.HasSuffix(info.Name(), ".yaml") || strings.HasSuffix(info.Name(), ".yml")) {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			// Parse the YAML data
			var yamlData map[string][]db.Dependency
			err = yaml.Unmarshal(data, &yamlData)
			if err != nil {
				fmt.Printf("Error parsing YAML file %s: %s\n", path, err)
				return nil
			}
			// Loop through each dependency entry and call CreateDependencyInterface
			for _, dep := range yamlData["dependency-interface"] {
				dep.ID = uuid.New()
				_, err := g.CreateDependencyInterface(&dep)
				if err != nil {
					fmt.Printf("Error updating the database: %s\n", err)
				} else {
					fmt.Println("Data inserted successfully!")
				}
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
