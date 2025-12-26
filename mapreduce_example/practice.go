package mapreduce_example

import (
	"fmt"
	"log"
	"time"

	"github.com/zeromicro/go-zero/core/mr"
)

// Product 商品结构体
type Product struct {
	ID       int
	Name     string
	Price    float64
	Stock    int
	Comments int
}

// DemoPractice 实战案例：批量商品数据聚合
// 场景：
// 给定一批商品ID，需要从不同的服务（价格服务、库存服务、评论服务）获取数据，
// 最终组装成完整的商品信息列表。
func DemoPractice() {
	fmt.Println("=== DemoPractice Start ===")

	// 1. 待处理的商品ID列表
	productIDs := []int{101, 102, 103, 104, 105, 106, 107, 108, 109, 110}

	// 2. 使用 MapReduce 并发处理
	// 返回类型 []*Product
	result, err := mr.MapReduce(func(source chan<- int) {
		// Generator: 发送任务
		for _, pid := range productIDs {
			source <- pid
		}
	}, func(item int, writer mr.Writer[*Product], cancel func(error)) {
		// Mapper: 处理单个任务
		pid := item

		// 模拟获取各个服务的数据
		// 注意：这里是在 Worker 中执行，mr 默认并发度为 16
		// 对于单个商品，我们顺序调用各个依赖服务（或者也可以在这里再开 goroutine 并发，但通常没必要，除非依赖服务响应很慢且无依赖）

		// 模拟可能发生的错误
		if pid == 105 {
			// 假设 ID 为 105 的商品数据有问题，演示错误处理
			// 如果调用 cancel，整个 MapReduce 会立即失败返回
			// cancel(errors.New("product 105 data corrupted"))
			// return

			// 如果只是想跳过该商品，不写入 writer 即可
			log.Printf("Skipping product %d due to error", pid)
			return
		}

		p := &Product{
			ID:   pid,
			Name: fmt.Sprintf("Product-%d", pid),
		}

		// 获取价格
		p.Price = getPrice(pid)
		// 获取库存
		p.Stock = getStock(pid)
		// 获取评论数
		p.Comments = getComments(pid)

		writer.Write(p)

	}, func(pipe <-chan *Product, writer mr.Writer[[]*Product], cancel func(error)) {
		// Reducer: 聚合结果
		var products []*Product
		for item := range pipe {
			products = append(products, item)
		}
		writer.Write(products)
	})

	if err != nil {
		log.Printf("Batch processing failed: %v", err)
	} else {
		products := result
		fmt.Printf("Successfully processed %d products:\n", len(products))
		for _, p := range products {
			fmt.Printf("ID: %d, Name: %s, Price: %.2f, Stock: %d, Comments: %d\n",
				p.ID, p.Name, p.Price, p.Stock, p.Comments)
		}
	}
	fmt.Println("=== DemoPractice End ===")
}

// 模拟外部服务调用

func getPrice(pid int) float64 {
	time.Sleep(time.Millisecond * 10) // 模拟延迟
	return float64(pid) * 1.5
}

func getStock(pid int) int {
	time.Sleep(time.Millisecond * 5)
	return pid % 10 * 100
}

func getComments(pid int) int {
	time.Sleep(time.Millisecond * 8)
	return pid * 2
}
