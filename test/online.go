// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package test

//funcation OnlineTest(number int) {
//	t1 := time.Now()
//	RegisterOnlineAccounts(number)
//	log.Info("", "RegisterAccounts  cost time: ", time.Since(t1), "Online client number ", number)
//	t2 := time.Now()
//	var wg sync.WaitGroup
//	wg.Add(number)
//	for i := 0; i < number; i++ {
//		go funcation(t int) {
//			GenWsConn(t)
//			log.Info("GenWsConn, the: ", t, " user")
//			wg.Done()
//		}(i)
//	}
//	wg.Wait()
//	log.Info("", "OnlineTest finish cost time: ", time.Since(t2), "Online client number ", number)
//}

//funcation GenWsConn(id int) {
//	userID := GenUid(id, "online")
//	token := RunGetToken(userID)
//	wsRespAsyn := interaction.NewWsRespAsyn()
//	wsConn := interaction.NewWsConn(new(testInitLister), token, userID, false)
//	cmdWsCh := make(chan common.Cmd2Value, 10)
//	pushMsgAndMaxSeqCh := make(chan common.Cmd2Value, 1000)
//	interaction.NewWs(wsRespAsyn, wsConn, cmdWsCh, pushMsgAndMaxSeqCh, nil)
//}
