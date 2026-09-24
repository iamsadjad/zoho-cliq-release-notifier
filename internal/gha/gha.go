package gha

import (
	"fmt"
	"os"
	"strings"
)

// Info writes an info log recognized by GitHub Actions.
func Info(format string, args ...any) {
	fmt.Printf(format+"\n", args...)
}

// Warning writes a warning annotation.
func Warning(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("::warning::%s\n", escape(msg))
}

// Error writes an error annotation.
func Error(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("::error::%s\n", escape(msg))
}

// SetSecret registers a secret so Actions masks it in logs.
func SetSecret(value string) {
	if strings.TrimSpace(value) == "" {
		return
	}
	fmt.Printf("::add-mask::%s\n", value)
}

// SetOutput writes a workflow output.
func SetOutput(name, value string) error {
	path := os.Getenv("GITHUB_OUTPUT")
	if path == "" {
		// Local/dev fallback.
		fmt.Printf("::set-output name=%s::%s\n", name, escape(value))
		return nil
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Use heredoc delimiter form to safely handle multiline values.
	delimiter := "ghadelim_" + strings.ReplaceAll(name, "-", "_")
	_, err = fmt.Fprintf(f, "%s<<%s\n%s\n%s\n", name, delimiter, value, delimiter)
	return err
}

func escape(value string) string {
	r := strings.NewReplacer("%", "%25", "\r", "%0D", "\n", "%0A")
	return r.Replace(value)
}
