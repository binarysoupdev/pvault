package logger

func Log(msg string) {
	if logger == nil {
		return
	}

	logger.Println(msg)
}

func LogError(err error) {
	Log("[X] " + err.Error())
}

func LogCreate(msg string) {
	Log("[+] " + msg)
}

func LogDelete(msg string) {
	Log("[-] " + msg)
}

func LogInfo(msg string) {
	Log("[=] " + msg)
}
