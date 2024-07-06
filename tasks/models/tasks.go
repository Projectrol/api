package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/lehoangvuvt/projectrol/common"
	pb "github.com/lehoangvuvt/projectrol/common/protos"
)

type RecurringType string

const (
	MINUTELY RecurringType = "MINUTELY"
	HOURLY   RecurringType = "HOURLY"
	DAILY    RecurringType = "DAILY"
	WEEKLY   RecurringType = "WEEKLY"
	MONTHLY  RecurringType = "MONTHLY"
	YEARLY   RecurringType = "YEARLY"
)

type TaskType string

const (
	MEETING TaskType = "MEETING"
)

type Task struct {
	Id          int            `json:"id"`
	NanoId      string         `json:"nanoid"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Duration    int            `json:"duration"` //in seconds
	Type        TaskType       `json:"type"`
	WorkspaceId int            `json:"workspace_id"`
	CreatedBy   int            `json:"created_by"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   sql.NullString `json:"updated_at"`
}

type CreateTaskInput struct {
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Duration    int                 `json:"duration"`
	DtStart     int64               `json:"dtstart"`
	Recurring   *pb.RecurringConfig `json:"recurring"`
	Type        string              `json:"type"`
	WorkspaceId int                 `json:"workspace_id"`
}

type TaskInstance struct {
	Id           int            `json:"id"`
	TaskEntityId int            `json:"task_entity_id"`
	DtStart      string         `json:"dtstart"`
	CreatedAt    string         `json:"created_at"`
	UpdatedAt    sql.NullString `json:"updated_at"`
}

type CreateTaskInstanceInput struct {
	TaskEntityId int    `json:"task_entity_id"`
	DtStart      string `json:"dt_start"`
}

type TaskModel struct {
	DB *sql.DB
}

func NewTaskModel(DB *sql.DB) *TaskModel {
	return &TaskModel{
		DB: DB,
	}
}

func (m *TaskModel) Insert(ctx context.Context, in *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	var taskId int
	nanoid := common.GenerateNanoid()
	err := m.DB.QueryRow("INSERT INTO tasks(nanoid, title, description, duration, type, workspace_id, created_by) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		nanoid,
		in.Title,
		in.Description,
		in.Duration,
		in.Type,
		in.WorkspaceId,
		in.CreatedBy,
	).Scan(&taskId)

	if err != nil {
		return nil, err
	}

	if in.Recurring != nil {
		inputs := ProcessRecurring(in, taskId)
		for _, input := range inputs {
			m.DB.Exec("INSERT INTO task_instances(task_entity_id, dtstart) VALUES($1, $2)",
				input.TaskEntityId,
				input.DtStart)
		}
	}
	return nil, nil
}

func GetDateOfNthDay(year int, monthIndex int, dayOfWeekIndex int, n int) (time.Time, error) {
	if n < 1 || n > 5 {
		return time.Now(), errors.New("Invalid value of n. It should be between 1 and 5.")
	}
	month, _ := strconv.Atoi(fmt.Sprintf("0%d", monthIndex+1))
	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	firstDayOfWeek := int(firstDay.Weekday())
	offset := dayOfWeekIndex - firstDayOfWeek
	if offset < 0 {
		offset += 7
	}
	nthDay := 1 + (n-1)*7 + offset
	return time.Date(year, time.Month(month), nthDay, 0, 0, 0, 0, time.UTC), nil
}

func ProcessRecurring(input *pb.CreateTaskRequest, taskId int) []CreateTaskInstanceInput {
	recurring := input.Recurring
	if recurring == nil {
		return make([]CreateTaskInstanceInput, 0)
	}
	dtstartTS := input.Dtstart
	dtstart := time.Unix(dtstartTS, 0).UTC().String()
	layout := "2006-01-02 15:04:05 -0700 MST"
	startDT, _ := time.Parse(layout, dtstart)
	interval := int(recurring.Interval)
	count := int(recurring.Count)
	rType := recurring.Type
	byweekdayRule := recurring.ByweekdayRule
	currentDT := startDT
	var inputs []CreateTaskInstanceInput
	if rType != "WEEKLY" && byweekdayRule != nil {
		inputs = append(inputs, CreateTaskInstanceInput{TaskEntityId: taskId, DtStart: dtstart})
		if rType == "MONTHLY" {
			layout := "2006-01-02 15:04:05 -0700 MST"
			formattedDTStart, _ := time.Parse(layout, dtstart)
			month := int(formattedDTStart.Month()) + interval
			for i := 0; i < count; i++ {
				foundedDate, _ := GetDateOfNthDay(
					formattedDTStart.Year(),
					month-1,
					int(byweekdayRule.Day),
					int(byweekdayRule.Every),
				)
				foundedDate = time.Date(foundedDate.Year(),
					foundedDate.Month(),
					foundedDate.Day(),
					formattedDTStart.Hour(),
					formattedDTStart.Minute(), 0, 0, time.UTC)

				inputs = append(inputs, CreateTaskInstanceInput{TaskEntityId: taskId, DtStart: foundedDate.UTC().String()})
				month = month + interval
			}
		}
	} else {
		for i := 0; i < count; i++ {
			inputs = append(inputs, CreateTaskInstanceInput{
				DtStart:      currentDT.UTC().String(),
				TaskEntityId: taskId,
			})
			// currentDT = currentDT.add(interval, intervalType)
		}
	}
	return inputs
}
