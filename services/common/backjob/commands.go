package backjob

type command int

const (
	command_stop command = iota
	command_force_do_job
)
