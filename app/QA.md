####Q
1.每次修改db数据时，都会修改meta page么？顺序是如何，一致性又如何保证呢？

2.node.spill先看node是否需要spill，然后将node重新申请page再写入，也就是说在内存中的node都会被重新写入新的page，这不是很浪费么

3.node.rebalance时，如果当前节点size<threshold就找到subling去合并，但是没有判断subling size和合并后的node的size，这样合并不会导致节点size过大么

4.freelist有什么用，为什么要设计他

5.不做wal，数据修改过程中失败如何处理

6.node相当于page的内存缓存，一些其他db使用了cache(lru)实现缓存以便于缓存无限增长，bolt对node做了什么处理呢

7.为什么inline的bucket有page字段呢

8.bucket的nodes变量map[pgid]*node在什么时候会缓存page呢


####A
3.node.rebalance只被bucket.rebalance调用，而bucket.rebalance只在tx.commit调用，调用rebalance后又调用了spill解决node size过大问题

8.bucket.node函数(bucket.go L643)会缓存page至nodes中，而bucket.node()函数在cursor.node()被调用，cursor.node()在增删查改中被调用