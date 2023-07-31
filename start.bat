cd /d %~dp0./testv3new
start ./pressure_test.exe --test.run TestPressureTester_PressureSendMsgs -m 1000 -t 100 -s 8774164829 -r 5798710778 -g 3906853784
start ./pressure_test.exe --test.run TestPressureTester_PressureSendGroupMsgs -m 1000 -t 100 -s register_test_1 -r register_test_10 -g 3906853784