package cmd

import (
	. "coca/src/adapter/api"
	"coca/src/adapter/models"
	"coca/src/domain"
	. "coca/src/utils"
	"encoding/json"
	"github.com/spf13/cobra"
	"log"
	"os/exec"
)

var apiCmd *cobra.Command = &cobra.Command{
	Use:   "api",
	Short: "scan api",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		path := cmd.Flag("path").Value.String()
		dependence := cmd.Flag("dependence").Value.String()

		if path != "" {
			app := new(JavaApiApp)
			restApis := app.AnalysisPath(path, dependence)

			cModel, _ := json.MarshalIndent(restApis, "", "\t")
			WriteToFile("apis.json", string(cModel))

			var parsedDeps []models.JClassNode
			file := ReadFile(dependence)
			if file == nil {
				log.Fatal("lost file:" + dependence)
			}
			_ = json.Unmarshal(file, &parsedDeps)

			var apiCallers []string
			for _, api := range restApis {
				apiCaller := api.PackageName + "." + api.ClassName + "." + api.MethodName
				apiCallers = append(apiCallers, apiCaller)
			}

			analyser := domain.NewCallGraph()
			dotContent := analyser.AnalysisByFiles(apiCallers, parsedDeps)

			WriteToFile("api.dot", dotContent)

			cmd := exec.Command("dot", []string{"-Tsvg", "api.dot", "-o", "api.svg"}...)
			_, err := cmd.CombinedOutput()
			if err != nil {
				log.Fatalf("cmd.Run() failed with %s\n", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)

	apiCmd.PersistentFlags().StringP("path", "p", "", "path")
	apiCmd.PersistentFlags().StringP("dependence", "d", "", "get dependence file")
}