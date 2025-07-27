package cli

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"starter/go_starter/chatClient"
	"starter/go_starter/docUpload"
	"starter/go_starter/knowledge"
	"strconv"
	"strings"
	"time"
)

const Blue = "\033[1;34m"
const Green = "\033[1;32m"
const Red = "\033[1;31m"
const Bold = "\033[1m"
const Underline = "\033[4m"
const colorReset = "\033[0m"

const CommandPrefix = "."

type CLI struct {
	chat      *chatclient.ChatClient
	reader    *bufio.Reader
	commands  map[string]Command
	Uploader  *docUpload.Uploader
	Store     *knowledge.Store
	Extractor *knowledge.Extractor
}

type Command struct {
	Desc   string
	Action func(*CLI)
}

func New(chat *chatclient.ChatClient, uploader *docUpload.Uploader, store *knowledge.Store, extractor *knowledge.Extractor) *CLI {
	commands := map[string]Command{
		".exit": {
			Desc: "Exit the application",
			Action: func(c *CLI) {
				c.Exit()
			},
		},
		".help": {
			Desc: "Show available commands",
			Action: func(c *CLI) {
				c.PrintHelp()
			},
		},
		".history": {
			Desc: "Show conversation history",
			Action: func(c *CLI) {
				c.PrintHistory()
			},
		},
		".upload": {
			Desc: "Get a doc review",
			Action: func(c *CLI) {
				c.UploadAndReviewDoc()
			},
		},
	}
	return &CLI{
		chat:      chat,
		commands:  commands,
		reader:    bufio.NewReader(os.Stdin),
		Uploader:  uploader,
		Store:     store,
		Extractor: extractor,
	}
}

func (c *CLI) Run() {
	c.printWelcome()
	c.PrintHelp()

	for {
		fmt.Print(Colored("you >  ", Blue))
		input, err := c.reader.ReadString('\n')

		if err != nil {
			log.Fatal(err)
		}

		input = strings.TrimSpace(input)

		if strings.HasPrefix(input, CommandPrefix) {
			c.handleCommand(input)
		} else {
			c.simpleChat(input)
			c.addInputToKnowledge(input)
		}

	}
}

func (c *CLI) simpleChat(input string) {
	ctx, done := context.WithCancel(c.chat.Ctx)
	c.ShowSpinner(ctx, Colored("🚀 generating...", Green))
	defer done()

	res := c.chat.SendMessage(input)
	fmt.Println(Colored("model> ", Green), res)
}

func (c *CLI) handleCommand(input string) {
	command, ok := c.commands[strings.ToLower(input)]
	if !ok {
		fmt.Printf("Error: '%s' is not a recognized command. \n", input)
	} else {
		command.Action(c)
	}
}

func (c *CLI) PrintHistory() {
	history := c.chat.Chat.History(false)
	if len(history) == 0 {
		fmt.Println("No conversation history yet.")
		return
	}

	fmt.Println("\n========== Conversation History ==========")
	for _, msg := range history {
		fmt.Print(Colored(msg.Role+"> ", Blue))
		fmt.Print(msg.Parts[0].Text)
	}
	fmt.Println("==============================================")
}

func (c *CLI) UploadAndReviewDoc() {
	c.printWelcomeToReview()

	file_name, err := c.GetFileName()
	if err != nil {
		fmt.Println(Colored("Error selecting file", Red))
		c.Exit()
	}

	ctx, done := context.WithCancel(c.chat.Ctx)
	c.ShowSpinner(ctx, Colored("🚀 getting Review...", Green))
	defer done()

	file_name = strings.TrimSpace(file_name)
	res := c.Uploader.UploadAndReviewDoc(file_name)

	fmt.Println(Colored("📋 Review Results:", Blue, Underline))
	fmt.Println(res)
}

func (c *CLI) GetFileName() (string, error) {
	dir := "./docs"

	fmt.Println(Colored("📂 Available Documents:", Green, Bold))
	files, _ := os.ReadDir(dir)
	for i, f := range files {
		fmt.Printf("%d: %s\n", i+1, f.Name())
	}

	fmt.Printf("%s", Colored("Select file number > ", Blue))
	fileNumStr, _ := c.reader.ReadString('\n')
	fileNumStr = strings.TrimSpace(fileNumStr)
	fileNum, err := strconv.Atoi(fileNumStr)
	if err != nil || fileNum < 1 || fileNum > len(files) {
		fmt.Println(Colored("🚫 Invalid selection", Red))
		return "", err
	}

	return filepath.Join(dir, files[fileNum-1].Name()), nil
}

func (c *CLI) PrintHelp() {

	keys := make([]string, 0, len(c.commands))
	for k := range c.commands {
		keys = append(keys, k)
	}

	fmt.Println("\nCommands:")
	for _, key := range sort.StringSlice(keys) {
		fmt.Printf("%s - %s\n", key, c.commands[key].Desc)
	}
	fmt.Println()
}

func (_ *CLI) printWelcome() {
	fmt.Println("==============================================")
	fmt.Println("         Welcome to Eventures!")
	fmt.Println("==============================================")
}

func (c *CLI) printWelcomeToReview() {
	fmt.Println("--------------------------------------------------")
	fmt.Println("           🚀 Welcome to AI Doc Review! 🚀      ")
	fmt.Println("--------------------------------------------------")
}

func (c *CLI) ShowSpinner(ctx context.Context, msg string) {
	spinnerFrames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	go func() {
		i := 0
		for {
			select {
			case <-ctx.Done():
				return
			default:
				fmt.Printf("\r%s  %s", msg, spinnerFrames[i%len(spinnerFrames)])
				fmt.Printf("\r")
				time.Sleep(100 * time.Millisecond)
				i++
			}
		}
	}()
}

func (c *CLI) addInputToKnowledge(input string) {
	chunks := c.Extractor.ExtractFromUserInput(input)
	if chunks != nil {
		c.Store.AddKnowledge(*chunks)
	}
}

func (_ *CLI) Exit() {
	fmt.Print("Good Bye")
	os.Exit(0)
}

func Colored(text string, styles ...string) string {
	styleStr := strings.Join(styles, " ")
	return fmt.Sprintf("%s%s%s", styleStr, text, colorReset)
}
