# 分布式存储漫游指南 1: 2025年了，存储硬件啥样了？

5202 年，我们能享受到怎样的存储硬件性能？


本漫游指南从数据中心硬件、单机性能出发，到构建分布式的对象存储产品的见闻和思考。
系列文章是笔者从事对象存储研发经历的小结。从自己曾有的疑惑出发，收集资料和思考形成，希望能和读者多多交流学习。


## 1 前言

我是一名 2050 年的开发者，莫名其妙穿越回了 2025 年。发现团队正准备从零开始构建自研对象存储产品，承载万亿对象和百PiB级别的数据。
我挠了挠头： 2050 年代这些事情都是 AI 自动搞定的，人类早都成最底层奴隶了。2025 年，存储服务器和数据中心硬件的性能究竟是啥水平？

哈哈，笔者只是开篇开了个玩笑。在构建分布式存储之前，我们不妨先来看看当下的硬件能力。实际硬件的选择，一定是成本和性能平衡的。

主要是明确需求：延迟 (us 级别 / ms 级别 / 秒级)、数据规模 (TiB / PiB / EiB) 和吞吐 (io型 iops / 带宽型 GBps or TBps)。

让我们先参考云服务提供商的一份性能基准，对目标产品的性能数量级大概有个印象。


| 特性                  | SSD 块存储(高性能) [1]                     | SSD 块存储(性价比) [1]                  | 对象存储 [2,3]                          |
|-----------------------|------------------------------------------|----------------------------------------|----------------------------------------|
| **容量范围**           | 1,261~65,536 GiB (单盘)                  | 1~65,536 GiB (单盘)                    | ＞PiB 级别                            |
| **IOPS性能**          | min{1,800+50*容量, 1,000,000}           | min{1,800+8*容量, 6,000}               | 单账户 30K qps                         |
| **吞吐**              | min{120+0.5*容量, 4,000} MB/s           | min{100+0.15*容量, 150} MB/s           | 单账户下载 100Gbps, 上传 20 Gbps (典型) |
| **4K随机读写延迟**    | 0.2ms                                   | 1~3ms                                  | 100ms~2000ms                           |

当然了，这里只是典型业务需求的参考值。如果扩大集群磁盘数量、使用更好的 SSD，都会提高集群的性能上限。

遗憾的是，地球人在 2025 年广泛使用的计算机架构仍然是冯诺依曼机，离不开 CPU、内存、磁盘、网络等几大件。
无论是服务器、消费级 PC 和智能手机，差别在于性能的高低（以及有多少预算）。

## 2 超高性能: 高速 SSD 与 RDMA 网络

LLM 训练和推理的存储、虚拟机集群的块存储等需求，对存储系统的性能要求很高。这些应用舍得花费成本，追求极低延迟和巨大吞吐。

考虑我们平时能够买到的云虚拟机产品。其中的硬盘由块存储服务提供，直接通过网络挂载到位于计算节点上的虚拟机上。任何微小的丢包和延迟波动都会在客户端成倍放大，甚至 Hang 住虚拟机中的应用，严重影响用户吞吐。因此这些应用是不能接受性能波动的。

另外对于 AI 训练场景，相较于 GPU 的成本和获取难度，追求存储系统的极限吞吐是很有性价比的。

![](https://static.zdfmc.net/imgs/2025/03/94e6511eebbf08bcf9620ac0f35a8bc2.png)
图: Dell PowerEdge R960 机架式服务器, 可选配 24 x NVMe SSD 和 2 x 100Gbps 网卡 [4]

2025 年开源的 3FS 存储系统就是原生基于现代 NVMe SSD 和 RDMA 网络构建。其使用的集群为 180 台机器，单机规格为 2×200Gbps InfiniBand 网卡 + 16×14TiB NVMe SSD。压测环境吞吐达到了 6.6 TiB/s。[5]


### 2.1 PCIe 5.0 与 NVMe SSD
既然是旗舰性能，自然少不了 PCIe 5.0 SSD 的加持。自己组装过主机的玩家可能知道，PCIe 每一代的前进基本都代表着总线带宽翻倍提升。

那么问题来了，什么是 PCIe，什么又是 NVMe？

我们的漫游指南不会去具体深究总线和接口的定义，不过可以这么试着理解
- **PCIe 是高速公路**，决定车道数量（x1/x2/x4）和车速（Gen3/4/5）。
- **NVMe 是交通规则**，优化车辆（数据包）的调度效率，避免拥堵（如并行队列、多核优化）。
- **NVMe 是协议，PCIe 是通道**，二者协同工作（NVMe over PCIe）



下表展示了现代企业级高性能 SSD 的性能水平。

| Peak Throughput and Bandwidth      | Solidigm PS1010 7.68TB | KIOXIA CM7-R 7.68TB | Samsung PM1743 7.68TB |
|------------------------------------|------------------------|----------------------|-----------------------|
| 256K sequential read (1T/64Q)      | 14,848MB/s             | 12,092MB/s           | 14,495MB/s            |
| 256K sequential write (1T/64Q)     | 7,117MB/s              | 5,796MB/s            | 6,052MB/s             |
| 4K random read (8T/32Q)            | 2,084,960 IOPS         | 1,963,066 IOPS       | 1,900,838 IOPS        |
| 4K random write (8T/32Q)           | 408,721 IOPS           | 301,061 IOPS         | 319,758 IOPS          |

表: 企业级 PCIe 5.0 NVMe SSD 性能级别 [6]

可以看到，4K 随机读来到了恐怖的 200万 iops，连续读也到达了 14GB/s。这个速度完全可以打满单网卡 100Gbps 的带宽。

![](https://static.zdfmc.net/imgs/2025/03/17e840857eb1678227a97456c4a7dcae.png)
图: U.2 接口的企业级 SSD，不会使用消费级常见的 `m.2` 规格，为了保证散热和可靠性，一般是 `E3.S 7.5mm` / `U.2 15mm`。

### 2.2 NVMe 也代表着存储软件栈的革新
NVMe 不仅是概念上的变化，kernel 中的驱动也是原生面向 SSD 设备设计的。下表是 与 AHCI (SATA) 的特性对比。

![](https://static.zdfmc.net/imgs/2025/03/1d7a9e6135bcd551762c761bfd84ba66.png)表: AHCI 和 NVMe 对比 [7]

AHCI 面向旋转的机械硬盘设计，单队列、深度 32 对于 15000 RPM 的机械硬盘是足够的。物理层面随机读取最多约为 500 IOPS。

而进化到高速 SSD，4K 随机读甚至来到了恐怖的百万 IOPS。因此，必须大幅度提升带宽上线、大幅提高队列和深度、支持多核中断、无锁设计，才能在软件驱动层面最大限度地发挥 SSD 的硬件能力。

有趣的是，2010 年 kernel 块设备/文件系统正准备迎接存储设备从机械硬盘到 SSD 的巨变，其实也仅仅是 15 年前的事情 [8]。

因此，高性能存储基础设施的发展，包含了硬件的软件的共同发展，为上层构建超高性能分布式存储提供了舞台。

### 2.3 高速网络与 RDMA

网卡规格达到 100 Gbps 后，为了追求极限性能，需要配合 RDMA 技术使用。RDMA（Remote Direct Memory Access，远程直接内存访问）

如果你问颇有经验的开发者，为啥非要使用 RDMA 不可？他们多半的回答是为了绕过操作系统内核和CPU，直接在应用程序内存之间进行数据传输。这个系统内核态的消耗究竟有多大呢？

我们使用一个没那么恰当的例子，观察用户态函数切换和内核态系统调用的消耗。
```
#include <iostream>
#include <chrono>
#include <unistd.h>

const int ITERATIONS = 1000000;

// 阻止编译器优化用户态函数
__attribute__((noinline)) void user_space_call() {
    // 模拟与 getpid 相当的轻量级操作
    static volatile int counter = 0; // volatile 防止优化
    counter += 1;                   // 简单算术运算
    asm volatile("" : "+r" (counter) : : "memory"); // 确保内存屏障
}

void kernel_mode_test() {
    volatile pid_t dummy; // 防止编译器优化系统调用
    
    auto start = std::chrono::steady_clock::now();
    for (int i = 0; i < ITERATIONS; ++i) {
        dummy = getpid(); // 系统调用
    }
    auto end = std::chrono::steady_clock::now();
    
    auto duration = std::chrono::duration_cast<std::chrono::microseconds>(end - start).count();
    std::cout << "内核态调用总耗时: " << duration << " μs | 单次开销: " 
              << static_cast<double>(duration)/ITERATIONS << " μs\n";
}

void user_mode_test() {
    auto start = std::chrono::steady_clock::now();
    for (int i = 0; i < ITERATIONS; ++i) {
        user_space_call(); // 用户态调用
    }
    auto end = std::chrono::steady_clock::now();
    
    auto duration = std::chrono::duration_cast<std::chrono::microseconds>(end - start).count();
    std::cout << "用户态调用总耗时: " << duration << " μs | 单次开销: " 
              << static_cast<double>(duration)/ITERATIONS << " μs\n";
}

int main() {
    std::cout << "--- 用户态/内核态切换开销测试 ---\n";
    user_mode_test();
    kernel_mode_test();
    
    return 0;
}
```

```
g++ -O0 -o sys_call_bench sys_call_bench.cc
```

得到的耗时对比如下：

```
--- 用户态/内核态切换开销测试 ---
用户态调用总耗时: 7139 μs | 单次开销: 0.007139 μs
内核态调用总耗时: 199932 μs | 单次开销: 0.199932 μs
```

这非常直观地说明了如果通过完全内核去收发网络包、多了几次内存拷贝和中断，可能还真会 hang 死用户的虚拟机磁盘。所以还真是不得不使用 RDMA 技术了。

常见的网络技术有：
- **InfiniBand**：原生支持RDMA，提供低延迟和高吞吐。
- **RoCE**（RDMA over Converged Ethernet）：在以太网上实现RDMA，兼容现有网络设施。
- **iWARP**：通过TCP实现RDMA，适合广域网场景。

除了时延和性能的提升，RDMA 技术还有一个很大的好处是节省了 CPU。对于大型集群，节省的 CPU 可以直接用来跑一些离线计算任务，提高整个集群的利用率。

那么挑战是什么呢？

需要专有的硬件网卡和交换机。比如 RoCE，需要无损网络和流量控制算法。这就意味着开发者不能像基于 kernel tcp/ip 编程那么畅快。网络性能观测、甚至基础设施交换机硬件的兼容性都要开发者面面俱到。

Azure Storage 大范围使用 RDMA 流量的论文 [9] 中，除了描述其带来的收益外，也花了大幅篇幅描述了他们遇到的挑战：
- 需要编写基于 RDMA 的通信框架，支持失败时候回退到 tcp
- 紧密观测基础设施 RDMA 流量网络交换的性能指标，防御暂停帧风暴等事件
- 推动网络硬件供应商更新固件解决现有的兼容和性能问题

![](https://static.zdfmc.net/imgs/2025/03/536302c08a4891c948c89c96a0c48b30.png)
图: RDMA vs TCP - Azure with RDMA on NSDI '23 (老哥确实想展示一下自己刚提的新车) [9]


## 3 一般性能: 成本、性能和规模的权衡

对象存储一般用于存储媒体资源，在这个看中成本的年代，首要追求低成本，高容量、追求性能和成本的平衡。
比如多媒体资源、用户 KYC 归档。单集群数据量应当是 PiB 起步。大公司的多租户集群甚至能到 EiB 级别 [10]。

### 3.1 HDD 容量为王: 你现在知道谁是老大了吼

### 3.2 适当的 SSD

### 3.3 线性扩展能力
媒体资源一般和 CDN 配合使用，因此一般不会对少量文件进行高 IOPS 热点读写。
- 性能取决于数据分片所在的硬件能力
- 数据存储：副本与 EC

### 3.4 Raid 卡不只能做 Raid


### 2.3 成本优先: 机器老了点，又不是不能存

故障率偏高，用于冷存。大型 EC 分片、减少修复 LRC 算法

### 2.4 大规模冷存储: 磁带


## 硬盘
### SSD 和 HDD 基准性能

### NVMe 和 SATA 驱动 
Kernel / LWN 链接

### 企业级和消费级有啥区别
fsync 性能比较

### 单台机器，存储密度越大越好吗？
这里的存储密度是指单台机器的存储容量。还要考虑 CPU 核心数 ：内存：存储容量

### 单块硬盘，存储容量越大越好吗？


### JBOD 与 Raid

### Raid 卡不光能做 Raid

### 更冷的存储：磁带


## 网络
### 为什么要用 RDMA 网络

### DPDK/SPDK


## CPU 和 内存

NUMA

L3 Cache 

服务器内存的 DDR4/5

底层硬件的发展：软件工程师只要安心调用就行了，硬件工程师考虑的就多了

CPU 加速：AVX 加速，HTTPS 卸载，EC 编码

Linus: AVX is shit

其他奇怪的设备：图形加速，FPGA 等


## 参考资料
[4] [Dell PowerEdge R960 Product Page](https://www.dell.com/zh-cn/shop/cty/pdp/spd/poweredge-r960/asper960)
[5] [deepseek-ai/3FS](https://github.com/deepseek-ai/3FS)
[6] [solidigm ps1010 ssd review - storagereview.com](https://www.storagereview.com/review/solidigm-ps1010-ssd-review)
[7] [AHCI vs NVMe - Phison Blog](https://phisonblog.com/ahci-vs-nvme-the-future-of-ssds-2/)
[8] [Solid-state storage devices and the block layer](https://lwn.net/Articles/408428/)
[9] [Empowering Azure Storage with RDMA](https://www.usenix.org/conference/nsdi23/presentation/bai)
[10] [Facebook’s Tectonic Filesystem: Efficiency from Exascale](https://www.usenix.org/system/files/fast21-pan.pdf)