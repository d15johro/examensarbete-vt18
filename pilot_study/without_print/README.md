# ANOVA Analysis

This is an ANOVA analysis of access time, response time, de-/serialization time, and data size for the WebSocket and HTTP implementations. The data is collected from 500 iterations preceded by 500 warmup iterations on 13 maps where the amount of data on each map varies. The implementations does not print anything on the client once the response from the server is deserialised which give some interesting results when measuring deserialization time on FlatBuffers.

## Access Time

### Anova Table

![Access Time Anova Table](img/access_time_anova_table.png)

### Clustered Columns Diagram

![Access Time Clustered Columns Diagram](img/access_time_clustered_columns_diagram)

