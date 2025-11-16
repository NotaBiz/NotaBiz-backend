package utils

import (
	"NotaBiz-backend/config"
	"NotaBiz-backend/model/entity"
	"log"
	"time"

	"github.com/jasonlvhit/gocron"
)

func RunScheduler() {
	gocron.Every(30).Minutes().From(gocron.NextTick()).Do(deleteNonVerifiedUsers)
	<- gocron.Start()
}

func deleteNonVerifiedUsers() {
	result := config.DB.Where("is_verified = ? AND updated_at >= ?", false, time.Now().Add(24*time.Hour)).Delete(&entity.MasterUser{})
	if result.Error != nil {
		log.Println("error deleting non-verified user", result.Error.Error())
		return
	}
	log.Println("successfully deleted non-verified users:", result.RowsAffected)
}
