<h1 align="center">Guide</h1>
<p align="center">
    <a href="https://github.com/rios0rios0/guide/releases/latest">
        <img src="https://img.shields.io/github/release/rios0rios0/guide.svg?style=for-the-badge&logo=github" alt="Latest Release"/></a>
    <a href="https://github.com/rios0rios0/guide/blob/main/LICENSE">
        <img src="https://img.shields.io/github/license/rios0rios0/guide.svg?style=for-the-badge&logo=github" alt="License"/></a>
    <a href="https://github.com/rios0rios0/guide/actions/workflows/update-wiki.yml">
        <img src="https://img.shields.io/github/actions/workflow/status/rios0rios0/guide/update-wiki.yml?branch=main&style=for-the-badge&logo=github" alt="Build Status"/></a>
    <a href="https://sonarcloud.io/summary/overall?id=rios0rios0_guide">
        <img src="https://img.shields.io/sonar/coverage/rios0rios0_guide?server=https%3A%2F%2Fsonarcloud.io&style=for-the-badge&logo=sonarqubecloud" alt="Coverage"/></a>
    <a href="https://sonarcloud.io/summary/overall?id=rios0rios0_guide">
        <img src="https://img.shields.io/sonar/quality_gate/rios0rios0_guide?server=https%3A%2F%2Fsonarcloud.io&style=for-the-badge&logo=sonarqubecloud" alt="Quality Gate"/></a>
    <a href="https://www.bestpractices.dev/projects/12025">
        <img src="https://img.shields.io/cii/level/12025?style=for-the-badge&logo=opensourceinitiative" alt="OpenSSF Best Practices"/></a>
</p>

A comprehensive software development guide covering agile culture, architecture, code style, testing, CI/CD, and team onboarding. Its objective is not to limit the development to a specific form, but to provide a common basis so that the synergy is better and starts with learning from past contexts.

## Summary

<!-- in GitHub Wiki each file name (page name) is an anchor. Make sure you have no duplicates -->
- [Home](Home.md)
- [Onboarding](Onboarding.md)
- [Agile & Culture](Agile-&-Culture.md)
  - [PDCA](Agile-&-Culture/PDCA.md)
- [Life Cycle](Life-Cycle.md)
  - [Git Flow](Life-Cycle/Git-Flow.md)
    - [Merge Guide](Life-Cycle/Git-Flow/Merge-Guide.md)
  - [Architecture](Life-Cycle/Architecture.md)
    - [Backend](Life-Cycle/Architecture/Backend-Design.md)
    - [Frontend](Life-Cycle/Architecture/Frontend-Design.md)
  - [Tests](Life-Cycle/Tests.md)
  - [CI & CD](Life-Cycle/CI-&-CD.md)
  - [Security](Life-Cycle/Security.md)
  - [Documentation & Change Control](Life-Cycle/Documentation-&-Change-Control.md)
    - [README Template](Life-Cycle/Documentation-&-Change-Control/README-Template.md)
    - [CONTRIBUTING Template](Life-Cycle/Documentation-&-Change-Control/CONTRIBUTING-Template.md)
    - [CHANGELOG Formatting](Life-Cycle/Documentation-&-Change-Control/CHANGELOG-Formatting.md)
- [Code Style](Code-Style.md)
  - [Language Guide Template](Code-Style/Language-Guide-Template.md)
  - [GoLang](Code-Style/GoLang.md)
    - [Conventions](Code-Style/GoLang/GoLang-Conventions.md)
    - [Formatting and Linting](Code-Style/GoLang/GoLang-Formatting-and-Linting.md)
    - [Type System](Code-Style/GoLang/GoLang-Type-System.md)
    - [Logging](Code-Style/GoLang/GoLang-Logging.md)
    - [Testing](Code-Style/GoLang/GoLang-Testing.md)
    - [Project Structure](Code-Style/GoLang/GoLang-Project-Structure.md)
  - [JavaScript](Code-Style/JavaScript.md)
    - [JavaScript Testing](Code-Style/JavaScript/JavaScript-Testing.md)
  - [YAML](Code-Style/YAML.md)
  - [Java](Code-Style/Java.md)
    - [Conventions](Code-Style/Java/Java-Conventions.md)
    - [Formatting and Linting](Code-Style/Java/Java-Formatting-and-Linting.md)
    - [Type System](Code-Style/Java/Java-Type-System.md)
    - [Logging](Code-Style/Java/Java-Logging.md)
    - [Testing](Code-Style/Java/Java-Testing.md)
    - [Project Structure](Code-Style/Java/Java-Project-Structure.md)
  - [Python](Code-Style/Python.md)
    - [Conventions](Code-Style/Python/Python-Conventions.md)
    - [Formatting and Linting](Code-Style/Python/Python-Formatting-and-Linting.md)
    - [Type System](Code-Style/Python/Python-Type-System.md)
    - [Logging](Code-Style/Python/Python-Logging.md)
    - [Testing](Code-Style/Python/Python-Testing.md)
    - [Project Structure](Code-Style/Python/Python-Project-Structure.md)
- [Cookbooks](Cookbooks.md)
  - [Tools & Setup](Cookbooks/Tools-&-Setup.md)
    - [WSL Setup](Cookbooks/Tools-&-Setup/WSL-Setup.md)
    - [Install Azure CLI](Cookbooks/Tools-&-Setup/Install-Azure-CLI.md)
    - [Azure Functions Setup](Cookbooks/Tools-&-Setup/Azure-Functions-Setup.md)
    - [Database Sync](Cookbooks/Tools-&-Setup/Database-Sync.md)
  - [Forking Technique](Cookbooks/Forking-Technique.md)
  - [Historical Repository Cleaner](Cookbooks/Historical-Repository-Cleaner.md)
  - [Mapper Design Pattern](Cookbooks/Mapper-Design-Pattern.md)
  - [Bulk Operations](Cookbooks/Bulk-Operations.md)
  - [AI-Assisted Workflows](Cookbooks/AI-Assisted-Workflows.md)

## AI Assistant Rules

This repository automatically generates rule files for AI coding assistants (Claude Code, Cursor, Codex, GitHub Copilot) from the documentation. Generated files live on the [`generated`](https://github.com/rios0rios0/guide/tree/generated) branch.

### Install with aisync (recommended)

[aisync](https://github.com/rios0rios0/aisync) syncs AI tool configurations across devices with encryption and multi-source support:

```bash
aisync init
aisync source add guide --source-repo rios0rios0/guide --branch generated
aisync pull
```

### Install as Claude Code plugin

```bash
/plugin marketplace add rios0rios0/guide
```

### Recommended External Sources

These community sources complement the guide's rules. Add any combination to your aisync config:

| Source | Repository | Description |
|--------|-----------|-------------|
| **Guide** | `rios0rios0/guide@generated` | 14 engineering standard rules, 7 agents, 8 commands, 5 skills |
| **Agents** | `wshobson/agents@main` | 112+ specialized agents (docs, TDD, CI/CD) |
| **Everything Claude Code** | `affaan-m/everything-claude-code@main` | Anthropic Hackathon winner. 28 agents, 125 skills, 60 commands |
| **Power Platform** | `microsoft/power-platform-skills@main` | Official Microsoft Power Platform plugins |
| **HVE Core** | `microsoft/hve-core@main` | 49 agents, 102 instructions, 63 prompts, 11 skills |

## References

- [Vizir's StyleGuide](https://github.com/Vizir/styleguide)
- [Awesome OSS Alternatives](https://github.com/VolhaBakanouskaya/awesome-oss-alternatives-public)

## Contributing

Contributions are welcome. See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the terms specified in the [LICENSE](LICENSE) file.
