#!/usr/bin/env bash
killall -9 open_im_test_client
uidCount=7
messageCount=100

cmd1num=1
cmd2num=1
cmd3num=1
cmd4num=1
cmd5num=1
cmd6num=1
cmd7num=1


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

for ((i = 1; i <= ${cmd7num}; i++)); do
echo 7 >> cmd.txt
done





for ((i = 1; i <= ${uidCount}; i++)); do
    nohup ./open_im_test_client -uid $i -uid_count ${uidCount} -message_count ${messageCount} >>./openIM.log.$i 2>&1 &
done

