package api

import (
	"context"
	"fmt"
	"mall-api/user-web/middleware"
	"mall-api/user-web/model"

	"github.com/go-redis/redis/v8"

	"github.com/dgrijalva/jwt-go"

	"github.com/go-playground/validator/v10"

	"mall-api/user-web/global"
	"mall-api/user-web/global/response"
	pb "mall-api/user-web/proto"
	"net/http"
	"strconv"

	"time"

	"google.golang.org/grpc/codes"

	"mall-api/user-web/forms"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
)

// HandleGrpcErrorToHttp 将grpc的错误码转换为http的状态码
func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其他错误" + e.Message(),
				})

			}
		}
	}
}

// HandleValidationError 将表单验证信息以中文返回
func HandleValidationError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": errs.Translate(global.Trans),
	})
	return
}

func GetUserList(c *gin.Context) {
	// 权限验证
	claims, ok := c.Get("claims")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "获取 claims 失败",
		})
		return
	}
	curClaims := claims.(*model.CustomClaims)
	zap.S().Infof("用户权限：%v", curClaims.AuthorityID)
	if curClaims.AuthorityID != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "无操作权限",
		})
	}

	// 接收数据 Pn 和 PSize
	page := c.DefaultQuery("page", "0")
	Pn, _ := strconv.Atoi(page)
	size := c.DefaultQuery("size", "0")
	PSize, _ := strconv.Atoi(size)

	//调用接口
	rsp, err := global.UserSrvClient.GetUserList(c, &pb.PageInfo{
		Pn:    uint32(Pn),
		PSize: uint32(PSize),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询失败")
		HandleGrpcErrorToHttp(err, c)
		return
	}

	result := make([]response.UserResponse, 0)
	for _, v := range rsp.Data {
		user := response.UserResponse{
			Id:       v.Id,
			NickName: v.NickName,
			Mobile:   v.Mobile,
			Birthday: time.Unix(int64(v.Birthday), 0).Format("2006-01-01"),
			Gender:   v.Gender,
		}
		result = append(result, user)
	}
	c.JSON(http.StatusOK, result)
}

func LoginByPassword(c *gin.Context) {
	// 表单验证
	loginByPasswordForm := forms.LoginByPasswordForm{}
	if err := c.ShouldBind(&loginByPasswordForm); err != nil {
		HandleValidationError(c, err)
	}

	// 检查验证码
	if !store.Verify(loginByPasswordForm.CaptchaId, loginByPasswordForm.Captcha, true) {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "验证码错误",
		})
		return
	}

	// 查询用户是否存在
	if rsp, err := global.UserSrvClient.GetUserByMobile(c, &pb.MobileRequest{Mobile: loginByPasswordForm.Mobile}); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": "用户不存在",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "登陆失败，GetUserByMobile 服务错误",
				})

			}
			return
		}
	} else {
		// 检查密码
		if check, err := global.UserSrvClient.CheckPassword(context.Background(), &pb.PasswordCheckInfo{
			Password:          loginByPasswordForm.Password,
			EncryptedPassword: rsp.Password,
		}); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "登陆失败，检查密码出错",
			})
		} else if check.Success {
			// 生成 jwt token
			j := middleware.NewJWT()
			claims := model.CustomClaims{
				ID:          uint(rsp.Id),
				NickName:    rsp.NickName,
				AuthorityID: uint(rsp.Role),
				StandardClaims: jwt.StandardClaims{
					NotBefore: time.Now().Unix(),               // 生效时间
					ExpiresAt: time.Now().Unix() + 60*60*24*30, // 过期时间 - 30天
					Issuer:    "hans",
				},
			}
			token, err := j.CreateToken(claims)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "token生成失败",
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"id":         rsp.Id,
				"nick_name":  rsp.NickName,
				"token":      token,
				"expired_at": (time.Now().Unix() + 60*60*24*7) * 1000,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "密码错误",
			})
		}
	}
}

func Register(c *gin.Context) {
	// 表单验证
	RegisterForm := forms.RegisterForm{}
	if err := c.ShouldBind(&RegisterForm); err != nil {
		HandleValidationError(c, err)
	}

	// 验证码校验
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.SrvConfig.RedisConfig.Host, global.SrvConfig.RedisConfig.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	r, err := rdb.Get(context.Background(), RegisterForm.Mobile).Result()
	if err == redis.Nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "redis-key不存在",
		})
		return
	} else if r != RegisterForm.SmsCode {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "验证码错误",
		})
		return
	}

	rsp, err := global.UserSrvClient.CreatUser(context.Background(), &pb.CreatRequest{
		NickName: RegisterForm.Mobile,
		Mobile:   RegisterForm.Mobile,
		Password: RegisterForm.Password,
	})
	if err != nil {
		zap.S().Error("[CreatUser] 调用失败：", err.Error())
		HandleGrpcErrorToHttp(err, c)
		return
	}

	// 注册后直接登陆（生成jwt免登录）
	j := middleware.NewJWT()
	claims := model.CustomClaims{
		ID:          uint(rsp.Id),
		NickName:    rsp.NickName,
		AuthorityID: uint(rsp.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),              // 生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*7, // 过期时间 - 7天
			Issuer:    "hans",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "token生成失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":         rsp.Id,
		"nick_name":  rsp.NickName,
		"token":      token,
		"expired_at": (time.Now().Unix() + 60*60*24*7) * 1000,
	})
}
