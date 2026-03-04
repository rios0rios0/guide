package main

// RuleGroup defines a mapping from source markdown files to a single rule output.
type RuleGroup struct {
	Name        string   // output filename (without extension)
	Description string   // human-readable description for Cursor frontmatter
	Sources     []string // relative paths from repo root, in concatenation order
	Globs       string   // file glob for language targeting; empty for always-apply
}

// ruleGroups returns the complete list of rule group definitions.
// Language sub-pages follow the order defined in Code-Style/Language-Guide-Template.md:
// index → Conventions → Formatting → Type System → Logging → Testing → Project Structure.
func ruleGroups() []RuleGroup {
	return []RuleGroup{
		// Language-specific rules
		{
			Name:        "golang",
			Description: "Go language coding standards and conventions",
			Sources: []string{
				"Code-Style/GoLang.md",
				"Code-Style/GoLang/GoLang-Conventions.md",
				"Code-Style/GoLang/GoLang-Formatting-and-Linting.md",
				"Code-Style/GoLang/GoLang-Type-System.md",
				"Code-Style/GoLang/GoLang-Logging.md",
				"Code-Style/GoLang/GoLang-Testing.md",
				"Code-Style/GoLang/GoLang-Project-Structure.md",
			},
			Globs: "**/*.go",
		},
		{
			Name:        "python",
			Description: "Python language coding standards and conventions",
			Sources: []string{
				"Code-Style/Python.md",
				"Code-Style/Python/Python-Conventions.md",
				"Code-Style/Python/Python-Formatting-and-Linting.md",
				"Code-Style/Python/Python-Type-System.md",
				"Code-Style/Python/Python-Logging.md",
				"Code-Style/Python/Python-Testing.md",
				"Code-Style/Python/Python-Project-Structure.md",
			},
			Globs: "**/*.py",
		},
		{
			Name:        "java",
			Description: "Java language coding standards and conventions",
			Sources: []string{
				"Code-Style/Java.md",
				"Code-Style/Java/Java-Conventions.md",
				"Code-Style/Java/Java-Formatting-and-Linting.md",
				"Code-Style/Java/Java-Type-System.md",
				"Code-Style/Java/Java-Logging.md",
				"Code-Style/Java/Java-Testing.md",
				"Code-Style/Java/Java-Project-Structure.md",
			},
			Globs: "**/*.java",
		},
		{
			Name:        "javascript",
			Description: "JavaScript and TypeScript coding standards and conventions",
			Sources: []string{
				"Code-Style/JavaScript.md",
				"Code-Style/JavaScript/JavaScript-Testing.md",
			},
			Globs: "**/*.{js,jsx,ts,tsx}",
		},
		{
			Name:        "yaml",
			Description: "YAML coding standards and conventions",
			Sources: []string{
				"Code-Style/YAML.md",
			},
			Globs: "**/*.{yml,yaml}",
		},
		// Cross-cutting concerns
		{
			Name:        "code-style",
			Description: "General code style and naming conventions",
			Sources: []string{
				"Code-Style.md",
			},
		},
		{
			Name:        "git-flow",
			Description: "Git workflow, branching, and commit conventions",
			Sources: []string{
				"Life-Cycle/Git-Flow.md",
				"Life-Cycle/Git-Flow/Merge-Guide.md",
			},
		},
		{
			Name:        "testing",
			Description: "Testing standards and patterns",
			Sources: []string{
				"Life-Cycle/Tests.md",
			},
		},
		{
			Name:        "architecture",
			Description: "Architecture principles and design patterns",
			Sources: []string{
				"Life-Cycle/Architecture.md",
				"Life-Cycle/Architecture/Backend-Design.md",
				"Life-Cycle/Architecture/Frontend-Design.md",
			},
		},
		{
			Name:        "security",
			Description: "Security practices and SAST pipeline",
			Sources: []string{
				"Life-Cycle/Security.md",
			},
		},
		{
			Name:        "ci-cd",
			Description: "CI/CD pipeline standards",
			Sources: []string{
				"Life-Cycle/CI-&-CD.md",
			},
		},
		{
			Name:        "documentation",
			Description: "Documentation and change control standards",
			Sources: []string{
				"Life-Cycle/Documentation-&-Change-Control.md",
				"Life-Cycle/Documentation-&-Change-Control/README-Template.md",
				"Life-Cycle/Documentation-&-Change-Control/CONTRIBUTING-Template.md",
			},
		},
		{
			Name:        "design-patterns",
			Description: "Design patterns and coding techniques",
			Sources: []string{
				"Cookbooks/Mapper-Design-Pattern.md",
				"Cookbooks/Forking-Technique.md",
			},
		},
		{
			Name:        "bulk-operations",
			Description: "Bulk operations across multiple repositories",
			Sources: []string{
				"Cookbooks/Bulk-Operations.md",
			},
		},
	}
}
