// error
package errors

import (
	"fmt"
)

const (
	NoError   = iota
	DeserializationError
	RequiredError
	TypeError
	ParamError
	InvalidError
	ExistError
	InternalError
	UnauthorizedError
	ForbiddenError

	AuthError = 1000 + iota
	UserExistError
	MailExistError
	AccessError
	DbError
	_
	JsonError
	NotFoundError
	PasswordError
	InvalidFileError
	HttpError
	FileNotFoundError
	_
	NotExistsError
	InvalidAddrError
	InvalidMsgError
	DeviceTokenError
	ReviewNotFoundError
	InviteCodeError
	FileTooLargeError
	FileUploadError
	UnimplementedError
    NsqSendError
	GrpcError
)

var errMap map[int]string = map[int]string{
	NoError:                "success",
	DeserializationError:   "转换出错",
	RequiredError:          "请求出错",
	TypeError:              "类型出错",
	ParamError:             "参数出错",
	InvalidError:           "不可用",
	ExistError:             "不存在",
	InternalError:          "网络出错",
	UnauthorizedError:      "未授权",
	ForbiddenError:         "被禁止的操作",
	AuthError:              "用户名或密码错误",
	UserExistError:         "用户已注册",
	MailExistError:         "邮箱已注册",
	AccessError:            "无效的访问请求",
	DbError:                "服务器出错啦！",
	JsonError:              "json data error",
	NotFoundError:          "not found",
	PasswordError:          "password invalid",
	InvalidFileError:       "file invalid",
	HttpError:              "http error",
	FileNotFoundError:      "file not found",
	NotExistsError:         "用户不存在",
	InvalidAddrError:       "address invalid",
	InvalidMsgError:        "message invalid",
	DeviceTokenError:       "device token invalid",
	ReviewNotFoundError:    "review not found",
	InviteCodeError:        "invite code invalid",
	FileTooLargeError:      "file too large",
	FileUploadError:        "file upload error",
	UnimplementedError:     "unimplemented",
    NsqSendError:           "nsq send error",
	GrpcError:				"grpc error",
}

type Error struct {
	Id   int    `json:"error_id"`
	Desc string `json:"error_desc"`
}

func NewError(id int, desc ...string) *Error {
	s := errMap[id]
	if len(desc) > 0 {
		s = desc[0]
	}
	return &Error{Id: id, Desc: s}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d %s", e.Id, e.Desc)
}
