package errors

import (
	"fmt"
	"net/http"

	"github.com/ares-cloud/ares-ddd-admin/pkg/hserver/herrors"
)

const (
	ReasonFolderNotFound     = "FOLDER_NOT_FOUND"     // 文件夹不存在
	ReasonFolderExists       = "FOLDER_EXISTS"        // 文件夹已存在
	ReasonFolderNotEmpty     = "FOLDER_NOT_EMPTY"     // 文件夹不为空
	ReasonFileNotFound       = "FILE_NOT_FOUND"       // 文件不存在
	ReasonFileExists         = "FILE_EXISTS"          // 文件已存在
	ReasonInvalidFileName    = "INVALID_FILE_NAME"    // 文件名无效
	ReasonInvalidFolderName  = "INVALID_FOLDER_NAME"  // 文件夹名无效
	ReasonInvalidOperation   = "INVALID_OPERATION"    // 无效操作
	ReasonInvalidMoveOp      = "INVALID_MOVE_OP"      // 无效的移动操作
	ReasonStorageError       = "STORAGE_ERROR"        // 存储错误
	ReasonShareNotFound      = "SHARE_NOT_FOUND"      // 分享不存在
	ReasonShareExpired       = "SHARE_EXPIRED"        // 分享已过期
	ReasonSharePasswordError = "SHARE_PASSWORD_ERROR" // 分享密码错误
)

// FolderNotFound 文件夹不存在
func FolderNotFound(id string) herrors.Herr {
	return herrors.New(http.StatusNotFound, ReasonFolderNotFound,
		fmt.Sprintf("folder not found: %s", id))
}

// FileNotFound 文件不存在
func FileNotFound(id string) herrors.Herr {
	return herrors.New(http.StatusNotFound, ReasonFileNotFound,
		fmt.Sprintf("file not found: %s", id))
}

// InvalidFileName 文件名无效
func InvalidFileName(name string) herrors.Herr {
	return herrors.New(http.StatusBadRequest, ReasonInvalidFileName,
		fmt.Sprintf("invalid file name: %s", name))
}

// InvalidFolderName 文件夹名无效
func InvalidFolderName(name string, reason string) herrors.Herr {
	return herrors.New(http.StatusBadRequest, ReasonInvalidFolderName,
		fmt.Sprintf("invalid folder name %s: %s", name, reason))
}

// FolderNotEmpty 文件夹不为空
func FolderNotEmpty(id string, count int) herrors.Herr {
	return herrors.New(http.StatusBadRequest, ReasonFolderNotEmpty,
		fmt.Sprintf("folder %s is not empty, contains %d items", id, count))
}

// InvalidOperation 无效操作
func InvalidOperation(reason string) herrors.Herr {
	return herrors.New(http.StatusBadRequest, ReasonInvalidOperation,
		fmt.Sprintf("invalid operation: %s", reason))
}

// InvalidMoveOperation 无效的移动操作
func InvalidMoveOperation(reason string) herrors.Herr {
	return herrors.New(http.StatusBadRequest, ReasonInvalidMoveOp,
		fmt.Sprintf("invalid move operation: %s", reason))
}

// StorageError 存储错误
func StorageError(err error) herrors.Herr {
	return herrors.New(http.StatusInternalServerError, ReasonStorageError,
		fmt.Sprintf("storage error: %v", err))
}

// FolderExists 文件夹已存在
func FolderExists(name string) herrors.Herr {
	return herrors.New(http.StatusBadRequest, ReasonFolderExists,
		fmt.Sprintf("folder already exists: %s", name))
}

// FileExists 文件已存在
func FileExists(name string) herrors.Herr {
	return herrors.New(http.StatusBadRequest, ReasonFileExists,
		fmt.Sprintf("file already exists: %s", name))
}

// ErrShareNotFound 分享不存在错误
var ErrShareNotFound = herrors.New(
	http.StatusNotFound,
	ReasonShareNotFound,
	"share not found",
)

// ErrShareExpired 分享已过期错误
var ErrShareExpired = herrors.New(
	http.StatusBadRequest,
	ReasonShareExpired,
	"share has expired",
)

// ErrSharePasswordIncorrect 分享密码错误
var ErrSharePasswordIncorrect = herrors.New(
	http.StatusBadRequest,
	ReasonSharePasswordError,
	"incorrect share password",
)

// ShareNotFound 分享不存在
func ShareNotFound(shareCode string) herrors.Herr {
	return herrors.New(
		http.StatusNotFound,
		ReasonShareNotFound,
		fmt.Sprintf("share not found: %s", shareCode),
	)
}

// ShareExpired 分享已过期
func ShareExpired(shareCode string) herrors.Herr {
	return herrors.New(
		http.StatusBadRequest,
		ReasonShareExpired,
		fmt.Sprintf("share has expired: %s", shareCode),
	)
}

// SharePasswordIncorrect 分享密码错误
func SharePasswordIncorrect(shareCode string) herrors.Herr {
	return herrors.New(
		http.StatusBadRequest,
		ReasonSharePasswordError,
		fmt.Sprintf("incorrect password for share: %s", shareCode),
	)
}
