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

	groupDetail, err := s.wahaClient.GetGroupDetail(newGroupRequest.WhatsappChatId)

	if err != nil {
		s.infoLog.Println(err)
		s.Fail(c, http.StatusBadRequest, "Error getting the group detail")
		return
	}

	if newGroupRequest.Title == "" {
		newGroupRequest.Title = groupDetail.Name
	}

	newGroup, err := s.models.GroupModel.Create(c.Request.Context(), newGroupRequest, groupDetail.Participants)

	if err != nil {
		s.infoLog.Println(err)
		s.Fail(c, http.StatusBadRequest, "Error creating the new group")
		return
	}

	s.Created(c, envelop{"new_group": newGroup})
}
