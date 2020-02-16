package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ant0ine/go-json-rest/rest"
)

var lock = sync.RWMutex{}

//Httpでもらう設定値
var UserID string
var Password string
var SendAddress string
var AreaIdString string
var ApiCert string
var SessionID string

//SpotInfo スクレイピング結果を格納する構造体
type SpotInfo struct {
	Time                              time.Time
	Area, Spot, Count, Name, Lat, Lon string
}

//JSpotinfo JSONマージャリング構造体
type JSpotinfo struct {
	Spotinfo []InnerSpotinfo `json:"spotinfo"`
}

//InnerSpotinfo 台数情報
type InnerSpotinfo struct {
	Time  string `json:"time"`
	Area  string `json:"area"`
	Spot  string `json:"spot"`
	Count string `json:"count"`
}

//Add SpotInfo構造体をJSON用にパースして加える
func (s *JSpotinfo) Add(time time.Time, area string, spot string, count string) {
	s.Spotinfo = append(s.Spotinfo, InnerSpotinfo{Time: time.Format("2006/01/02 15:04:05"), Area: area, Spot: spot, Count: count})
}

//GetSessionID ログインしてセッションIDを取得する
func GetSessionID() (string, error) {
	//リクエストBody作成
	values := url.Values{}
	values.Set("EventNo", "21401")
	values.Add("GarblePrevention", "ＰＯＳＴデータ")
	values.Add("MemberID", UserID)
	values.Add("Password", Password)
	values.Add("MemAreaID", "1")

	req, err := http.NewRequest(
		"POST",
		"https://tcc.docomo-cycle.jp/cycle/TYO/cs_web_main.php",
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return "", err
	}

	// リクエストHead作成
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "ja,en-US;q=0.9,en;q=0.8,pt;q=0.7")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Length", "144")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Host", "tcc.docomo-cycle.jp")
	req.Header.Set("Origin", "https://tcc.docomo-cycle.jp")
	req.Header.Set("Referer", "https://tcc.docomo-cycle.jp/cycle/TYO/cs_web_main.php")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.106 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, e := goquery.NewDocumentFromResponse(resp)
	if e != nil {
		return "", e
	}

	SessionID, success := doc.Find("input[name='SessionID']").Attr("value")
	if !success {
		return "", fmt.Errorf("error")
	} else {
		return SessionID, nil
	}
}

//GetSpotInfoMain スクレイピングメイン関数
func GetSpotInfoMain(AreaID string, retry bool) ([]SpotInfo, error) {
	var list []SpotInfo
	//リクエストBody作成
	values := url.Values{}
	values.Set("EventNo", "25706")
	values.Add("SessionID", SessionID)
	values.Add("UserID", "TYO")
	values.Add("MemberID", UserID)
	values.Add("GetInfoNum", "200")
	values.Add("GetInfoTopNum", "1")
	values.Add("MapType", "1")
	values.Add("MapCenterLat", "")
	values.Add("MapCenterLon", "")
	values.Add("MapZoom", "13")
	values.Add("EntServiceID", "TYO0001")
	values.Add("Location", "")
	values.Add("AreaID", AreaID)

	req, err := http.NewRequest(
		"POST",
		"https://tcc.docomo-cycle.jp/cycle/TYO/cs_web_main.php",
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return list, err
	}

	// リクエストHead作成
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "ja,en-US;q=0.9,en;q=0.8,pt;q=0.7")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Length", "219")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Host", "tcc.docomo-cycle.jp")
	req.Header.Set("Origin", "https://tcc.docomo-cycle.jp")
	req.Header.Set("Referer", "https://tcc.docomo-cycle.jp/cycle/TYO/cs_web_main.php")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.106 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return list, err
	}
	defer resp.Body.Close()

	doc, e := goquery.NewDocumentFromResponse(resp)
	if e != nil {
		return list, e
	}

	//エラーならログインし直して再チャレンジ
	if err := CheckErrorPage(doc); err != nil {
		if retry {
			SessionID, err = GetSessionID()
			if err != nil {
				return nil, err
			}
			time.Sleep(1 * time.Second)
			//再帰呼び出し（次はリトライしない）
			return GetSpotInfoMain(AreaID, false)
		} else {
			//２回目は諦める
			return nil, err
		}
	}

	//スポットリスト解析
	doc.Find("form[name^=tab_]").Each(func(i int, s *goquery.Selection) {
		spotinfo := SpotInfo{Time: time.Now()}
		err := ParseSpotInfoByText(s.Find("a").Text(), &spotinfo)
		if err != nil {
			return
		}
		if val, exist := s.Find("input[name=ParkingLat]").Attr("value"); exist {
			spotinfo.Lat = val
		}
		if val, exist := s.Find("input[name=ParkingLon]").Attr("value"); exist {
			spotinfo.Lon = val
		}
		list = append(list, spotinfo)
	})

	fmt.Printf("AreaID = %s (%d件)\n", AreaID, len(list))

	return list, nil
}

//ParseSpotInfoByText テキスト解析
// "B1-01.十思公園B1-01.Jisshi Park8台"の形式のテキストからarea,spot,name,countを取得する
func ParseSpotInfoByText(text string, s *SpotInfo) error {
	//駿河台とかが引っかからないように最後から検索する
	indexCount := strings.LastIndex(text, "台")
	if indexCount < 1 {
		return fmt.Errorf("error")
	}

	indexDot := strings.Index(text, ".")
	if indexDot < 1 {
		return fmt.Errorf("error")
	}
	// "D1-10"のコード
	code := text[:indexDot]
	if arr := strings.Split(code, "-"); len(arr) == 2 {
		s.Area = arr[0]
		s.Spot = arr[1]
	} else {
		return fmt.Errorf("error")
	}

	//逆順のループで数値じゃ無くなるところまでを台数とする
	indexNum := indexCount - 1
	for {
		if indexNum < 1 {
			return fmt.Errorf("error")
		}
		_, err := strconv.Atoi(text[indexNum : indexNum+1])
		if err != nil {
			s.Count = text[indexNum+1 : indexCount]
			break
		} else {
			indexNum--
		}
	}

	//名前を取得 基本的にはコードが2回目に現れるまでが日本語名だが、コードが1つしかないパターンもわずかにあるため分岐
	if strings.LastIndex(text, code) > indexDot {
		s.Name = text[indexDot+1 : strings.LastIndex(text, code)]
	} else {
		//このパターンでは日本語名と英語名が区別できない
		s.Name = text[indexDot+1 : indexNum]
	}
	return nil
}

//RegAllSpotInfo 全スポット登録関数
func RegAllSpotInfo() (err error) {
	SessionID, err = GetSessionID()
	if err != nil {
		return err
	}

	//特に指定してない場合は全スポット
	if AreaIdString == "" {
		AreaIdString = "1,2,3,5,6,4,10,12,7,8"
	}
	AreaIDs := strings.Split(AreaIdString, ",")
	for _, AreaID := range AreaIDs {
		if AreaID == "" {
			continue
		}
		//待ち時間いれる
		time.Sleep(5 * time.Second)
		//台数取得
		var list []SpotInfo
		list, err = GetSpotInfoMain(AreaID, true)
		if err != nil {
			continue
		}
		SendSpotInfo(list)
	}

	return nil
}

//CheckErrorPage エラーページかをチェックする
func CheckErrorPage(doc *goquery.Document) error {
	if title := doc.Find(".tittle_h1").Text(); strings.Index(title, "エラー") > -1 {
		return fmt.Errorf(strings.TrimSpace(doc.Find(".main_inner_message").Text()))
	}
	return nil
}

//SendSpotInfo DBに送信する
func SendSpotInfo(list []SpotInfo) {
	var jsonStruct JSpotinfo
	for _, s := range list {
		jsonStruct.Add(s.Time, s.Area, s.Spot, s.Count)
	}
	marshalized, _ := json.Marshal(jsonStruct)
	req, err := http.NewRequest(
		"POST",
		SendAddress,
		bytes.NewBuffer(marshalized),
	)
	if err != nil {
		return
	}

	// リクエストHead作成
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("cert", ApiCert)

	//SSL証明書を無視したクライアント作成
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
	}
	//送信
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}

//TestGetSpotInfoMain 単体テスト
func TestGetSpotInfoMain(html string) ([]SpotInfo, error) {
	var list []SpotInfo
	// ファイルからドキュメントを作成（テスト用）
	f, e := os.Open(html)
	if e != nil {
		log.Fatal(e)
	}
	defer f.Close()

	doc, e := goquery.NewDocumentFromReader(f)
	if e != nil {
		log.Fatal(e)
	}

	//スポットリスト解析
	doc.Find("form[name^=tab_]").Each(func(i int, s *goquery.Selection) {
		spotinfo := SpotInfo{Time: time.Now()}
		err := ParseSpotInfoByText(s.Find("a").Text(), &spotinfo)
		if err != nil {
			return
		}
		if val, exist := s.Find("input[name=ParkingLat]").Attr("value"); exist {
			spotinfo.Lat = val
		}
		if val, exist := s.Find("input[name=ParkingLon]").Attr("value"); exist {
			spotinfo.Lon = val
		}
		list = append(list, spotinfo)
	})

	return list, nil

}

//Start スクレイピング開始
func Start(w rest.ResponseWriter, r *rest.Request) {
	r.ParseForm()
	params := r.Form
	UserID = params.Get("id")
	Password = params.Get("password")
	SendAddress = params.Get("address")
	if UserID == "" || Password == "" || SendAddress == "" {
		w.WriteHeader(http.StatusForbidden)
		w.WriteJson("[ERROR] lack of parameter")
		return
	}
	AreaIdString = params.Get("areaID")
	if env := params.Get("env"); env != "" {
		os.Setenv("API_CERT", env)
	}
	//引数で渡さなくても環境変数にある場合がある
	if val := os.Getenv("API_CERT"); val != "" {
		ApiCert = val
	}
	//スクレイピング実行（非同期）
	go RegAllSpotInfo()
	//先にOKを返しておく
	w.WriteHeader(http.StatusOK)
	w.WriteJson("OK")
}

//Recover SQLiteからリカバリー
func Recover(w rest.ResponseWriter, r *rest.Request) {

}

// func main() {
// 	TestGetSpotInfoMain("./エラーページ.html")
// }

func main() {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/start", Start),
		rest.Get("/recover", Recover),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	port := "5005"
	if val := os.Getenv("PORT"); val != "" {
		port = val
	}
	log.Fatal(http.ListenAndServe(":"+port, api.MakeHandler()))
}
