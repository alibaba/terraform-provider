#!/bin/bash

touch ~/.gitcookies
chmod 0600 ~/.gitcookies

git config --global http.cookiefile ~/.gitcookies

tr , \\t <<\__END__ >>~/.gitcookies
.googlesource.com,TRUE,/,TRUE,2147483647,o,git-shuwei.yin.alibaba-inc.com=1/JPFTLZ4ijSXFTD5C0wNT-khTbbdDE6WewBwkBZ23Fr4
__END__
