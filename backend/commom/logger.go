package commom

import (
	"fmt"
	"os"
	"time"
)

var Logfile, _ = os.OpenFile(fmt.Sprintf("./log/log-%s.log", time.Now().Format("2006-01-02")), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
