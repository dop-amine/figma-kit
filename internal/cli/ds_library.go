package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/dop-amine/figma-kit/internal/codegen"
	"github.com/dop-amine/figma-kit/internal/restapi"
)

func newDSLibraryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "library",
		Short: "Browse, inspect, and import published library components, styles, and variables",
		Long: `Library commands let you discover and import published Figma assets.

Discovery commands (list, info) use the REST API and require a Personal Access
Token (FIGMA_TOKEN / FIGMA_PAT). Import commands generate Plugin API JavaScript
that runs in use_figma — no PAT needed for imports.

  figma-kit ds library list --team 12345
  figma-kit ds library info a1b2c3d4
  figma-kit ds library import a1b2c3d4 --name "Hero"
  figma-kit ds library import-set x9y8z7 --variant "Size=Large"
  figma-kit ds library import-style s1t2y3
  figma-kit ds library variables`,
	}
	cmd.AddCommand(newDSLibraryListCmd())
	cmd.AddCommand(newDSLibraryInfoCmd())
	cmd.AddCommand(newDSLibraryImportCmd())
	cmd.AddCommand(newDSLibraryImportSetCmd())
	cmd.AddCommand(newDSLibraryImportStyleCmd())
	cmd.AddCommand(newDSLibraryVariablesCmd())
	return cmd
}

// ---------------------------------------------------------------------------
// ds library list
// ---------------------------------------------------------------------------

func newDSLibraryListCmd() *cobra.Command {
	var (
		teamID     string
		fileKey    string
		assetType  string
		asJSON     bool
		limit      int
		cursor     string
	)
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List published components, component sets, or styles from a team or file library",
		Example: `  figma-kit ds library list --team 123456
  figma-kit ds library list --file abc123def
  figma-kit ds library list --file abc123def --type styles
  figma-kit ds library list --team 123456 --type component-sets --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if teamID == "" && fileKey == "" {
				fileKey = resolveFileKey()
				if fileKey == "" {
					return fmt.Errorf("provide --team or --file (or set FIGMA_FILE_KEY)")
				}
			}
			client, err := resolveRESTClient()
			if err != nil {
				return err
			}

			switch assetType {
			case "", "components":
				return listComponents(client, teamID, fileKey, limit, cursor, asJSON)
			case "component-sets":
				return listComponentSets(client, teamID, fileKey, limit, cursor, asJSON)
			case "styles":
				return listStyles(client, teamID, fileKey, limit, cursor, asJSON)
			default:
				return fmt.Errorf("unknown --type %q (use: components, component-sets, styles)", assetType)
			}
		},
	}
	cmd.Flags().StringVar(&teamID, "team", "", "Team ID to query")
	cmd.Flags().StringVar(&fileKey, "file", "", "File key to query")
	cmd.Flags().StringVar(&assetType, "type", "", "Asset type: components (default), component-sets, styles")
	cmd.Flags().BoolVar(&asJSON, "json", false, "Output raw JSON")
	cmd.Flags().IntVar(&limit, "limit", 30, "Page size")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor (after)")
	return cmd
}

func listComponents(c *restapi.Client, teamID, fileKey string, limit int, cursor string, asJSON bool) error {
	var resp *restapi.ComponentsResponse
	var err error
	if teamID != "" {
		resp, err = c.GetTeamComponents(teamID, limit, cursor)
	} else {
		resp, err = c.GetFileComponents(fileKey)
	}
	if err != nil {
		return err
	}
	if asJSON {
		return writeJSON(resp)
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "KEY\tNAME\tDESCRIPTION\tUPDATED")
	for _, comp := range resp.Meta.Components {
		desc := truncateDesc(comp.Description, 40)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", comp.Key, comp.Name, desc, comp.UpdatedAt)
	}
	w.Flush()
	if resp.Meta.Cursor.After != "" {
		fmt.Fprintf(os.Stderr, "\nNext page: --cursor %s\n", resp.Meta.Cursor.After)
	}
	return nil
}

func listComponentSets(c *restapi.Client, teamID, fileKey string, limit int, cursor string, asJSON bool) error {
	var resp *restapi.ComponentSetsResponse
	var err error
	if teamID != "" {
		resp, err = c.GetTeamComponentSets(teamID, limit, cursor)
	} else {
		resp, err = c.GetFileComponentSets(fileKey)
	}
	if err != nil {
		return err
	}
	if asJSON {
		return writeJSON(resp)
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "KEY\tNAME\tDESCRIPTION\tUPDATED")
	for _, cs := range resp.Meta.ComponentSets {
		desc := truncateDesc(cs.Description, 40)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", cs.Key, cs.Name, desc, cs.UpdatedAt)
	}
	w.Flush()
	if resp.Meta.Cursor.After != "" {
		fmt.Fprintf(os.Stderr, "\nNext page: --cursor %s\n", resp.Meta.Cursor.After)
	}
	return nil
}

func listStyles(c *restapi.Client, teamID, fileKey string, limit int, cursor string, asJSON bool) error {
	var resp *restapi.StylesResponse
	var err error
	if teamID != "" {
		resp, err = c.GetTeamStyles(teamID, limit, cursor)
	} else {
		resp, err = c.GetFileStyles(fileKey)
	}
	if err != nil {
		return err
	}
	if asJSON {
		return writeJSON(resp)
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "KEY\tNAME\tTYPE\tDESCRIPTION\tUPDATED")
	for _, s := range resp.Meta.Styles {
		desc := truncateDesc(s.Description, 40)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", s.Key, s.Name, s.StyleType, desc, s.UpdatedAt)
	}
	w.Flush()
	if resp.Meta.Cursor.After != "" {
		fmt.Fprintf(os.Stderr, "\nNext page: --cursor %s\n", resp.Meta.Cursor.After)
	}
	return nil
}

// ---------------------------------------------------------------------------
// ds library info
// ---------------------------------------------------------------------------

func newDSLibraryInfoCmd() *cobra.Command {
	var (
		asJSON    bool
		assetType string
	)
	cmd := &cobra.Command{
		Use:   "info <key>",
		Short: "Show detailed info for a single published component, component set, or style",
		Example: `  figma-kit ds library info a1b2c3d4e5
  figma-kit ds library info a1b2c3d4e5 --type style --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := resolveRESTClient()
			if err != nil {
				return err
			}
			key := args[0]

			switch assetType {
			case "", "component":
				comp, err := client.GetComponent(key)
				if err != nil {
					return err
				}
				if asJSON {
					return writeJSON(comp)
				}
				printComponentInfo(comp)
			case "component-set":
				cs, err := client.GetComponentSet(key)
				if err != nil {
					return err
				}
				if asJSON {
					return writeJSON(cs)
				}
				printComponentSetInfo(cs)
			case "style":
				s, err := client.GetStyle(key)
				if err != nil {
					return err
				}
				if asJSON {
					return writeJSON(s)
				}
				printStyleInfo(s)
			default:
				return fmt.Errorf("unknown --type %q (use: component, component-set, style)", assetType)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "Output raw JSON")
	cmd.Flags().StringVar(&assetType, "type", "", "Asset type: component (default), component-set, style")
	return cmd
}

func printComponentInfo(c *restapi.Component) {
	fmt.Printf("Name:            %s\n", c.Name)
	fmt.Printf("Key:             %s\n", c.Key)
	fmt.Printf("Description:     %s\n", c.Description)
	fmt.Printf("File Key:        %s\n", c.FileKey)
	fmt.Printf("Node ID:         %s\n", c.NodeID)
	fmt.Printf("Containing Page: %s\n", c.ContainingFrame.PageName)
	fmt.Printf("Frame:           %s\n", c.ContainingFrame.Name)
	fmt.Printf("Thumbnail:       %s\n", c.ThumbnailURL)
	fmt.Printf("Created:         %s\n", c.CreatedAt)
	fmt.Printf("Updated:         %s\n", c.UpdatedAt)
}

func printComponentSetInfo(cs *restapi.ComponentSet) {
	fmt.Printf("Name:            %s\n", cs.Name)
	fmt.Printf("Key:             %s\n", cs.Key)
	fmt.Printf("Description:     %s\n", cs.Description)
	fmt.Printf("File Key:        %s\n", cs.FileKey)
	fmt.Printf("Node ID:         %s\n", cs.NodeID)
	fmt.Printf("Containing Page: %s\n", cs.ContainingFrame.PageName)
	fmt.Printf("Frame:           %s\n", cs.ContainingFrame.Name)
	fmt.Printf("Thumbnail:       %s\n", cs.ThumbnailURL)
	fmt.Printf("Created:         %s\n", cs.CreatedAt)
	fmt.Printf("Updated:         %s\n", cs.UpdatedAt)
}

func printStyleInfo(s *restapi.Style) {
	fmt.Printf("Name:            %s\n", s.Name)
	fmt.Printf("Key:             %s\n", s.Key)
	fmt.Printf("Type:            %s\n", s.StyleType)
	fmt.Printf("Description:     %s\n", s.Description)
	fmt.Printf("File Key:        %s\n", s.FileKey)
	fmt.Printf("Node ID:         %s\n", s.NodeID)
	fmt.Printf("Containing Page: %s\n", s.ContainingFrame.PageName)
	fmt.Printf("Frame:           %s\n", s.ContainingFrame.Name)
	fmt.Printf("Thumbnail:       %s\n", s.ThumbnailURL)
	fmt.Printf("Created:         %s\n", s.CreatedAt)
	fmt.Printf("Updated:         %s\n", s.UpdatedAt)
}

// ---------------------------------------------------------------------------
// ds library import — composable, Plugin API
// ---------------------------------------------------------------------------

func newDSLibraryImportCmd() *cobra.Command {
	var (
		parent string
		name   string
		posX   float64
		posY   float64
	)
	useLast := new(bool)
	cmd := &cobra.Command{
		Use:   "import <key> [key...]",
		Short: "Import published component(s) by key and create instance(s)",
		Long: `Generates Plugin API JavaScript to import a published component via
figma.importComponentByKeyAsync() and create an instance.

Supports multiple keys for bulk import. Works standalone or in compose
workflows with --parent / --last for chaining.`,
		Example: `  figma-kit ds library import a1b2c3d4e5
  figma-kit ds library import a1b2c3d4e5 --name "Main Hero" --parent 123:456
  figma-kit ds library import a1b2 b2c3 c3d4

  # In compose:
  figma-kit compose "ds library import a1b2 --name Hero" "fx glow --last"`,
		Annotations: map[string]string{"composable": "true"},
		Args:        cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b := newBuilder()

			var ids []string
			for i, key := range args {
				varName := fmt.Sprintf("inst%d", i)
				b.ImportComponent(key, varName)
				if name != "" && i == 0 {
					b.Linef("%s.name = %q;", varName, name)
				}
				if posX != 0 || posY != 0 {
					b.Linef("%s.x = %s; %s.y = %s;", varName, codegen.FmtFloat(posX), varName, codegen.FmtFloat(posY))
				}
				if parent != "" {
					emitParentAppend(b, parent, varName)
				} else if *useLast {
					b.Linef("const _par%d = _results[_results.length - 1];", i)
					b.Linef("if (_par%d && 'appendChild' in _par%d) _par%d.appendChild(%s);", i, i, i, varName)
				}
				ids = append(ids, varName+".id")
				b.Blank()
			}
			b.ReturnIDs(ids...)
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&parent, "parent", "", "Parent node ID (or _results[N] in compose)")
	cmd.Flags().StringVar(&name, "name", "", "Rename the first instance")
	cmd.Flags().Float64Var(&posX, "x", 0, "X position")
	cmd.Flags().Float64Var(&posY, "y", 0, "Y position")
	cmd.Flags().BoolVar(useLast, "last", false, "Append to last _results[] node (compose chaining)")
	return cmd
}

// ---------------------------------------------------------------------------
// ds library import-set — composable, Plugin API
// ---------------------------------------------------------------------------

func newDSLibraryImportSetCmd() *cobra.Command {
	var (
		parent  string
		name    string
		variant string
		posX    float64
		posY    float64
	)
	useLast := new(bool)
	cmd := &cobra.Command{
		Use:   "import-set <key>",
		Short: "Import a component set and create an instance with a specific variant",
		Example: `  figma-kit ds library import-set x9y8z7 --variant "Size=Large,State=Default"
  figma-kit ds library import-set x9y8z7 --variant "Size=Small" --parent 123:456`,
		Annotations: map[string]string{"composable": "true"},
		Args:        cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			b := newBuilder()

			b.ImportComponentSet(key, variant, "inst")
			if name != "" {
				b.Linef("inst.name = %q;", name)
			}
			if posX != 0 || posY != 0 {
				b.Linef("inst.x = %s; inst.y = %s;", codegen.FmtFloat(posX), codegen.FmtFloat(posY))
			}
			if parent != "" {
				emitParentAppend(b, parent, "inst")
			} else if *useLast {
				b.Line("const _par = _results[_results.length - 1];")
				b.Line("if (_par && 'appendChild' in _par) _par.appendChild(inst);")
			}
			b.ReturnIDs("inst.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&parent, "parent", "", "Parent node ID (or _results[N] in compose)")
	cmd.Flags().StringVar(&name, "name", "", "Rename the instance")
	cmd.Flags().StringVar(&variant, "variant", "", "Variant properties: Key=Value,Key=Value")
	cmd.Flags().Float64Var(&posX, "x", 0, "X position")
	cmd.Flags().Float64Var(&posY, "y", 0, "Y position")
	cmd.Flags().BoolVar(useLast, "last", false, "Append to last _results[] node (compose chaining)")
	return cmd
}

// ---------------------------------------------------------------------------
// ds library import-style — composable, Plugin API
// ---------------------------------------------------------------------------

func newDSLibraryImportStyleCmd() *cobra.Command {
	var apply string
	cmd := &cobra.Command{
		Use:   "import-style <key>",
		Short: "Import a published style and optionally apply it to a node",
		Example: `  figma-kit ds library import-style s1t2y3l4e5
  figma-kit ds library import-style s1t2y3l4e5 --apply 123:456`,
		Annotations: map[string]string{"composable": "true"},
		Args:        cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			b := newBuilder()

			b.ImportStyle(key, "style")
			if apply != "" {
				if strings.HasPrefix(apply, "_results[") {
					b.Linef("const _target = %s;", apply)
				} else {
					b.Linef("const _target = await figma.getNodeByIdAsync(%q);", apply)
				}
				b.Line("if (_target) {")
				b.Line("  const styleType = style.type;")
				b.Line("  if (styleType === 'PAINT') {")
				b.Line("    _target.fillStyleId = style.id;")
				b.Line("  } else if (styleType === 'TEXT') {")
				b.Line("    _target.textStyleId = style.id;")
				b.Line("  } else if (styleType === 'EFFECT') {")
				b.Line("    _target.effectStyleId = style.id;")
				b.Line("  } else if (styleType === 'GRID') {")
				b.Line("    _target.gridStyleId = style.id;")
				b.Line("  }")
				b.Line("}")
			}
			b.ReturnExpr("style.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&apply, "apply", "", "Node ID to apply the style to (or _results[N])")
	return cmd
}

// ---------------------------------------------------------------------------
// ds library variables — Plugin API (figma.teamLibrary)
// ---------------------------------------------------------------------------

func newDSLibraryVariablesCmd() *cobra.Command {
	var collection string
	var asJSON bool
	cmd := &cobra.Command{
		Use:   "variables",
		Short: "List available library variable collections and their variables",
		Long: `Generates Plugin API JavaScript to list variable collections from enabled
team libraries using figma.teamLibrary APIs.

The library must be enabled in the current Figma file via the UI.`,
		Example: `  figma-kit ds library variables
  figma-kit ds library variables --collection <collectionKey>`,
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			b := newBuilder()

			if collection != "" {
				b.Comment("Fetch variables from a specific library collection")
				b.Linef("const vars = await figma.teamLibrary.getVariablesInLibraryCollectionAsync(%q);", collection)
				b.Line("const result = vars.map(v => ({ name: v.name, resolvedType: v.resolvedType, key: v.key }));")
				b.ReturnExpr("JSON.stringify(result, null, 2)")
			} else {
				b.Comment("List all available library variable collections")
				b.Line("const collections = await figma.teamLibrary.getAvailableLibraryVariableCollectionsAsync();")
				b.Line("const result = collections.map(c => ({ name: c.name, key: c.key, libraryName: c.libraryName }));")
				b.ReturnExpr("JSON.stringify(result, null, 2)")
			}
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&collection, "collection", "", "Collection key to list variables from")
	cmd.Flags().BoolVar(&asJSON, "json", false, "Output raw JSON (default for this command)")
	return cmd
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func writeJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func truncateDesc(s string, max int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) > max {
		return s[:max-3] + "..."
	}
	return s
}
