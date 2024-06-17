package main

import (
	"client.server.api/domain/model"
	"client.server.api/domain/ports"
	"client.server.api/domain/repository"
)

func main(){

	db, err := repository.AbrirConexao()
	if err != nil{
		panic(err)
	}
	repository.CriarTabela(db, model.Cotacao{})

	go ports.CriarServer()
	ports.ConsultaCotacaoDolar()

}