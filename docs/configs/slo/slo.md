### sloth-slo
> [slo.yaml](https://github.com/Ayoralol/doc-gen/tree/main/configs/slo/slo.yaml)

- version:
    ***prometheus/v1***
- service:
    ***tester***
- labels:
    - ci:
        ***CI000***
    - owner:
        ***oner***
- slos:
    - name:
        ***prober***
    - objective:
        ***99***
    - description:
        ***SLO for uptime***
    - sli:
        - events:
            - error_query:
                ```sql
                sum_over_time(probe_success{ci="CI000"} == 0[{{.window}}])
                ```
            - total_query:
                ```sql
                sum_over_time(probe_success{ci="CI000"}[{{.window}}])
                ```
    - alerting:
        - page_alert:
            - disabled:
                ***true***
        - ticker_alert:
            - disabled:
                ***true***
