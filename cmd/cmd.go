package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"termtyper/database"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/wish/v2"
	"charm.land/wish/v2/activeterm"
	"charm.land/wish/v2/bubbletea"
	lm "charm.land/wish/v2/logging"
	"github.com/charmbracelet/ssh"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// TODO: Gamify the software. Add levels, achievements, stats, etc.
// TODO: Add a daily/weekly challenge with a global leaderboard. Maybe also add a local leaderboard for each user.
// TODO: Add user levels. The leveling system should be similar to monkeytype, where you level up by earning experience points from typing.
// TODO: Keep track of user test history. Allow users to view their past test result stats.
// TODO: Use the test history to show users their progress over time. Maybe add some basic analytics, like average WPM over the past x number of tests.
// TODO: Keep track of time spent in the app, and show it to the user in their profile.
// TODO: Add a really good readme.md with screenshots, gifs, and maybe even a demo video. The readme should also include instructions on how to use the software.

var (
	host           = "localhost"
	port           = 22222
	privateKeyPath string
)

type Session struct {
	mu            sync.Mutex
	User          *database.ApplicationUser
	RemoteAddr    string
	Authenticated bool
	LastActivity  time.Time
}

var (
	Version       = "dev"
	sshServerFlag bool
	RootCmd       = &cobra.Command{
		Use:  "TermTyper",
		Long: "TermTyper - Terminal Typing Test",
		RunE: func(cmd *cobra.Command, args []string) error {
			termWidth, termHeight, err := term.GetSize(int(os.Stdout.Fd()))

			if err != nil {
				fmt.Println("Error getting terminal size", err)
				return err
			}

			p := tea.NewProgram(
				initModel(
					termenv.ColorProfile(),
					termenv.ForegroundColor(),
					termWidth, termHeight,
					&Session{
						LastActivity: time.Now(),
						User: &database.ApplicationUser{
							Id:       -1,
							Username: "Guest",
							Config:   &database.DefaultConfig,
						},
					},
				),
			)

			_, err = p.Run()
			return err
		},
	}
	serveCmd = &cobra.Command{
		Use:  "serve",
		Long: "Serve as an SSH server",
		RunE: func(cmd *cobra.Command, args []string) error {
			resolvedHost, err := resolveHost(host)
			if err != nil {
				return fmt.Errorf("failed to resolve host: %w", err)
			}

			s, err := wish.NewServer(
				wish.WithAddress(fmt.Sprintf("%s:%d", resolvedHost, port)),
				wish.WithHostKeyPath(privateKeyPath),
				wish.WithMiddleware(
					bubbletea.Middleware(teaHandler),
					activeterm.Middleware(),
					lm.Middleware(),
				),
			)

			if err != nil {
				return err
			}

			done := make(chan os.Signal, 1)
			signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

			log.Printf("Starting server on %s:%d", resolvedHost, port)
			if resolvedHost == "0.0.0.0" {
				ips := getLocalNetworkIPs()
				if len(ips) > 0 {
					log.Printf("Connect using: %v", ips)
				}
			}
			go func() {
				if err := s.ListenAndServe(); err != nil {
					log.Fatalln(err)
				}
			}()

			<-done

			log.Printf("Stopping SSH server on %s:%d", resolvedHost, port)
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer func() { cancel() }()
			if err := s.Shutdown(ctx); err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	RootCmd.PersistentFlags().BoolVar(&sshServerFlag, "ssh-server", false, "Serve as an SSH server")
	serveCmd.Flags().StringVarP(&privateKeyPath, "key", "k", "id_rsa", "path to the server key")
	serveCmd.Flags().StringVarP(&host, "host", "", "localhost", "address to serve on (localhost or network)")
	serveCmd.Flags().IntVarP(&port, "port", "p", port, "port to serve on")
	RootCmd.AddCommand(serveCmd)
}

func getLocalNetworkIPs() []string {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}

	return ips
}

func resolveHost(hostFlag string) (string, error) {
	switch hostFlag {
	case "localhost":
		return "localhost", nil
	case "network":
		return "0.0.0.0", nil
	default:
		return hostFlag, nil
	}
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()
	sess := &Session{
		RemoteAddr:   s.RemoteAddr().String(),
		LastActivity: time.Now(),
		User: &database.ApplicationUser{
			Id:       -1,
			Username: "Guest",
			Config:   &database.DefaultConfig,
		},
	}

	m := initModel(
		termenv.ANSI256,
		termenv.ANSIWhite,
		pty.Window.Width,
		pty.Window.Height,
		sess,
	)

	return m, nil
}
