# OTEL Habits

Pull daily checklists from Bear note app, convert to OTEL Metric format, submit to collector, export to data store, 
visualize in grafana.

![diagram](./media/diagram.png)

# Setup

You will need to set an ENV var to point at the collector:

`OTEL_EXPORTER_OTLP_METRICS_ENDPOINT=http://localhost:4317`


```bash
# Startup Clickhouse DB, OTEL collector, and Grafana
docker compose up -d

# Make a copy of the Bear App SQLlite db
make importdb

# Run App to pull data and submit to collector
make run
```



## Clickhouse

Queries

```sql
/* show database list */
SHOW DATABASES;

/* identify all tables */
SELECT
    table,
    sum(rows) AS rows,
    max(modification_time) AS latest_modification,
    formatReadableSize(sum(bytes)) AS data_size,
    formatReadableSize(sum(primary_key_bytes_in_memory)) AS primary_keys_size,
    any(engine) AS engine,
    sum(bytes) AS bytes_size
FROM clusterAllReplicas(default, system.parts)
WHERE active
GROUP BY
    database,
    table
ORDER BY bytes_size DESC;

/* See schema */
SHOW CREATE TABLE otel.otel_metrics_sum

/* Query all*/
SELECT MetricName, Attributes, StartTimeUnix, Value FROM otel.otel_metrics_sum;

/* Query by attribute */
SELECT
    $__timeInterval(StartTimeUnix) as time,
    count(Value)
FROM
    "otel"."otel_metrics_sum"
WHERE
    $__timeFilter(StartTimeUnix) AND
    (Attributes = '{\'name\':\'Workout\'}') AND
    (StartTimeUnix < now())
GROUP BY
    time
ORDER BY
    time
    LIMIT 1000;


select StartTimeUnix as time, round(Value) as value
from otel_metrics_sum
where Attributes['name']='Read Book'
GROUP by StartTimeUnix, Attributes, value
ORDER by StartTimeUnix;



SELECT
    MetricName,
    Attributes,
    StartTimeUnix,
    Value
FROM otel.otel_metrics_sum
WHERE Attributes = '{\'name\':\'Read Book\'}'
ORDER BY StartTimeUnix DESC
    LIMIT 100
```