### prometheus-scrape-config
> [flat-prometheus-scrape.yaml](https://github.com/Ayoralol/doc-gen/tree/main/configs/prometheus/flat-prometheus-scrape.yaml)

- job_name:
    ***blackbox***
- scrape_interval:
    ***30s***
- static_configs:
    - targets:
        - [***https://example.com***](https://example.com)
        - [***https://google.com***](https://google.com)
