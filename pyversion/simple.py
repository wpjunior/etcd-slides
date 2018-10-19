import etcd3

etcd = etcd3.client()

etcd.get('foo')
etcd.put('bar', 'doot')
etcd.delete('bar')
