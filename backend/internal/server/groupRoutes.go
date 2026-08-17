package server

import (
	"errors"
	"github.com/Piccio-Code/AlfredEventsOranizer/backend/internal/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) registerGroupRoute(rg *gin.RouterGroup) {
	groups := rg.Group("/groups")
	{
		groups.POST("", s.CreateGroup)
		groups.GET("", s.GetGroups)
		groups.GET("/:id", s.GetGroup)
		groups.PATCH("/:id", s.UpdateGroup)
		groups.PUT("/:id", s.UpdateGroup)
		groups.DELETE("/:id", s.DeleteGroup)
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

func (s *Server) GetGroups(c *gin.Context) {
	groups, err := s.models.GroupModel.GetAll(c.Request.Context())

	if err != nil {
		s.infoLog.Println(err)
		s.Fail(c, http.StatusBadRequest, "Error listing the groups")
		return
	}

	s.Ok(c, envelop{"groups": groups}, nil)
}

func (s *Server) GetGroup(c *gin.Context) {
	id := c.Param("id")

	group, err := s.models.GroupModel.GetByID(c.Request.Context(), id)

	if errors.Is(err, database.NotFoundError) {
		s.Fail(c, http.StatusNotFound, "Group not found")
		return
	}

	if err != nil {
		s.infoLog.Println(err)
		s.Fail(c, http.StatusBadRequest, "Error getting the group")
		return
	}

	s.Ok(c, envelop{"group": group}, nil)
}

func (s *Server) UpdateGroup(c *gin.Context) {
	id := c.Param("id")

	var updateRequest database.GroupRequest

	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		s.infoLog.Println(err)
		s.Fail(c, http.StatusBadRequest, "Malformed JSON")
		return
	}

	if updateRequest.Title == "" && updateRequest.WhatsappChatId == "" {
		s.Fail(c, http.StatusBadRequest, "Malformed JSON")
		return
	}

	group, err := s.models.GroupModel.GetByID(c.Request.Context(), id)

	if errors.Is(err, database.NotFoundError) {
		s.Fail(c, http.StatusNotFound, "Group not found")
		return
	}

	if err != nil {
		s.infoLog.Println(err)
		s.Fail(c, http.StatusBadRequest, "Error getting the group")
		return
	}

	title := group.Title
	if updateRequest.Title != "" {
		title = updateRequest.Title
	}

	whatsappChatId := group.WhatsappChatId
	if updateRequest.WhatsappChatId != "" {
		whatsappChatId = updateRequest.WhatsappChatId
	}

	if whatsappChatId != group.WhatsappChatId {
		groupDetail, wahaErr := s.wahaClient.GetGroupDetail(whatsappChatId)

		if wahaErr != nil {
			s.infoLog.Println(wahaErr)
			s.Fail(c, http.StatusBadRequest, "Error getting the group detail")
			return
		}

		err = s.models.GroupModel.Update(c.Request.Context(), id, title, whatsappChatId, groupDetail.Participants)
	} else {
		err = s.models.GroupModel.Update(c.Request.Context(), id, title, whatsappChatId, nil)
	}

	if errors.Is(err, database.NotFoundError) {
		s.Fail(c, http.StatusNotFound, "Group not found")
		return
	}

	if err != nil {
		s.infoLog.Println(err)
		s.Fail(c, http.StatusBadRequest, "Error updating the group")
		return
	}

	s.Updated(c)
}

func (s *Server) DeleteGroup(c *gin.Context) {
	id := c.Param("id")

	err := s.models.GroupModel.Delete(c.Request.Context(), id)

	if errors.Is(err, database.NotFoundError) {
		s.Fail(c, http.StatusNotFound, "Group not found")
		return
	}

	if err != nil {
		s.infoLog.Println(err)
		s.Fail(c, http.StatusBadRequest, "Error deleting the group")
		return
	}

	s.Delete(c)
}
