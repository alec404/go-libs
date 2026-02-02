package dingtalk

import (
	"net/http"
	"testing"
)

var dingToken = []string{"63655d30f4d1a05b01335577951262cd3bb83f1c75b987"}

var dingTalkCli = InitDingTalk(dingToken, ".")
var dingTalkCliWithSecret = InitDingTalkWithSecret("f380b7dfb0d7697b61f68af7c2651b5c9506ce7518138320", "SECe70ee0c1014bcc903084edc096504bfdba83be356a4c0ee9d9e96da1dbffa874") // 加签

var testImg = "https://golang.google.cn/lib/godoc/images/footer-gopher.jpg"
var testUrl = "https://golang.google.cn/"
var testPhone = "10086"

func init() {
	dingTalkCli = InitDingTalk(dingToken, ".")
}

func TestTextMsgWithSecret(t *testing.T) {
	err := dingTalkCliWithSecret.SendTextMessage("加签测试", WithAtMobiles([]string{testPhone}))
	if err != nil {
		t.Errorf("TestTextMsgWithSecret expected be nil, but %v got", err)
	}
}

func TestTextMsg(t *testing.T) {
	err := dingTalkCli.SendTextMessage("Text 测试", WithAtMobiles([]string{testPhone}))
	if err != nil {
		t.Errorf("TestTextMsg expected be nil, but %v got", err)
	}
}

func TestLinkMsg(t *testing.T) {
	err := dingTalkCli.SendLinkMessage("Link title", "Link test.", testImg, testUrl)
	if err != nil {
		t.Errorf("TestLinkMsg expected be nil, but %v got", err)
	}
}

func TestMarkDownMsg(t *testing.T) {
	err := dingTalkCli.SendMarkDownMessage("Markdown title", "### Link test\n --- \n- <font color=#ff0000 size=6>红色文字</font> \n - content2.", WithAtAll())
	if err != nil {
		t.Errorf("TestMarkDownMsg expected be nil, but %v got", err)
	}
}

func TestSendDTMDMessage(t *testing.T) {
	// 有序map
	dtmdOrderMap := DingMap().
		Set("dtmdOrderMap1", "dtmdValue1").
		Set("dtmdOrderMap2", "dtmdValue2").
		Set("dtmdOrderMap3", "dtmdValue3")
	err := dingTalkCli.SendDTMDMessage("DTMD title", dtmdOrderMap)
	if err != nil {
		t.Errorf("TestSendDTMDMessage expected be nil, but %v got", err)
	}
}

func TestSendMarkDownMessageByList(t *testing.T) {
	msg := []string{
		"### Link test",
		"---",
		"- <font color=#ff0000 size=2>红色文字</font>",
		"- content2",
	}
	err := dingTalkCli.SendMarkDownMessageBySlice("Markdown title", msg, WithAtMobiles([]string{testPhone}))
	if err != nil {
		t.Errorf("TestSendMarkDownMessageByList expected be nil, but %v got", err)
	}
}

func TestActionCardMultiMsg(t *testing.T) {
	Btns := []ActionCardMultiBtnModel{{
		Title:     "test1",
		ActionURL: testUrl,
	}, {
		Title:     "test2",
		ActionURL: testUrl,
	},
	}
	//err := dingTalkCli.SendActionCardMessage("ActionCard title", "ActionCard text.", WithCardSingleTitle("title"), WithCardSingleURL(testUrl))
	err := dingTalkCli.SendActionCardMessage("ActionCard title", "- ActionCard text.", WithCardBtns(Btns), WithCardBtnVertical())
	if err != nil {
		t.Errorf("TestActionCardMultiMsg expected be nil, but %v got", err)
	}
}

func TestActionCardMultiMsgBySlice(t *testing.T) {
	Btns := []ActionCardMultiBtnModel{{
		Title:     "test1",
		ActionURL: testUrl,
	}, {
		Title:     "test2",
		ActionURL: testUrl,
	},
	}
	dm := DingMap()
	dm.Set("颜色测试", H2)
	dm.Set("失败：$$ 同行不同色 $$", RED)
	dm.Set("---", "")
	dm.Set("金色", GOLD)
	dm.Set("成功", GREEN)
	dm.Set("警告", BLUE)
	dm.Set("普通文字", N)
	//err := dingTalkCli.SendActionCardMessage("ActionCard title", "ActionCard text.", WithCardSingleTitle("title"), WithCardSingleURL(testUrl))
	err := dingTalkCli.SendActionCardMessageBySlice("ActionCard title", dm.Slice(), WithCardBtns(Btns), WithCardBtnVertical())
	if err != nil {
		t.Errorf("TestActionCardMultiMsgBySlice expected be nil, but %v got", err)
	}
}

func TestFeedCardMsg(t *testing.T) {
	links := []FeedCardLinkModel{
		{
			Title:      "FeedCard1.",
			MessageURL: testUrl,
			PicURL:     testImg,
		},
		{
			Title:      "FeedCard2",
			MessageURL: testUrl,
			PicURL:     testImg,
		},
		{
			Title:      "FeedCard3",
			MessageURL: testUrl,
			PicURL:     testImg,
		},
	}
	err := dingTalkCli.SendFeedCardMessage(links)
	if err != nil {
		t.Errorf("TestFeedCardMsg expected be nil, but %v got", err)
	}
}

func TestDingMap(t *testing.T) {
	dm := DingMap()
	dm.Set("颜色测试", H2)
	dm.Set("失败：$$ 同行不同色 $$", RED)
	dm.Set("---", "")
	dm.Set("金色", GOLD)
	dm.Set("成功", GREEN)
	dm.Set("警告", BLUE)
	dm.Set("普通文字", N)
	err := dingTalkCli.SendMarkDownMessageBySlice("color test", dm.Slice())
	if err != nil {
		t.Errorf("TestDingMap expected be nil, but %v got", err)
	}
}

func TestOutGoing(t *testing.T) {
	outgoingFunc := func(args []string) []byte {
		// do what you want to
		return NewTextMsg("hello").Marshaler()
	}
	RegisterCommand("hello", outgoingFunc, 1, true)
	http.Handle("/outgoing", &OutGoingHandler{})
	_ = http.ListenAndServe(":8000", nil)
}
