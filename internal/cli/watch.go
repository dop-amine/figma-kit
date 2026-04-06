package cli

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
)

func newWatchCmd() *cobra.Command {
	var debounce int
	cmd := &cobra.Command{
		Use:   "watch <recipe.yaml>",
		Short: "Watch a recipe YAML and re-emit JS blocks on every save",
		Long: `Watches a batch recipe YAML file for changes.
On every save it re-runs the recipe and prints the JS blocks to stdout,
the same output as 'figma-kit batch'. Press Ctrl-C to stop.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("resolve path: %w", err)
			}

			// Initial emit.
			if err := emitRecipe(path); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "⚠  %s\n", err)
			}

			watcher, err := fsnotify.NewWatcher()
			if err != nil {
				return fmt.Errorf("create watcher: %w", err)
			}
			defer watcher.Close()

			if err := watcher.Add(path); err != nil {
				return fmt.Errorf("watch %s: %w", path, err)
			}

			_, _ = fmt.Fprintf(os.Stderr, "👁  watching %s  (Ctrl-C to stop)\n", path)

			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

			debounceTimer := time.NewTimer(0)
			<-debounceTimer.C // drain initial fire

			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return nil
					}
					if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
						// Debounce rapid saves (e.g. editors writing then truncating).
						debounceTimer.Reset(time.Duration(debounce) * time.Millisecond)
					}
				case <-debounceTimer.C:
					_, _ = fmt.Fprintf(os.Stderr, "\n── %s ──\n", time.Now().Format("15:04:05"))
					if err := emitRecipe(path); err != nil {
						_, _ = fmt.Fprintf(os.Stderr, "⚠  %s\n", err)
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						return nil
					}
					_, _ = fmt.Fprintf(os.Stderr, "watch error: %s\n", err)
				case <-quit:
					_, _ = fmt.Fprintln(os.Stderr, "\nstopped.")
					return nil
				}
			}
		},
	}
	cmd.Flags().IntVar(&debounce, "debounce", 150, "Debounce delay in ms (avoid double-emit on save)")
	return cmd
}

// emitRecipe reads, parses, and prints a recipe file in the same format as `batch`.
func emitRecipe(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	var rec batchRecipe
	if err := yaml.Unmarshal(data, &rec); err != nil {
		return fmt.Errorf("parse YAML: %w", err)
	}
	if len(rec.Steps) == 0 {
		return fmt.Errorf("recipe has no steps")
	}
	var out strings.Builder
	if rec.Title != "" {
		fmt.Fprintf(&out, "// Recipe: %s\n\n", rec.Title)
	}
	for i, step := range rec.Steps {
		label := step.Title
		if label == "" {
			label = fmt.Sprintf("Step %d", i+1)
		}
		js := strings.TrimSpace(step.JS)
		if js == "" {
			return fmt.Errorf("step %d (%q) has empty js", i+1, label)
		}
		fmt.Fprintf(&out, "// --- Block %d: %s ---\n", i+1, label)
		out.WriteString(js)
		if !strings.HasSuffix(js, "\n") {
			out.WriteByte('\n')
		}
		out.WriteByte('\n')
	}
	output(out.String())
	return nil
}
