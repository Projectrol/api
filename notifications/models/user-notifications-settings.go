package models

import (
	"context"
	"database/sql"
	"log"

	pb "github.com/lehoangvuvt/projectrol/common/protos"
)

type UpsertUserNotiSettingsInput struct {
	IsViaInbox          bool   `json:"is_via_inbox"`
	IsViaEmail          bool   `json:"is_via_email"`
	TaskNotiSettings    string `json:"task_noti_settings"`
	ProjectNotiSettings string `json:"project_noti_settings"`
	EventNotiSettings   string `json:"event_noti_settings"`
	EventNoticeBefore   int    `json:"event_notice_before"`
}

type UserNotificationsSettingsModel struct {
	DB *sql.DB
}

func NewUserNotificationsSettingsModel(DB *sql.DB) *UserNotificationsSettingsModel {
	return &UserNotificationsSettingsModel{
		DB: DB,
	}
}

func (m *UserNotificationsSettingsModel) Insert(ctx context.Context, in *pb.UpsertUserNotificationsSettingsRequest) (*pb.UpsertUserNotificationsSettingsResponse, error) {
	settingsInput := in.Settings
	_, err := m.DB.Exec(`INSERT INTO user_notifications_settings 
						(user_id, is_via_inbox, is_via_email, task_noti_settings, project_noti_settings, event_noti_settings, event_notice_before) 
						VALUES($1, $2, $3, $4, $5, $6, $7)`,
		settingsInput.UserId,
		settingsInput.IsViaInbox,
		settingsInput.IsViaEmail,
		settingsInput.TaskNotiSettings,
		settingsInput.ProjectNotiSettings,
		settingsInput.EventNotiSettings,
		settingsInput.EventNoticeBefore,
	)
	if err != nil {
		return nil, err
	}
	return &pb.UpsertUserNotificationsSettingsResponse{
		Settings: settingsInput,
	}, nil
}

func (m *UserNotificationsSettingsModel) Update(ctx context.Context, in *pb.UpsertUserNotificationsSettingsRequest) (*pb.UpsertUserNotificationsSettingsResponse, error) {
	settingsInput := in.Settings
	_, err := m.DB.Exec(`UPDATE user_notifications_settings 
						SET is_via_inbox = $1, is_via_email = $2, task_noti_settings = $3, project_noti_settings = $4, 
						event_noti_settings = $5, event_notice_before = $6 WHERE user_id = $7`,
		settingsInput.IsViaInbox,
		settingsInput.IsViaEmail,
		settingsInput.TaskNotiSettings,
		settingsInput.ProjectNotiSettings,
		settingsInput.EventNotiSettings,
		settingsInput.EventNoticeBefore,
		settingsInput.UserId,
	)
	if err != nil {
		return nil, err
	}
	return &pb.UpsertUserNotificationsSettingsResponse{
		Settings: settingsInput,
	}, nil
}

func (m *UserNotificationsSettingsModel) GetSettingsByUserId(ctx context.Context, in *pb.GetUserNotificationsSettingsRequest) (*pb.GetUserNotificationsSettingsResponse, error) {
	settings := &pb.UserNotificationsSettings{}
	log.Print(in.UserId)
	err := m.DB.QueryRow(`SELECT user_id, is_via_inbox, is_via_email, task_noti_settings, project_noti_settings, event_noti_settings, event_notice_before 
						FROM user_notifications_settings WHERE user_id=$1`, in.UserId,
	).Scan(&settings.UserId, &settings.IsViaInbox, &settings.IsViaEmail, &settings.TaskNotiSettings,
		&settings.ProjectNotiSettings, &settings.EventNotiSettings, &settings.EventNoticeBefore)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserNotificationsSettingsResponse{Settings: settings}, nil
}
