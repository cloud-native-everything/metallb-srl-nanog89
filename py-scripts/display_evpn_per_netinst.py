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

import yaml
import time
import logging
from srl_evpn_class import SrlDevice
from srl_evpn_class import MergeEvpnToArray
from srl_evpn_class import HighlightAlternateGroups
from tabulate import tabulate

# Define logging levels as constants
LOGGING_LEVEL_CRITICAL = logging.CRITICAL
LOGGING_LEVEL_WARNING = logging.WARNING

# Define magic numbers as constants
DEFAULT_SKIP_VERIFY = True
DEFAULT_MODEL = 'ixrd3'
DEFAULT_RELEASE = '21.6.4'

# Define the router info filename as a command line argument
import argparse
parser = argparse.ArgumentParser()
parser.add_argument('filename', help='the filename of the router info YAML file')
args = parser.parse_args()

start = time.time()

def main():
    try:
        with open(args.filename, 'r') as fh:
            router_info = yaml.safe_load(fh)
    except FileNotFoundError:
        print(f"File {args.filename} not found.")
        return
    except yaml.YAMLError as exc:
        print(f"Error in configuration file: {exc}")
        return

    try:
        switches = router_info['switches']
        routers = switches['srl']
        username = router_info['username']
        password = router_info['password']
        port = router_info['gnmi_port']
        skip_verify = router_info['skip_verify']
    except KeyError as e:
        print(f"Key {e} not found in configuration file.")
        return

    srl_devices = []
    for router in routers:
        srl_devices.append(SrlDevice(router, port, DEFAULT_MODEL, DEFAULT_RELEASE, username, password, skip_verify))

    rows = MergeEvpnToArray(srl_devices)

    if not rows:
        print("No data to display.")
        return


    sorted_rows = sorted(rows, key=lambda x: x[1])
    print("Table: Sorted by Network Instance")          
    highlighted_rows = HighlightAlternateGroups(sorted_rows, 5)  # Assuming Network Instance is the 1st column (0-indexed)
    table = tabulate(highlighted_rows, headers=['Router', 'Network instance', 'ID', 'EVPN Admin state', 
                                                'VXLAN interface', 'EVI', 'ECMP', 'Oper state', 
                                                'RD', 'import-rt', 'export-rt'], tablefmt="pretty")
    print(table)


if __name__ == '__main__':
    main()
    end = time.time()
    float_format = '{0:.2F}'
    print(f'Total time: {float_format.format(end - start)} seconds')