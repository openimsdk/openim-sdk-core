#!/usr/bin/env bash
source ./common.sh




while ((1))
do

 time3=$(date "+%Y-%m-%d %H:%M:%S")
 echo $time3




echo -e "test client num: \c"
ps -ef |grep open_im_test_client | grep -v grep | wc -l


echo -e "login&recv client num: \c"
grep "login do test, only login" openIM.log* | wc -l

echo -e "login&send&recv  client num: \c"
grep "login do test, login and send" openIM.log* | wc -l

echo -e "login&send&recv&random sleep  client num: \c"
grep "random sleep and send" openIM.log* | wc -l

echo -e "expect send num:\c"
let var=`expr ${messageCount}*${cmd2num}+${messageCount}*${cmd3num}+${messageCount}*${cmd4num}+10`
echo $var

echo -e "send num: \c"
grep "func send" openIM.log* | wc -l

echo -e "send msg success: \c"
grep "test_openim: send msg success" openIM.log*  | wc -l

echo -e "send msg failed: \c"
grep "test_openim: send msg failed" openIM.log* | wc -l

echo -e "recv msg: \c"
grep "test_openim: " openIM.log*  |grep "recv time" | wc -l

echo -e "openim ws  recv push msg: \c"
grep "openim ws  recv push msg" openIM.log* | wc -l

echo -e "pull msg num: \c"
grep "open_im pull one msg:" openIM.log* | wc -l


sleep 5
echo ""

done
