package app_session

import (
	"github.com/kataras/iris/v12/sessions"
	"github.com/kataras/iris/v12/sessions/sessiondb/boltdb"
	"os"
)

//定义后台通用的session
var Sess = sessions.New(sessions.Config{Cookie: "seesionId"})

var Sessdb, _ = boltdb.New("./sessions.db", os.FileMode(0750))

//Sess.UseDatabase(Sessdb)
