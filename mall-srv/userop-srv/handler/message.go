package handler

import (
	"context"
	"mall-srv/userop-srv/global"
	"mall-srv/userop-srv/model"
	pb "mall-srv/userop-srv/proto"
)

func (*UserOpServer) MessageList(ctx context.Context, req *pb.MessageRequest) (*pb.MessageListResponse, error) {
	var rsp pb.MessageListResponse
	var messages []model.LeavingMessages
	var messageList []*pb.MessageResponse

	result := global.DB.Where(&model.LeavingMessages{User: req.UserId}).Find(&messages)
	rsp.Total = int32(result.RowsAffected)

	for _, message := range messages {
		messageList = append(messageList, &pb.MessageResponse{
			Id:          message.ID,
			UserId:      message.User,
			MessageType: message.MessageType,
			Subject:     message.Subject,
			Message:     message.Message,
			File:        message.File,
		})
	}

	rsp.Data = messageList
	return &rsp, nil
}

func (*UserOpServer) CreateMessage(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {
	var message model.LeavingMessages

	message.User = req.UserId
	message.MessageType = req.MessageType
	message.Subject = req.Subject
	message.Message = req.Message
	message.File = req.File

	global.DB.Save(&message)

	return &pb.MessageResponse{Id: message.ID}, nil
}
