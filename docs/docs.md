# Doc-Gen

> ### [Repository](https://github.com/Ayoralol/doc-gen/)

> Documentation generator from .yaml files, written in Go

## prometheus-scrape-config

#### flat-prometheus-scrape.md

- job_name:
    ***blackbox***
- scrape_interval:
    ***30s***
- static_configs:
    - targets:
        - [***https://example.com***](https://example.com)
        - [***https://google.com***](https://google.com)


#### nested-prometheus-scrape.md

- scrape_configs:
    - job_name:
        ***prometheus-self***
    - static_configs:
        - targets:
            - ***localhost:9090***

    - job_name:
        ***blackbox-probe***
    - static_configs:
        - targets:
            - [***https://www.example.com***](https://www.example.com)
            - [***https://www.google.com***](https://www.google.com)
            - [***https://www.facebook.com***](https://www.facebook.com)


#### prometheus-copy.md

- job_name:
    ***blackbox***
- scrape_interval:
    ***30s***
- static_configs:
    - targets:
        - [***https://example.com***](https://example.com)
        - [***https://google.com***](https://google.com)


## sloth-slo

#### slo-copy.md

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


#### slo.md

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


