package server

import (
	"github.com/Piccio-Code/AlfredEventsOranizer/backend/internal/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) registerGroupRoute(rg *gin.RouterGroup) {
	groups := rg.Group("/groups")
	{
		groups.POST("", s.CreateGroup)
	}
}

func (s *Server) CreateGroup(c *gin.Context) {
	var newGroupRequest database.GroupRequest

	if err := c.ShouldBindJSON(&newGroupRequest); err != nil {
		s.infoLog.Println(err)
		s.Fail(c, http.StatusBadRequest, "Malformed JSON")
		return
	}

	if err := s.wahaClient.IsValidID(newGroupRequest.WhatsappChatId); err != nil {
		s.infoLog.Println(err)
		s.Fail(c, http.StatusBadRequest, "Error verifying the group ID")
		return
	}

	newGroup, err := s.models.GroupModel.Create(c.Request.Context(), newGroupRequest)

	if err != nil {
		s.infoLog.Println(err)
		s.Fail(c, http.StatusBadRequest, "Error creating the new group")
		return
	}

	s.Created(c, envelop{"new_group": newGroup})
}
