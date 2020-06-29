# MapReduce
mapreduce分布式计算框架

step1:
go run clientNode.go map节点个数            // 客户端节点用于将数据进行分片

step2:
go run mapNode.go mapID reduce节点的个数    // map节点用于将数据进行分区

step3:
go run reduceNode.go map节点个数 reduceID  // 将分区数据进行规约

step4:
go run mergeNode.go reduce节点的个数       // 将数据进行合并

设计部分参考我的博客：https://blog.csdn.net/qq_34276797/article/details/106989421

