package bond

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"pigeon/config"
	"pigeon/pkg/framework"
	"strconv"
	"strings"
	"time"

	"github.com/axgle/mahonia"
	"github.com/robfig/cron/v3"
)

const Name = "bond"

// StockMonitor 股票监控器
type StockMonitor struct {
	stocks map[string]*StockInfo
	Cron   string `json:"cron"`
}

// StockInfo 股票信息
type StockInfo struct {
	Name       string
	AlertPrice float64
}

// StockData 实时股票数据
type StockData struct {
	Name  string
	Price float64
	Time  string
}

func New(arguments config.Arguments) framework.Plugin {
	// log.Println("bond plugin init", Arguments)

	sm := &StockMonitor{
		Cron:   arguments["cron"].(string),
		stocks: make(map[string]*StockInfo),
	}

	bondList := arguments["bonds"].([]interface{})
	for _, bondItem := range bondList {
		bond := bondItem.(map[interface{}]interface{})
		code := bond["code"].(string)

		sm.stocks[code] = &StockInfo{
			Name:       bond["name"].(string),
			AlertPrice: bond["price"].(float64),
		}
	}

	return sm
}

func (sm StockMonitor) Name() string {
	return Name
}

// SendMessage send msg to channel
func (sm StockMonitor) SendMessage(bindMsg interface{}, Msg chan config.Msg) error {
	msgjson := bindMsg.(string)
	msg := config.Msg{
		Title:       Name,
		Description: msgjson,
		Channel:     9,
	}

	Msg <- msg
	return nil
}

func (sm StockMonitor) Run(Msg chan config.Msg, config config.Config, c *cron.Cron) error {
	_, err := c.AddFunc(sm.Cron, func() {
		log.Printf("%s 开始执行任务", Name)
		msg := sm.CheckAlerts()
		if msg == "" {
			log.Println("no need bond remainder")
			return
		}
		if err := sm.SendMessage(msg, Msg); err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		return err
	}
	return nil
}

// NewStockMonitor 创建新的股票监控器
func NewStockMonitor() *StockMonitor {
	return &StockMonitor{
		stocks: make(map[string]*StockInfo),
	}
}

// AddStockAlert 添加股票提醒
func (sm *StockMonitor) AddStockAlert(stockCode string, alertPrice float64, stockName string) {
	sm.stocks[stockCode] = &StockInfo{
		Name:       stockName,
		AlertPrice: alertPrice,
	}
	fmt.Printf("✅ 已添加监控: %s(%s), 提醒价格: %.2f元\n", stockName, stockCode, alertPrice)
}

// GetStockPrice 获取股票实时价格（新浪财经API）
func (sm *StockMonitor) GetStockPrice(stockCode string) (*StockData, error) {
	// 构造股票代码
	var code string
	if strings.HasPrefix(stockCode, "6") {
		code = "sh" + stockCode
	} else {
		code = "sz" + stockCode
	}

	// 构建请求URL
	url := fmt.Sprintf("http://hq.sinajs.cn/list=%s", code)

	// 创建HTTP客户端
	client := &http.Client{Timeout: 10 * time.Second}

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	req.Header.Set("Referer", "http://finance.sina.com.cn/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应数据
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析数据
	// 数据格式: var hq_str_sh600000="浦发银行,12.35,12.36,...";
	// 将 GBK 编码转换为 UTF-8
	decoder := mahonia.NewDecoder("gbk")
	if decoder == nil {
		return nil, fmt.Errorf("不支持的编码: gbk")
	}

	dataStr := decoder.ConvertString(string(body))
	start := strings.Index(dataStr, `="`)
	end := strings.Index(dataStr, `";`)

	if start == -1 || end == -1 {
		return nil, fmt.Errorf("无效的响应数据")
	}

	stockInfo := dataStr[start+2 : end]
	fields := strings.Split(stockInfo, ",")
	if len(fields) < 3 {
		return nil, fmt.Errorf("数据格式错误")
	}

	// 解析当前价格（第3个字段）
	currentPrice, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return nil, fmt.Errorf("解析价格失败: %v", err)
	}

	return &StockData{
		Name:  fields[0],
		Price: currentPrice,
		Time:  time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

// CheckAlerts 检查所有股票提醒
func (sm *StockMonitor) CheckAlerts() string {
	log.Printf("\n📊 %s 检查股票价格...\n", time.Now().Format("15:04:05"))
	var msgs string
	for stockCode, stockInfo := range sm.stocks {
		stockData, err := sm.GetStockPrice(stockCode)
		if err != nil {
			log.Printf("❌ 获取 %s(%s) 价格失败: %v\n", stockInfo.Name, stockCode, err)
			continue
		}

		log.Printf("   %s(%s): %.2f元\n", stockData.Name, stockCode, stockData.Price)

		// 检查是否达到提醒条件
		if stockData.Price <= stockInfo.AlertPrice {
			msgs += sm.sendAlert(stockData, stockCode, stockInfo.AlertPrice) + "\n\n"
		}

		// 避免请求过于频繁
		time.Sleep(1 * time.Second)
	}
	return msgs
}

// sendAlert 发送提醒
func (sm *StockMonitor) sendAlert(stockData *StockData, stockCode string, alertPrice float64) string {
	// priceStr := fmt.Sprintf("%.2f", stockData.Price)
	// alertPriceStr := fmt.Sprintf("%.2f", alertPrice)
	// msg := strings.Repeat("=", 50) + "\n\n" + "🚨 股票: " + stockData.Name + "(" + stockCode + ")" + "\n\n" + "🚨 当前价格:" + string(priceStr) + "\n\n" + "🚨 已低于设定价格: " + string(alertPriceStr) + "\n\n" + strings.Repeat("=", 50)
	msg := fmt.Sprintf(
		"🚨%s🚨\n\n🚨 股票价格提醒!\n\n🚨 股票: %s(%s)\n\n🚨 当前价格: %.2f元\n\n🚨 已低于设定价格: %.2f元\n\n🚨 时间: %s\n\n🚨%s🚨",
		strings.Repeat("=", 50),
		stockData.Name, stockCode,
		stockData.Price,
		alertPrice,
		stockData.Time,
		strings.Repeat("=", 50),
	)

	// 同时打印到控制台
	log.Println(msg)

	return msg
}

// StartMonitoring 开始监控
func (sm *StockMonitor) StartMonitoring(interval time.Duration) {
	fmt.Printf("🎯 开始监控股票价格，检查间隔: %v\n", interval)
	fmt.Println("按 Ctrl+C 停止监控")

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		sm.CheckAlerts()
	}
}

func main() {
	// 创建股票监控器
	monitor := NewStockMonitor()

	// 添加要监控的股票
	// 参数：股票代码, 提醒价格, 股票名称
	monitor.AddStockAlert("000001", 100.5, "平安银行")
	monitor.AddStockAlert("600036", 300.0, "招商银行")
	monitor.AddStockAlert("000858", 150.0, "五粮液")
	monitor.AddStockAlert("600519", 16000.0, "贵州茅台")

	// 开始监控，每30秒检查一次
	monitor.StartMonitoring(2 * time.Second)
}
