package handler

import (
	"fmt"
	"fsb-mse-iad591-server/datasource/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type Handler struct {
	DB *gorm.DB
}

func New(db *gorm.DB) Handler {
	return Handler{db}
}

func (h Handler) GetTemperatureRecords(ctx *gin.Context) {
	var records []models.TemperatureRecord

	if result := h.DB.Find(&records); result.Error != nil {
		fmt.Println(result.Error)
	}

	//w.Header().Add("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)
	//err := json.NewEncoder(w).Encode(records)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//ctx.IndentedJSON(http.StatusOK, records)
	ctx.IndentedJSON(http.StatusOK, records)
}
