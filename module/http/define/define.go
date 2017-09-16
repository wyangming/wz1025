package define

//filter define
const (
	//member filter url
	FILTER_MEMBER     = "/member"
	FILTER_MEMBER_ALL = "/member/**"
)

//url define
const (
	//main controller url
	URL_INDEX = "/"
	URL_LOGIN = "/login"
	URL_REG   = "/reg"

	//member url
	URL_MEMBER       = "/member"
	URL_MEMBER_INFO  = "/member/info"
	URL_MEMBER_VIDEO = "/member/video"
)

//session const
const (
	SESSION_MEMBER_INFO = "memberInfo"
)

//main controller
const (
	//使用的常量
	CON_MAIN_LOGIN_STATUS = "login_status"
	CON_MAIN_REG_STATUS   = "reg_status"

	//页面信息
	CON_MAIN_LOGIN_PAGE = "login.html"
	CON_MAIN_REG_PAGE   = "reg.html"
	CON_MAIN_INDEX_PAGE = "index.html"
)

//member controller
const (
	CON_MEMBER_MAIN_PAGE  = "member/main.html"
	CON_MEMBER_INFO_PAGE  = "member/info.html"
	CON_MEMBER_VIDEO_PAGE = "member/video.html"
)
