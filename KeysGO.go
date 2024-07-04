package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

func clearScreen() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "linux", "darwin":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		fmt.Println("Cannot clear terminal screen for this OS")
	}
}

func setTitle(title string) {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "title", title)
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "linux":
		fmt.Printf("\033]0;%s\007", title)
	case "darwin":
		fmt.Printf("\033]0;%s\007", title)
	default:
		fmt.Println("Cannot set terminal title for this OS")
	}
}

func parseNumber(input string) (*big.Int, error) {
	input = strings.ToLower(input)
	input = strings.ReplaceAll(input, ",", ".")

	re := regexp.MustCompile(`(\d*\.?\d*)([kmgt]?)`)
	matches := re.FindStringSubmatch(input)
	if len(matches) < 3 {
		return nil, fmt.Errorf("entrada inválida")
	}

	baseStr := matches[1]
	suffix := matches[2]

	base, err := strconv.ParseFloat(baseStr, 64)
	if err != nil {
		return nil, fmt.Errorf("entrada inválida")
	}

	switch suffix {
	case "k":
		base *= 1e3
	case "m":
		base *= 1e6
	case "b":
		base *= 1e9
	case "t":
		base *= 1e12
	}

	return new(big.Int).SetInt64(int64(base)), nil
}

func calculateTimeToCrack(bits int, keysPerSecond *big.Int) string {
	totalCombinations := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(bits)), nil)
	timeInSeconds := new(big.Int).Div(totalCombinations, keysPerSecond)
	secondsInAYear := big.NewInt(60 * 60 * 24 * 365)
	timeInYears := new(big.Float).Quo(new(big.Float).SetInt(timeInSeconds), new(big.Float).SetInt(secondsInAYear))

	years := new(big.Float).Quo(timeInYears, big.NewFloat(1))
	yearsInt, _ := years.Int(nil)

	if yearsInt.Sign() == 0 {
		// Menos de um ano
		months := new(big.Float).Mul(years, big.NewFloat(12))
		monthsInt, _ := months.Int(nil)

		seconds := new(big.Float).SetFloat64(0)
		if monthsInt.Sign() > 0 {
			monthsSeconds := new(big.Float).Mul(months, big.NewFloat(365*24*60*60/12))
			seconds = new(big.Float).Sub(timeInYears, new(big.Float).SetFloat64(float64(monthsInt.Int64())/12)).Mul(monthsSeconds, big.NewFloat(12))
		} else {
			seconds = timeInYears
		}

		return fmt.Sprintf("%d meses %.2f segundos", monthsInt, seconds)
	} else {
		// Um ano ou mais
		return fmt.Sprintf("%d anos", yearsInt)
	}
}

func calculateProgress(bits int, keysPerSecond, keysChecked *big.Int) string {
	totalCombinations := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(bits)), nil)
	keysLeft := new(big.Int).Sub(totalCombinations, keysChecked)
	timeInSeconds := new(big.Int).Div(keysLeft, keysPerSecond)
	secondsInAYear := big.NewInt(60 * 60 * 24 * 365)
	timeInYears := new(big.Float).Quo(new(big.Float).SetInt(timeInSeconds), new(big.Float).SetInt(secondsInAYear))

	years := new(big.Float).Quo(timeInYears, big.NewFloat(1))
	yearsInt, _ := years.Int(nil)

	if yearsInt.Sign() == 0 {
		// Menos de um ano
		months := new(big.Float).Mul(years, big.NewFloat(12))
		monthsInt, _ := months.Int(nil)

		seconds := new(big.Float).SetFloat64(0)
		if monthsInt.Sign() > 0 {
			monthsSeconds := new(big.Float).Mul(months, big.NewFloat(365*24*60*60/12))
			seconds = new(big.Float).Sub(timeInYears, new(big.Float).SetFloat64(float64(monthsInt.Int64())/12)).Mul(monthsSeconds, big.NewFloat(12))
		} else {
			seconds = timeInYears
		}

		return fmt.Sprintf("%d meses %.2f segundos", monthsInt, seconds)
	} else {
		// Um ano ou mais
		return fmt.Sprintf("%d anos", yearsInt)
	}
}

func main() {
	setTitle("KeysGo")
	clearScreen()

	reader := bufio.NewReader(os.Stdin)

	green := "\033[32m"
	reset := "\033[0m"

	fmt.Println(green + `
▀████▀ ▀███▀                           ▄▄█▀▀▀█▄█  ▄▄█▀▀██▄  
  ██   ▄█▀                           ▄██▀     ▀█▄██▀    ▀██▄
  ██ ▄█▀      ▄▄█▀██▀██▀   ▀██▀▄██▀████▀       ▀██▀      ▀██
  █████▄     ▄█▀   ██ ██   ▄█  ██   ▀▀█         ██        ██
  ██  ███    ██▀▀▀▀▀▀  ██ ▄█   ▀█████▄█▄    ▀█████▄      ▄██
  ██   ▀██▄  ██▄    ▄   ███    █▄   ████▄     ██▀██▄    ▄██▀
▄████▄   ███▄ ▀█████▀   ▄█     ██████▀ ▀▀███████  ▀▀████▀▀  
                      ▄█                                    
                    ██▀                                     
` + reset)

	fmt.Println(green + "Creditos: Dog_Gabriel" + reset)

	for {
		fmt.Println("1. Calcular tempo necessário para quebrar uma chave")
		fmt.Println("2. Calcular progresso com base no número de chaves já verificadas")
		fmt.Println("Pressione Ctrl + C para sair do programa")

		fmt.Print("Escolha uma opção: ")

		choiceStr, _ := reader.ReadString('\n')
		choiceStr = strings.TrimSpace(choiceStr)
		choice, err := strconv.Atoi(choiceStr)
		if err != nil || (choice != 1 && choice != 2) {
			fmt.Println("Opção inválida")
			continue
		}

		switch choice {
		case 1:
			fmt.Print("Digite o número de bits da chave: ")
			bitsStr, _ := reader.ReadString('\n')
			bitsStr = strings.TrimSpace(bitsStr)
			bits, err := strconv.Atoi(bitsStr)
			if err != nil {
				fmt.Println("Entrada inválida para o número de bits")
				continue
			}

			fmt.Print("Digite o número de chaves verificadas por segundo: ")
			keysPerSecondStr, _ := reader.ReadString('\n')
			keysPerSecondStr = strings.TrimSpace(keysPerSecondStr)
			keysPerSecond, err := parseNumber(keysPerSecondStr)
			if err != nil {
				fmt.Println("Entrada inválida para o número de chaves por segundo")
				continue
			}

			timeToCrack := calculateTimeToCrack(bits, keysPerSecond)
			fmt.Printf("Tempo estimado para testar todas as combinações de uma carteira de %d bits a %.2f chaves por segundo: %s\n", bits, float64(keysPerSecond.Int64())/1e6, timeToCrack)

		case 2:
			fmt.Print("Digite o número de bits da chave: ")
			bitsStr, _ := reader.ReadString('\n')
			bitsStr = strings.TrimSpace(bitsStr)
			bits, err := strconv.Atoi(bitsStr)
			if err != nil {
				fmt.Println("Entrada inválida para o número de bits")
				continue
			}

			fmt.Print("Digite o número de chaves verificadas por segundo: ")
			keysPerSecondStr, _ := reader.ReadString('\n')
			keysPerSecondStr = strings.TrimSpace(keysPerSecondStr)
			keysPerSecond, err := parseNumber(keysPerSecondStr)
			if err != nil {
				fmt.Println("Entrada inválida para o número de chaves por segundo")
				continue
			}

			fmt.Print("Digite o número de chaves já verificadas: ")
			keysCheckedStr, _ := reader.ReadString('\n')
			keysCheckedStr = strings.TrimSpace(keysCheckedStr)
			keysChecked, err := parseNumber(keysCheckedStr)
			if err != nil {
				fmt.Println("Entrada inválida para o número de chaves verificadas")
				continue
			}

			progress := calculateProgress(bits, keysPerSecond, keysChecked)
			fmt.Printf("Tempo estimado para testar as combinações restantes de uma carteira de %d bits a %.2f chaves por segundo: %s\n", bits, float64(keysPerSecond.Int64())/1e6, progress)
		}
	}
}
