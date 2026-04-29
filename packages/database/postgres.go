package database

import (
	"cryptox/internal/modules/auth"
	cashwallet "cryptox/internal/modules/cah_wallet"
	cryptowallet "cryptox/internal/modules/crypto_wallet"
	ecard "cryptox/internal/modules/e_card"
	"cryptox/internal/modules/kyc"
	"cryptox/packages/config"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresConnection(cfg *config.Config)(*gorm.DB,error){
	//connect postgres
	db,err:=gorm.Open(postgres.Open(cfg.DBUrl),&gorm.Config{})
	if err!=nil{
		return nil,err
	}
	db.Debug()

	sqlDB,err:=db.DB()
	if err!=nil{
		return nil,err
	}
	//set connection pooling for reuse and limit
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Connect to PostgreSQL")

	err = db.AutoMigrate(
		&auth.User{},
		&kyc.KYC{},
		&ecard.Card{},
		&cashwallet.Wallet{},
		&cashwallet.WalletTransaction{},
		&cryptowallet.CryptoWallet{},
		&cryptowallet.CryptoAsset{},
		&cryptowallet.CryptoWalletTransaction{},
	)
	if err != nil {
		log.Fatal(err)
	}

	return db,nil
}
