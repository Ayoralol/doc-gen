### prometheus-scrape-config
> [nested-prometheus-scrape.yaml](https://github.com/Ayoralol/doc-gen/tree/main/configs/prometheus/nested-prometheus-scrape.yaml)

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
