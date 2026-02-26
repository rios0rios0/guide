# Code Style

> **TL;DR:** Follow the naming conventions and file structures defined here as the baseline. Language-specific guides override these defaults where applicable. Use the standard operations vocabulary (`List`, `Get`, `Insert`, `Update`, `Delete`, and their batch variants) for naming files and classes.

## Overview

This document establishes the baseline coding standards for file creation, function naming, and variable naming across all projects. If a language-specific guide exists (see the sub-pages for Go, JavaScript, Java, and Python), its conventions take precedence over this document.

The architectural layers referenced throughout are defined in the [Backend Design](Life-Cycle/Architecture/Backend-Design.md) section. All language-specific examples follow **Hexagonal Architecture (Ports & Adapters)** with **Domain-Driven Design (DDD)** and **CQRS** -- see [Backend Design](Life-Cycle/Architecture/Backend-Design.md) for details.

## File Structure

The default file structure follows the [Backend Design](Life-Cycle/Architecture/Backend-Design.md) specification, which separates code into `domain` (contracts) and `infrastructure` (implementations) layers.

## Service Naming

When creating service files, the naming must reflect both the **entity** being operated on and whether the operation targets a **single record or a batch**.

### Operations Vocabulary

Use the following standard operation prefixes consistently across all projects:

| Operation     | Description               |
|---------------|---------------------------|
| `List`        | Retrieve multiple records |
| `Get`         | Retrieve a single record  |
| `Insert`      | Create a single record    |
| `Update`      | Modify a single record    |
| `Delete`      | Remove a single record    |
| `BatchInsert` | Create multiple records   |
| `BatchUpdate` | Modify multiple records   |
| `BatchDelete` | Remove multiple records   |
| `DeleteAll`   | Remove all records        |
