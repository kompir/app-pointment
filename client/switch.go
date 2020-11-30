package client

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

/** Values passed from CLI */
type idsFlag []string

func (list *idsFlag) String() string {
	return strings.Join(*list, ",")
}

/** Override the value ot the list */
func (list *idsFlag) Set(v string) error {
	*list = append(*list, v)
	return nil
}

/** HTTP client for communicating with the Backend API */
type BackendHTTPClient interface {
	Create(title, message string, duration time.Duration) ([]byte, error)
	Edit(id string, title string, message string, duration time.Duration) ([]byte, error)
	List(ids []string) ([]byte, error)
	Delete(ids []string) error
	Healthy(host string) bool
}

/** CLI command switch */
type Switch struct {
	client        BackendHTTPClient
	backendAPIURI string
	commands      map[string]func() func(string) error
}

/** Creates a new instance of command Switch */
func NewSwitch(uri string) Switch { //constructor
	httpClient := NewHTTPClient(uri)
	s := Switch{client: httpClient, backendAPIURI: uri}
	s.commands = map[string]func() func(string) error{
		"create": s.create,
		"edit":   s.edit,
		"list":   s.list,
		"delete": s.delete,
		"health": s.health,
	}
	return s
}

/** Executes the given command by args */
func (s Switch) Switch() error {
	cmdName := os.Args[1]
	cmd, ok := s.commands[cmdName]
	if !ok {
		return fmt.Errorf("Invalid command: '%s\n' ", cmdName)
	}
	return cmd()(cmdName)
}

func (s Switch) Help() {
	var help string
	for name := range s.commands {
		help += name + "\t --help\n"
	}
	fmt.Printf("Usage of: %s:\n <command> [<args>\n%s]", os.Args[0], help)

}

/** Parse command flags */
func (s Switch) parseCmd(cmd *flag.FlagSet) error {
	err := cmd.Parse(os.Args[2:])
	if err != nil {
		return wrapError("Could not parse '"+cmd.Name()+"' command flags.", err)
	}
	return nil
}

/** Create new reminder */
func (s Switch) create() func(string) error {
	return func(cmd string) error {
		createCmd := flag.NewFlagSet(cmd, flag.ExitOnError)
		t, m, d := s.reminderFlags(createCmd)
		if err := s.checkArgs(3); err != nil {
			return err
		}
		if err := s.parseCmd(createCmd); err != nil {
			return err
		}

		res, err := s.client.Create(*t, *m, *d)
		if err != nil {
			return wrapError("Could not create reminder.", err)
		}

		fmt.Printf("Reminder created successfuly:\n%s.", string(res))
		return nil
	}
}

/** Edit a reminder */
func (s Switch) edit() func(string) error {
	return func(cmd string) error {
		ids := idsFlag{}
		editCmd := flag.NewFlagSet(cmd, flag.ExitOnError)
		editCmd.Var(&ids, "id", "The ID of the reminder to edit.")
		t, m, d := s.reminderFlags(editCmd)
		if err := s.checkArgs(2); err != nil {
			return err
		}
		if err := s.parseCmd(editCmd); err != nil {
			return err
		}

		lastID := ids[len(ids)-1]
		res, err := s.client.Edit(lastID, *t, *m, *d)
		if err != nil {
			return wrapError("Could not edit reminder.", err)
		}

		fmt.Printf("Reminder edited successfuly:\n%s.", string(res))
		return nil
	}
}

/** Get aall reminders */
func (s Switch) list() func(string) error {
	return func(cmd string) error {
		ids := idsFlag{}
		listCmd := flag.NewFlagSet(cmd, flag.ExitOnError)
		listCmd.Var(&ids, "id", "The ID of the reminder to list.")

		if err := s.checkArgs(1); err != nil {
			return err
		}
		if err := s.parseCmd(listCmd); err != nil {
			return err
		}

		res, err := s.client.List(ids)
		if err != nil {
			return wrapError("Could not list reminder.", err)
		}

		fmt.Printf("Appoitments listed successfuly:\n%s.", string(res))
		return nil
	}
}

/** Delete a reminder */
func (s Switch) delete() func(string) error {
	return func(cmd string) error {
		ids := idsFlag{}
		deleteCmd := flag.NewFlagSet(cmd, flag.ExitOnError)
		deleteCmd.Var(&ids, "id", "The ID of the reminder to delete.")

		if err := s.checkArgs(1); err != nil {
			return err
		}
		if err := s.parseCmd(deleteCmd); err != nil {
			return err
		}

		err := s.client.Delete(ids)
		if err != nil {
			return wrapError("Could not delete reminder.", err)
		}

		fmt.Printf("Reminder deleted successfuly:\n%v\n.", ids)
		return nil
	}
}

/** Ping the host */
func (s Switch) health() func(string) error {
	return func(cmd string) error {
		var host string
		healthCmd := flag.NewFlagSet(cmd, flag.ExitOnError)
		healthCmd.StringVar(&host, "host", s.backendAPIURI, "Host to ping for helath.")
		if err := s.parseCmd(healthCmd); err != nil {
			return err
		}
		if !s.client.Healthy(host) {
			fmt.Printf("Host %s is down.\n", host)
		} else {
			fmt.Printf("Host %s is up and running.\n", host)

		}
		return nil
	}
}

/** A specific flags */
func (s Switch) reminderFlags(f *flag.FlagSet) (*string, *string, *time.Duration) {
	t, m, d := "", "", time.Duration(0)
	f.StringVar(&t, "title", "", "Reminder title.")
	f.StringVar(&t, "t", "", "Reminder title.")
	f.StringVar(&m, "message", "", "Reminder message.")
	f.StringVar(&m, "m", "", "Reminder message.")
	f.DurationVar(&d, "duration", 0, "Reminder duration.")
	f.DurationVar(&d, "d", 0, "Reminder duration.")
	return &t, &m, &d
}

/** Checks if the number of passed in args is greater or equal to min args */
func (s Switch) checkArgs(minArgs int) error {
	if len(os.Args) == 3 && os.Args[2] == "--help" {
		return nil
	}
	if len(os.Args)-2 < minArgs {
		fmt.Printf("Incorrect use of %s\n%s %s --help\n",
			os.Args[1],
			os.Args[0],
			os.Args[1])
		return fmt.Errorf("%s expects atleast %d args, %d provided",
			os.Args[1],
			minArgs,
			len(os.Args)-2)
	}
	return nil
}
