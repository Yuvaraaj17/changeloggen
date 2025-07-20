package main

import (
	"fmt"
	"os"
	"strings"
	"os/exec"
	"log"
)

var helpMsg string = `
Welcome to ChangeLogGen - Change Log Generator

Usage commands

-h - Show the list of cmds available for usage
-f <commit-id-1> -t <commit-id-2> - Generates change from the given commits
-v - Shows the version
`
var groups = map[string][]string{}

func categorize(msg string) string {
	msg = strings.ToLower(msg)
	switch {
	case strings.Contains(msg, "fix"), strings.Contains(msg, "bug"):
		return "Fixes"
	case strings.Contains(msg, "add"), strings.Contains(msg, "create"), strings.Contains(msg, "new"), strings.Contains(msg, "update"):
		return "Features"
	case strings.Contains(msg, "refactor"), strings.Contains(msg, "clean"), strings.Contains(msg, "remove"):
		return "Cleanup"
	case strings.Contains(msg, "readme"), strings.Contains(msg, "doc"):
		return "Docs"
	default:
		return "Other"
	}
}

func writeChLog(mp map[string][]string, from string, to string){
	file, err := os.Create("changelog.txt")
	if err != nil {
		log.Fatalf("❌ Failed to create file: %v", err)
	}
	defer file.Close()
	_,_ = file.WriteString(fmt.Sprintf("##### Change Log For versions %s --- %s #####\n", from, to))
	for cat,lines :=range mp{
		_, err := file.WriteString(fmt.Sprintf("## %s\n", cat))
		if err != nil {
			log.Fatalf("❌ Failed to write category: %v", err)
		}
		for _, line := range lines {
			_, err := file.WriteString(fmt.Sprintf("- %s\n", line))
			if err != nil {
				log.Fatalf("❌ Failed to write line: %v", err)
			}
		}
		_, _ = file.WriteString("\n") // Add space between categories
	}
}

func generateChLog(args []string) {

	if len(args) == 4 && args[0] == "-f" && args[2] == "-t" {
		pretty := "--pretty=format:%s"
		result := exec.Command("git","log",fmt.Sprintf(`%s...%s`,args[1],args[3]),pretty)
		op, err := result.Output()
		if err != nil {
			log.Fatal(err)
			return
		}
		for msg := range strings.SplitSeq(string(op), "\n"){
			cat := categorize(msg)
			groups[cat] = append(groups[cat], msg)
		}
		writeChLog(groups, args[1],args[3])
		return 
	}

	fmt.Println("Invalid arguments provided check -h for usage")
}

func main() {
	argument := strings.ToLower(os.Args[1])
	switch argument {
	case "-h":
		fmt.Println(helpMsg)
	case "-f":
		generateChLog(os.Args[1:])
	case "-v":
		fmt.Println("ChangeLogGen Version: v1.0.0")
	}

}
