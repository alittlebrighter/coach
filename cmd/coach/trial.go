package main

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

var (
	TRIAL_DATE_LENGTH string
	expiration        time.Time

	trialEnded = errors.New("The trial period has ended.\n" +
		"You can purchase coach at https://coach.alittlebrighter.io#download")
	expireNotice = ""
)

func checkTrial() error {
	if strings.TrimSpace(TRIAL_DATE_LENGTH) == "" {
		return nil
	}

	parts := strings.Split(TRIAL_DATE_LENGTH, "|")
	startDate, err := time.Parse(time.UnixDate, parts[0])
	if err != nil || !startDate.After(time.Unix(0, 0)) {
		return nil
	}

	var days int
	days, err = strconv.Atoi(parts[1])
	if err != nil || days == 0 {
		return nil
	}

	expiration = startDate.Add(time.Duration(days*24) * time.Hour)

	if time.Now().After(expiration) {
		return trialEnded
	}

	expireNotice = "This trial will end on " + expiration.Add(time.Duration(-1*24)*time.Hour).Format("2006-01-02") + "!\n"
	version = "TRIAL"
	return nil
}
