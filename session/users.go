package session

import (
	"github.com/AlkBur/GoIDE/conf"
	"github.com/AlkBur/GoIDE/event"
	"github.com/AlkBur/GoIDE/log"
	"github.com/AlkBur/GoIDE/util"
	"sync"
	"time"
)

var ActiveUsers = &Users{Users: make([]*UserSession, 0)}

type UserSession struct {
	Usr   *conf.User
	Event *event.UserEventQueue
	Sid   string // IDE session id related
}

type Users struct {
	sync.RWMutex
	Users []*UserSession
}

func StartUserMonitor() {
	go func() {
		defer util.Recover()
		timeChan := time.NewTimer(time.Minute).C
		for {
			select {
			case <-timeChan:
				SaveActiveUsers()
			case q := <-event.EventQueue:
				log.Debug("Received a global event [code=%d]", q.Code)
				ActiveUsers.RLock()
				for _, user := range ActiveUsers.Users {
					ev := q.Copy()
					ev.Sid = user.Sid
					user.Event.Queue <- ev
				}
				ActiveUsers.RUnlock()
			}
		}
	}()
}

func SaveActiveUsers() {
	ActiveUsers.RLock()
	defer ActiveUsers.RUnlock()
	for _, u := range ActiveUsers.Users {
		if u.Usr.Save() {
			log.Debug("Saved online user [%s]'s configurations", u.Usr.Name)
		}
	}
}

func GetUserSessions(usr *conf.User) []*UserSession {
	ActiveUsers.RLock()
	defer ActiveUsers.RUnlock()

	ret := make([]*UserSession, 0)

	for _, s := range ActiveUsers.Users {
		if s.Usr == usr {
			ret = append(ret, s)
		}
	}
	return ret
}
