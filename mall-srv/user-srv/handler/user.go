package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"mall-srv/user-srv/global"
	"mall-srv/user-srv/model"
	pb "mall-srv/user-srv/proto"
	"strings"
	"time"

	"github.com/anaskhan96/go-password-encoder"
	"github.com/golang/protobuf/ptypes/empty"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

//type UserServer interface {
//	GetUserList(context.Context, *PageInfo) (*UserListResponse, error)
//	GetUserByMobile(context.Context, *MobileRequest) (*UserInfoResponse, error)
//	GetUserById(context.Context, *IdRequest) (*UserInfoResponse, error)
//	CreatUser(context.Context, *CreatRequest) (*UserInfoResponse, error)
//	UpdateUser(context.Context, *UpdateRequest) (*emptypb.Empty, error)
//	CheckPassword(context.Context, *PasswordCheckInfo) (*CheckResponse, error)
//	mustEmbedUnimplementedUserServer()
//}

type UserServer struct {
	pb.UnimplementedUserServer
}

// Paginate 分页逻辑
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func ModelToResponse(m model.User) pb.UserInfoResponse {
	//
	userInfoRsp := pb.UserInfoResponse{
		Id:       m.ID,
		Password: m.Password,
		Mobile:   m.Mobile,
		NickName: m.NickName,
		Gender:   m.Gender,
		Role:     int32(m.Role),
	}
	if m.Birthday != nil {
		userInfoRsp.Birthday = uint64(m.Birthday.Unix())
	}
	return userInfoRsp
}

// GetUserList 获取用户列表
func (s *UserServer) GetUserList(ctx context.Context, req *pb.PageInfo) (*pb.UserListResponse, error) {
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	rsp := &pb.UserListResponse{
		Total: 0,
		Data:  nil,
	}
	rsp.Total = int32(result.RowsAffected)

	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)

	for _, user := range users {
		userInfoRsp := ModelToResponse(user)
		rsp.Data = append(rsp.Data, &userInfoRsp)
	}

	return rsp, nil
}

func (s *UserServer) GetUserByMobile(ctx context.Context, req *pb.MobileRequest) (*pb.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil
}

func (s *UserServer) GetUserById(ctx context.Context, req *pb.IdRequest) (*pb.UserInfoResponse, error) {
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil
}

func (s *UserServer) CreatUser(ctx context.Context, req *pb.CreatRequest) (*pb.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}

	user.Mobile = req.Mobile
	user.NickName = req.NickName

	// 密码加密 by go-password-encoder
	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	salt, encodedPwd := password.Encode(req.Password, options)
	user.Password = fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)

	// 持久化到DB
	result = global.DB.Create(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "持久化到 DB 失败：%v", result.Error.Error())
	}

	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil
}

func (s *UserServer) UpdateUser(ctx context.Context, req *pb.UpdateRequest) (*empty.Empty, error) {
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	reqBirthday := time.Unix(int64(req.Birthday), 0)
	user.Birthday = &reqBirthday
	user.Gender = req.Gender
	user.NickName = req.NickName

	result = global.DB.Save(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &empty.Empty{}, nil
}

func (s *UserServer) CheckPassword(ctx context.Context, req *pb.PasswordCheckInfo) (*pb.CheckResponse, error) {
	// 校验密码
	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	passwordInfo := strings.Split(req.EncryptedPassword, "$")
	check := password.Verify(req.Password, passwordInfo[2], passwordInfo[3], options)
	return &pb.CheckResponse{Success: check}, nil
}
