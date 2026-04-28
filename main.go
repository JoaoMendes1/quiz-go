package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
)

type Question struct {
	Text    string
	Options []string
	Answer  int
}

type GameState struct {
	Name      string
	Points    int
	Questions []Question
}

func (g *GameState) Init() {
	fmt.Println("Seja bem vindo(a) ao Quiz")
	fmt.Println("Escreva o seu nome:")
	reader := bufio.NewReader(os.Stdin)

	name, err := reader.ReadString('\n') // Ler o que escreve no terminal

	if err != nil {
		panic("Errp ao ler a string")
	}

	g.Name = name

	fmt.Printf("Vamos ao jogo %s", g.Name)
}

func (g *GameState) ProccessCSV() {
	f, err := os.Open("quiz-go.csv")
	if err != nil {
		panic("Erro ao ler arquivo")
	}

	defer f.Close()

	reader := csv.NewReader(f)
	record, err := reader.ReadAll()
	if err != nil {
		panic("Erro ao ler csv")
	}

	for index, record := range record {
		if index > 0 {
			correctAnswer, _ := toInt(record[5])
			question := Question{
				Text:    record[0],
				Options: record[1:5],
				Answer: correctAnswer,
			}

			g.Questions = append(g.Questions, question)
		}
	}
}

func (g *GameState) Run() {
	// Exibir a pergunta pro usuário
	for index, question := range g.Questions {
		fmt.Printf("\033[33m %d. %s \033[0m\n", index+1, question.Text)

		// Iterar sobre as opções no game state e exibir no terminal para o usuário 

		for j, option := range question.Options {
			fmt.Printf("[%d] %s\n", j+1, option)
		}

		fmt.Println("Digite uma alternativa:")

		var answer int 
		var err error 

		for {
			reader := bufio.NewReader(os.Stdin)
			read, _ := reader.ReadString('\n')

			answer, err = toInt(read[:len(read)-1])

			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			break 
		
		}

		//Validar a resposta 
		// Exibir a resposta correta para o usuário
		// Calcular a pontuação 
		if answer == question.Answer {
			fmt.Println("Parabéns vocÇe acertou!!")
			g.Points += 10
		} else {
			fmt.Println("Ops! Errou!")
			fmt.Println("----------------------------------")
		}

	}
}

func main() {
	game := &GameState{Points: 0}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		game.ProccessCSV()
		wg.Done()
	}()
//	game.Init()
	wg.Wait()
	game.Run()

	fmt.Printf("Fim de jogo, Você fez %d pontos\n", game.Points)
}

func toInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.New("Não é permitido caractere diferente de número")
	}

	return i, nil 
}
