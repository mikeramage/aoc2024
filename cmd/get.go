/*
Copyright © 2024 Mike Ramage <mike.ramage@gmail.com>
*/
package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var day int
var outDir string
var sessionCookie string

const sessionFilename = "./session.txt"

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get input for specified day",
	Long:  `Get problem input for specified day from Advent of Code website and write it to disk`,
	Run: func(cmd *cobra.Command, args []string) {

		if err := getSessionCookie(&sessionCookie); err != nil {
			log.Fatalln("Command failed: invalid session")
		}

		body, err := getProblemInputFromServer()

		if err != nil {
			log.Fatalln("Command failed: could not get input from server")
		}

		if err := writeToFile(body); err != nil {
			log.Fatalln("Command failed: could not write input to file")
		}

		fmt.Printf("Successfully written day%v program input to file\n", day)
	},
}

func getSessionCookie(sessionCookie *string) error {
	if len(*sessionCookie) == 0 {
		bytes, err := os.ReadFile(sessionFilename)
		if err != nil {
			log.Printf("Could not read session from file: %v. Please use --session / -s to specify a session cookie\n", err)
			return err
		}

		*sessionCookie = string(bytes)
	} else {
		sessionFile, err := os.Create(sessionFilename)
		if err != nil {
			log.Printf("Failed to create session cookie cache file: %v: received error %v\n", sessionFilename, err)
			return err
		}
		defer func() {
			if err := sessionFile.Close(); err != nil {
				log.Println("Warning: failed to close session cookie file:", err)
			}
		}()

		_, err = sessionFile.WriteString(*sessionCookie)
		if err != nil {
			log.Printf("Failed to write session to session cookie cache file: %v: received error %v\n", sessionFilename, err)
			return err
		}
	}

	return nil
}

func getProblemInputFromServer() ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://adventofcode.com/2024/day/%v/input", day), nil)
	if err != nil {
		log.Printf("Failed to build request: Error: %v\n", err)
		return nil, err
	}

	req.Header.Add("Cookie", fmt.Sprintf("session=%v", sessionCookie))

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to get input for day %v from AoC website: received error %v\n", day, err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read body from AoC website: received error:", err)
		return nil, err
	}

	return body, nil
}

func writeToFile(body []byte) error {
	f, err := os.Create(fmt.Sprintf("%v/day%v.txt", outDir, day))
	defer func() {
		if err := f.Close(); err != nil {
			log.Println("Warning: failed to close file:", err)
		}
	}()
	if err != nil {
		log.Printf("Failed to create file %v/day%v.txt: received error %v\n", outDir, day, err)
		return err
	}

	_, err = f.Write(body)
	if err != nil {
		log.Printf("Failed to write to file %v/day%v.txt: received error %v\n", outDir, day, err)
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().IntVarP(&day, "day", "d", 1, "Day for which to retrieve input data")
	getCmd.Flags().StringVarP(&outDir, "outdir", "o", "./input", "Output destination directory for problem input")
	getCmd.Flags().StringVarP(&sessionCookie, "session", "s", "", "Session cookie (only required on first use or if expires)")
}
