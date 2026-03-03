// masterfabric_go — code-generation CLI for the MasterFabric platform.
//
// Usage:
//
//	masterfabric_go generate dart   — generate sdk/dart_go_api Dart package
//	masterfabric_go generate swift  — generate sdk/swift_go_api Swift package
package main

import (
	"fmt"
	"os"

	"github.com/masterfabric/masterfabric_go_basic/internal/codegen/dart"
	"github.com/masterfabric/masterfabric_go_basic/internal/codegen/swift"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "masterfabric_go",
		Short: "MasterFabric code-generation CLI",
		Long: `masterfabric_go is the code-generation tool for the MasterFabric Go backend.

It reads the GraphQL schema files and generates typed client SDK packages
for supported target platforms.`,
	}

	generate := &cobra.Command{
		Use:   "generate",
		Short: "Generate SDK packages from the GraphQL schema",
	}

	// ── Phase 1: Dart ─────────────────────────────────────────────────────────
	var (
		schemaDir string
		outputDir string
	)

	dartCmd := &cobra.Command{
		Use:   "dart",
		Short: "Generate a Dart/Flutter GraphQL client package (sdk/dart_go_api)",
		Long: `Reads all *.graphqls files from the schema directory and emits a complete
Dart package under sdk/dart_go_api ready to be dropped into any Flutter project.

Generated package includes:
  • pubspec.yaml
  • lib/src/models/       — Dart model classes with fromJson / toJson
  • lib/src/queries/      — gql() DocumentNode constants
  • lib/src/client/       — typed GraphQL API client (graphql package)
  • lib/dart_go_api.dart  — barrel export`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("masterfabric_go: generating Dart package...")
			if err := dart.Generate(schemaDir, outputDir); err != nil {
				return fmt.Errorf("dart generation failed: %w", err)
			}
			fmt.Printf("masterfabric_go: Dart package written to %s\n", outputDir)
			return nil
		},
	}

	dartCmd.Flags().StringVar(&schemaDir, "schema", "internal/infrastructure/graphql/schema", "Directory containing *.graphqls files")
	dartCmd.Flags().StringVar(&outputDir, "output", "sdk/dart_go_api", "Output directory for the generated Dart package")

	// ── Phase 2: Swift ────────────────────────────────────────────────────────
	var (
		swiftSchemaDir string
		swiftOutputDir string
	)

	swiftCmd := &cobra.Command{
		Use:   "swift",
		Short: "Generate a Swift/iOS GraphQL client package (sdk/swift_go_api)",
		Long: `Reads all *.graphqls files from the schema directory and emits a complete
Swift Package Manager package under sdk/swift_go_api ready to be added to any
iOS/macOS Xcode project via File > Add Package Dependencies.

Generated package includes:
  • Package.swift                   — SPM manifest (Apollo iOS dependency)
  • Sources/MasterFabricAPI/Models/ — Codable structs, enums, input types
  • Sources/MasterFabricAPI/Queries/ — GraphQL operation string constants
  • Sources/MasterFabricAPI/Client/ — typed async/await API client`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("masterfabric_go: generating Swift package...")
			if err := swift.Generate(swiftSchemaDir, swiftOutputDir); err != nil {
				return fmt.Errorf("swift generation failed: %w", err)
			}
			fmt.Printf("masterfabric_go: Swift package written to %s\n", swiftOutputDir)
			return nil
		},
	}

	swiftCmd.Flags().StringVar(&swiftSchemaDir, "schema", "internal/infrastructure/graphql/schema", "Directory containing *.graphqls files")
	swiftCmd.Flags().StringVar(&swiftOutputDir, "output", "sdk/swift_go_api", "Output directory for the generated Swift package")

	generate.AddCommand(dartCmd, swiftCmd)
	root.AddCommand(generate)

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
