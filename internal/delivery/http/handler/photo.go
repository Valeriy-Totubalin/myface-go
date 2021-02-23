package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Valeriy-Totubalin/myface-go/internal/delivery/http/request"
	"github.com/Valeriy-Totubalin/myface-go/internal/service/photo_service"
	"github.com/gin-gonic/gin"
)

func (h *Handler) upload(c *gin.Context) {
	var data request.Upload
	if err := c.ShouldBindJSON(&data); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userId := c.MustGet("user_id")
	if nil == userId {
		return
	}
	id, _ := strconv.Atoi(userId.(string))

	service, err := photo_service.NewPhotoService()
	if nil != err {
		log.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": UNKNOW_ERROR})
		return
	}

	err = service.CheckCorrectData(data.Photo)
	if nil != err {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = service.Upload(id, data.Photo)
	if nil != err {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func (h *Handler) get(c *gin.Context) {
	var photoInput request.Photo
	if err := c.ShouldBindUri(&photoInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	id, err := strconv.Atoi(photoInput.Id)
	if nil != err {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	userId := c.MustGet("user_id")
	if nil == userId {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user id",
		})
		return
	}
	userIdInt, _ := strconv.Atoi(userId.(string))

	service, err := photo_service.NewPhotoService()
	if nil != err {
		log.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": UNKNOW_ERROR})
		return
	}

	canGet, err := service.CanGet(userIdInt, id)
	if nil != err {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if false == canGet {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Photo not found for this user",
		})
		return
	}

	base64, err := service.GetById(id)
	if nil != err {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"photo_id": photoInput.Id,
		"photo":    base64,
	})
	return
}