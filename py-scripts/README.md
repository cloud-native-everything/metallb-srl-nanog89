# pyGNMI scripts to check EVPN configurations

## SRL BGP Device Fetcher Class
SRL Device BGP Information Fetcher

This Python application connects to a list of specified routers using the gNMI protocol 
and retrieves both BGP EVPN and BGP VPN information. The information is then formatted 
into a PrettyTable for easy viewing.

The application primarily consists of the SrlDevice class, which represents a router. 
This class is initialized with the router's basic information and uses the gNMI client 
to fetch BGP EVPN and BGP VPN information.

The fetched BGP EVPN and BGP VPN information is stored as instances of the nested 
BgpEvpn and BgpVpn classes. Each of these instances encapsulates the relevant information 
for a single BGP EVPN or BGP VPN instance.

```python
    def __init__(self, router, port, model, release, username, password, skip_verify=DEFAULT_SKIP_VERIFY):
        self.router = router
        self.port = port
        self.password = password
        self.username = username
        self.skip_verify = skip_verify
        self.model = model
        self.release = release
        self.bgp_evpn = self.get_bgp_evpn_info()
        self.bgp_vpn = self.get_bgp_vpn_info()
```

Defined methods extract BGP-EVPN and BG-VPN data use a helper function like this:
```python
    def _get_gnmi_info(self, gnmi_path):
        info = []
        logging.getLogger('pygnmi').setLevel(LOGGING_LEVEL_CRITICAL)
        result = None
        try:
            with gNMIclient(target=(self.router, self.port), username=self.username, password=self.password, skip_verify=True) as gc:
                result = gc.get(path=gnmi_path)
        except Exception as e:
            print(f"Failed to connect to router or fetch data: {e}")
        finally:
            logging.getLogger('pygnmi').setLevel(LOGGING_LEVEL_WARNING)

        if result is not None:
            try:
                for notification in result['notification']:
                    for update in notification['update']:
                        if 'srl_nokia-network-instance:network-instance' in update['val']:
                            network_instances = update['val']['srl_nokia-network-instance:network-instance']
                            for network_instance in network_instances:
                                info.append(network_instance)
            except KeyError as ke:
                print(f"Failed to process data due to missing key: {ke}")
            except Exception as e:
                print(f"Unexpected error occurred while processing data: {e}")

        return info
```

## Display EVPN data from all devices per router

The application then creates a list of these SrlDevice instances based on a YAML 
configuration file ('datacenter-nodes.yml'). It generates a tables sorted by 
router name.
```bash
[root@rbc-r2-hpe4 py-scripts]# python3 display_evpn_per_router.py datacenter-nodes.yml
Table: Sorted by Router
+-----------------------+------------------+----+------------------+-----------------+------+------+------------+--------------+-------------------+-------------------+
| Router                | Network instance | ID | EVPN Admin state | VXLAN interface | EVI  | ECMP | Oper state | RD           | import-rt         | export-rt         |
+-----------------------+------------------+----+------------------+-----------------+------+------+------------+--------------+-------------------+-------------------+
| clab-dc-k8s-LEAF-DC-1 | kube-ipvrf       | 1  | enable           | vxlan1.4        | 4    | 4    | up         | 1.1.1.1:4    | target:65123:4    | target:65123:4    |
| clab-dc-k8s-LEAF-DC-1 | kube_macvrf      | 1  | enable           | vxlan1.1        | 1    | 1    | up         | 1.1.1.1:1    | target:65123:1    | target:65123:1    |
| clab-dc-k8s-LEAF-DC-1 | l2evpn1001       | 2  | enable           | vxlan2.1001     | 1001 | 1    | no state   | 1.1.1.1:1001 | target:65123:1001 | target:65123:1001 |
| clab-dc-k8s-LEAF-DC-1 | l2evpn1002       | 2  | enable           | vxlan2.1002     | 1002 | 1    | no state   | 1.1.1.1:1002 | target:65123:1002 | target:65123:1002 |
| clab-dc-k8s-LEAF-DC-1 | l2evpn1003       | 2  | enable           | vxlan2.1003     | 1003 | 1    | no state   | 1.1.1.1:1003 | target:65123:1003 | target:65123:1003 |
| clab-dc-k8s-LEAF-DC-1 | l2evpn1004       | 2  | enable           | vxlan2.1004     | 1004 | 1    | no state   | 1.1.1.1:1004 | target:65123:1004 | target:65123:1004 |
| clab-dc-k8s-LEAF-DC-1 | l2evpn1005       | 2  | enable           | vxlan2.1005     | 1005 | 1    | no state   | 1.1.1.1:1005 | target:65123:1005 | target:65123:1005 |
| clab-dc-k8s-LEAF-DC-1 | l2evpn1006       | 2  | enable           | vxlan2.1006     | 1006 | 1    | no state   | 1.1.1.1:1006 | target:65123:1006 | target:65123:1006 |
| clab-dc-k8s-LEAF-DC-1 | l3evpn           | 1  | enable           | vxlan1.2        | 2    | 4    | up         | 1.1.1.1:2    | target:65123:2    | target:65123:2    |
| clab-dc-k8s-LEAF-DC-2 | kube-ipvrf       | 1  | enable           | vxlan1.4        | 4    | 4    | up         | 1.1.1.2:4    | target:65123:4    | target:65123:4    |
| clab-dc-k8s-LEAF-DC-2 | kube_macvrf      | 1  | enable           | vxlan1.1        | 1    | 1    | up         | 1.1.1.2:1    | target:65123:1    | target:65123:1    |
| clab-dc-k8s-LEAF-DC-2 | l2evpn1001       | 2  | enable           | vxlan2.1001     | 1001 | 1    | up         | 1.1.1.2:1001 | target:65123:1001 | target:65123:1001 |
| clab-dc-k8s-LEAF-DC-2 | l2evpn1002       | 2  | enable           | vxlan2.1002     | 1002 | 1    | up         | 1.1.1.2:1002 | target:65123:1002 | target:65123:1002 |
| clab-dc-k8s-LEAF-DC-2 | l2evpn1003       | 2  | enable           | vxlan2.1003     | 1003 | 1    | no state   | 1.1.1.2:1003 | target:65123:1003 | target:65123:1003 |
| clab-dc-k8s-LEAF-DC-2 | l2evpn1004       | 2  | enable           | vxlan2.1004     | 1004 | 1    | no state   | 1.1.1.2:1004 | target:65123:1004 | target:65123:1004 |
| clab-dc-k8s-LEAF-DC-2 | l2evpn1005       | 2  | enable           | vxlan2.1005     | 1005 | 1    | no state   | 1.1.1.2:1005 | target:65123:1005 | target:65123:1005 |
| clab-dc-k8s-LEAF-DC-2 | l2evpn1006       | 2  | enable           | vxlan2.1006     | 1006 | 1    | no state   | 1.1.1.2:1006 | target:65123:1006 | target:65123:1006 |
| clab-dc-k8s-LEAF-DC-2 | l3evpn           | 1  | enable           | vxlan1.2        | 2    | 4    | up         | 1.1.1.2:2    | target:65123:2    | target:65123:2    |
| clab-dc-k8s-BORDER-DC | kube-ipvrf       | 1  | enable           | vxlan1.4        | 4    | 4    | up         | 1.1.1.10:4   | target:65123:4    | target:65123:4    |
+-----------------------+------------------+----+------------------+-----------------+------+------+------------+--------------+-------------------+-------------------+
Total time: 1.52 seconds
```

## Display EVPN data from all devices per network-instance
The application then creates a list of these SrlDevice instances based on a YAML 
configuration file ('datacenter-nodes.yml'). It generates a tables sorted by 
network instance name. Also, it would highlight any EVI difference between the routers using the same instance.

```bash
[root@rbc-r2-hpe4 py-scripts]# python3 display_evpn_per_netinst.py datacenter-nodes.yml
Table: Sorted by Network Instance
+-----------------------+------------------+----+------------------+-----------------+------+------+------------+--------------+-------------------+-------------------+
|        Router         | Network instance | ID | EVPN Admin state | VXLAN interface | EVI  | ECMP | Oper state |      RD      |     import-rt     |     export-rt     |
+-----------------------+------------------+----+------------------+-----------------+------+------+------------+--------------+-------------------+-------------------+
| clab-dc-k8s-LEAF-DC-1 |    kube-ipvrf    | 1  |      enable      |    vxlan1.4     |  4   |  4   |     up     |  1.1.1.1:4   |  target:65123:4   |  target:65123:4   |
| clab-dc-k8s-LEAF-DC-2 |    kube-ipvrf    | 1  |      enable      |    vxlan1.4     |  4   |  4   |     up     |  1.1.1.2:4   |  target:65123:4   |  target:65123:4   |
| clab-dc-k8s-BORDER-DC |    kube-ipvrf    | 1  |      enable      |    vxlan1.4     |  4   |  4   |     up     |  1.1.1.10:4  |  target:65123:4   |  target:65123:4   |
| clab-dc-k8s-LEAF-DC-1 |   kube_macvrf    | 1  |      enable      |    vxlan1.1     |  1   |  1   |     up     |  1.1.1.1:1   |  target:65123:1   |  target:65123:1   |
| clab-dc-k8s-LEAF-DC-2 |   kube_macvrf    | 1  |      enable      |    vxlan1.1     |  1   |  1   |     up     |  1.1.1.2:1   |  target:65123:1   |  target:65123:1   |
| clab-dc-k8s-LEAF-DC-1 |    l2evpn1001    | 2  |      enable      |   vxlan2.1001   | 1001 |  1   |  no state  | 1.1.1.1:1001 | target:65123:1001 | target:65123:1001 |
| clab-dc-k8s-LEAF-DC-2 |    l2evpn1001    | 2  |      enable      |   vxlan2.1001   | 1001 |  1   |     up     | 1.1.1.2:1001 | target:65123:1001 | target:65123:1001 |
| clab-dc-k8s-LEAF-DC-1 |    l2evpn1002    | 2  |      enable      |   vxlan2.1002   | 1002 |  1   |  no state  | 1.1.1.1:1002 | target:65123:1002 | target:65123:1002 |
| clab-dc-k8s-LEAF-DC-2 |    l2evpn1002    | 2  |      enable      |   vxlan2.1002   | 1002 |  1   |     up     | 1.1.1.2:1002 | target:65123:1002 | target:65123:1002 |
| clab-dc-k8s-LEAF-DC-1 |    l2evpn1003    | 2  |      enable      |   vxlan2.1003   | 1003 |  1   |  no state  | 1.1.1.1:1003 | target:65123:1003 | target:65123:1003 |
| clab-dc-k8s-LEAF-DC-2 |    l2evpn1003    | 2  |      enable      |   vxlan2.1003   | 1003 |  1   |  no state  | 1.1.1.2:1003 | target:65123:1003 | target:65123:1003 |
| clab-dc-k8s-LEAF-DC-1 |    l2evpn1004    | 2  |      enable      |   vxlan2.1004   | 1004 |  1   |  no state  | 1.1.1.1:1004 | target:65123:1004 | target:65123:1004 |
| clab-dc-k8s-LEAF-DC-2 |    l2evpn1004    | 2  |      enable      |   vxlan2.1004   | 1004 |  1   |  no state  | 1.1.1.2:1004 | target:65123:1004 | target:65123:1004 |
| clab-dc-k8s-LEAF-DC-1 |    l2evpn1005    | 2  |      enable      |   vxlan2.1005   | 1005 |  1   |  no state  | 1.1.1.1:1005 | target:65123:1005 | target:65123:1005 |
| clab-dc-k8s-LEAF-DC-2 |    l2evpn1005    | 2  |      enable      |   vxlan2.1005   | 1005 |  1   |  no state  | 1.1.1.2:1005 | target:65123:1005 | target:65123:1005 |
| clab-dc-k8s-LEAF-DC-1 |    l2evpn1006    | 2  |      enable      |   vxlan2.1006   | 1006 |  1   |  no state  | 1.1.1.1:1006 | target:65123:1006 | target:65123:1006 |
| clab-dc-k8s-LEAF-DC-2 |    l2evpn1006    | 2  |      enable      |   vxlan2.1006   | 1006 |  1   |  no state  | 1.1.1.2:1006 | target:65123:1006 | target:65123:1006 |
| clab-dc-k8s-LEAF-DC-1 |      l3evpn      | 1  |      enable      |    vxlan1.2     |  2   |  4   |     up     |  1.1.1.1:2   |  target:65123:2   |  target:65123:2   |
| clab-dc-k8s-LEAF-DC-2 |      l3evpn      | 1  |      enable      |    vxlan1.2     |  2   |  4   |     up     |  1.1.1.2:2   |  target:65123:2   |  target:65123:2   |
+-----------------------+------------------+----+------------------+-----------------+------+------+------------+--------------+-------------------+-------------------+
Total time: 1.51 seconds

```