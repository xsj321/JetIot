package errorCode

const (
	ERR_NO_ERROR                    = 0
	ERR_UNSUPPORT_VER               = 1  // 版本不支持
	ERR_AUTH_FAILED                 = 2  // 用户名密码错误
	ERR_VA_CODE_MISMATCH            = 3  // VA CODE mimatched.
	ERR_VA_CODE_TIMEOUT             = 4  // 验证码超时
	ERR_USER_NON_EXIST              = 5  // mysql中用户不存在
	ERR_USER_EXISTS                 = 6  // mysql中用户存在
	ERR_MYSQL_FAILED                = 7  // mysql查询错误
	ERR_REDIS_FAILED                = 8  // redis查询错误
	ERR_LAST_VA_CODE_STILL_VALID    = 9  // 上一个VA_CODE 依然有效,不能产生新的
	ERR_ADD_MSG_FAILED              = 10 // 邮件或者短信加入队列失败
	ERR_APP_NO_NEW_VER              = 11 // No other new app version
	ERR_ONE_HOUR_SMS_OVERFLOW       = 12 // 一小时内发送短信已达上限
	ERR_OTHER_APP_ONLINE            = 13 // 其他app已经登录该帐号
	ERR_ACCOUNT_WARNING             = 14 // 帐号不安全,重新设置密码
	ERR_ONE_DAY_SMS_OVERFLOW        = 15 // 二十四小时内发送短信已达上限
	ERR_SHARE_ACCOUNT_CNT_OVER_FLOW = 16 // 授权帐号个数已达上线
	ERR_ACCOUNT_HAD_PERMISSION      = 17 // 帐号已经有权限
	ERR_NEW_APP_MANAGE              = 18 // 有新帐号来管理,失去管理权限
	ERR_OLD_APP_MANAGE              = 19 // 其老帐号正在管理,不能进行管理
	ERR_TOKEN_EXPIRED               = 20 // token超时，请重新获取token
	ERR_APP_NEW_VER                 = 21 // 当前app有更新版本
	ERR_SVR_INTERNAL                = 22 // 服务器内部错误
	ERR_DATA_ERR                    = 23 // 来源数据错误
	ERR_SIG_ERR                     = 24 // 签名错误
	ERR_USER_REDIS_NON_EXIST        = 25 // redis中用户不存在
	ERR_THIRD_NAME_BINDED           = 26 // 第三方账号已绑定
	ERR_OLD_PASSWD                  = 27 // 旧密码错误
	ERR_ACCOUNT_STILL_BOUND_DEVICE  = 28 // 账号下仍有绑定设备
	ERR_MEMBER_REQ_NON_EXIST        = 29 // 好友请求不存在
	ERR_ORDER_NO_PRODUCE_FAILED     = 30 // 订单号生成失败
	ERR_MEMBERSHIP_ALREADY_EXIST    = 31 // 好友关系已存在
	ERR_HAS_NO_FRIENDS              = 32 // 没有好友
	ERR_TELBOOK_MATCH_NOTHING       = 33 // 没有匹配的电话簿
	ERR_ORDER_PAID                  = 34 // 订单已支付
	ERR_ORDER_NOT_PAY               = 35 // 订单未支付
	ERR_ORDER_CLOSED                = 36 // 订单已关闭
	ERR_ORDER_PAYERROR              = 37 // 订单支付失败，非内部原因
	ERR_DEV_BINDED                  = 38 // 设备已被绑定
	ERR_ACCOUNT_LOGOUT              = 39 // 账号已被注销
	ERR_MAX                         = 40
	ERR_DEVICE_NOT_FIND             = 41
)
