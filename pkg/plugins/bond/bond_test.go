package bond

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"pigeon/config"
)

func TestGetStockPrice(t *testing.T) {
	monitor := NewStockMonitor()

	// 定义测试股票列表
	testStocks := []struct {
		Code        string
		Name        string
		TargetPrice float64 // 仅用于参考，不用于断言
	}{
		{Code: "002304", Name: "洋河酒业", TargetPrice: 54.0},
		{Code: "600938", Name: "中国海油", TargetPrice: 22.0},
		{Code: "000858", Name: "五粮液", TargetPrice: 95.0},
		{Code: "600519", Name: "贵州茅台", TargetPrice: 1300.0},
		{Code: "002594", Name: "比亚迪", TargetPrice: 75.0},
		{Code: "hk00175", Name: "吉利汽车", TargetPrice: 12.189}, // 空代码，预期失败
		{Code: "600941", Name: "中国移动", TargetPrice: 90.14},
		{Code: "600900", Name: "长江电力", TargetPrice: 23.45},
		{Code: "600566", Name: "济川药业", TargetPrice: 18.0},
		{Code: "600398", Name: "海澜之家", TargetPrice: 5.28},
		{Code: "hk03709", Name: "赢家时尚", TargetPrice: 6.772},
		{Code: "603605", Name: "珀莱雅", TargetPrice: 59.8},
		{Code: "000538", Name: "云南白药", TargetPrice: 33.0},
		{Code: "601006", Name: "大秦铁路", TargetPrice: 4.5},
		{Code: "600600", Name: "青岛啤酒A股", TargetPrice: 58.7},
		{Code: "001872", Name: "招商港口", TargetPrice: 16.133},
		{Code: "hk02319", Name: "蒙牛乳业", TargetPrice: 15.0},
		{Code: "600887", Name: "伊利股份", TargetPrice: 23.8},
		{Code: "hk00168", Name: "青岛啤酒H股", TargetPrice: 49.19},
		{Code: "hk06049", Name: "保利物业", TargetPrice: 30.438},
		{Code: "hk02669", Name: "中海物业", TargetPrice: 4.15},
		{Code: "hk01448", Name: "福寿园", TargetPrice: 2.58},
		{Code: "hk03613", Name: "同仁堂国药", TargetPrice: 8.479},
		{Code: "hk00392", Name: "北京控股", TargetPrice: 32.643},
		{Code: "hk00696", Name: "中国民航信息网络", TargetPrice: 10.4},
		{Code: "600377", Name: "宁沪高速", TargetPrice: 10.6},
		{Code: "hk00177", Name: "江苏宁沪高速H股", TargetPrice: 9.53},
		{Code: "hk00548", Name: "深高速", TargetPrice: 4.754},
		{Code: "000651", Name: "格力电器", TargetPrice: 35.57},
	}

	fmt.Println("开始测试股票价格获取...")
	fmt.Println("========================================")

	for i, stock := range testStocks {
		if stock.Code == "" {
			fmt.Printf("[%02d] %s (代码为空，跳过)\n", i+1, stock.Name)
			continue
		}

		fmt.Printf("[%02d] 正在查询: %s (%s) 参考目标价: %.2f\n",
			i+1, stock.Name, stock.Code, stock.TargetPrice)

		// 调用 GetStockPrice 获取价格
		stockData, err := monitor.GetStockPrice(stock.Code)

		if err != nil {
			fmt.Printf("   ❌ 查询失败: %v\n", err)
		} else {
			// 计算与目标价的差异
			diff := stockData.Price - stock.TargetPrice
			diffPercent := (diff / stock.TargetPrice) * 100

			fmt.Printf("   ✅ 成功获取: %s\n", stockData.Name)
			fmt.Printf("      当前价: %.3f\n", stockData.Price)
			fmt.Printf("      目标价: %.2f\n", stock.TargetPrice)
			fmt.Printf("      差异值: %.3f (%.2f%%)\n", diff, diffPercent)
			fmt.Printf("      时间: %s\n", stockData.Time)
		}
		fmt.Println("----------------------------------------")

		// 添加短暂延迟，避免请求过快被限
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("========================================")
	fmt.Println("测试完成")
}

// TestToFloat64 覆盖 yaml 把 price 解析成各种类型的情况
func TestToFloat64(t *testing.T) {
	cases := []struct {
		name    string
		input   interface{}
		want    float64
		wantErr bool
	}{
		{name: "int", input: 5, want: 5},
		{name: "int 零", input: 0, want: 0},
		{name: "int64", input: int64(5), want: 5},
		{name: "uint64", input: uint64(5), want: 5},
		{name: "float64", input: 5.28, want: 5.28},
		{name: "float64 整数值", input: 1250.0, want: 1250},
		{name: "float32", input: float32(4.5), want: 4.5},
		{name: "string", input: "12.189", want: 12.189},
		{name: "string 带空格", input: " 6.15 ", want: 6.15},
		{name: "string 非数字", input: "abc", wantErr: true},
		{name: "bool", input: true, wantErr: true},
		{name: "nil", input: nil, wantErr: true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := toFloat64(c.input)
			if c.wantErr {
				if err == nil {
					t.Fatalf("toFloat64(%#v) 期望报错，实际返回 %v", c.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("toFloat64(%#v) 意外报错: %v", c.input, err)
			}
			if got != c.want {
				t.Fatalf("toFloat64(%#v) = %v, 期望 %v", c.input, got, c.want)
			}
		})
	}
}

// TestNewWithIntPrice 回归测试：config.yaml 里 price 写成整数（如 price: 5）不能报错
func TestNewWithIntPrice(t *testing.T) {
	yamlContent := `
plugins:
  bond:
    cron: "0 10 * * 1-5"
    bonds:
      - code: "600519"
        name: "贵州茅台"
        price: 5
      - code: "000858"
        name: "五粮液"
        price: 68.0
      - code: "600398"
        name: "海澜之家"
        price: "5.28"
        reason: "低估"
      - code: "000001"
        name: "坏数据"
        price: "不是数字"
      - code: "000002"
        name: "缺价格"
`

	yamlPath := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("写入临时配置失败: %v", err)
	}

	conf, err := config.GetConf(yamlPath)
	if err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}

	arguments, ok := conf.Plugins["bond"]
	if !ok {
		t.Fatal("配置里没找到 bond 插件")
	}

	// price 是 int / float / string 时都应该能正常构建，不能 panic
	monitor := New(arguments).(*StockMonitor)

	want := map[string]float64{
		"600519": 5,
		"000858": 68.0,
		"600398": 5.28,
	}
	for code, wantPrice := range want {
		stock, ok := monitor.stocks[code]
		if !ok {
			t.Errorf("%s 没有被加载", code)
			continue
		}
		if stock.AlertPrice != wantPrice {
			t.Errorf("%s 的 AlertPrice = %v, 期望 %v", code, stock.AlertPrice, wantPrice)
		}
	}

	if stock := monitor.stocks["000858"]; stock.Name != "五粮液" {
		t.Errorf("000858 的 Name = %q, 期望 %q", stock.Name, "五粮液")
	}
	if stock := monitor.stocks["600398"]; stock.BuyReason != "低估" {
		t.Errorf("600398 的 BuyReason = %q, 期望 %q", stock.BuyReason, "低估")
	}

	// 价格非法或缺失的条目应被跳过，而不是让整个插件挂掉
	for _, code := range []string{"000001", "000002"} {
		if _, ok := monitor.stocks[code]; ok {
			t.Errorf("%s 价格非法/缺失，应该被跳过", code)
		}
	}
}
