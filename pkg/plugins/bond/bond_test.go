package bond

import (
	"fmt"
	"testing"
	"time"
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
