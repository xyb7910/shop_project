package api

import (
	"mxshop-api/user-web/forms"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/middlewares"
	"mxshop-api/user-web/models"
	"mxshop-api/user-web/proto"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleGrpcErrToHttp(err error, c *gin.Context) {
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
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其他错误",
				})
			}
			return
		}
	}
}

func GetUserList(ctx *gin.Context) {
	userConn, err := grpc.Dial(":50001", grpc.WithInsecure())
	if err != nil {
		zap.S().Error("【GetUserList】链接用户服务失败：", err)
		HandleGrpcErrToHttp(err, ctx)
		return
	}
	userClient := proto.NewUserClient(userConn)

	var page string
	var size string
	page = ctx.DefaultQuery("page", "1")
	size = ctx.DefaultQuery("size", "10")

	pageInt, _ := strconv.Atoi(page)
	sizeInt, _ := strconv.Atoi(size)
	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{Page: uint32(pageInt), Size: uint32(sizeInt)})
	//if err != nil {
	//	zap.S().Error("【GetUserList】请求接口失败：", err.Error())
	//	HandleGrpcErrToHttp(err, ctx)
	//	return
	//}
	//result := make([]interface{}, 0)
	//for _, value := range rsp.Data {
	//	data := make(map[string]interface{})
	//	data["id"] = value.Id
	//	data["nickname"] = value.NickName
	//	data["gender"] = value.Gender
	//	data["mobile"] = value.Mobile
	//	result = append(result, data)
	//
	//}
	ctx.JSON(http.StatusOK, rsp.Data)
}

func PassWordLogin(c *gin.Context) {
	passwordLoginForm := forms.PassWordLoginForm{}
	if err := c.ShouldBindJSON(&passwordLoginForm); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": global.RemoveTopStruct(errs.Translate(global.Translator)),
		})
		return
	}
	userConn, err := grpc.Dial(":50001", grpc.WithInsecure())
	if err != nil {
		zap.S().Error("【GetUserList】链接用户服务失败：", err)
		HandleGrpcErrToHttp(err, c)
		return
	}
	userClient := proto.NewUserClient(userConn)
	rsp, err := userClient.GetUserMobile(c, &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest, map[string]string{
					"mobile": "用户不存在",
				})
			default:
				c.JSON(http.StatusBadRequest, map[string]string{
					"mobile": "登陆失败",
				})
			}
			return
		}
	}
	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:          uint(rsp.Id),
		NickName:    rsp.NickName,
		AuthorityId: uint(rsp.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Unix() + 60*60*24*30,
			Issuer:    "imooc",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成token失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":       "登陆成功",
		"user":      rsp,
		"token":     token,
		"expireAt0": (time.Now().Unix() + 60*60*24*30),
	})

}
