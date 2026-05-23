package rooms

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRoomBody struct {
	Name string `json:"name"`
}

func CreateRoom(ctx *gin.Context) {
	var body CreateRoomBody

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	room := Room{
		ID:   uuid.NewString(),
		Name: body.Name,
	}

	Rooms = append(Rooms, room)

	ctx.JSON(http.StatusCreated, room)
}

func GetRooms(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Rooms)
}
