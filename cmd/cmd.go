package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"termtyper/database"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

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
				tea.WithAltScreen(),
			)

			_, err = p.Run()
			return err
		},
	}
	serveCmd = &cobra.Command{
		Use:  "serve",
		Long: "Serve as an SSH server",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := wish.NewServer(
				wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
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

			log.Printf("Starting server on %s:%d", host, port)
			go func() {
				if err := s.ListenAndServe(); err != nil {
					log.Fatalln(err)
				}
			}()

			<-done

			log.Printf("Stopping SSH server on %s:%d", host, port)
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
	serveCmd.Flags().StringVarP(&host, "host", "", "localhost", "address to serve on")
	serveCmd.Flags().IntVarP(&port, "port", "p", port, "port to serve on")
	RootCmd.AddCommand(serveCmd)
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

	return m, []tea.ProgramOption{tea.WithAltScreen()}
}
