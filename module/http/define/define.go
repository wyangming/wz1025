package define

//filter define
const (
	//member filter url
	FILTER_MEMBER     = "/member"
	FILTER_MEMBER_ALL = "/member/**"

	//admin filter url
	FILTER_ADMIN     = "/admin"
	FILTER_ADMIN_ALL = "/admin/**"
)

//url define
const (
	//main controller url
	URL_INDEX      = "/"
	URL_LOGIN      = "/login"
	URL_REG        = "/reg"
	URL_ADMINLOGIN = "/adminlogin"

	//member url
	URL_MEMBER              = "/member"
	URL_MEMBER_INFO         = "/member/info"
	URL_MEMBER_VIDEO        = "/member/video"
	URL_MEMBER_EXPLAIN_INFO = "/member/ajax_explain_info"

	//admin url
	URL_ADMIN                      = "/admin"
	URL_ADMIN_EXPLAIN              = "/admin/explain"
	URL_ADMIN_EXPLAIN_LIST         = "/admin/explain_list"
	URL_ADMIN_EXPLAIN_ACTIVEUPDATE = "/admin/explain_activeUpdate"
	URL_ADMIN_EXPLAIN_SPIDERUPDATE = "/admin/explain_spiderUpdate"
	URL_ADMIN_MEMBER               = "/admin/member"
	URL_ADMIN_MEMBER_LIST          = "/admin/member_list"
	URL_ADMIN_MEMBER_ACTIVEUPDATE  = "/admin/member_activeUpdate"
	URL_ADMIN_MEMBER_EXPIREUPDATE  = "/admin/member_expireUpdate"
)

//session const
const (
	SESSION_MEMBER_INFO = "memberInfo"
	SESSION_ADMIN_INFO  = "adminInfo"
)

//main controller
const (
	//使用的常量
	CON_MAIN_LOGIN_STATUS = "login_status"
	CON_MAIN_REG_STATUS   = "reg_status"

	//页面信息
	CON_MAIN_LOGIN_PAGE      = "login.html"
	CON_MAIN_REG_PAGE        = "reg.html"
	CON_MAIN_INDEX_PAGE      = "index.html"
	CON_MAIN_ADMINLOGIN_PAGE = "adminlogin.html"

	//方法名
	CON_MAIN_GET_METHOD        = "*:Get"
	CON_MAIN_LOGIN_METHOD      = "get,post:Login"
	CON_MAIN_REG_METHOD        = "get,post:Reg"
	CON_MAIN_ADMINLOGIN_METHOD = "get,post:Adminlogin"
)

//member controller
const (
	//页面信息
	CON_MEMBER_MAIN_PAGE  = "member/main.html"
	CON_MEMBER_INFO_PAGE  = "member/info.html"
	CON_MEMBER_VIDEO_PAGE = "member/video.html"

	//方法名
	CON_MEMBER_GET_METHOD             = "*:Get"
	CON_MEMBER_INFO_METHOD            = "*:Info"
	CON_MEMBER_VIDEO_METHOD           = "*:Video"
	CON_MEMBER_AJAXEXPLAININFO_METHOD = "*:AjaxExplainInfo"
)

//admin controller
const (
	//页面信息
	CON_ADMIN_MAIN_PAGE    = "admin/main.html"
	CON_ADMIN_EXPLAIN_PAGE = "admin/explain.html"
	CON_ADMIN_MEMBER_PAGE  = "admin/member.html"

	//方法名
	CON_ADMIN_GET_METHOD                 = "*:Get"
	CON_ADMIN_EXPLAIN_METHOD             = "*:Explain"
	CON_ADMIN_EXPLAINLIST_METHOD         = "*:Explain_List"
	CON_ADMIN_EXPLAINACTIVEUPDATE_METHOD = "*:Explain_ActiveUpdate"
	CON_ADMIN_EXPLAINSPIDERUPDATE_METHOD = "*:Explain_spiderUpdate"
	CON_ADMIN_MEMBER_METHOD              = "*:Member"
	CON_ADMIN_MEMBERlIST_METHOD          = "*:Member_List"
	CON_ADMIN_MEMBERACTIVEUPDATE_METHOD  = "*:Member_ActiveUpdate"
	CON_ADMIN_MEMBEREXPIREUPDATE_METHOD  = "*:Member_ExpireUpdate"
)
