package test

//func OnlineTest(number int) {
//	t1 := time.Now()
//	RegisterOnlineAccounts(number)
//	log.Info("", "RegisterAccounts  cost time: ", time.Since(t1), "Online client number ", number)
//	t2 := time.Now()
//	var wg sync.WaitGroup
//	wg.Add(number)
//	for i := 0; i < number; i++ {
//		go func(t int) {
//			GenWsConn(t)
//			log.Info("GenWsConn, the: ", t, " user")
//			wg.Done()
//		}(i)
//	}
//	wg.Wait()
//	log.Info("", "OnlineTest finish cost time: ", time.Since(t2), "Online client number ", number)
//}

//func GenWsConn(id int) {
//	userID := GenUid(id, "online")
//	token := RunGetToken(userID)
//	wsRespAsyn := interaction.NewWsRespAsyn()
//	wsConn := interaction.NewWsConn(new(testInitLister), token, userID, false)
//	cmdWsCh := make(chan common.Cmd2Value, 10)
//	pushMsgAndMaxSeqCh := make(chan common.Cmd2Value, 1000)
//	interaction.NewWs(wsRespAsyn, wsConn, cmdWsCh, pushMsgAndMaxSeqCh, nil)
//}
