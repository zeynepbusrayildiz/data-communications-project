package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

const (
	port      = ":8080"
	wordsFile = "words.txt"
)

func main() {
	words, err := loadWordsFromFile(wordsFile)
	if err != nil {
		fmt.Println("Could not load words from file:", err)
		os.Exit(1)
	}

	listen, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Could not start server:", err)
		os.Exit(1)
	}
	defer listen.Close()

	fmt.Println("Server has started, waiting for the client...")

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Could not receive client connection:", err)
			continue
		}
		fmt.Println("New client has connected:", conn.RemoteAddr())

		go handleClient(conn, words)
	}
}

func handleClient(conn net.Conn, words []string) {
	defer conn.Close()

	rand.Seed(time.Now().UnixNano())

	streak := 0

	for {

		secretWord := words[rand.Intn(len(words))]
		fmt.Println("New game started! Secret Word:", secretWord)

		conn.Write([]byte(fmt.Sprintf("\nHello! Try guessing the 4-letter word. You have 5 attempts. Current streak: %d\n", streak)))

		reader := bufio.NewReader(conn)
		attempts := 5

		for attempts > 0 {
			conn.Write([]byte(fmt.Sprintf("\nYou have %d attempts left. Current streak: %d. Enter your guess:\n", attempts, streak)))

			gue, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Could not receive data from client:", err)
				return
			}

			gue = strings.TrimSpace(gue)

			if gue == secretWord {
				streak++
				conn.Write([]byte(fmt.Sprintf("Congrats! Your guess is correct! Streak: %d\n", streak)))
				break
			}

			attempts--

			if attempts == 0 {
				conn.Write([]byte(fmt.Sprintf("Game over! The correct word was: %s. Your streak is reset to 0.\n", secretWord)))
				streak = 0
				break
			}

			answ := analyze(gue, secretWord)
			conn.Write([]byte(answ + "\n"))
		}
	}
}

func loadWordsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if len(word) == 4 {
			words = append(words, word)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

func analyze(gue, secretWord string) string {
	answ := ""
	for i := 0; i < len(secretWord); i++ {
		if i < len(gue) {
			if gue[i] == secretWord[i] {
				answ += string(gue[i]) + " is in its correct place, "
			} else if strings.Contains(secretWord, string(gue[i])) {
				answ += string(gue[i]) + " is in the wrong place, "
			} else {
				answ += string(gue[i]) + " is not in this word, "
			}
		} else {
			answ += "_ is missing, "
		}
	}
	return strings.TrimSuffix(answ, ", ")
}
