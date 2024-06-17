package ports

import (
	"context"
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

	file, err := os.Create("cotacao.txt")
	if err != nil{
		panic(err)
	}
	defer file.Close()
	_, err = file.WriteString("Dólar:" + string(v.GetStringBytes("bid")) )
	if err != nil{
		panic(err)
	}
}