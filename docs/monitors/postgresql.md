<!--- GENERATED BY gomplate from scripts/docs/monitor-page.md.tmpl --->

# postgresql

Monitor Type: `postgresql` ([Source](https://github.com/signalfx/signalfx-agent/tree/master/internal/monitors/postgresql))

**Accepts Endpoints**: **Yes**

**Multiple Instances Allowed**: Yes

## Overview

This monitor pulls metrics from all PostgreSQL databases from a specific
Postgres server instance.  It pulls basic information that is applicable
to any database.  It gathers these metrics via SQL queries.

<!--- OVERVIEW --->
## Metrics about Queries

In order to get metrics about query execution time, you must enable the
`pg_stat_statements` extension.  This extension must be specified in the
`shared_preload_libraries` config option in the main PostgreSQL
configuration at server start up.  Then the extension must be enabled for
each database by running `CREATE EXTENSION IF NOT EXISTS
pg_stat_statements;` on each database.

Note that in order to get consistent and accurate query execution time
metrics, you must set the [pg_stat_statements.max config
option](https://www.postgresql.org/docs/9.3/pgstatstatements.html#AEN160631)
to larger than the number of distinct queries on the server.

Here is a [sample configuration of Postgres to enable statement tracking](https://www.postgresql.org/docs/9.3/pgstatstatements.html#AEN160631).

Tested with PostgreSQL 9.2+.

If you want to collect additional metrics about PostgreSQL, use the [sql monitor](./sql.md).

<!--- SETUP --->
## Example Configuration

This example uses the [Vault remote config
source](https://github.com/signalfx/signalfx-agent/blob/master/docs/remote-config.md#nested-values-vault-only)
to connect to PostgreSQL using the `params` map that allows you to pull
out the username and password individually from Vault and interpolate
them into the `connectionString` config option.

```yaml
monitors:
 - type: postgresql
   connectionString: 'sslmode=disable user={{.username}} password={{.password}}'
   params: &psqlParams
     username: {"#from": "vault:secret/my-database[username]"}
     password: {"#from": "vault:secret/my-database[password]"}
   discoveryRule: 'container_image =~ "postgres" && port == 5432'

 # This monitor will monitor additional queries from PostgreSQL using the
 # provided SQL queries.
 - type: sql
   dbDriver: postgres
   connectionString: 'sslmode=disable user={{.username}} password={{.password}}'
   # This is a YAML reference to avoid duplicating the above config.
   params: *psqlParams
   queries:
     - query: 'SELECT COUNT(*) as count, country, status FROM customers GROUP BY country, status;'
       metrics:
         - metricName: "customers"
           valueColumn: "count"
           dimensionColumns: ["country", "status"]
```


## Configuration

To activate this monitor in the Smart Agent, add the following to your
agent config:

```
monitors:  # All monitor config goes under this key
 - type: postgresql
   ...  # Additional config
```

**For a list of monitor options that are common to all monitors, see [Common
Configuration](../monitor-config.md#common-configuration).**


| Config option | Required | Type | Description |
| --- | --- | --- | --- |
| `host` | no | `string` |  |
| `port` | no | `integer` |  (**default:** `0`) |
| `connectionString` | no | `string` | See https://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters. |
| `params` | no | `map of strings` | Parameters to the connection string that can be templated into the connection string with the syntax `{{.key}}`. |
| `databases` | no | `list of strings` | List of databases to send database-specific metrics about.  If omitted, metrics about all databases will be sent.  This is an [overridable set](https://docs.signalfx.com/en/latest/integrations/agent/filtering.html#overridable-filters). (**default:** `[*]`) |
| `databasePollIntervalSeconds` | no | `integer` | How frequently to poll for new/deleted databases in the DB server. Defaults to the same as `intervalSeconds` if not set. (**default:** `0`) |


## Metrics

These are the metrics available for this monitor.
Metrics that are categorized as
[container/host](https://docs.signalfx.com/en/latest/admin-guide/usage.html#about-custom-bundled-and-high-resolution-metrics)
(*default*) are ***in bold and italics*** in the list below.


 - ***`postgres_block_hit_ratio`*** (*gauge*)<br>    The proportion (between 0 and 1, inclusive) of block reads that used the cache and did not have to go to the disk.  Is sent for `table`, `index`, and the `database` as a whole.
 - ***`postgres_database_size`*** (*gauge*)<br>    Size in bytes of the database on disk
 - ***`postgres_deadlocks`*** (*cumulative*)<br>    Total number of deadlocks detected by the system
 - ***`postgres_index_scans`*** (*cumulative*)<br>    Total number of index scans on the `table`.
 - ***`postgres_live_rows`*** (*gauge*)<br>    Number of rows live (not deleted) in the `table`.
 - ***`postgres_query_count`*** (*cumulative*)<br>    Total number of queries executed on the `database`, broken down by `user`.  Note that the accuracy of this metric depends on the PostgreSQL [pg_stat_statements.max config option](https://www.postgresql.org/docs/9.3/pgstatstatements.html#AEN160631) being large enough to hold all queries.

 - ***`postgres_query_time`*** (*cumulative*)<br>    Total time taken to execute queries on the `database`, broken down by `user`.
 - ***`postgres_rows_deleted`*** (*cumulative*)<br>    Number of rows deleted from the `table`.
 - ***`postgres_rows_inserted`*** (*cumulative*)<br>    Number of rows inserted into the `table`.
 - ***`postgres_rows_updated`*** (*cumulative*)<br>    Number of rows updated in the `table`.
 - ***`postgres_sequential_scans`*** (*cumulative*)<br>    Total number of sequential scans on the `table`.
 - ***`postgres_sessions`*** (*gauge*)<br>    Number of sessions currently on the server instance.  The `state` dimension will specify which which type of session (see `state` row of [pg_stat_activity](https://www.postgresql.org/docs/9.2/monitoring-stats.html#PG-STAT-ACTIVITY-VIEW)).

 - ***`postgres_table_size`*** (*gauge*)<br>    The size in bytes of the `table` on disk.

### Non-default metrics (version 4.7.0+)

**The following information applies to the agent version 4.7.0+ that has
`enableBuiltInFiltering: true` set on the top level of the agent config.**

To emit metrics that are not _default_, you can add those metrics in the
generic monitor-level `extraMetrics` config option.  Metrics that are derived
from specific configuration options that do not appear in the above list of
metrics do not need to be added to `extraMetrics`.

To see a list of metrics that will be emitted you can run `agent-status
monitors` after configuring this monitor in a running agent instance.

### Legacy non-default metrics (version < 4.7.0)

**The following information only applies to agent version older than 4.7.0. If
you have a newer agent and have set `enableBuiltInFiltering: true` at the top
level of your agent config, see the section above. See upgrade instructions in
[Old-style whitelist filtering](../legacy-filtering.md#old-style-whitelist-filtering).**

If you have a reference to the `whitelist.json` in your agent's top-level
`metricsToExclude` config option, and you want to emit metrics that are not in
that whitelist, then you need to add an item to the top-level
`metricsToInclude` config option to override that whitelist (see [Inclusion
filtering](../legacy-filtering.md#inclusion-filtering).  Or you can just
copy the whitelist.json, modify it, and reference that in `metricsToExclude`.

## Dimensions

The following dimensions may occur on metrics emitted by this monitor.  Some
dimensions may be specific to certain metrics.

| Name | Description |
| ---  | ---         |
| `database` | The name of the database within a PostgreSQL server to which the metric pertains. |
| `index` | For index metrics, the name of the index |
| `schemaname` | The name of the schema within which the object being monitored resides (e.g. `public`). |
| `table` | The name of the table to which the metric pertains. |
| `tablespace` | For table metrics, the tablespace in which the table belongs, if not null. |
| `type` | Whether the object (table, index, function, etc.) belongs to the `system` or `user`. |
| `user` | For query metrics, the user name of the user that executed the queries. |



