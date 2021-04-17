package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitoramento = 3
const delay = 5

//https://golang.org/src/time/format.go
var dateNow = time.Now().Format("02/01/2006 15:04:05")
var dateLog = time.Now().Format("02012006")
var logFile = "log_" + dateLog + ".txt"

func main() {
	exibeIntroducao()
	//leSitesDoArquivo()
	//registraLog("site-falso", false)

	for {

		exibeMenu()
		comando := lerComando()

		switch comando {
		case 1:
			iniciarMonitoramento()
		case 2:
			imprimeLogs()
		case 0:
			fmt.Println("Saindo do programa.")
			os.Exit(0) // informar uma saída bem sucedida
		default:
			fmt.Println("Comando não reconhecido.")
			os.Exit(-1)
		}
	}
}

func exibeIntroducao() {
	version := 1.2
	fmt.Println("Bem Vindo!", dateNow)
	fmt.Println("Este programa esta na versão", version)
}

func exibeMenu() {
	fmt.Println("1- Iniciar Monitoramento")
	fmt.Println("2- Exibir Logs")
	fmt.Println("0- Sair do Programa")
}

func lerComando() int {
	var comandoLido int
	fmt.Scan(&comandoLido)
	fmt.Println("O comando escolhido foi", comandoLido)
	fmt.Println("")

	return comandoLido
}

//& é o endereço da variável que queremos salvar a entrada,
//pois a função Scanf não espera uma variável, e sim o seu endereço, um ponteiro para a variável

func iniciarMonitoramento() {
	fmt.Println("Monitorando...")

	//Declarando slice ("array")
	/* 	sites := []string{
		"https://random-status-code.herokuapp.com/",
	} */

	sites := leSitesDoArquivo()

	for i := 0; i < monitoramento; i++ {
		for i, site := range sites {
			fmt.Println("Testando site", i+1, ":", site)
			testaSite(site)
		}
		time.Sleep(delay * time.Second)
		fmt.Println("")
	}

	fmt.Println("")
}

func testaSite(site string) {
	resp, err := http.Get(site)

	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}

	if resp.StatusCode == 200 {
		fmt.Println("Site", site, "foi carregado com sucesso!")
		registraLog(site, true)
	} else {
		fmt.Println("Site: '", site, "' está com problemas. Status code:", resp.StatusCode)
		registraLog(site, false)
	}
}

func leSitesDoArquivo() []string {

	var sites []string

	arquivo, err := os.Open("sites.txt")
	//arquivo, err := ioutil.ReadFile("sites.txt")

	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}

	leitor := bufio.NewReader(arquivo)
	for {
		linha, err := leitor.ReadString('\n')
		linha = strings.TrimSpace(linha)

		sites = append(sites, linha)

		if err == io.EOF {
			// end of file
			break
		}
	}

	arquivo.Close() // antes de finalizar o programa, por boa prática, precisa fechar o arquivo

	return sites
}

func registraLog(site string, status bool) {
	date := time.Now().Format("02/01/2006 15:04:05")

	arquivo, err := os.OpenFile(logFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println(err)
	}

	arquivo.WriteString(
		date +
			" - " +
			site +
			" - online: " +
			strconv.FormatBool(status) +
			"\n")

	arquivo.Close()
}

func imprimeLogs() {
	fmt.Println("Exibindo logs..")

	arquivo, err := ioutil.ReadFile(logFile)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(arquivo))
}
