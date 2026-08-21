package server

import (
	"errors"
	"net/http"

	"github.com/Piccio-Code/AlfredEventsOranizer/backend/internal/database"
	"github.com/gin-gonic/gin"
)

func (s *Server) registerPollTemplateRoute(rg *gin.RouterGroup) {
	pollTemplates := rg.Group("/poll-templates")
	{
		pollTemplates.POST("", s.CreatePollTemplate)
		pollTemplates.GET("", s.GetPollTemplates)

		pollTemplates.GET("/:id", s.GetPollTemplate)
		pollTemplates.PUT("/:id", s.UpdatePollTemplate)
		pollTemplates.DELETE("/:id", s.DeletePollTemplate)

		pollTemplates.POST("/:id/options", s.CreatePollTemplateOptions)
		pollTemplates.PUT("/:id/options/:optionId", s.UpdatePollTemplateOption)
		pollTemplates.DELETE("/:id/options/:optionId", s.DeletePollTemplateOption)
	}
}

func (s *Server) CreatePollTemplate(c *gin.Context) {
	var newTemplateRequest database.PollTemplateRequest

	if err := c.ShouldBindJSON(&newTemplateRequest); err != nil {
		s.infoLog.Println(err)
		s.Fail(c, http.StatusBadRequest, "Malformed JSON")
		return
	}

	if newTemplateRequest.Name == "" ||
		newTemplateRequest.Title == "" ||
		newTemplateRequest.GroupID == "" {
		s.Fail(c, http.StatusBadRequest, "Malformed JSON")
		return
	}

	newTemplate, err := s.models.PollTemplateModel.Create(c.Request.Context(), newTemplateRequest)
	if err != nil {
		s.infoLog.Println(err)
		s.Fail(c, http.StatusBadRequest, "Error creating the new poll template")
		return
	}

	s.Created(c, envelop{"new_poll_template": newTemplate})
}

func (s *Server) GetPollTemplates(c *gin.Context) {
	pollTemplates, err := s.models.PollTemplateModel.GetAll(c.Request.Context())
	if err != nil {
		s.infoLog.Println(err)
		s.Fail(c, http.StatusBadRequest, "Error listing the poll templates")
		return
	}

	s.Ok(c, envelop{"poll_templates": pollTemplates}, nil)
}

func (s *Server) GetPollTemplate(c *gin.Context) {
	id := c.Param("id")

	pollTemplate, err := s.models.PollTemplateModel.GetByID(c.Request.Context(), id)
	if errors.Is(err, database.NotFoundError) {
		s.Fail(c, http.StatusNotFound, "Poll template not found")
		return
	}
	if err != nil {
		s.infoLog.Println(err)
		s.Fail(c, http.StatusBadRequest, "Error getting the poll template")
		return
	}

	s.Ok(c, envelop{"poll_template": pollTemplate}, nil)
}

func (s *Server) UpdatePollTemplate(c *gin.Context) {
	id := c.Param("id")

	var updateRequest database.PollTemplateRequest
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		s.infoLog.Println(err)
		s.Fail(c, http.StatusBadRequest, "Malformed JSON")
		return
	}

	if updateRequest.Name == "" &&
		updateRequest.Title == "" &&
		updateRequest.GroupID == "" &&
		updateRequest.MultipleChoice == nil {
		s.Fail(c, http.StatusBadRequest, "Malformed JSON")
		return
	}

	pollTemplate, err := s.models.PollTemplateModel.GetByID(c.Request.Context(), id)
	if errors.Is(err, database.NotFoundError) {
		s.Fail(c, http.StatusNotFound, "Poll template not found")
		return
	}
	if err != nil {
		s.infoLog.Println(err)
		s.Fail(c, http.StatusBadRequest, "Error getting the poll template")
		return
	}

	name := pollTemplate.Name
	if updateRequest.Name != "" {
		name = updateRequest.Name
	}

	title := pollTemplate.Title
	if updateRequest.Title != "" {
		title = updateRequest.Title
	}

	groupID := pollTemplate.GroupID
	if updateRequest.GroupID != "" {
		groupID = updateRequest.GroupID
	}

	multipleChoice := pollTemplate.MultipleChoice
	if updateRequest.MultipleChoice != nil {
		multipleChoice = *updateRequest.MultipleChoice
	}

	err = s.models.PollTemplateModel.Update(
		c.Request.Context(),
		id,
		name,
		title,
		groupID,
		multipleChoice,
	)
	if errors.Is(err, database.NotFoundError) {
		s.Fail(c, http.StatusNotFound, "Poll template not found")
		return
	}
	if err != nil {
		s.infoLog.Println(err)
		s.Fail(c, http.StatusBadRequest, "Error updating the poll template")
		return
	}

	s.Updated(c)
}

func (s *Server) DeletePollTemplate(c *gin.Context) {
	id := c.Param("id")

	err := s.models.PollTemplateModel.Delete(c.Request.Context(), id)
	if errors.Is(err, database.NotFoundError) {
		s.Fail(c, http.StatusNotFound, "Poll template not found")
		return
	}
	if err != nil {
		s.infoLog.Println(err)
		s.Fail(c, http.StatusBadRequest, "Error deleting the poll template")
		return
	}

	s.Delete(c)
}

func (s *Server) CreatePollTemplateOptions(c *gin.Context) {
	templateID := c.Param("id")

	var newOptionsRequest []database.PollTemplateOptionRequest
	if err := c.ShouldBindJSON(&newOptionsRequest); err != nil {
		s.infoLog.Println(err)
		s.Fail(c, http.StatusBadRequest, "Malformed JSON")
		return
	}

	if len(newOptionsRequest) == 0 {
		s.Fail(c, http.StatusBadRequest, "Malformed JSON")
		return
	}

	for _, option := range newOptionsRequest {
		if option.Label == "" {
			s.Fail(c, http.StatusBadRequest, "Malformed JSON")
			return
		}
	}

	newOptions, err := s.models.PollTemplateModel.AddOptions(
		c.Request.Context(),
		templateID,
		newOptionsRequest,
	)
	if errors.Is(err, database.NotFoundError) {
		s.Fail(c, http.StatusNotFound, "Poll template not found")
		return
	}
	if err != nil {
		s.infoLog.Println(err)
		s.Fail(c, http.StatusBadRequest, "Error creating the poll template options")
		return
	}

	s.Created(c, envelop{"new_poll_template_options": newOptions})
}

func (s *Server) UpdatePollTemplateOption(c *gin.Context) {
	templateID := c.Param("id")
	optionID := c.Param("optionId")

	var updateRequest database.PollTemplateOptionRequest
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		s.infoLog.Println(err)
		s.Fail(c, http.StatusBadRequest, "Malformed JSON")
		return
	}

	if updateRequest.Label == "" &&
		updateRequest.MotivationNeeded == nil &&
		updateRequest.CongratulationNeeded == nil &&
		updateRequest.SpecificationNeeded == nil {
		s.Fail(c, http.StatusBadRequest, "Malformed JSON")
		return
	}

	option, err := s.models.PollTemplateModel.GetOptionByID(
		c.Request.Context(),
		templateID,
		optionID,
	)
	if errors.Is(err, database.NotFoundError) {
		s.Fail(c, http.StatusNotFound, "Poll template option not found")
		return
	}
	if err != nil {
		s.infoLog.Println(err)
		s.Fail(c, http.StatusBadRequest, "Error getting the poll template option")
		return
	}

	label := option.Label
	if updateRequest.Label != "" {
		label = updateRequest.Label
	}

	motivationNeeded := option.MotivationNeeded
	if updateRequest.MotivationNeeded != nil {
		motivationNeeded = *updateRequest.MotivationNeeded
	}

	congratulationNeeded := option.CongratulationNeeded
	if updateRequest.CongratulationNeeded != nil {
		congratulationNeeded = *updateRequest.CongratulationNeeded
	}

	specificationNeeded := option.SpecificationNeeded
	if updateRequest.SpecificationNeeded != nil {
		specificationNeeded = *updateRequest.SpecificationNeeded
	}

	err = s.models.PollTemplateModel.UpdateOption(
		c.Request.Context(),
		templateID,
		optionID,
		label,
		motivationNeeded,
		congratulationNeeded,
		specificationNeeded,
	)
	if errors.Is(err, database.NotFoundError) {
		s.Fail(c, http.StatusNotFound, "Poll template option not found")
		return
	}
	if err != nil {
		s.infoLog.Println(err)
		s.Fail(c, http.StatusBadRequest, "Error updating the poll template option")
		return
	}

	s.Updated(c)
}

func (s *Server) DeletePollTemplateOption(c *gin.Context) {
	templateID := c.Param("id")
	optionID := c.Param("optionId")

	err := s.models.PollTemplateModel.DeleteOption(c.Request.Context(), templateID, optionID)
	if errors.Is(err, database.NotFoundError) {
		s.Fail(c, http.StatusNotFound, "Poll template option not found")
		return
	}
	if err != nil {
		s.infoLog.Println(err)
		s.Fail(c, http.StatusBadRequest, "Error deleting the poll template option")
		return
	}

	s.Delete(c)
}
