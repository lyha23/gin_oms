package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wejectchen/ginblog/model"
)

func GetATASelection(c *gin.Context) {
	var ctl model.ATAEvent
	if err := ctl.GetAll(); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
	}
	c.JSON(
		http.StatusOK, gin.H{
			"code":   http.StatusOK,
			"result": ctl.ATAList,
		},
	)
}

func GetPartSelection(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("equipment_id"))
	var ctl model.EquipmentEvent
	if err := ctl.GetListByEquipment(id); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"code":   http.StatusOK,
		"result": ctl.Equipment.Parts,
	})
}

func LoadATAEquipment(c *gin.Context) {
	var data model.LoadATAEquipmentParam
	_ = c.ShouldBindJSON(&data)
	var ctl model.PartLoadLogEvent
	if err := ctl.SaveLog(data); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"code":   http.StatusOK,
		"result": "success",
	})
}

func GetAllLoadStatus(c *gin.Context){
	var ctl model.PartLoadLogEvent
	if err := ctl.GetAll(); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"code":   http.StatusOK,
		"result": ctl.PartLoadLogList,
	})
}
