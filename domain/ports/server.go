package ports

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"client.server.api/domain/model"
	"client.server.api/domain/repository"
	"github.com/valyala/fastjson"
)



func CriarServer()  {

	mux := http.NewServeMux()
	mux.HandleFunc("GET /cotacao", ConsultaCotacao)
	http.ListenAndServe(":8080",mux)

}

func ConsultaCotacao(w http.ResponseWriter, r *http.Request){

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	select{
	case <-time.After(200*time.Millisecond) :
		log.Fatal("Cancelado por timeout na request do Server Cotacao!")
	case <- ctx.Done():
	}

	db, err := repository.AbrirConexao()
	if err != nil{
		panic(err)
	}
	var p fastjson.Parser	

	body, err := io.ReadAll(res.Body)
	if err != nil{
		panic(err)
	}

	v, err := p.ParseBytes(body)
	if err != nil{
		panic(err)
	}	

	cotacao:= model.Cotacao{}
	json.Unmarshal([]byte(v.GetObject("USDBRL").String()), &cotacao)

	repository.InsertCotacao(db, &cotacao)

	bid := model.Bid{Bid: cotacao.Bid}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bid)

}

