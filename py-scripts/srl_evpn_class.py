"""
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

The application then creates a list of these SrlDevice instances based on a YAML 
configuration file ('datacenter-nodes.yml'). It generates two tables: one sorted by 
router name and another sorted by network instance.

Author: Mauricio Rojas
Last Updated: June 2023
"""

from pygnmi.client import gNMIclient
import logging
from itertools import groupby

# Define logging levels as constants
LOGGING_LEVEL_CRITICAL = logging.CRITICAL
LOGGING_LEVEL_WARNING = logging.WARNING

# Define magic numbers as constants
DEFAULT_SKIP_VERIFY = True



class SrlDevice:
    """
    A class representing a router device.

    Attributes:
    router (str): The router's hostname or IP address.
    port (int): The port number used for the gNMI connection.
    model (str): The router's model.
    release (str): The router's software release version.
    username (str): Username for the gNMI connection.
    password (str): Password for the gNMI connection.
    skip_verify (bool): Whether to skip SSL verification.
    """    
    def __init__(self, router, port, model, release, username, password, skip_verify=DEFAULT_SKIP_VERIFY):
        self.router = router
        self.port = port
        self.password = password
        self.username = username
        self.skip_verify = skip_verify
        self.model = model
        self.release = release
        self.bgp_evpn = self._get_bgp_evpn_info()
        self.bgp_vpn = self._get_bgp_vpn_info()

    class BgpEvpn:
        def __init__(self, network_instance, id, admin_state, vxlan_interface, evi, ecmp, oper_state):
            self.network_instance = network_instance
            self.id = id
            self.admin_state = admin_state
            self.vxlan_interface = vxlan_interface
            self.evi = evi
            self.ecmp = ecmp
            self.oper_state = oper_state

    class BgpVpn:
        def __init__(self, network_instance, id, rd, export_rt, import_rt):
            self.network_instance = network_instance
            self.id = id
            self.rd = rd
            self.export_rt = export_rt
            self.import_rt = import_rt               

    def _get_gnmi_info(self, gnmi_path):
        """
        Fetches gNMI information from the specified path.

        Args:
        gnmi_path (list): The gNMI path to fetch information from.

        Returns:
        list: A list of network instances.
        """ 
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


    def _get_bgp_evpn_info(self): 
        """
        Fetches BGP EVPN information from the router.
        Info like EVI, VXLAN Interface

        Returns:
        list: A list of BgpVpn instances.
        """                
        info = self._get_gnmi_info(['network-instance/protocols/bgp-evpn/bgp-instance'])
        bgp_evpn = []
        for network_instance in info:
            for bgp_instance in network_instance['protocols']['bgp-evpn']['srl_nokia-bgp-evpn:bgp-instance']:
                if not('oper-state' in bgp_instance):
                    bgp_evpn.append(self.BgpEvpn(network_instance['name'], bgp_instance['id'], bgp_instance['admin-state'], bgp_instance['vxlan-interface'], bgp_instance['evi'], bgp_instance['ecmp'], "up"))
                else:
                    bgp_evpn.append(self.BgpEvpn(network_instance['name'], bgp_instance['id'], bgp_instance['admin-state'], bgp_instance['vxlan-interface'], bgp_instance['evi'], bgp_instance['ecmp'], bgp_instance['oper-state']))
        return bgp_evpn

    def _get_bgp_vpn_info(self):
        """
        Fetches BGP VPN information from the router.

        Returns:
        list: A list of BgpVpn instances.
        """        
        info = self._get_gnmi_info(['network-instance/protocols/bgp-vpn/bgp-instance'])
        bgp_vpn = []
        for network_instance in info:
            for bgp_instance in network_instance['protocols']['srl_nokia-bgp-vpn:bgp-vpn']['bgp-instance']:
                bgp_vpn.append(self.BgpVpn(network_instance['name'], bgp_instance['id'], bgp_instance.get('route-distinguisher', {}).get('rd'), bgp_instance.get('route-target', {}).get('export-rt'),bgp_instance.get('route-target', {}).get('import-rt')))
        return bgp_vpn

def MergeEvpnToArray(srl_devices):
    rows = []
    for device in srl_devices:
        bgp_Evpn_dict = {item.network_instance: item for item in device.bgp_evpn}
        bgp_Vpn_dict = {item.network_instance: item for item in device.bgp_vpn}
        for key in bgp_Evpn_dict.keys():
            if key in bgp_Vpn_dict:
                rows.append([device.router, key, bgp_Evpn_dict[key].id, bgp_Evpn_dict[key].admin_state, 
                            bgp_Evpn_dict[key].vxlan_interface, bgp_Evpn_dict[key].evi, bgp_Evpn_dict[key].ecmp,
                            bgp_Evpn_dict[key].oper_state, bgp_Vpn_dict[key].rd, bgp_Vpn_dict[key].import_rt,
                            bgp_Vpn_dict[key].export_rt])
    if not rows:
        print("No data to display.")
        return    
    return rows   

def HighlightAlternateGroups(sorted_rows, column_to_check):
    """
    sorted_rows has been sorted out based on network instance already
    function will display a difference of evi between domains in different routers
    """
    lighted_rows = []
    grouped_rows = groupby(sorted_rows, key=lambda x: x[1])
    for network, group in grouped_rows:
        previous_value = None
        color_switch = False
        for row in list(group):
            if previous_value is not None and previous_value != row[column_to_check]:
                color_switch = not color_switch
            if color_switch:
                row[column_to_check] = f"\033[43m{row[column_to_check]}\033[0m"
            previous_value = row[column_to_check]
            lighted_rows.append(row)
    return lighted_rows

