package meetinghandler

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	meetingmodel "github.com/theRAGEhero/Democracy-Routes/feature/discussion/server/meetinghandler/model"
)

type Handler struct {
	db *sql.DB
}

func New(db *sql.DB) (*Handler, error) {
	if db == nil {
		return nil, fmt.Errorf("no db")
	}

	return &Handler{
		db: db,
	}, nil
}

func (h *Handler) Create(params meetingmodel.CreateMeeting) (*meetingmodel.Meeting, error) {
	var meeting meetingmodel.Meeting
	meeting.ID = uuid.NewString()
	meeting.Title = params.Title

	_, err := h.db.Exec("INSERT INTO meeting (id, title) VALUES ($1, $2)", meeting.ID, meeting.Title)
	if err != nil {
		return nil, fmt.Errorf("create meeting: %w", err)
	}

	return &meeting, nil
}
