package cli

import (
	"github.com/spf13/cobra"
)

func newPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "page",
		Short: "Page management (list, create, rename, delete, duplicate)",
		Example: `  # "Create separate pages for landing, components, and design system"
  figma-kit page create "Landing Page"
  figma-kit page create "Components"
  figma-kit page create "Design System"

  # "List all pages"
  figma-kit page list`,
	}
	cmd.AddCommand(newPageListCmd())
	cmd.AddCommand(newPageCreateCmd())
	cmd.AddCommand(newPageRenameCmd())
	cmd.AddCommand(newPageDeleteCmd())
	cmd.AddCommand(newPageDuplicateCmd())
	return cmd
}

func newPageListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all pages in the current file",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			b := newBuilder()
			b.Comment("List all pages")
			b.Line("return figma.root.children.map((p, i) => ({ index: i, id: p.id, name: p.name }));")
			output(b.String())
			return nil
		},
	}
}

func newPageCreateCmd() *cobra.Command {
	var after int
	cmd := &cobra.Command{
		Use:         "create <name>",
		Short:       "Add a new page to the file",
		Args:        cobra.ExactArgs(1),
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			b := newBuilder()
			b.Comment("Create page")
			b.Linef("const pg = figma.createPage();")
			b.Linef("pg.name = %q;", args[0])
			if after >= 0 {
				b.Linef("figma.root.insertChild(%d, pg);", after+1)
			}
			b.ReturnIDs("pg.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().IntVar(&after, "after", -1, "Insert after this page index (-1 appends)")
	return cmd
}

func newPageRenameCmd() *cobra.Command {
	return &cobra.Command{
		Use:         "rename <index> <new-name>",
		Short:       "Rename a page by index",
		Args:        cobra.ExactArgs(2),
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			b := newBuilder()
			b.Comment("Rename page")
			b.Linef("const pg = figma.root.children[%s];", args[0])
			b.Linef("if (!pg) throw new Error('No page at index %s');", args[0])
			b.Linef("pg.name = %q;", args[1])
			b.ReturnDone()
			output(b.String())
			return nil
		},
	}
}

func newPageDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:         "delete <index>",
		Short:       "Delete a page by index",
		Args:        cobra.ExactArgs(1),
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			b := newBuilder()
			b.Comment("Delete page")
			b.Line("if (figma.root.children.length <= 1) throw new Error('Cannot delete the only page');")
			b.Linef("const pg = figma.root.children[%s];", args[0])
			b.Linef("if (!pg) throw new Error('No page at index %s');", args[0])
			b.Line("pg.remove();")
			b.ReturnDone()
			output(b.String())
			return nil
		},
	}
}

func newPageDuplicateCmd() *cobra.Command {
	return &cobra.Command{
		Use:         "duplicate <index>",
		Short:       "Duplicate a page by index",
		Args:        cobra.ExactArgs(1),
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			b := newBuilder()
			b.Comment("Duplicate page")
			b.Linef("const src = figma.root.children[%s];", args[0])
			b.Linef("if (!src) throw new Error('No page at index %s');", args[0])
			b.Line("const dup = src.clone();")
			b.Line("dup.name = src.name + ' (copy)';")
			b.ReturnIDs("dup.id")
			output(b.String())
			return nil
		},
	}
}
