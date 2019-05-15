package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"flag"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"github.com/valyala/fasthttp"
	"time"
	"strings"
	"strconv"
)

var(
	rateLimitTime = 5 // In seconds
	userTime = make(map[string]string)
	client = fasthttp.Client{ReadTimeout: time.Second * 10, WriteTimeout: time.Second * 10}
	BotID string
	token string
	Bot *discordgo.User
	session *discordgo.Session
	messageC *discordgo.MessageCreate

	Algos = []string{
		"Scrypt",
		"SHA256",
		"ScryptNf",
		"X11",
		"X13",
		"Keccak",
		"X15",
		"Nist5",
		"eoScrypt",
		"Lyra2RE",
		"WhirlpoolX",
		"Qubit",
		"Quark",
		"Axiom",
		"Lyra2REv2",
		"ScryptJaneNf16",
		"Blake256r8",
		"Blake256r14",
		"Blake256r8vnl",
		"Hodl",
		"DaggerHashimoto",
		"Decred",
		"CryptoNight",
		"Lbry",
		"Equihash",
		"Pascal",
		"X11Gost",
		"Sia",
		"Blake2s",
		"Skunk",
	}

)

const(
	binanceMethodBase = "https://api.binance.com"
	urlMethodBase = "https://api.nicehash.com/api?method="
	statsProvider = "stats.provider&"
	statsProviderEx = "stats.provider.ex&"
	statsProviderWorkers = "stats.provider.workers&"
)

func init(){
	flag.StringVar(&token, "t", "", "Bot Token")
}

func main(){
	go forever()

	if token == "" {
		fmt.Println("No token provided")
		return
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating new bot session.")
		return
	}

	self, err := dg.User("@me")
	if err != nil {
		fmt.Println("Error creating bot user.")
		return
	}

	BotID = self.ID
	Bot = self

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error creating Discord session.")
	}

	fmt.Println("Running.")
	select{}
}

func forever(){}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate){
	session = s
	messageC = m
	//For use in other functions ^^
	c := strings.ToLower(m.Content)

	if strings.HasPrefix(c, "?balance") || strings.HasPrefix(c, "?bal") && !isRateLimited(){
		if strings.Contains(c, " ") ||  m.Author.ID == "71371505652477952" || m.Author.ID == "157630049644707840" || m.Author.ID == "134799002129530880"{
			var wallet string
			if strings.HasPrefix(c, "?balance") {

				wallet = strings.TrimPrefix(m.Content, "?balance ")

				if c == "?balance" {
					if m.Author.ID == "71371505652477952" {
						wallet = "3LDpUUQp19McZqkBSagT1qkNHyZTz5RMNy"
					}
					if m.Author.ID == "157630049644707840" {
						wallet = "3MaJ5Tj4GgwgSzFxkgRu2VUxDp5GnazMdP"
					}
					if m.Author.ID == "237779090759745536" {
						wallet = "383r1QLUjStH65DGmPkzbVyVge37D4DtR6"
					}
					if m.Author.ID == "134799002129530880" {
						wallet = "133sF5nBxasooMTRP9cofRjnyk3NQ2Wps6"
					}
				}
			}

			if strings.HasPrefix(c, "?bal") {

				wallet = strings.TrimPrefix(m.Content, "?bal ")

				if c == "?bal" {
					if m.Author.ID == "71371505652477952" {
						wallet = "3LDpUUQp19McZqkBSagT1qkNHyZTz5RMNy"
					}
					if m.Author.ID == "157630049644707840" {
						wallet = "3MaJ5Tj4GgwgSzFxkgRu2VUxDp5GnazMdP"
					}
					if m.Author.ID == "237779090759745536" {
						wallet = "383r1QLUjStH65DGmPkzbVyVge37D4DtR6"
					}
					if m.Author.ID == "134799002129530880" {
						wallet = "133sF5nBxasooMTRP9cofRjnyk3NQ2Wps6"
					}
				}
			}

			walletParams := params {
				Wallet: wallet,
			}

			results, err := GetInfo(walletParams)
			if err != nil {
				fmt.Println(err)
				return
			}

				if len(results.Result.Stats) >= 1 {
					var amount = 0.0

					for i := 0; i < len(results.Result.Stats); i++ {

						num, err := strconv.ParseFloat(results.Result.Stats[i].Balance, 64)

						if err != nil {
							fmt.Println(err)
							return

						}
						amount = amount + num
					}

					amountUSDStr, err := getPriceUSD("")
					if err != nil{
						sendMessage("Unpaid balance in BTC: " + string(strconv.FormatFloat(amount, 'f', -1, 64)))
					}

					amountUSDFloat, err := strconv.ParseFloat(amountUSDStr.Data.Rates.USD, 64)

					amountUSDFloatBal := amountUSDFloat * amount

					sendMessage("Unpaid balance in BTC: " + string(strconv.FormatFloat(amount, 'f', -1, 64)) + "\nUnpaid balance in USD: " + string(strconv.FormatFloat(amountUSDFloatBal, 'f', -1, 64)))
					return
				} else {
					sendMessage("API returned nothing.")
				}
			} else {
			sendMessage("Invalid arguments")
			return
		}
		/*}  else if strings.HasPrefix(c, "?price ") && !isRateLimited() {
			rate := strings.TrimPrefix(c, "?price ")
			amountUSDStr, err := getPriceUSD(rate)
			if err != nil {
				fmt.Println(err)
				sendMessage("Error getting price")
				return
			}
			if amountUSDStr.Data.Rates.USD == "" {
				sendMessage("No matching tickers")
				return
			}
			sendMessage("1 " + strings.ToUpper(rate) + " = $" + amountUSDStr.Data.Rates.USD + " USD")
			return
		} else if c == "?price" && !isRateLimited() {
			amountUSDStr, err := getPriceUSD("")
			if err != nil {
				fmt.Println(err)
				sendMessage("Error getting price")
				return
			}

			sendMessage("1 BTC = $" + amountUSDStr.Data.Rates.USD + " USD")
		*/} else if strings.HasPrefix(c, "?price ") && !isRateLimited(){
		tick := strings.TrimPrefix(c, "?price ")
		amountBTC, err := getTicker(tick)
		if err != nil {
			sendMessage("An error has occurred retrieving price.")
			return
		}
		if amountBTC.Price != "" {
			sendMessage(amountBTC.Price + " " + amountBTC.Symbol)
		} else {
			sendMessage("An error has occurred finding ticker data.")
		}

		return
	}


}

func sendMessage(message string) {
	session.ChannelMessageSend(messageC.ChannelID, message)
}

func getPriceUSD(rate string) (*CBTop, error){
	if rate == "" {
		rate = "BTC"
	}
	urn, err := fetchUrl("https://api.coinbase.com/v2/exchange-rates?currency=" + strings.ToUpper(rate))

	var contentResult *CBTop
	if err != nil {
		return nil, err
	}
	return contentResult, json.Unmarshal(urn, &contentResult)
}

func fetchUrl(u string) ([]byte, error) {
	response, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return contents, nil
}

func getAllTickers() {

}

func getTicker(t string) (*Price, error) {
	symData, err := fetchUrl(binanceMethodBase + "/api/v3/ticker/price?symbol=" + strings.ToUpper(t))
	if err != nil {
		return nil, err
	}
	var symResult *Price

	return symResult, json.Unmarshal(symData, &symResult)

}
func GetInfo(params params) (*Main, error) {
	urn, err := fetchUrl(urlMethodBase + statsProvider + "addr=" + params.Wallet)



	var contentResult *Main
	if err != nil {
		return nil, err
	}

	return contentResult, json.Unmarshal(urn, &contentResult)
}

//This whole damn API is a mess
type Main struct {
	Result Result`json:"result"`
}
type Result struct {
	Stats []Stats `json:"stats"`
	Addr string `json:"addr"`//btc wallet address
	Algo int `json:"algo"`//https://old.nicehash.com/?p=api
	Workers []Workers `json:"workers"`
	Payments []Payments `json:"payments"`
}

type Stats struct {
	Method string `json:"method"` //Not 100% sure what you'd need this for, but included anyway
	Balance string `json:"balance"`
	RejectedSpeed string `json:"rejected_speed"`
	Algo int `json:"algo"`
	AcceptedSpeed string `json:"accepted_speed"`
}

type Payments struct {
	Amount string `json:"amount"`
	Fee string `json:"fee"`
	TXID string `json:"TXID"`
	Time string `json:"time"` //UNIX timestamp. Not sure that's the correct type
	Type int `json:"type"`
}

type Workers struct {
	WorkerName string
	A string `json:"a"`//Accepted speed
	Rs string `json:"rs"`//Rejected speed (stale)
	TimeConnected int //in minutes
	XnsubStatus int //0 = false; 1 = true
	Difficulty string
	Location int //0 = EU; 1 = US; 2 = HK; 3 = JP
}

type Current struct {
	Algo int `json:"algo"`
	Name string `json:"name"`
	Suffix string `json:"suffix"`
	Profitability string `json:"profitability"`
	Data []Data
}

type params struct {
	Wallet string
	Algo int
}

type Data struct {
	A string `json:"a"`
	Rs string `json:"rs"`
	Data string `json:"data"`//Unpaid balance
}

type Rates struct { //TODO: Actually finish this
	AED string `json:"AED"`
	AFN string `json:"AFN"`
	ALL string `json:"ALL"`
	AMD string `json:"AMD"`
	ANG string `json:"ANG"`
	AOA string `json:"AOA"`
	ARS string `json:"ARS"`
	AUD string `json:"AUD"`
	AWG string `json:"AWG"`
	AZN string `json:"AZN"`
	BAM string `json:"BAM"`
	BBD string `json:"BBD"`
	BDT string `json:"BDT"`
	BGN string `json:"BGN"`
	BHD string `json:"BHD"`
	BIF string `json:"BIF"`
	BMD string `json:"BMD"`
	BND string `json:"BND"`
	BOB string `json:"BOB"`
	BRL string `json:"BRL"`
	BSD string `json:"BSD"`
	BTC string `json:"BTC"` //Not sure why coinbase feels the need to return this
	BTN string `json:"BTN"`
	BWP string `json:"BWP"`
	BYN string `json:"BYN"`
	BYR string `json:"BYR"`
	BZD string `json:"BZD"`
	CAD string `json:"CAD"`
	CDF string `json:"CDF"`
	CHF string `json:"CHF"`
	CLF string `json:"CLF"`
	CLP string `json:"CLP"`
	CNY string `json:"CNY"`
	COP string `json:"COP"`
	CRC string `json:"CRC"`
	CUC string `json:"CUC"`
	CVE string `json:"CVE"`
	CZK string `json:"CZK"`
	DJF string `json:"DJF"`
	DKK string `json:"DKK"`
	DOP string `json:"DOP"`
	DZD string `json:"DZD"`
	EEK string `json:"EEK"`
	EGP string `json:"EGP"`
	ERN string `json:"ERN"`
	ETB string `json:"ETB"`
	ETH string `json:"ETH"`
	EUR string `json:"EUR"`
	FJD string `json:"FJD"`
	FKP string `json:"FKP"`
	GBP string `json:"GBP"`
	GEL string `json:"GEL"`
	GGP string `json:"GGP"`
	GHS string `json:"GHS"`
	GIP string `json:"GIP"`
	GMD string `json:"GMD"`
	GNF string `json:"GNF"`
	GTQ string `json:"GTQ"`
	GYD string `json:"GYD"`
	HKD string `json:"HKD"`
	HNL string `json:"HNL"`
	HRK string `json:"HRK"`
	HTG string `json:"HTG"`
	HUF string `json:"HUF"`
	IDR string `json:"IDR"`
	ILS string `json:"ILS"`
	IMP string `json:"IMP"`
	INR string `json:"INR"`
	IQD string `json:"IQD"`
	ISK string `json:"ISK"`
	JEP string `json:"JEP"`
	JMD string `json:"JMD"`
	JOD string `json:"JOD"`
	JPY string `json:"JPY"`
	KES string `json:"KES"`
	KGS string `json:"KGS"`
	KHR string `json:"KHR"`
	KMF string `json:"KMF"`
	KRW string `json:"KRW"`
	KWD string `json:"KWD"`
	KYD string `json:"KYD"`
	KZT string `json:"KZT"`
	LAK string `json:"LAK"`
	LBP string `json:"LBP"`
	LKR string `json:"LKR"`
	LRD string `json:"LRD"`
	LSL string `json:"LSL"`
	LTC string `json:"LTC"`
	LTL string `json:"LTL"`
	LVL string `json:"LVL"`
	LYD string `json:"LYD"`
	MAD string `json:"MAD"`
	MDL string `json:"MDL"`
	MGA string `json:"MGA"`
	MKD string `json:"MKD"`
	MMK string `json:"MMK"`
	MNT string `json:"MNT"`
	MOP string `json:"MOP"`
	MRO string `json:"MRO"`
	MTL string `json:"MTL"`
	MUR string `json:"MUR"`
	MVR string `json:"MVR"`
	MWK string `json:"MWK"`
	MXN string `json:"MXN"`
	MYR string `json:"MYR"`
	MZN string `json:"MZN"`
	NAD string `json:"NAD"`
	NGN string `json:"NGN"`
	NIO string `json:"NIO"`
	NOK string `json:"NOK"`
	NPR string `json:"NPR"`
	NZD string `json:"NZD"`
	OMR string `json:"OMR"`
	PAB string `json:"PAB"`
	PEN string `json:"PEN"`
	PGK string `json:"PGK"`
	PHP string `json:"PHP"`
	PKR string `json:"PKR"`
	PLN string `json:"PLN"`
	PYG string `json:"PYG"`
	QAR string `json:"QAR"`
	RON string `json:"RON"`
	RSD string `json:"RSD"`
	RUB string `json:"RUB"`
	RWF string `json:"RWF"`
	SAR string `json:"SAR"`
	SBD string `json:"SBD"`
	SCR string `json:"SCR"`
	SEK string `json:"SEK"`
	SGD string `json:"SGD"`
	SHP string `json:"SHP"`
	SLL string `json:"SLL"`
	SOS string `json:"SOS"`
	SRD string `json:"SRD"`
	SSP string `json:"SSP"`
	STD string `json:"STD"`
	SVC string `json:"SVC"`
	SZL string `json:"SZL"`
	THB string `json:"THB"`
	TJS string `json:"TJS"`
	TMT string `json:"TMT"`
	TND string `json:"TND"`
	TOP string `json:"TOP"`
	TRY string `json:"TRY"`
	TTD string `json:"TTD"`
	TWD string `json:"TWD"`
	TZS string `json:"TZS"`
	UAH string `json:"UAH"`
	UGX string `json:"UGX"`
	USD string `json:"USD"`
	UYU string `json:"UYU"`
	UZS string `json:"UZS"`
	VEF string `json:"VEF"`
	VND string `json:"VND"`
	VUV string `json:"VUV"`
	WST string `json:"WST"`
	XAF string `json:"XAF"`
	XAG string `json:"XAG"`
	XAU string `json:"XAU"`
	XCD string `json:"XCD"`
	XDR string `json:"XDR"`
	XOF string `json:"XOF"`
	XPD string `json:"XPD"`
	XPF string `json:"XPF"`
	XPT string `json:"XPT"`
	YER string `json:"YER"`
	ZAR string `json:"ZAR"`
	ZMK string `json:"ZMK"`
	ZMW string `json:"ZMW"`
	ZWL string `json:"ZWL"`
}

type PriceSlice struct {
	Top []string `json:"[]"`
}

type Price struct {
	Symbol string `json:"symbol"`
	Price string `json:"price"`
}

//Coinbase for conversion
type CBTop struct {
	Data CBData `json:"data"`
}

type CBData struct {
	Base string `json:"base"`
	Currency string `json:"currency"`
	Price string `json:"amount"`
	Rates Rates `json:"rates"`
}

//Yoinked from dogbot because I'm too lazy to rewrite it

func isRateLimited() bool {
	//if m.Author.ID == "157630049644707840" {
	//	return false //hehe
	//}

	//A really janky way to rate limit
	var m = messageC
	tNow := time.Now() 																							//Time of type Time
	tForm := tNow.Format("Mon Jan 2 15:04:05 -0700 MST 2006")											//Time of type string
	if userTime[m.Author.ID] == "" {		 																	//Check to see if last message is nil
		userTime[m.Author.ID] = tForm 																			//Pack string into map
		return false
	} else {
		tExpand, err := time.Parse("Mon Jan 2 15:04:05 -0700 MST 2006", userTime[m.Author.ID])			//Change tForm back into tNow
		if err != nil {
			fmt.Println(err)
		}
		tSince := time.Now().Sub(tExpand)
		noPrevEntrTime, err := time.ParseDuration("2562047h47m16.854775807s")
		if err != nil {
			fmt.Println(err)
		}
		if tSince != noPrevEntrTime {
			timeSec := strconv.Itoa(rateLimitTime) + "s"
			rLimitS, err := time.ParseDuration(timeSec)
			if err != nil {
				fmt.Println("Error occurred parsing duration:")
				fmt.Println(err)
				return false
			}
			tNow := time.Now() 														//Time of type Time
			tForm := tNow.Format("Mon Jan 2 15:04:05 -0700 MST 2006")		//Time of type string
			userTime[m.Author.ID] = tForm											//Pack string into map

			if tSince > rLimitS {
				//s.ChannelMessageSend(m.ChannelID, "You are being rate-limited. Please try again later.")
				return false
			} else {
				return true
			}
		}
		return false
	}
}