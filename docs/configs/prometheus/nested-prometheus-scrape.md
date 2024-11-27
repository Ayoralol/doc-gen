### prometheus-scrape-config

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
            - ***https://www.example.com***
            - ***https://www.google.com***
            - ***https://www.facebook.com***