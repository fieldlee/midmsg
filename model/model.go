package model

import "github.com/micro/go-micro/errors"

type Min_Net_MsgBody struct {
	m_lAsktype  		uint64 		//请求的服务类型
	m_lServerSequence   uint64		//服务端响应序列号
	m_lAskSequence		uint64		//客户端请求序列号
	m_cMsgAckType		int32		//0 无需回复, 1 回复到发送方, 2 回复到离线服务器
	m_cMsgType			int32		//消息类型
	m_sSendCount		int32   	//同一请求，请求服务端的次数
	m_lExpireTime		uint32		//过期时间，0表示永不过期
	m_iSendTimeApp		uint64		//请求发送的时间，单位秒
	m_lResult			int32		//0: SUCCESS  !0:FAILURE
	m_lBack				uint64		//备份数据，默认为0
	m_iDiscard			int32		//消息是否可以丢弃 0 表示可以丢弃 1表示不可以
}

type  Net_Pack struct {
	m_Msg	    []byte					//每个服务类型定义的protobuf结构，打包成流缓存在该字段
	m_MsgBody   Min_Net_MsgBody
}

type GJ_Net_Pack struct{
	m_Net_Pack	map[uint32]Net_Pack		//可缓存多个客户端请求
}

var (
	ErrHeaderLength = errors.New("errheaderlength","the header length error",40001)
	ErrClientType = errors.New("errclienttype","the client type error",40002)
	ErrCompressionType = errors.New("errcompressiontype","the compression type error",40003)
	ErrEncrptyType = errors.New("errencryptiontype","the encryption type error",40004)
	ErrCompressedLength = errors.New("errcompressedlength","the compressed length error",40005)
	ErrUNCompressedLength = errors.New("erruncompressedlength","the uncompressed length error",40006)
	ErrCompresseduncompressedLength = errors.New("errcompresseduncompressedlength","the compressed length more the uncompressed length error",40007)
)

var (
    GJ_PUBLIC_START        		=   0							//公共部分的请求
    GJ_PUBLIC_NET_OPERATION		=   GJ_PUBLIC_START+10000		//操作部分
    ETN_ASK_LOAIN_SERVER     	=	GJ_PUBLIC_NET_OPERATION+1 	// 登录
    ETN_SERVER_NET_CONNET    	=	GJ_PUBLIC_NET_OPERATION+2 	//有客户端连接成功
    ETN_SERVER_NET_CLOSE     	=	GJ_PUBLIC_NET_OPERATION+3	//服务端网络层关闭
    ETN_ASK_USER_LEAVE       	=	GJ_PUBLIC_NET_OPERATION+4	//用户登录退出
    ETN_SERVER_PUSH_NOTICE_MSG  =   GJ_PUBLIC_NET_OPERATION+5	//服务器推送通知
    ETN_HEARTBEAT_PACK          =   GJ_PUBLIC_NET_OPERATION+6	//心跳包
    ETN_SERVER_SUBSRCTIBE_MSG   =   GJ_PUBLIC_NET_OPERATION+7	//广播消息
)

type MSG_TYPE int32
var (
	MSG_TYPE_ACK				MSG_TYPE	= 0				//普通请求类型
	MSG_TYPE_LOGIN_REQ			MSG_TYPE	= 1				//注册
	MSG_TYPE_LOGIN_ACK			MSG_TYPE	= 2				//注册响应
	MSG_TYPE_KEEPALIVE_REQ		MSG_TYPE	= 3				//心跳检测
	MSG_TYPE_KEEPALIVE_ACK		MSG_TYPE	= 4				//心跳检测响应
	MSG_TYPE_PUSHMSG_REQ		MSG_TYPE	= 5				//下发消息
	MSG_TYPE_PUSHMSG_ACK		MSG_TYPE	= 6				//下发消息响应
	MSG_TYPE_UPLOADMSG_REQ		MSG_TYPE	= 7				//上传消息
	MSG_TYPE_UPLOADMSG_ACK		MSG_TYPE	= 8				//上传消息响应
	MSG_TYPE_BROADCAST			MSG_TYPE	= 9				//广播消息
	MSG_TYPE_SUBSCRIBE_REQ		MSG_TYPE	= 10			//订阅消息
	MSG_TYPE_SUBSCRIBE_ACK		MSG_TYPE	= 11			//订阅消息响应
	MSG_TYPE_ERROR				MSG_TYPE	= 12			//错误信息应答
	MSG_TYPE_NOTICE             MSG_TYPE    = 13			//通知消息
	MSG_TYPEMAX					MSG_TYPE    = 14			//用于判断合法性预留，以后该枚举需扩展，则在该枚举值上面进行扩展
	MSG_TYPE_CONNECT_COUNT_MSG	MSG_TYPE	= 0xFFFF		//查询连接数
)

type CLIENT_TYPE int16
var (
	Window_pc				CLIENT_TYPE		= 0				//请求来自PC
	IOS_mobile				CLIENT_TYPE		= 1				//请求来自苹果手机
	Android_mobilewindow_pc	CLIENT_TYPE		= 3				//请求来自安卓手机
	Web_side				CLIENT_TYPE		= 4				//请求来自WEB端
	ClientTypeMax			CLIENT_TYPE		= 5				//用于判断合法性预留，以后该枚举需扩展，则在该枚举值上面进行扩展
)
type COMPRESS_TYPE int16
var (
	Compression_no			COMPRESS_TYPE		= 0				//表示数据未压缩
	Compression_zip			COMPRESS_TYPE		= 1				//表示数据被压缩zlib压缩格式
	CompressionWayMax		COMPRESS_TYPE		= 2				//用于判断合法性预留，以后该枚举需扩展，则在该枚举值上面进行扩展
)

type ENCRPTION_TYPE int16
var (
	Encryption_No			ENCRPTION_TYPE		= 0				//请求数据流未被加密
	Encryption_Des			ENCRPTION_TYPE		= 1				//请求数据是通过Des算法加密
	Encryption_AES			ENCRPTION_TYPE		= 2				//请求数据是通过AES算法加密
	Encryption_RSA			ENCRPTION_TYPE		= 3				//请求数据是通过RSA算法加密
	Encryption_Max			ENCRPTION_TYPE		= 4				//用于判断合法性预留，以后该枚举需扩展，则在该枚举值上面进行扩展
)

type DATAFORMAT_TYPE int16

var (
	DataFormat_Probufo		DATAFORMAT_TYPE		= 0			//数据是通过protobuf进行格式化
	DataFormatMax			DATAFORMAT_TYPE		= 1			//用于判断合法性预留，以后该枚举需扩展，则在该枚举值上面进行扩展
)

type CALL_CLIENT_TYPE uint8
var (
	CALL_CLIENT_SYNC 	CALL_CLIENT_TYPE = 0
	CALL_CLIENT_ASYNC   CALL_CLIENT_TYPE = 1
)