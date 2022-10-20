# LC_Nexus_Namesilo

A small prictice project of Golang. Works like a DDNS script, but for IPv6, and can change a group of hosts at once.

As a http server, this program would listen on port 5066.

When recieved a `Get` query, format like `http://THISSERVER:5066/namesiloapi?prefix=aaaa:bbbb:cccc:dddd`, from a client (typically from a router), this program will check current DNS record in Namesilo, and change all hosts' dns to new value when your network get a new IPv6 prefix ,or do nothing if there is no need to change.

For personal use only.