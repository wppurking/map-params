## mp: Map Params
因为在 Golang 中经常使用 `map[string]interface{}` 这个类型, 但每次频繁的进行参数值转换是在麻烦,
看到 `gocraft/work` 中对 `map[string]interface{}` 类型的抽象的辅助方法比较好使用, 然后就单独将
这段代码抽取出来公用, 放在自己的 repo 中让不同的 golang 项目之间共享..
