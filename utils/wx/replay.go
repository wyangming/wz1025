package wx

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

const (
	//WeiXin conf info not User
	Token = "c1ca15e2cac8482cacb3fc07216e58db"

	// Event type
	EventSubscribe   = "subscribe"
	EventUnsubscribe = "unsubscribe"
	EventQrscene_    = "qrscene_"
	EventScan        = "SCAN"
	EventView        = "VIEW"
	EventClick       = "CLICK"
	EventLocation    = "LOCATION"

	// Message type
	MsgTypeText       = "text"
	MsgTypeImage      = "image"
	MsgTypeVoice      = "voice"
	MsgTypeVideo      = "video"
	MsgTypeShortVideo = "shortvideo"
	MsgTypeLocation   = "location"
	MsgTypeLink       = "link"
	MsgTypeEvent      = "event"

	// Reply format
	replyText    = "<xml>%s<MsgType><![CDATA[text]]></MsgType><Content><![CDATA[%s]]></Content></xml>"
	replyImage   = "<xml>%s<MsgType><![CDATA[image]]></MsgType><Image><MediaId><![CDATA[%s]]></MediaId></Image></xml>"
	replyVoice   = "<xml>%s<MsgType><![CDATA[voice]]></MsgType><Voice><MediaId><![CDATA[%s]]></MediaId></Voice></xml>"
	replyVideo   = "<xml>%s<MsgType><![CDATA[video]]></MsgType><Video><MediaId><![CDATA[%s]]></MediaId><Title><![CDATA[%s]]></Title><Description><![CDATA[%s]]></Description></Video></xml>"
	replyMusic   = "<xml>%s<MsgType><![CDATA[music]]></MsgType><Music><Title><![CDATA[%s]]></Title><Description><![CDATA[%s]]></Description><MusicUrl><![CDATA[%s]]></MusicUrl><HQMusicUrl><![CDATA[%s]]></HQMusicUrl><ThumbMediaId><![CDATA[%s]]></ThumbMediaId></Music></xml>"
	replyNews    = "<xml>%s<MsgType><![CDATA[news]]></MsgType><ArticleCount>%d</ArticleCount><Articles>%s</Articles></xml>"
	replyComm    = "<ToUserName><![CDATA[%s]]></ToUserName><FromUserName><![CDATA[%s]]></FromUserName><CreateTime>%d</CreateTime>"
	replyArticle = "<item><Title><![CDATA[%s]]></Title> <Description><![CDATA[%s]]></Description><PicUrl><![CDATA[%s]]></PicUrl><Url><![CDATA[%s]]></Url></item>"
)

type ReceiveMsg struct {
	ToUserName   string
	FromUserName string
	CreateTime   int
	MsgType      string
	MsgId        int64
	Content      string
	PicUrl       string
	MediaId      string
	Format       string
	ThumbMediaId string
	LocationX    float32 `xml:"Location_X"`
	LocationY    float32 `xml:"Location_Y"`
	Scale        float32
	Label        string
	Title        string
	Description  string
	Url          string
	Event        string
	EventKey     string
	Ticket       string
	Latitude     float32
	Longitude    float32
	Precision    float32
	Recognition  string
	Status       string
}
type Music struct {
	Title        string
	Description  string
	MusicUrl     string
	HQMusicUrl   string
	ThumbMediaId string
}
type Article struct {
	Title       string
	Description string
	PicUrl      string
	Url         string
}

func WxAuth(timestamp, nonce, signatureIn, echostr, wxToken string) string {
	signatureGen := makeSignature(timestamp, nonce, wxToken)
	if signatureGen != signatureIn {
		fmt.Printf("signatureGen != signatureIn signatureGen=%s,signatureIn=%s\n", signatureGen, signatureIn)
		return ""
	} else {
		//如果请求来自于微信，则原样返回echostr参数内容 以上完成后，接入验证就会生效，开发者配置提交就会成功。
		return echostr
	}
}

//算法排序
func makeSignature(timestamp, nonce, wxToken string) string {
	//1. 将 plat_token、timestamp、nonce三个参数进行字典序排序
	sl := []string{wxToken, timestamp, nonce}
	sort.Strings(sl)
	//2. 将三个参数字符串拼接成一个字符串进行sha1加密
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))

	return fmt.Sprintf("%x", s.Sum(nil))
}
func replyHeader(reMsg ReceiveMsg) string {
	return fmt.Sprintf(replyComm, reMsg.FromUserName, reMsg.ToUserName, time.Now().Unix())
}

// Reply text message
func ReplyText(text string, reMsg ReceiveMsg) string {
	return fmt.Sprintf(replyText, replyHeader(reMsg), text)
}

// Reply image message
func ReplyImage(mediaId string, reMsg ReceiveMsg) string {
	return fmt.Sprintf(replyImage, replyHeader(reMsg), mediaId)
}

// Reply voice message
func ReplyVoice(mediaId string, reMsg ReceiveMsg) string {
	return fmt.Sprintf(replyVoice, replyHeader(reMsg), mediaId)
}

// Reply video message
func ReplyVideo(mediaId string, title string, description string, reMsg ReceiveMsg) string {
	return fmt.Sprintf(replyVideo, replyHeader(reMsg), mediaId, title, description)
}

// Reply music message
func ReplyMusic(m *Music, reMsg ReceiveMsg) string {
	return fmt.Sprintf(replyMusic, replyHeader(reMsg), m.Title, m.Description, m.MusicUrl, m.HQMusicUrl, m.ThumbMediaId)
}

// Reply news message (max 10 news)
func ReplyNews(articles []Article, reMsg ReceiveMsg) string {
	var ctx string
	for _, article := range articles {
		ctx += fmt.Sprintf(replyArticle, article.Title, article.Description, article.PicUrl, article.Url)
	}
	return fmt.Sprintf(replyNews, replyHeader(reMsg), len(articles), ctx)
}
func InitRecMsg(xml_str string) ReceiveMsg {
	reText := ReceiveMsg{}
	xml.Unmarshal([]byte(xml_str), &reText)
	return reText
}
