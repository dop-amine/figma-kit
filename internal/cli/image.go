package cli

import (
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dop-amine/figma-kit/internal/codegen"
)

const maxInlineBytes = 33000 // ~44K base64 chars, leaves room for JS within 50K limit

func newImageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "image",
		Short: "Place local images or URLs into Figma",
		Long: `Upload images to Figma from local files or URLs.

Local files are base64-encoded and embedded directly in the generated
JavaScript — no server, no public URL needed. Files up to ~33 KB work
inline. For larger files, use a URL or run 'figma-kit image serve' to
start a local HTTP server and fetch from it.`,
		Example: `  # "Add our logo to the Figma file"
  figma-kit image place ./assets/logo.png --width 200 --height 60

  # "Use this hero image from Unsplash"
  figma-kit image place https://images.unsplash.com/photo-xxx --width 1440 --height 900

  # "Serve my local assets so I can use them in Figma"
  figma-kit image serve ./assets`,
	}

	cmd.AddCommand(newImagePlaceCmd())
	cmd.AddCommand(newImageFillCmd())
	cmd.AddCommand(newImageServeCmd())
	return cmd
}

func newImagePlaceCmd() *cobra.Command {
	var (
		width     int
		height    int
		name      string
		scaleMode string
	)
	cmd := &cobra.Command{
		Use:   "place <path-or-url>",
		Short: "Place an image as a new frame in Figma",
		Long: `Create a new image frame from a local file or URL.

Local files (png, jpg, gif, webp, svg) are base64-encoded into the
generated JavaScript. The Figma Plugin API decodes and renders them
directly — no server or public URL needed.

For files larger than ~33 KB, use a URL instead or run
'figma-kit image serve' to host them locally.`,
		Example: `  # From a local file
  figma-kit image place ./logo.png --name "Brand Logo" --width 200 --height 60

  # From a URL
  figma-kit image place https://example.com/hero.jpg --width 1440 --height 900

  # With custom scale mode
  figma-kit image place ./photo.jpg --scale-mode FIT --width 800 --height 600`,
		Args: cobra.ExactArgs(1),
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			src := args[0]
			mode := strings.ToUpper(scaleMode)
			switch mode {
			case "FILL", "FIT", "CROP", "TILE":
			default:
				return fmt.Errorf("invalid --scale-mode %q (use FILL, FIT, CROP, or TILE)", scaleMode)
			}

			page := resolvePage()
			b := newBuilder()
			b.PageSetup(page)

			if err := emitImageLoad(b, src); err != nil {
				return err
			}

			frameName := name
			if frameName == "" {
				frameName = inferName(src)
			}

			b.Line("const frame = figma.createFrame();")
			b.Linef("frame.name = %q;", frameName)
			b.Linef("frame.resize(%d, %d);", width, height)
			b.Line("frame.clipsContent = true;")
			b.Linef("frame.fills = [{ type: 'IMAGE', imageHash: img.hash, scaleMode: %q }];", mode)
			b.Line("pg.appendChild(frame);")
			b.ReturnIDs("frame.id")

			output(b.String())
			return nil
		},
	}
	cmd.Flags().IntVar(&width, "width", 400, "Frame width in pixels")
	cmd.Flags().IntVar(&height, "height", 300, "Frame height in pixels")
	cmd.Flags().StringVar(&name, "name", "", "Frame name (defaults to filename)")
	cmd.Flags().StringVar(&scaleMode, "scale-mode", "FILL", "Image scale mode: FILL, FIT, CROP, TILE")
	return cmd
}

func newImageFillCmd() *cobra.Command {
	var (
		nodeID    string
		scaleMode string
	)
	cmd := &cobra.Command{
		Use:   "fill <path-or-url>",
		Short: "Fill an existing Figma node with an image",
		Long: `Replace a node's fill with an image from a local file or URL.

The target node must already exist. Use 'figma-kit find' or
'figma-kit tree' to discover node IDs.`,
		Example: `  # Fill a hero section with a background image
  figma-kit image fill ./hero.jpg --node "2:5"

  # Fill from a URL with FIT mode
  figma-kit image fill https://example.com/bg.png --node "12:34" --scale-mode FIT`,
		Args: cobra.ExactArgs(1),
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if nodeID == "" {
				return fmt.Errorf("--node is required (e.g. --node \"2:5\")")
			}
			src := args[0]
			mode := strings.ToUpper(scaleMode)
			switch mode {
			case "FILL", "FIT", "CROP", "TILE":
			default:
				return fmt.Errorf("invalid --scale-mode %q (use FILL, FIT, CROP, or TILE)", scaleMode)
			}

			b := newBuilder()

			if err := emitImageLoad(b, src); err != nil {
				return err
			}

			b.Linef("const node = await figma.getNodeByIdAsync(%q);", nodeID)
			b.Line("if (!node) throw new Error('Node not found: ' + " + fmt.Sprintf("%q", nodeID) + ");")
			b.Linef("node.fills = [{ type: 'IMAGE', imageHash: img.hash, scaleMode: %q }];", mode)
			b.ReturnExpr("'Filled ' + node.name + ' with image'")

			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&nodeID, "node", "", "Target node ID (e.g. \"2:5\")")
	cmd.Flags().StringVar(&scaleMode, "scale-mode", "FILL", "Image scale mode: FILL, FIT, CROP, TILE")
	_ = cmd.MarkFlagRequired("node")
	return cmd
}

func newImageServeCmd() *cobra.Command {
	var port int
	cmd := &cobra.Command{
		Use:   "serve [directory]",
		Short: "Start a local HTTP server for image files",
		Long: `Serve a directory over HTTP so Figma can fetch images from it.

This is useful when images are too large for base64 embedding (~33 KB limit)
or when you want to serve many images at once. Start the server, then use
the printed URLs with 'figma-kit image place <url>' or 'figma-kit card image'.

Note: This works when Figma runs locally (desktop app). For browser-based
Figma, use a tunneling service or upload images to a public URL.`,
		Example: `  # Serve current directory
  figma-kit image serve .

  # Serve a specific assets folder on port 9090
  figma-kit image serve ./assets --port 9090

  # Then use the URLs in other commands:
  # figma-kit image place http://localhost:8741/logo.png --width 200`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) > 0 {
				dir = args[0]
			}

			abs, err := filepath.Abs(dir)
			if err != nil {
				return err
			}
			info, err := os.Stat(abs)
			if err != nil {
				return fmt.Errorf("directory not found: %s", abs)
			}
			if !info.IsDir() {
				return fmt.Errorf("not a directory: %s", abs)
			}

			if port == 0 {
				ln, err := net.Listen("tcp", "127.0.0.1:0")
				if err != nil {
					return fmt.Errorf("finding free port: %w", err)
				}
				port = ln.Addr().(*net.TCPAddr).Port
				_ = ln.Close()
			}

			addr := fmt.Sprintf("127.0.0.1:%d", port)
			fmt.Fprintf(os.Stderr, "Serving %s on http://%s\n\n", abs, addr)

			imageExts := map[string]bool{
				".png": true, ".jpg": true, ".jpeg": true,
				".gif": true, ".webp": true, ".svg": true, ".ico": true,
			}
			entries, _ := os.ReadDir(abs)
			found := 0
			for _, e := range entries {
				if e.IsDir() {
					continue
				}
				ext := strings.ToLower(filepath.Ext(e.Name()))
				if imageExts[ext] {
					info, _ := e.Info()
					size := ""
					if info != nil {
						size = formatBytes(info.Size())
					}
					fmt.Fprintf(os.Stderr, "  http://%s/%s  (%s)\n", addr, e.Name(), size)
					found++
				}
			}
			if found == 0 {
				fmt.Fprintln(os.Stderr, "  (no image files found in directory)")
			}
			fmt.Fprintf(os.Stderr, "\nUse with: figma-kit image place http://%s/<filename>\n", addr)
			fmt.Fprintln(os.Stderr, "Press Ctrl+C to stop.")

			return http.ListenAndServe(addr, http.FileServer(http.Dir(abs)))
		},
	}
	cmd.Flags().IntVar(&port, "port", 0, "Port to listen on (default: random free port)")
	return cmd
}

// emitImageLoad generates JS that loads an image into a variable called `img`.
// For local files it base64-encodes inline; for URLs it uses fetch.
func emitImageLoad(b *codegen.Builder, src string) error {
	if isURL(src) {
		b.Linef("const res = await fetch(%q);", src)
		b.Line("if (!res.ok) throw new Error('Image fetch failed: ' + res.status);")
		b.Line("const buf = new Uint8Array(await res.arrayBuffer());")
		b.Line("const img = figma.createImage(buf);")
		b.Blank()
		return nil
	}

	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("reading %s: %w", src, err)
	}

	if len(data) > maxInlineBytes {
		return fmt.Errorf(
			"%s is %s — too large for inline embedding (max ~33 KB)\n\n"+
				"Options:\n"+
				"  1. Use a URL:    figma-kit image place https://...\n"+
				"  2. Serve locally: figma-kit image serve %s\n"+
				"     then:         figma-kit image place http://localhost:<port>/%s",
			src, formatBytes(int64(len(data))),
			filepath.Dir(src), filepath.Base(src),
		)
	}

	encoded := base64.StdEncoding.EncodeToString(data)
	b.Linef("const b64 = %q;", encoded)
	b.Line("const raw = atob(b64);")
	b.Line("const buf = new Uint8Array(raw.length);")
	b.Line("for (let i = 0; i < raw.length; i++) buf[i] = raw.charCodeAt(i);")
	b.Line("const img = figma.createImage(buf);")
	b.Blank()
	return nil
}

func isURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

func inferName(src string) string {
	if isURL(src) {
		parts := strings.Split(src, "/")
		if len(parts) > 0 {
			return parts[len(parts)-1]
		}
		return "Image"
	}
	return strings.TrimSuffix(filepath.Base(src), filepath.Ext(src))
}

func formatBytes(b int64) string {
	switch {
	case b >= 1<<20:
		return fmt.Sprintf("%.1f MB", float64(b)/float64(1<<20))
	case b >= 1<<10:
		return fmt.Sprintf("%.1f KB", float64(b)/float64(1<<10))
	default:
		return fmt.Sprintf("%d B", b)
	}
}
