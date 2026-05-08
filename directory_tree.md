# Agent Directory Tree

```text
nevarix-agent/
|-- api/
|   `-- openapi.yaml
|-- cmd/
|   |-- cli/
|   `-- server/
|       |-- agent_commands.go
|       |-- constants.go
|       |-- main.go
|       |-- router.go
|       |-- runtime.go
|       |-- send_to_hub.go
|       `-- server.exe
|-- internal/
|   |-- domain/
|   |   |-- post_to_hub/
|   |   |   `-- agent_push.go
|   |   `-- prober/
|   |       |-- utils/
|   |       |   |-- agent_collect.go
|   |       |   |-- config_hosts.go
|   |       |   |-- database_operations.go
|   |       |   |-- database_schema.go
|   |       |   |-- db_checks.go
|   |       |   |-- dns_operations.go
|   |       |   |-- host_operations.go
|   |       |   |-- logging.go
|   |       |   |-- probe_operations.go
|   |       |   |-- runtime_options.go
|   |       |   |-- sqlite_db.go
|   |       |   |-- state_dns.go
|   |       |   |-- status.go
|   |       |   `-- types.go
|   |       |-- http_prober.go
|   |       |-- icmp_prober.go
|   |       |-- runtime_config.go
|   |       `-- utils.go
|   |-- http/
|   |   `-- api/
|   |       |-- handlers.go
|   |       |-- middleware.go
|   |       `-- router.go
|   `-- hub/
|       `-- client.go
|-- .gitattributes
|-- agent.log
|-- config.toml
|-- directory_tree.md
|-- go.mod
|-- go.sum
|-- ncs.state
|-- README.md
|-- server
|-- server.exe
|-- server.exe~
`-- update_directory_trees.py
```
