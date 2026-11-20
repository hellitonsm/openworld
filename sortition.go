package main

import (
	"bufio"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const tipURL = "https://mempool.space/api/blocks/tip/hash"

func fetchTipHash() (string, error) {
	resp, err := http.Get(tipURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	buf := make([]byte, 256)
	n, _ := resp.Body.Read(buf)
	h := strings.TrimSpace(string(buf[:n]))
	return h, nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Sistema de Sorteio (Bitcoin hash)")
	fmt.Print("Capturando hash atual... ")
	last, err := fetchTipHash()
	if err != nil {
		fmt.Println("erro:", err)
		return
	}
	fmt.Println(last)

	fmt.Print("Informe X (numero maximo): ")
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	x, err := strconv.Atoi(line)
	if err != nil || x < 1 {
		fmt.Println("X invalido")
		return
	}

	fmt.Println("Aguardando novo hash (novo bloco)...")
	var newHash string
	for {
		time.Sleep(15 * time.Second)
		h, err := fetchTipHash()
		if err != nil {
			fmt.Println("erro:", err)
			continue
		}
		if h != last {
			newHash = h
			fmt.Println("Novo hash:", newHash)
			break
		}
		fmt.Print(".")
	}

	n := new(big.Int)
	n.SetString(newHash, 16)
	xBig := big.NewInt(int64(x))
	mod := new(big.Int).Mod(n, xBig)
	winner := int(mod.Int64()) + 1
	fmt.Printf("Numero sorteado: %d\n", winner)
}
