package ports

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/valyala/fastjson"
)




func ConsultaCotacaoDolar(){


	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	select{
	case <-time.After(300*time.Millisecond) :
		log.Fatal("Cancelado por timeout na request do Server!")
	case <- ctx.Done():
	}

	body, err := io.ReadAll(res.Body)
	if err != nil{
		panic(err)
	}

	gravarArquivo(body)
	


}

func gravarArquivo(body []byte){

	var p fastjson.Parser

	v, err := p.ParseBytes(body)
	if err != nil{
		panic(err)
	}

	file, err := os.OpenFile("cotacao.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil{
		panic(err)
	}
	defer file.Close()

	linha := fmt.Sprintf("Dólar: %v \n", string(v.GetStringBytes("bid")))

	_, err = file.WriteString(linha)
	if err != nil{
		panic(err)
	}
}