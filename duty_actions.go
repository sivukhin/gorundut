package gorundut

import (
	"path"
	"time"
)

type dutyAction struct {
	name string
	f    func(at time.Time) error
}

func NewDutyAction(name string, action func(at time.Time) error) DutyAction {
	return dutyAction{name: name, f: action}
}

func (s dutyAction) Do(at time.Time) error { return s.f(at) }
func (s dutyAction) Name() string          { return s.name }

type dumpProfileAction struct {
	directory string
	profile   RuntimeProfileName
}

func NewDumpProfileAction(directory string, profile RuntimeProfileName) DutyAction {
	return &dumpProfileAction{directory: directory, profile: profile}
}

func (d dumpProfileAction) Name() string { return "DumpProfile:" + string(d.profile) }
func (d dumpProfileAction) Do(at time.Time) error {
	prefix := at.Format("2006-01-02_15-04-05")
	_, err := MakeProfileDump(path.Join(d.directory, prefix), d.profile)
	return err
}
