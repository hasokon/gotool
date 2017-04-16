## Goで作るコマンドラインツール

# d2b

10進数の値を2進数で表示するだけのコマンド

**Example**

    $ d2b 345
    $ 101011001

**Options**

- r : 2進数を引数に取って10進数に変換する

# goping

 Goで実装したPingコマンド
 
**Example**

    $ goping google.com                                                                                       
    PING google.com (172.217.26.46) 23 bytes of data.
    43 bytes from nrt12s17-in-f14.1e100.net. (172.217.26.46) : icmp_seq=1 ttl=52 time=72.730 ms
    43 bytes from nrt12s17-in-f14.1e100.net. (172.217.26.46) : icmp_seq=2 ttl=52 time=69.636 ms
    43 bytes from nrt12s17-in-f46.1e100.net. (172.217.26.46) : icmp_seq=3 ttl=52 time=75.482 ms
    43 bytes from nrt12s17-in-f14.1e100.net. (172.217.26.46) : icmp_seq=4 ttl=52 time=78.564 ms
    43 bytes from nrt12s17-in-f46.1e100.net. (172.217.26.46) : icmp_seq=5 ttl=52 time=83.943 ms
    ^C
    --- google.com ping statistics ---
    5 packets transmitted, 5 received, 0% packet loss, time 5389ms
    rtt min/avg/max/mdev = 69.636/76.071/83.943/4.146 ms

[Go で ping](http://tyamagu2.xyz/articles/go_ping/)を参考にさせて頂きました
