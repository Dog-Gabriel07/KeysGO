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
		fmt.Println("Não foi possível limpar a tela para este sistema operacional")
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
		fmt.Println("Não foi possível definir o título do terminal para este sistema operacional")
	}
}

func calculateTimeToCrack(bits int, keysPerSecond *big.Int) float64 {
	totalCombinations := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(bits)), nil)
	timeInSeconds := new(big.Int).Div(totalCombinations, keysPerSecond)
	secondsInAYear := big.NewInt(31536000) // segundos em um ano padrão (60 * 60 * 24 * 365)
	timeInYears := new(big.Float).Quo(new(big.Float).SetInt(timeInSeconds), new(big.Float).SetInt(secondsInAYear))
	timeToCrack, _ := timeInYears.Float64()
	return timeToCrack
}

func calculateProgress(bits int, keysPerSecond, keysChecked *big.Int) float64 {
	totalCombinations := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(bits)), nil)
	keysLeft := new(big.Int).Sub(totalCombinations, keysChecked)
	timeInSeconds := new(big.Int).Div(keysLeft, keysPerSecond)
	secondsInAYear := big.NewInt(31536000) // segundos em um ano padrão (60 * 60 * 24 * 365)
	timeInYears := new(big.Float).Quo(new(big.Float).SetInt(timeInSeconds), new(big.Float).SetInt(secondsInAYear))
	timeToFinish, _ := timeInYears.Float64()
	return timeToFinish
}

func parseKeysPerSecond(input string) (*big.Int, error) {
	input = strings.ToLower(input)
	if strings.Contains(input, "m") {
		value, err := strconv.ParseFloat(strings.TrimRight(input, "m"), 64)
		if err != nil {
			return nil, err
		}
		return new(big.Int).Mul(big.NewInt(int64(value)), big.NewInt(1000000)), nil
	} else if strings.Contains(input, "b") {
		value, err := strconv.ParseFloat(strings.TrimRight(input, "b"), 64)
		if err != nil {
			return nil, err
		}
		return new(big.Int).Mul(big.NewInt(int64(value)), big.NewInt(1000000000)), nil
	} else if strings.Contains(input, "t") {
		value, err := strconv.ParseFloat(strings.TrimRight(input, "t"), 64)
		if err != nil {
			return nil, err
		}
		return new(big.Int).Mul(big.NewInt(int64(value)), big.NewInt(1000000000000)), nil
	} else if strings.Contains(input, ",") {
		value, err := strconv.ParseFloat(strings.Replace(input, ",", "", -1), 64)
		if err != nil {
			return nil, err
		}
		return new(big.Int).Mul(big.NewInt(int64(value)), big.NewInt(1000)), nil
	} else if strings.Contains(input, ".") {
		value, err := strconv.ParseFloat(strings.Replace(input, ".", "", -1), 64)
		if err != nil {
			return nil, err
		}
		return new(big.Int).Mul(big.NewInt(int64(value)), big.NewInt(10)), nil
	} else {
		return nil, fmt.Errorf("Formato inválido")
	}
}

func main() {
	setTitle("KeysGo")
	clearScreen()

	reader := bufio.NewReader(os.Stdin)

	fmt.Println(`
▀████▀ ▀███▀                           ▄▄█▀▀▀█▄█  ▄▄█▀▀██▄  
  ██   ▄█▀                           ▄██▀     ▀█▄██▀    ▀██▄
  ██ ▄█▀      ▄▄█▀██▀██▀   ▀██▀▄██▀████▀       ▀██▀      ▀██
  █████▄     ▄█▀   ██ ██   ▄█  ██   ▀▀█         ██        ██
  ██  ███    ██▀▀▀▀▀▀  ██ ▄█   ▀█████▄█▄    ▀█████▄      ▄██
  ██   ▀██▄  ██▄    ▄   ███    █▄   ████▄     ██▀██▄    ▄██▀
▄████▄   ███▄ ▀█████▀   ▄█     ██████▀ ▀▀███████  ▀▀████▀▀  
                      ▄█                                    
                    ██▀                                    `)

	for {
		fmt.Println("Escolha uma opção:")
		fmt.Println("1. Calcular tempo necessário para quebrar uma chave")
		fmt.Println("2. Calcular progresso com base no número de chaves já verificadas")
		fmt.Println("3. Estimativas de custo computacional (em desenvolvimento)")
		fmt.Println("4. Eficiência energética (em desenvolvimento)")
		fmt.Println("5. Ajuda")
		fmt.Println("6. Sair")
		fmt.Print("Opção: ")

		choiceStr, _ := reader.ReadString('\n')
		choiceStr = strings.TrimSpace(choiceStr)

		switch choiceStr {
		case "1":
			fmt.Print("Digite o número de bits da chave: ")
			bitsStr, _ := reader.ReadString('\n')
			bitsStr = strings.TrimSpace(bitsStr)
			bits, err := strconv.Atoi(bitsStr)
			if err != nil {
				fmt.Println("Entrada inválida para o número de bits")
				continue
			}

			fmt.Print("Digite a taxa de chaves verificadas por segundo (ex. 15M, 2.5k, 45,9k): ")
			keysPerSecondStr, _ := reader.ReadString('\n')
			keysPerSecondStr = strings.TrimSpace(keysPerSecondStr)

			keysPerSecond, err := parseKeysPerSecond(keysPerSecondStr)
			if err != nil {
				fmt.Println("Entrada inválida para a taxa de chaves verificadas por segundo")
				continue
			}

			timeToCrack := calculateTimeToCrack(bits, keysPerSecond)
			fmt.Printf("Levaria aproximadamente %.2f anos para testar todas as combinações de uma carteira de %d bits a uma taxa de %s chaves por segundo.\n", timeToCrack)

		case "2":
			fmt.Print("Digite o número de bits da chave: ")
			bitsStr, _ := reader.ReadString('\n')
			bitsStr = strings.TrimSpace(bitsStr)
			bits, err := strconv.Atoi(bitsStr)
			if err != nil {
				fmt.Println("Entrada inválida para o número de bits")
				continue
			}

			fmt.Print("Digite a taxa de chaves verificadas por segundo (ex. 15M, 2.5k, 45,9k): ")
			keysPerSecondStr, _ := reader.ReadString('\n')
			keysPerSecondStr = strings.TrimSpace(keysPerSecondStr)

			keysPerSecond, err := parseKeysPerSecond(keysPerSecondStr)
			if err != nil {
				fmt.Println("Entrada inválida para a taxa de chaves verificadas por segundo")
				continue
			}

			fmt.Print("Digite o número de chaves já verificadas: ")
			keysCheckedStr, _ := reader.ReadString('\n')
			keysCheckedStr = strings.TrimSpace(keysCheckedStr)
			keysChecked, ok := new(big.Int).SetString(keysCheckedStr, 10)
			if !ok {
				fmt.Println("Entrada inválida para o número de chaves verificadas")
				continue
			}

			timeToFinish := calculateProgress(bits, keysPerSecond, keysChecked)
			fmt.Printf("Levaria aproximadamente %.2f anos para testar as combinações restantes de uma carteira de %d bits a uma taxa de %s chaves por segundo.\n", timeToFinish)

		case "3":
			fmt.Println("Opção em desenvolvimento: Estimativas de custo computacional")
			// Implementar funcionalidade de estimativas de custo computacional aqui

		case "4":
			fmt.Println("Opção em desenvolvimento: Eficiência energética")
			// Implementar funcionalidade de eficiência energética aqui

		case "5", "ajuda":
			clearScreen()
			setTitle("KeysGo - Ajuda")

			fmt.Println("Bem-vindo ao KeysGo!")
			fmt.Println("Este programa oferece várias funcionalidades relacionadas à quebra de chaves criptográficas e cálculos relacionados.")
			fmt.Println("Aqui estão as opções disponíveis:")
			fmt.Println("1. Calcular tempo necessário para quebrar uma chave")
			fmt.Println("   - Esta opção calcula quanto tempo levaria para testar todas as combinações possíveis de uma chave criptográfica.")
			fmt.Println("2. Calcular progresso com base no número de chaves já verificadas")
			fmt.Println("   - Esta opção estima quanto tempo falta para testar as combinações restantes de uma chave criptográfica, com base em chaves já verificadas.")
			fmt.Println("3. Estimativas de custo computacional (em desenvolvimento)")
			fmt.Println("   - Opção futura para calcular o custo computacional estimado para quebrar uma chave criptográfica.")
			fmt.Println("4. Eficiência energética (em desenvolvimento)")
			fmt.Println("   - Opção futura para calcular a eficiência energética envolvida na quebra de chaves criptográficas.")
			fmt.Println("5. Ajuda")
			fmt.Println("   - Exibe esta mensagem de ajuda.")
			fmt.Println("6. Sair")
			fmt.Println("   - Encerra o programa.")

			fmt.Print("\nPressione Enter para voltar ao menu principal...")
			reader.ReadString('\n')
			clearScreen()
			setTitle("KeysGo")

		case "6", "sair":
			fmt.Println("Encerrando o programa...")
			return

		default:
			fmt.Println("Opção inválida. Escolha uma das opções disponíveis.")
		}
	}
}
