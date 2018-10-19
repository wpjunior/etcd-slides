import etcd3

etcd = etcd3.client()

# watch key
events_iterator, cancel = etcd.watch("pybr")
for event in events_iterator:
    print("*"* 20)
    print("key", event.key)
    print("value", event.value)
    print("revision", event.mod_revision)
    print("*"* 20)
