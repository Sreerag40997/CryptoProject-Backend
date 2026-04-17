package seed

import (
	"cryptox/internal/modules/auth"
	"cryptox/packages/utils"
	"fmt"

	"gorm.io/gorm"
)

func Addadmin(db *gorm.DB) {

  password := "admin123"

  hashed, _ := utils.Hashing(password)

  admin := auth.User{
    Name: "Admin",
    Email: "admincrytinox@gmail.com",
    Password: hashed,
		IsVerified: true,
		ProfilePicURL: "https://img.freepik.com/free-vector/follow-me-social-business-theme-design_24877-50426.jpg?semt=ais_user_personalization&w=740&q=80",
    Role: "admin",
  }

  err:=db.Where(auth.User{Email: "admincrytinox@gmail.com"}).FirstOrCreate(&admin).Error
	if err!=nil{
		fmt.Println("admin not added",err)
	}

  fmt.Println("Admin seeded successfully")
}