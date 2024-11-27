### sloth-slo
> [slo-copy.yaml](https://github.com/Ayoralol/doc-gen/tree/main/configs/slo/slo-copy.yaml)

- version:
    ***prometheus/v1***
- service:
    ***tester1***
- labels:
    - ci:
        ***CI0001***
    - owner:
        ***oners***
- slos:
    - name:
        ***another-prober***
    - objective:
        ***99***
    - description:
        ***SLO for uptime***
    - sli:
        - events:
            - error_query:
                ```sql
                sum_over_time(probe_success{ci="CI0001"} == 0[{{.window}}])
                ```
            - total_query:
                ```sql
                sum_over_time(probe_success{ci="CI0001"}[{{.window}}])
                ```
    - alerting:
        - page_alert:
            - disabled:
                ***true***
        - ticker_alert:
            - disabled:
                ***true***
