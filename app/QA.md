####Q
1.每次修改db数据时，都会修改meta page么？顺序是如何，一致性又如何保证呢？

2.node.spill先看node是否需要spill，然后将node重新申请page再写入，也就是说在内存中的node都会被重新写入新的page，这不是很浪费么

3.node.rebalance时，如果当前节点size<threshold就找到subling去合并，但是没有判断subling size和合并后的node的size，这样合并不会导致节点size过大么

4.freelist有什么用，为什么要设计他

5.不做wal，数据修改过程中失败如何处理

6.node相当于page的内存缓存，一些其他db使用了cache(lru)实现缓存以便于缓存无限增长，bolt对node做了什么处理呢

7.为什么inline的bucket有page字段呢

8.bucket的nodes变量map[pgid]*node在什么时候会缓存page呢

9.inline bucket容量达到限制时，如何分裂成一个普通的bucket呢？

10.node.spill是从下向上split，也就是从leaf向root split，为什么root node split时要清空children再对parent split一次呢

11.node.rebalance除了对当前node进行rebalance，对parent递归（spill会对children递归)，那么如何做到对整颗b+ tree rebalance呢

12.inline bucket分裂成普通bucket后，inline page被删除，rootNode对应的实际page会怎么改变呢

13.事务的提交和回滚是如何实现的

14.boltdb的MVCC是如何实现的


####A
3.node.rebalance只被bucket.rebalance调用，而bucket.rebalance只在tx.commit调用，调用rebalance后又调用了spill解决node size过大问题

6.bolt的crud从bucket开始，每个bucket有自己的缓存，所以缓存生命周期和bucket对象一致

7.因为inline bucket的page_id=0，而且该page在物理页面中是跟在bucket head后面，那么通过bucket函数获取bucket对象时，应该把该page一起读出

8.bucket.node函数(bucket.go L643)会缓存page至nodes中，而bucket.node()函数在cursor.node()被调用，cursor.node()在增删查改中被调用

9.inline bucket只有一个rootNode，rootNode split成两个leafNode，然后新生成一个parent作为rootNode，bucket的page_id记录该rootNode

10.因为root node split时其实创建了一个新的node，也就是当前node的parent，而递归只到当前node，清空children避免重复split，因为children只是缓存

11.bucket的nodes缓存了增删查改的node（叶节点)，那么rebalance实际上是从叶节点向root rebalance。而spill是从root向下 spill

12.和普通的key删除类似，其实增删查改都类似，内存中的node转换为page写入磁盘，所以boltdb的读写的单位都是page，读写一个key和一页的key消耗相同

13.数据修改先只存在于node，也就是内存中，在commit函数中的spill，实现node向page的转化，也就是写入磁盘。所以只要不commit，数据不会写入磁盘

14.数据修改后spill写入磁盘时，每次都是获取一个新page，那么旧page只要不删除，就不会影响别的事务读。删除page在freepage.release实现，在每个事务开始时，会将小于
   当前db正在运行的事务中最小的事务的pending page真正删除。也就是确定该页面没有任何事务所使用，才会加入到freepage（删除）