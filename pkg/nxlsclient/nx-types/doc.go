/*
Package nxtypes provides Go type definitions for Nx workspace data structures.

This package contains Go representations of Nx-specific data structures used by the
Nx Language Server Protocol. These types are used for marshaling and unmarshaling
JSON data exchanged with the Nx LSP server.

# Overview

The nxtypes package includes:

  - Project configuration models
  - Project graph structures
  - Generator schemas
  - Workspace configuration
  - Package manager definitions
  - Cloud integration types

# Usage

The types in this package are primarily used as return types for commands in the commands package:

	// Get project graph
	projectGraph, err := client.Commander.SendCreateProjectGraphRequest(ctx, params)
	if err != nil {
		// Handle error
	}

	// Access project graph nodes
	for name, node := range projectGraph.Nodes {
		fmt.Printf("Project: %s, Type: %s\n", name, node.Type)
	}

	// Access dependencies
	for source, deps := range projectGraph.Dependencies {
		fmt.Printf("Project %s has %d dependencies\n", source, len(deps))
	}

# Key Types

## Project Graph

The ProjectGraph type represents the dependency structure of an Nx workspace:

	type ProjectGraph struct {
		Nodes         map[string]ProjectGraphProjectNode  `json:"nodes"`
		ExternalNodes map[string]ProjectGraphExternalNode `json:"externalNodes,omitempty"`
		Dependencies  map[string][]ProjectGraphDependency `json:"dependencies"`
		Version       *string                             `json:"version,omitempty"`
	}

## Project Configuration

ProjectConfiguration contains the configuration for a single project:

	type ProjectConfiguration struct {
		Root        string              `json:"root"`
		SourceRoot  string              `json:"sourceRoot,omitempty"`
		ProjectType string              `json:"projectType,omitempty"`
		Targets     map[string]Target   `json:"targets,omitempty"`
		Tags        []string            `json:"tags,omitempty"`
		Implicitly  bool                `json:"implicitly,omitempty"`
		Name        string              `json:"name,omitempty"`
	}

## Generator Schemas

GeneratorSchema represents an Nx generator schema:

	type GeneratorSchema struct {
		Name        string                     `json:"name"`
		Factory     string                     `json:"factory"`
		Schema      map[string]SchemaProperty  `json:"schema"`
		Description string                     `json:"description"`
		Hidden      bool                       `json:"hidden"`
		// Additional fields omitted for brevity
	}
*/
package nxtypes
