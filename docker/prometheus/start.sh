docker run -d -p 9090:9090 \
            -v $PWD/prometheus.yml:/etc/prometheus/prometheus.yml \
            --name prometheus \
            prom/prometheus
