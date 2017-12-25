package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/axgle/mahonia"
)

type Transfer struct {
	Time    time.Time
	Amount  float32
	TradeNo string
	TelHead string
	TelTail string
	Email   string
	Examed  bool
}

var TransferMap map[string]*Transfer

var count uint64
var first bool = true

// var url string = `https://my.alipay.com/tile/service/portal:recent.tile\?t\=1493948220999\&_input_charset\=utf-8\&ctoken\=ao8zUmGAs_1yxlzd\&_output_charset\=utf-8"`

var url string = `https://my.alipay.com/tile/service/portal:recent.tile?`

func init() {
	TransferMap = make(map[string]*Transfer)
}

func GetTransfer(c string) string {

	// ctokenIndex := strings.Index(c, "ctoken=")
	// if ctokenIndex < 0 {
	//     return ""
	// }
	// ctoken := c[ctokenIndex:]
	// ctoken = ctoken[:strings.IndexByte(ctoken, ';')]
	// nurl := url + ctoken
	// nurl = nurl + `&t=` + fmt.Sprint(time.Now().Unix()*1000)
	//
	// cmd := exec.Command("curl", "--cookie", `"`+c+`"`, nurl)
	//
	// output, err := cmd.CombinedOutput()
	// if err != nil {
	//     log.Println(err)
	// }

	//http request by library

	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Cookie", c)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.89 Safari/537.36")
	// fmt.Println(req.Cookies())
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return ""
	}
	output, err := ioutil.ReadAll(resp.Body)

	// fmt.Println(string(output))

	r, _ := regexp.Compile(`<td class="amount">\s*<span class="amount-pay">\+ [0-9]*\.[0-9]{2}</span>\s*</td>\s*<td class="detail">\s*<a .{5,300}</a>`)
	outstr := string(output)
	outstr = strings.Replace(outstr, "\n", "", 99999)
	outstr = strings.Replace(outstr, "\r", "", 99999)

	b := r.FindAllString(outstr, 5)

	for _, v := range b {
		trans := &Transfer{}
		amountRegexp, _ := regexp.Compile(`\+ [0-9]+\.[0-9]{2}`)

		amount := amountRegexp.FindString(v)
		fmt.Sscanf(amount[2:], "%f", &(trans.Amount))

		tradeNoRegexp, _ := regexp.Compile(`bizInNo=[0-9]{20,50}`)
		tradeNo := tradeNoRegexp.FindString(v)[len("bizInNo="):]
		trans.TradeNo = tradeNo

		fmt.Println(*trans)

		//find new transfer
		if TransferMap[trans.TradeNo] == nil {

			TransferMap[trans.TradeNo] = trans

			// first to catch html
			if first == true {
				defer func() {
					first = false
				}()
				continue
			}

			url := `https://shenghuo.alipay.com/send/queryTransferDetail.htm?tradeNo=` + trans.TradeNo

			cmd := exec.Command("curl", "--cookie", `"`+c+`"`, url)
			output, err := cmd.CombinedOutput()

			if err != nil {
				log.Println(err)
				continue
			}

			//remove \n \r and convert string from gdk to utf8
			outputstr := string(output)
			outputstr = strings.Replace(outputstr, "\n", "", 99999)
			outputstr = strings.Replace(outputstr, "\r", "", 99999)
			outputstr = mahonia.NewDecoder("gbk").ConvertString(outputstr)

			//get infomation form string
			r, _ := regexp.Compile(`<th>对方信息：</th>\s*<td>.*\s*(([0-9]{3}\*{4}[0-9]{4})|([0-9a-z]{1,16}\*+@[0-9a-z]*.com))`)
			outstr = r.FindString(outputstr)

			if outstr != "" {
				log.Println(outstr)
				r, _ = regexp.Compile(`[0-9]{3}\*{4}[0-9]{4}`)
				if tel := r.FindString(outstr); tel != "" {
					log.Println(tel)
					trans.TelHead = tel[0:3]
					trans.TelTail = tel[7:11]
				} else {
					r, _ = regexp.Compile(`[0-9a-z]{1,16}\*+@[0-9a-z]*.com`)
					if email := r.FindString(outstr); email != "" {
						trans.Email = email
					}
				}

			}

			trans.Time = time.Now()
			log.Printf("New transfer! %+v\n", trans)
		}
	}

	fmt.Println("count: ", count)
	count++
	return ""
}
