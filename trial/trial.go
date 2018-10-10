package trial

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	// set this when building with
	// `-ldflags "-X 'github.com/alittlebrighter/coach/trial.TRIAL_DATE_LENGTH=$(date)|${trial_length_in_days}'"`
	TRIAL_DATE_LENGTH string
	expiration        time.Time

	trialEnded = errors.New("The trial period has ended.\n" +
		"You can purchase coach at https://coach.alittlebrighter.io#downloads")
	ExpireNotice = ""
	Version      = "PRO"
)

func init() {
	if err := checkTrial(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

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

	ExpireNotice = "This trial will end on " + expiration.Add(time.Duration(-1*24)*time.Hour).Format("2006-01-02") + "!\n"
	Version = "TRIAL"
	return nil
}
