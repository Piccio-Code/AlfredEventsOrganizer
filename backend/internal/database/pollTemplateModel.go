package database

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PollTemplateModel struct {
	DB *pgxpool.Pool
}

type PollTemplateDB struct {
	ID             string    `json:"id,omitempty"`
	Name           string    `json:"name,omitempty"`
	Title          string    `json:"title,omitempty"`
	GroupID        string    `json:"groupId,omitempty"`
	MultipleChoice bool      `json:"multipleChoice"`
	CreatedAt      time.Time `json:"createdAt"`
}

type PollTemplateRequest struct {
	Name           string `json:"name,omitempty"`
	Title          string `json:"title,omitempty"`
	GroupID        string `json:"groupId,omitempty"`
	MultipleChoice *bool  `json:"multipleChoice,omitempty"`
}

type PollTemplateOptionRequest struct {
	Label                string `json:"label,omitempty"`
	MotivationNeeded     *bool  `json:"motivationNeeded,omitempty"`
	CongratulationNeeded *bool  `json:"congratulationNeeded,omitempty"`
}

type PollTemplateOption struct {
	ID                   string `json:"id,omitempty"`
	TemplateID           string `json:"templateId,omitempty"`
	Label                string `json:"label,omitempty"`
	Position             int    `json:"position"`
	MotivationNeeded     bool   `json:"motivationNeeded"`
	CongratulationNeeded bool   `json:"congratulationNeeded"`
}

type PollTemplateResponse struct {
	ID             string               `json:"id,omitempty"`
	Name           string               `json:"name,omitempty"`
	Title          string               `json:"title,omitempty"`
	GroupID        string               `json:"groupId,omitempty"`
	MultipleChoice bool                 `json:"multipleChoice"`
	Options        []PollTemplateOption `json:"options,omitempty"`
	CreatedAt      time.Time            `json:"createdAt"`
}

func (m *PollTemplateModel) Create(ctx context.Context, newTemplate PollTemplateRequest) (PollTemplateDB, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	multipleChoice := false
	if newTemplate.MultipleChoice != nil {
		multipleChoice = *newTemplate.MultipleChoice
	}

	stmt := `INSERT INTO poll_templates(name, title, group_id, multiple_choice)
			 VALUES ($1, $2, $3, $4)
			 RETURNING id, name, title, group_id, multiple_choice, created_at`

	var createdTemplate PollTemplateDB
	err := m.DB.QueryRow(
		ctx,
		stmt,
		newTemplate.Name,
		newTemplate.Title,
		newTemplate.GroupID,
		multipleChoice,
	).Scan(
		&createdTemplate.ID,
		&createdTemplate.Name,
		&createdTemplate.Title,
		&createdTemplate.GroupID,
		&createdTemplate.MultipleChoice,
		&createdTemplate.CreatedAt,
	)
	if err != nil {
		return PollTemplateDB{}, err
	}

	return createdTemplate, nil
}

func (m *PollTemplateModel) GetAll(ctx context.Context) ([]PollTemplateResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	stmt := `SELECT id, name, title, group_id, multiple_choice, created_at
			 FROM poll_templates
			 ORDER BY created_at DESC`

	rows, err := m.DB.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	templates := []PollTemplateResponse{}
	templateIndex := map[string]int{}

	for rows.Next() {
		var template PollTemplateResponse
		err := rows.Scan(
			&template.ID,
			&template.Name,
			&template.Title,
			&template.GroupID,
			&template.MultipleChoice,
			&template.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		template.Options = []PollTemplateOption{}
		templateIndex[template.ID] = len(templates)
		templates = append(templates, template)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	rows.Close()

	stmt = `SELECT id, template_id, label, position, motivation_needed, congratulation_needed
			FROM poll_template_options
			ORDER BY template_id, position`

	optionRows, err := m.DB.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer optionRows.Close()

	for optionRows.Next() {
		var option PollTemplateOption
		err := optionRows.Scan(
			&option.ID,
			&option.TemplateID,
			&option.Label,
			&option.Position,
			&option.MotivationNeeded,
			&option.CongratulationNeeded,
		)
		if err != nil {
			return nil, err
		}

		idx, ok := templateIndex[option.TemplateID]
		if !ok {
			continue
		}

		templates[idx].Options = append(templates[idx].Options, option)
	}

	if err := optionRows.Err(); err != nil {
		return nil, err
	}

	return templates, nil
}

func (m *PollTemplateModel) GetByID(ctx context.Context, id string) (PollTemplateResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	stmt := `SELECT id, name, title, group_id, multiple_choice, created_at
			 FROM poll_templates
			 WHERE id = $1`

	var template PollTemplateResponse
	err := m.DB.QueryRow(ctx, stmt, id).Scan(
		&template.ID,
		&template.Name,
		&template.Title,
		&template.GroupID,
		&template.MultipleChoice,
		&template.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return PollTemplateResponse{}, NotFoundError
	}
	if err != nil {
		return PollTemplateResponse{}, err
	}

	stmt = `SELECT id, template_id, label, position, motivation_needed, congratulation_needed
			FROM poll_template_options
			WHERE template_id = $1
			ORDER BY position`

	rows, err := m.DB.Query(ctx, stmt, id)
	if err != nil {
		return PollTemplateResponse{}, err
	}
	defer rows.Close()

	template.Options = []PollTemplateOption{}
	for rows.Next() {
		var option PollTemplateOption
		err := rows.Scan(
			&option.ID,
			&option.TemplateID,
			&option.Label,
			&option.Position,
			&option.MotivationNeeded,
			&option.CongratulationNeeded,
		)
		if err != nil {
			return PollTemplateResponse{}, err
		}

		template.Options = append(template.Options, option)
	}

	if err := rows.Err(); err != nil {
		return PollTemplateResponse{}, err
	}

	return template, nil
}

func (m *PollTemplateModel) Update(
	ctx context.Context,
	id string,
	name string,
	title string,
	groupID string,
	multipleChoice bool,
) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	stmt := `UPDATE poll_templates
			 SET name = $1, title = $2, group_id = $3, multiple_choice = $4
			 WHERE id = $5`

	cmd, err := m.DB.Exec(ctx, stmt, name, title, groupID, multipleChoice, id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return NotFoundError
	}

	return nil
}

func (m *PollTemplateModel) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	stmt := `DELETE FROM poll_templates WHERE id = $1`

	cmd, err := m.DB.Exec(ctx, stmt, id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return NotFoundError
	}

	return nil
}

func (m *PollTemplateModel) AddOptions(
	ctx context.Context,
	templateID string,
	newOptions []PollTemplateOptionRequest,
) ([]PollTemplateOption, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tx, err := m.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var lockedTemplateID string
	err = tx.QueryRow(
		ctx,
		`SELECT id FROM poll_templates WHERE id = $1 FOR UPDATE`,
		templateID,
	).Scan(&lockedTemplateID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, NotFoundError
	}
	if err != nil {
		return nil, err
	}

	var lastPosition int
	err = tx.QueryRow(
		ctx,
		`SELECT COALESCE(MAX(position), -1)
		 FROM poll_template_options
		 WHERE template_id = $1`,
		templateID,
	).Scan(&lastPosition)
	if err != nil {
		return nil, err
	}

	createdOptions := make([]PollTemplateOption, 0, len(newOptions))
	stmt := `INSERT INTO poll_template_options(
				template_id,
				label,
				position,
				motivation_needed,
				congratulation_needed
			 )
			 VALUES ($1, $2, $3, $4, $5)
			 RETURNING id, template_id, label, position, motivation_needed, congratulation_needed`

	for index, newOption := range newOptions {
		motivationNeeded := newOption.MotivationNeeded != nil && *newOption.MotivationNeeded
		congratulationNeeded := newOption.CongratulationNeeded != nil && *newOption.CongratulationNeeded

		var createdOption PollTemplateOption
		err = tx.QueryRow(
			ctx,
			stmt,
			templateID,
			newOption.Label,
			lastPosition+index+1,
			motivationNeeded,
			congratulationNeeded,
		).Scan(
			&createdOption.ID,
			&createdOption.TemplateID,
			&createdOption.Label,
			&createdOption.Position,
			&createdOption.MotivationNeeded,
			&createdOption.CongratulationNeeded,
		)
		if err != nil {
			return nil, err
		}

		createdOptions = append(createdOptions, createdOption)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return createdOptions, nil
}

func (m *PollTemplateModel) GetOptionByID(
	ctx context.Context,
	templateID string,
	optionID string,
) (PollTemplateOption, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	stmt := `SELECT id, template_id, label, position, motivation_needed, congratulation_needed
			 FROM poll_template_options
			 WHERE template_id = $1 AND id = $2`

	var option PollTemplateOption
	err := m.DB.QueryRow(ctx, stmt, templateID, optionID).Scan(
		&option.ID,
		&option.TemplateID,
		&option.Label,
		&option.Position,
		&option.MotivationNeeded,
		&option.CongratulationNeeded,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return PollTemplateOption{}, NotFoundError
	}
	if err != nil {
		return PollTemplateOption{}, err
	}

	return option, nil
}

func (m *PollTemplateModel) UpdateOption(
	ctx context.Context,
	templateID string,
	optionID string,
	label string,
	motivationNeeded bool,
	congratulationNeeded bool,
) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	stmt := `UPDATE poll_template_options
			 SET label = $1, motivation_needed = $2, congratulation_needed = $3
			 WHERE template_id = $4 AND id = $5`

	cmd, err := m.DB.Exec(
		ctx,
		stmt,
		label,
		motivationNeeded,
		congratulationNeeded,
		templateID,
		optionID,
	)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return NotFoundError
	}

	return nil
}

func (m *PollTemplateModel) DeleteOption(ctx context.Context, templateID string, optionID string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	stmt := `DELETE FROM poll_template_options WHERE template_id = $1 AND id = $2`

	cmd, err := m.DB.Exec(ctx, stmt, templateID, optionID)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return NotFoundError
	}

	return nil
}
