package constant

const (
	KeyAccessToken = "access_token"
)

const (
	ReasonTokenEmpty      = "tokenEmpty"
	ReasonTokenVerifyFail = "tokenVerifyFail"
	IsNotAdminAccount     = "isNotAdminAccount"
	SQLInjectionDetected  = "sqlInjectionDetected"
	ReasonSuccess         = "Success"
	PleaseDoNotResubmit   = "PleaseDoNotResubmit"
	ReasonNoAccess        = "noAccess"
)
const (
	RespCode      = "code"
	RespMsg       = "msg"
	RespData      = "data"
	RespReason    = "reason"
	ErrMsg        = "errMsg"
	RespTimestamp = "timestamp"
	ReasonHttpOk  = "httpOk"
)

const (
	StatusSuccess      = 200
	StatusNotFindData  = 404
	StatusInvalidParam = 400
	StatusNoPermission = 405
	StatusServerError  = 500
	StatusInvalidToken = 401
)
const (
	RoleUser       = "user"
	RoleAdmin      = "admin"
	RoleSuperAdmin = "superAdmin"
	RoleAgent      = "Agent"
)

func GetConstRole(isAdmin int32) string {
	switch isAdmin {
	case 0:
		return RoleUser
	case 1:
		return RoleAdmin
	default:
		return RoleUser
	}
}

const (
	MetadataUserId     = "x-md-global-userId"
	MetadataRole       = "x-md-global-role"
	MetadataTenantId   = "x-md-global-tenant-id"
	MetadataDeviceId   = "x-md-global-device-Id"
	MetadataDeviceName = "x-md-global-device-name"
	MetadataIpAddress  = "x-md-global-ip-address"
)

type GroupMemberLevel byte

const (
	ReqPending = 0 //待处理
	ReqApprove = 1 //同意
	ReqRefuse  = 2 //拒绝

	GroupOrdinaryUser GroupMemberLevel = 20  //普通用户
	GroupAdmin        GroupMemberLevel = 60  //管理员
	GroupLeader       GroupMemberLevel = 100 //群主

	GroupOk              = 0 //正常状态
	GroupBanChat         = 1 //禁止聊天
	GroupStatusDismissed = 2 //群已经解散状态
	GroupStatusMuted     = 3 //禁言状态

	GroupFilterAll                   = 0
	GroupFilterOwner                 = 1
	GroupFilterAdmin                 = 2
	GroupFilterOrdinaryUsers         = 3
	GroupFilterAdminAndOrdinaryUsers = 4
	GroupFilterOwnerAndAdmin         = 5
)
