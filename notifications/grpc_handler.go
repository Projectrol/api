package main

import (
	"context"

	pb "github.com/lehoangvuvt/projectrol/common/protos"
	"github.com/lehoangvuvt/projectrol/notifications/models"
)

type server struct {
	pb.UnimplementedNotificationsServiceServer
	UserNotificationsSettingsModel *models.UserNotificationsSettingsModel
}

func (s *server) CreateUserNotificationsSettings(ctx context.Context, in *pb.UpsertUserNotificationsSettingsRequest) (*pb.UpsertUserNotificationsSettingsResponse, error) {
	return s.UserNotificationsSettingsModel.Insert(ctx, in)
}

func (s *server) UpdateUserNotificationsSettings(ctx context.Context, in *pb.UpsertUserNotificationsSettingsRequest) (*pb.UpsertUserNotificationsSettingsResponse, error) {
	return s.UserNotificationsSettingsModel.Update(ctx, in)
}

func (s *server) GetUserNotificationsSettings(ctx context.Context, in *pb.GetUserNotificationsSettingsRequest) (*pb.GetUserNotificationsSettingsResponse, error) {
	return s.UserNotificationsSettingsModel.GetSettingsByUserId(ctx, in)
}
