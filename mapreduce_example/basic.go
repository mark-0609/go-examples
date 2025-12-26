package mapreduce_example

import (
	"fmt"
	"log"
	"strings"

	"github.com/zeromicro/go-zero/core/mr"
)

// DemoMapReduce 演示 MapReduce 的基本用法
// 场景：计算一组整数的平方和
func DemoMapReduce() {
	fmt.Println("=== DemoMapReduce Start ===")

	// 模拟输入数据
	nums := []int{1, 2, 3, 4, 5}

	// MapReduce 函数
	// Generate: 产生数据
	// Mapper: 处理数据（并发执行）
	// Reducer: 聚合结果
	val, err := mr.MapReduce(func(source chan<- int) {
		for _, v := range nums {
			source <- v
		}
	}, func(item int, writer mr.Writer[int], cancel func(error)) {
		// 模拟耗时操作
		// time.Sleep(time.Millisecond * 10)
		writer.Write(item * item)
	}, func(pipe <-chan int, writer mr.Writer[int], cancel func(error)) {
		var sum int
		for v := range pipe {
			sum += v
		}
		writer.Write(sum)
	})

	if err != nil {
		log.Printf("MapReduce error: %v", err)
	} else {
		fmt.Printf("Sum of squares: %v\n", val)
	}
	fmt.Println("=== DemoMapReduce End ===")
}

// DemoMapReduceChan 演示 MapReduceChan 的基本用法
// 场景：处理一个已经存在的 Channel 数据
func DemoMapReduceChan() {
	fmt.Println("=== DemoMapReduceChan Start ===")

	// 创建并填充 channel
	source := make(chan int, 10)
	go func() {
		for i := 0; i < 5; i++ {
			source <- i
		}
		close(source)
	}()

	// MapReduceChan 直接消费 channel
	val, err := mr.MapReduceChan(source, func(item int, writer mr.Writer[int], cancel func(error)) {
		writer.Write(item * 2)
	}, func(pipe <-chan int, writer mr.Writer[[]int], cancel func(error)) {
		var result []int
		for v := range pipe {
			result = append(result, v)
		}
		writer.Write(result)
	})

	if err != nil {
		log.Printf("MapReduceChan error: %v", err)
	} else {
		fmt.Printf("Doubled values: %v\n", val)
	}
	fmt.Println("=== DemoMapReduceChan End ===")
}

// DemoForEach 演示 ForEach 的基本用法
// 场景：并发处理一组任务，不需要返回值，只需确保全部完成
func DemoForEach() {
	fmt.Println("=== DemoForEach Start ===")

	names := []string{"alice", "bob", "charlie", "david"}

	// ForEach 并发处理每个元素
	// 注意：go-zero v1.9+ ForEach 没有返回值
	mr.ForEach(func(source chan<- string) {
		for _, name := range names {
			source <- name
		}
	}, func(item string) {
		// 模拟处理
		upper := strings.ToUpper(item)
		fmt.Printf("Processed: %s -> %s\n", item, upper)
	})

	fmt.Println("All tasks completed (ForEach does not return error in new version)")
	fmt.Println("=== DemoForEach End ===")
}
