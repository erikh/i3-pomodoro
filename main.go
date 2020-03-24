package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/pkg/errors"
)

const (
	stateFile = "/tmp/pomodoro-state"

	iconRun   = "➤"
	iconPause = "Ⅱ"

	defaultTimerLength      = 25 * time.Minute
	defaultShortBreakLength = 5 * time.Minute
	defaultLongBreakLength  = 15 * time.Minute

	clockPomodoro   = 0
	clockPomodoro2  = iota
	clockShortBreak = iota
	clockLongBreak  = iota
)

// State is the state of the pomodoro timer.
type State struct {
	Running   bool
	Paused    bool
	Duration  time.Duration
	LastTime  time.Time
	Now       time.Time
	ClockType int
}

func (s *State) clockText() string {
	clockText := "POMODORO"

	switch s.ClockType {
	case clockPomodoro2:
		clockText = "POMODORO 2"
	case clockShortBreak:
		clockText = "SHORT BREAK"
	case clockLongBreak:
		clockText = "LONG BREAK"
	}

	return clockText
}

func (s *State) cycleClock() {
	switch s.ClockType {
	case clockPomodoro:
		s.ClockType = clockShortBreak
	case clockPomodoro2:
		s.ClockType = clockLongBreak
	case clockShortBreak:
		s.ClockType = clockPomodoro2
	case clockLongBreak:
		s.ClockType = clockPomodoro
	}
}

func (s *State) finish() {
	s.cycleClock()
	exec.Command("i3-nagbar", "-m", fmt.Sprintf("%s!", s.clockText())).Start()
}

func (s *State) output() {
	icon := iconPause
	if s.Running && !s.Paused {
		icon = iconRun
	}

	pomodoro := fmt.Sprintf("%s: %s %s", s.clockText(), icon, s.Duration.Round(time.Second))
	fmt.Println(pomodoro)
	fmt.Println(pomodoro)
}

func (s *State) write() error {
	content, err := json.Marshal(s)
	if err != nil {
		return errors.Wrap(err, "while marshaling json")
	}

	return ioutil.WriteFile(stateFile, content, 0600)
}

func loadState() (*State, error) {
	s := &State{}
	content, err := ioutil.ReadFile(stateFile)
	if os.IsNotExist(err) {
		return s, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "while reading state file")
	}

	if err := json.Unmarshal(content, &s); err != nil {
		return nil, errors.Wrap(err, "while unmarshaling json")
	}

	return s, nil
}

func (s *State) reset() {
	s.Running = false
	s.Paused = false
	s.LastTime = s.Now

	switch s.ClockType {
	case clockPomodoro, clockPomodoro2:
		s.Duration = defaultTimerLength
	case clockLongBreak:
		s.Duration = defaultLongBreakLength
	case clockShortBreak:
		s.Duration = defaultShortBreakLength
	}

	s.write()
}

func main() {
	s, err := loadState()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	s.Now = time.Now()
	if s.LastTime.IsZero() {
		s.reset()
	}

	switch os.Getenv("BLOCK_BUTTON") {
	case "1":
		if s.Running {
			s.Paused = !s.Paused
		} else {
			s.Running = true
		}
		s.LastTime = s.Now
	case "2":
		s.reset()
	case "3":
		s.cycleClock()
		s.reset()
	}

	if s.Running && !s.Paused {
		s.Duration -= s.Now.Sub(s.LastTime)
	}

	s.LastTime = s.Now

	if s.Duration < 0 {
		s.output()
		s.finish()
		s.reset()
	} else {
		s.output()
	}

	if s.Running {
		s.write()
	}
}
