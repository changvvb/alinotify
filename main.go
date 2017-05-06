package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

var cookies *string

func main() {
	cookies = flag.String("cookies", "", "the cookies used by requestion")
	flag.Parse()
	log.Println(*cookies)

	*cookies = `cna=XNOREQJzulgCATtDALzB4H7X; session.cookieNameId=ALIPAYJSESSIONID; unicard1.vm="K1iSL1Db2zfGhm2HjBdH8A=="; CHAIR_SESS=K6iO619fGWnMOmQO_wFsrCKhc2eBub829k_yUHiNPHHAINZX6FTrWAftFmBiZ5EgsBmxGLkxJvykRNHH_ngkvZ3kMUn5vmoTU4sGNWRtFEe1FVJ0T6Vi_lxtuMG1xQhgBlJuPZQeSkNawP49kzTrzw==; ctoken=ao8zUmGAs_1yxlzd; LoginForm=alipay_login_auth; alipay="K1iSL1Db2zfGhm2HjBdH8LWURTRFOfP6GHMNz7ANvw=="; CLUB_ALIPAY_COM=2088902277892901; iw.userid="K1iSL1Db2zfGhm2HjBdH8A=="; ali_apache_tracktmp="uid=2088902277892901"; mobileSendTime=-1; credibleMobileSendTime=-1; ctuMobileSendTime=-1; riskMobileBankSendTime=-1; riskMobileAccoutSendTime=-1; riskMobileCreditSendTime=-1; riskCredibleMobileSendTime=-1; riskOriginalAccountMobileSendTime=-1; zone=RZ24B; ALIPAYJSESSIONID=RZ24GOLSuRSVof8T6wC7wtrjwnJDBEauthGZ00RZ24; ALIPAYJSESSIONID.sig=cMjgFQdwIP2VgAiAdDN6hrKDOmNDlTPUw9muslVezqU; spanner=1VHGnC9Kos4Oy1dCyZ/e4GSrrzYBVlnw0cV1wxqJaYc=`

	// GetTransfer(*cookies)
	HttpSetup()
	Run()
}

func HttpSetup() {
	http.HandleFunc("/setcookie", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			b := []byte(`
				<html>
				<body>
				<form method="POST" target="/setcookie"> <input type=text name="cookies"/> <button type=submit />  </form>
				</body>
				</html>
			`)
			w.Write(b)

		} else if r.Method == http.MethodPost {
			newcookies := r.FormValue("cookies")
			log.Println(newcookies)
			*cookies = newcookies
		}
	})

	http.HandleFunc("/exam", func(w http.ResponseWriter, r *http.Request) {
		tel := r.URL.Query().Get("tel")
		email := r.URL.Query().Get("email")
		// amountStr := r.URL.Query().Get("amount")

		if len(tel) == 0 && len(email) == 0 {
			w.Write([]byte("Param error!"))
			return
		}

		if r.Method == http.MethodGet {
			for _, v := range TransferMap {
				if v.Examed == true {
					continue
				}
				if len(tel) >= 11 {
					if len(v.TelHead) != 0 && len(v.TelHead) != 0 {
						if v.TelHead == tel[:3] && v.TelTail == tel[7:11] {
							if time.Now().Sub(v.Time) < time.Hour*2 {
								w.Write([]byte(fmt.Sprint(v.Amount)))
								v.Examed = true
								return
							}
						}
					}
				} else if len(email) > 5 {
					if len(v.Email) == 0 {
						continue
					}
					m, n := strings.IndexByte(v.Email, '*'), strings.LastIndexByte(v.Email, '*')
					log.Println(m, n)
					part1 := v.Email[:m]
					part2 := v.Email[n+1:]

					log.Println(part1, part2)
					//campare v.Email and email
					if part1 == email[:len(part1)] && part2 == email[len(email)-len(part2):] {
						if time.Now().Sub(v.Time) < time.Hour*2 {
							w.Write([]byte(fmt.Sprint(v.Amount)))
							v.Examed = true
							return
						}
					}
				}

			}
		}

		w.Write([]byte("No result!"))
	})
	go http.ListenAndServe(":2048", nil)
}

func Run() {
	for {
		GetTransfer(*cookies)
		time.Sleep(time.Second * 2)
	}
}
