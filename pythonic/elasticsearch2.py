from elasticsearch_dsl import connections


class ES:
    def __init__(self) -> None:
        self._es = None

    def get_search(self, index):
        from elasticsearch_dsl import Search
        return Search(using=self._es, index=index)

    def connect(self, host, timeout=10):
        self._es = connections.create_connection(hosts=host, timeout=timeout)


def search_data(es, host, index_name, query_body):
    es.connect(host)
    es_search = es.get_search(index_name)
    es_query = es_search.from_dict(query_body)[:10]
    # query = query.filter(id="25012203503110100", amount__in=[2000000])
    for i in es_query:
        print(i.to_dict())


if __name__ == "__main__":
    es = ES()
    # host = "http://wl-global-elasticsearch01-1.offline-ops.net:9200"
    host = "http://localhost:9200"
    index = "mysql-cdc.test.users"
    query_body = {
        "query": {
            "bool": {
                "must": [
                    {"term": {"id": 1000}},
                ]
            }
        }
    }
    search_data(es, host, index, query_body)
