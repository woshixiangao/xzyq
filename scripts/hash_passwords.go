package main

import (
	"log"
	"xzyq/database"
	"xzyq/models"
	"xzyq/utils"
)

func main() {
	// 初始化数据库连接
	database.InitDB()

	// 获取所有用户
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		log.Fatalf("获取用户列表失败: %v", err)
	}

	// 更新每个用户的密码
	for _, user := range users {
		// 如果密码看起来不像是已经加密的（长度小于50），则进行加密
		if len(user.Password) < 50 {
			hashedPassword, err := utils.HashPassword(user.Password)
			if err != nil {
				log.Printf("加密用户 %s 的密码失败: %v", user.Username, err)
				continue
			}

			// 更新用户密码
			if err := database.DB.Model(&user).Update("password", hashedPassword).Error; err != nil {
				log.Printf("更新用户 %s 的密码失败: %v", user.Username, err)
				continue
			}

			log.Printf("成功更新用户 %s 的密码", user.Username)
		}
	}

	log.Println("密码更新完成")
}
