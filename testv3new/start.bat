cd /d %~dp0
start ./pressure_test.exe --test.run TestPressureTester_PressureSendMsgs -m 1 -t 1 -s 1000 -r 1 -g 0
start ./pressure_test.exe --test.run TestPressureTester_PressureSendGroupMsgs -m 1 -t 1 -s 1000 -r 0 -g 1000