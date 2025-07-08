# SwitchConfigSim

Simple network switch configuration simulator similar to NVIDIA's NVUE.

## What it does

Allows for basic switch management with both CLI and REST API interfaces. Built with Go, shell scripts, and includes OpenAPI documentation.

## Building

```bash
make build
```

## Running

CLI tool:
```bash
./switchctl show
./switchctl set hostname new-name
./switchctl set interface eth0 up
```

REST API server:
```bash
./api-server
```

API docs available at http://localhost:8080/docs

## What's included

- CLI tool (`cli/`) - command line interface
- REST API server (`api/`) - HTTP endpoints  
- Shell scripts (`shell/`) - backend automation
- OpenAPI spec (`openapi.yaml`) - API documentation
- Makefile - build automation

Both CLI and API call the same shell scripts for consistency.

## Requirements

- Go 1.19+
- Unix/Linux environment
