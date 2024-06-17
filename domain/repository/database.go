package repository

import (
	"context"
	"log"
	"time"

	"client.server.api/domain/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)


func AbrirConexao() (*gorm.DB, error){
	db, err := gorm.Open(sqlite.Open("cotacao.db"), &gorm.Config{})
	if err != nil{
		panic(err)
	}

	return db, nil	
	
}

func CriarTabela[T any](db *gorm.DB, t T ) error{

	err := db.AutoMigrate(&t)
	if err != nil{
		return  err
	}
	return nil	

	
	
}

func InsertCotacao(db *gorm.DB, cotacao *model.Cotacao ) error{

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	db.WithContext(ctx).Create(&cotacao)

	select{
	case <-time.After(10*time.Millisecond) :
		log.Fatal("Cancelado por timeout na insert Database!")
	case <- ctx.Done():
	}
	return nil

}
