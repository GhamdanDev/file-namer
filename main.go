package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

var (
	inputText   string
	ext         string
	sep         string
	keepNumbers bool
	showHelp    bool
	versionFlag bool
)

const (
	Version = "1.1.0"
)

func init() {
	flag.StringVar(&inputText, "t", "", "Input text for filename generation (required)")
	flag.StringVar(&ext, "ext", "dart", "File extension (e.g., dart, go, cs)")
	flag.StringVar(&sep, "sep", "_", "Word separator (e.g., _, -)")
	flag.BoolVar(&keepNumbers, "keep-numbers", false, "Preserve numbers in filename")
	flag.BoolVar(&showHelp, "help", false, "Show help")
	flag.BoolVar(&versionFlag, "version", false, "Show version")
}

func main() {
	flag.Parse()

	switch {
	case showHelp:
		printHelp()
		return
	case versionFlag:
		fmt.Printf("file-namer v%s\n", Version)
		return
	case inputText == "" && !flag.Parsed():
		runInteractiveMode()
	default:
		if inputText == "" {
			fmt.Println("❗ Error: Please provide input text using -t")
			os.Exit(1)
		}
		filename := generateFilename(inputText, ext, sep, keepNumbers)
		fmt.Println("✅ Generated filename:", filename)
	}
}

func generateFilename(input, ext, sep string, keepNumbers bool) string {
	cleanInput := sanitizeInput(input, keepNumbers)
	words := strings.Fields(cleanInput)
	filename := strings.Join(words, sep)
	filename = fmt.Sprintf("%s.%s", filename, ext)

	if !isValidFilename(filename) {
		fmt.Println("🚫 Invalid filename characters detected!")
		os.Exit(1)
	}

	return filename
}

func sanitizeInput(input string, keepNumbers bool) string {
	pattern := "[^a-z ]"
	if keepNumbers {
		pattern = "[^a-z0-9 ]"
	}

	reg := regexp.MustCompile(pattern)
	clean := reg.ReplaceAllString(strings.ToLower(input), "")
	return clean
}

func isValidFilename(filename string) bool {
	invalidChars := `<>:"/\|?*`
	return !strings.ContainsAny(filename, invalidChars)
}

func runInteractiveMode() {
	answers := struct {
		InputText   string
		Ext         string
		Sep         string
		KeepNumbers bool
	}{}

	// Input text prompt
	survey.AskOne(&survey.Input{
		Message: "Enter text for filename:",
	}, &answers.InputText, survey.WithValidator(survey.Required))

	// Extension selection
	extOptions := []string{"dart", "go", "cs", "js", "ts", "py"}
	survey.AskOne(&survey.Select{
		Message: "Choose file extension:",
		Options: extOptions,
		Default: "dart",
	}, &answers.Ext)

	// Separator selection
	sepOptions := []string{"_", "-", "camel", ""}
	survey.AskOne(&survey.Select{
		Message: "Choose word separator:",
		Options: sepOptions,
		Default: "_",
	}, &answers.Sep)

	// Keep numbers confirmation
	survey.AskOne(&survey.Confirm{
		Message: "Keep numbers in filename?",
		Default: false,
	}, &answers.KeepNumbers)

	// Handle camelCase separately
	if answers.Sep == "camel" {
		filename := generateCamelCase(answers.InputText, answers.Ext, answers.KeepNumbers)
		fmt.Println("\n✅ Generated filename:", filename)
		return
	}

	filename := generateFilename(answers.InputText, answers.Ext, answers.Sep, answers.KeepNumbers)
	fmt.Println("\n✅ Generated filename:", filename)
}

func generateCamelCase(input, ext string, keepNumbers bool) string {
	cleanInput := sanitizeInput(input, keepNumbers)
	words := strings.Fields(cleanInput)
	for i, word := range words {
		if i > 0 {
			words[i] = strings.Title(word)
		}
	}
	filename := strings.Join(words, "") + "." + ext
	return filename
}

func printHelp() {
	fmt.Println(`
Usage:
  file-namer [flags]

Flags:
  -t            Input text (required)
  -ext          File extension (default: dart)
  -sep          Word separator (default: _)
  -keep-numbers Preserve numbers (default: false)
  -help         Show help
  -version      Show version

Interactive Mode:
  Simply run the program without flags

Examples:
  file-namer -t "Auth Service" -ext go -sep _
  file-namer -t "Data2024" -ext cs -keep-numbers`)
}
