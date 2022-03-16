package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type userSessionDetail struct {
	userName       string
	sessionType    string
	sessionHours   string
	sessionMinutes string
	sessionSeconds string
}

type usrSessionDetails []userSessionDetail

type userSession struct {
	userName  string
	startTime time.Time
	endTime   time.Time
}

type userSessions []userSession

type sessionDetails struct {
	numberOfSessions int
	sessionTime      float64
}

var emptyTime = time.Time{}

func main() {
	sessions := usrSessionDetails{}
	if len(os.Args) <= 1 {
		log.Fatal("cannot start application without file path as command line argument")
	}
	fl, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("file opening failed with: %v", err)
	}
	defer fl.Close()
	sc := bufio.NewScanner(fl)
	for sc.Scan() {
		//validate log line, if invalid then skip logline
		if len(strings.Split(sc.Text(), " ")) != 3 {
			continue
		}
		logLine := strings.Split(sc.Text(), " ")
		//validate timestamp, if invalid then skip logline
		if len(strings.Split(logLine[0], ":")) != 3 {
			continue
		}
		//validate Start or End marker, if invalid then skip logline
		if logLine[2] != "Start" && logLine[2] != "End" {
			continue
		}
		logTime := strings.Split(logLine[0], ":")
		s := userSessionDetail{}
		s.sessionHours = logTime[0]
		s.sessionMinutes = logTime[1]
		s.sessionSeconds = logTime[2]
		s.userName = logLine[1]
		s.sessionType = logLine[2]
		sessions = append(sessions, s)
	}
	sessions.Process()

}

func (us usrSessionDetails) Process() {
	if len(us) < 1 {
		return
	}
	var logsStartTime time.Time
	var logsEndTime time.Time
	var err error

	//start time from logs file
	logsStartTime, err = time.Parse("15:04:05", strings.Join([]string{us[0].sessionHours, us[0].sessionMinutes, us[0].sessionSeconds}, ":"))
	if err != nil {
		log.Fatal(err)
	}

	//end time from logs file
	logsEndTime, err = time.Parse("15:04:05", strings.Join([]string{us[len(us)-1].sessionHours, us[len(us)-1].sessionMinutes, us[len(us)-1].sessionSeconds}, ":"))
	if err != nil {
		log.Fatal(err)
	}

	//prepare sessions
	uSessions := userSessions{}
	for _, sess := range us {
		processed := false

		if sess.sessionType == "Start" {
			n := userSession{userName: sess.userName}
			n.startTime, err = time.Parse("15:04:05", strings.Join([]string{sess.sessionHours, sess.sessionMinutes, sess.sessionSeconds}, ":"))
			if err != nil {
				log.Fatal(err)
			}
			uSessions = append(uSessions, n)
			processed = true
		}

		if sess.sessionType == "End" {
			for idx, ses := range uSessions {
				if ses.endTime == emptyTime && ses.userName == sess.userName {
					uSessions[idx].endTime, err = time.Parse("15:04:05", strings.Join([]string{sess.sessionHours, sess.sessionMinutes, sess.sessionSeconds}, ":"))
					processed = true
					break
				}
			}
		}

		if !processed {
			n := userSession{userName: sess.userName}
			n.endTime, err = time.Parse("15:04:05", strings.Join([]string{sess.sessionHours, sess.sessionMinutes, sess.sessionSeconds}, ":"))
			if err != nil {
				log.Fatal(err)
			}
			uSessions = append(uSessions, n)
		}
	}

	usrSessions := make(map[string]userSessions)

	for _, val := range uSessions {
		usess, ok := usrSessions[val.userName]
		if !ok {
			usess = userSessions{}
		}
		usess = append(usess, val)
		usrSessions[val.userName] = usess
	}

	finalDetails := make(map[string]sessionDetails)
	for _, v1 := range usrSessions {
		for _, v2 := range v1 {
			v3, ok := finalDetails[v2.userName]
			if !ok {
				v3 = sessionDetails{}
			}
			v3.numberOfSessions += 1
			if v2.startTime != emptyTime {
				if v2.endTime != emptyTime {
					v3.sessionTime += v2.endTime.Sub(v2.startTime).Seconds()
				} else {
					v3.sessionTime += logsEndTime.Sub(v2.startTime).Seconds()
				}
			} else {
				v3.sessionTime += v2.endTime.Sub(logsStartTime).Seconds()
			}
			finalDetails[v2.userName] = v3
		}
	}

	for k, v := range finalDetails {
		fmt.Printf("%v %v %v\n", k, v.numberOfSessions, v.sessionTime)
	}

}
