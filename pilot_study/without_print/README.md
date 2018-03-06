# Data Analysis

This is a Data analysis of access time, response time, de-/serialization time, and data size for the WebSocket and HTTP implementations. The data is collected from 500 iterations preceded by 500 warmup iterations on 13 maps where the amount of data on each map varies. The implementations does not print anything on the client once the response from the server is deserialised which give some interesting, but not surprising, results when measuring deserialization time on FlatBuffers.

---

## ~ Access Time ~

### Single Factor Anova Table 

![Access Time Anova Table](http://wwwlab.iit.his.se/d15johro/examensarbete/pilot_study/without_print/img/access_time_anova_table.PNG)

We use a confidence level of 95%. Since the P-value is not less than 0.05. There is no bias. We might need to improve the tests by running more iterations and and using vary the amount of data more.

---