package main

import (
	"gotalk/constants"
	"gotalk/models"
	"runtime"
	"testing"
	"time"
)

//TODO extend tests
func TestClientSession(t *testing.T) {
	testContent := testUi.newUi()
	quit := make(chan bool)

	go func() {
		for {
			select {
			case <-quit:
				return
			default:
				testMsg.Body = nil
				client.SetReadDeadline(time.Now().Add(timeoutDuration))
				n, err := client.Read(testBuf)
				if err == nil {
					err := testMsg.UnmarshalMSG(testBuf[:n])
					if err == nil {
						switch testMsg.Action {
						case constants.ACTION_SENDMESSAGE:
							if len(testMsg.Body) == 0 {
								t.Errorf("bad test user message")
								t.Fail()
							} else {
								t.Log("ACTION_SENDMESSAGE passed")
							}
						case constants.ACTION_SENDSTATUS:
							if len(testMsg.Body) == 0 {
								t.Errorf("bad test status mesage")
								t.Fail()
							} else {
								t.Log("ACTION_SENDSTATUS passed")
							}
						case constants.ACTION_REVISION:
							if len(testMsg.Body) != 1 {
								t.Errorf("bad test revision message")
								t.Fail()
							} else {
								t.Log("ACTION_REVISION passed")
							}
						}
					}
				}
			}

		}
	}()

	if runtime.GOOS == "windows" {
		newline = "\r\n"
	} else {
		newline = "\n"
	}

	testWindow.SetContent(testContent)

	testSnd.Body = nil
	testSnd.Body = append(testSnd.Body, "Testmessage")
	testSession.WriteMessage(testSnd.Body)

	testSnd.Body = nil
	testSnd.Body = append(testSnd.Body, "Test status")
	testSession.WriteStatus(testSnd.Body)

	models.SendJSONMessage(testSession.Conn, constants.ACTION_REVISION, []string{constants.REVISION})

	quit <- true

	client.Close()
	server.Close()
}
