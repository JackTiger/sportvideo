package modledata

import (
	"github.com/ginuerzh/sportvideo/common/errors"
)

var (
	ErrOK           = errors.NewError(errors.NoError, "", "ok")
	ErrInvalid      = errors.NewError(errors.InvalidError, "", "Invalid")
	ErrInternal     = errors.NewError(errors.InternalError, "", "Internal error")
	ErrNotFound     = errors.NewError(errors.NotFoundError, "", "Not found")
	ErrExist        = errors.NewError(errors.ExistError, "", "Exists")
	ErrUnauthorized = errors.NewError(errors.UnauthorizedError, "", "Unauthorized")
	ErrForbidden    = errors.NewError(errors.ForbiddenError, "", "Forbidden")
	ErrDB           = errors.NewError(errors.DbError, "", "database error")
)

const (
	dbName         = "ytbdownload"
    DownloadCollection  = "download"//Id FileName ImageURL FileSize VideoID DownloadUrl DownloadItag DownloadResolution DownloadExt Createtime
    KeywordCollection  = "keyword"//Id Word PageToken
	MongoAddr      = "172.24.222.44:27017"
)
