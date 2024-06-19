package v1

import (
	"net/http"

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
