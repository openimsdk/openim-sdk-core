#!/usr/bin/env bash
killall -9 open_im_test_client
uidCount=1310
messageCount=1000

cmd1num=10
cmd2num=500
cmd3num=500
cmd4num=100
cmd5num=100
cmd6num=100


for ((i = 1; i <= ${cmd1num}; i++)); do
echo 1 >> cmd.txt
done


for ((i = 1; i <= ${cmd2num}; i++)); do
echo 2 >> cmd.txt
done


for ((i = 1; i <= ${cmd3num}; i++)); do
echo 3 >> cmd.txt
done


for ((i = 1; i <= ${cmd4num}; i++)); do
echo 4 >> cmd.txt
done


for ((i = 1; i <= ${cmd5num}; i++)); do
echo 5 >> cmd.txt
done


for ((i = 1; i <= ${cmd6num}; i++)); do
echo 6 >> cmd.txt
done





for ((i = 1; i <= ${uidCount}; i++)); do
    nohup ./open_im_test_client -uid $i -uid_count ${uidCount} -message_count ${messageCount} >>./openIM.log.$i 2>&1 &
done

