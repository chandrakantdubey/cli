# CLI Mastery

A practical, builder-first guide to the Linux command line and the tools engineers use every day.

The goal is not to memorize commands. It is to understand how the shell, filesystem, processes, networking, permissions, scripting, Git, containers, and production operations fit together so you can work confidently on real systems.

## What you'll learn

### Foundations
- terminal and shell mental models
- Linux filesystem and paths
- users, groups, permissions, ownership
- files, streams, redirection, and pipelines

### Everyday engineering
- text processing and search
- environment variables and configuration
- packages and software installation
- networking and SSH
- processes and signals
- shell scripting
- Git workflows

### Deeper system skills
- archiving and advanced text processing
- Node.js and Python CLI tooling
- shell expansion and customization
- storage and filesystems
- systemd and services
- Linux troubleshooting and security

### Production and cloud-native work
- containers
- observability
- production operations
- automation
- Linux for cloud-native engineers
- shell portability
- Linux operations capstone

## How to use this site

Start with **CLI & Linux Intro**, then move through the topics in order or jump directly to the area you need.

Every page is designed to be read, tried in a terminal, and connected to real engineering work.

## Local development

```bash
npm install
npm run dev
```

Build and lint:

```bash
npm run build
npm run lint
```

## Project structure

```text
src/
├── data/modules/       # learning content
├── data/moduleRegistry.ts
├── layouts/            # navigation and page shell
├── pages/              # reader and loading states
└── App.tsx             # routing
```

## Learning principle

The command line is not a collection of commands. It is a way to inspect state, compose programs, automate work, and operate systems.
