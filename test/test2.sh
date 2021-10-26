#!/usr/bin/env bash
killall -9 open_im_test_client
source ./common.sh

for ((i = 1; i <= ${uidCount}; i++)); do
    nohup ./open_im_test_client -uid $i -uid_count ${uidCount} -message_count ${messageCount} >>./openIM.log.$i 2>&1 &
done

