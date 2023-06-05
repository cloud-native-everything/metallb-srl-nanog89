replaces:
{{ $target := index .Vars .TargetName }}
{{- range $netinstances := index $target "network-instances" }}
  - path: "/interface[name=irb0]/subinterface[index={{ index $netinstances "vni" }}]"
    encoding: "json_ietf"
    value:
      admin-state: enable
      anycast-gw: {}
      ipv4:
        address: 
          - ip-prefix: {{ index $netinstances "anycast-gw" }}
            anycast-gw: True
        arp:
          learn-unsolicited: True
          host-route:
            populate:
              - route-type: dynamic
          evpn:
            advertise:
              - route-type: dynamic             
  - path: "/tunnel-interface[name=vxlan2]/vxlan-interface[index={{ index $netinstances "vni" }}]"
    encoding: "json_ietf"
    value:
      type: bridged
      ingress:
        vni: {{ index $netinstances "vni" }}
  - path: "/network-instance[name={{ index $netinstances "name" }}]"
    encoding: "json_ietf"
    value: 
      admin-state: {{ index $netinstances "admin-state" | default "disable" }}
      type: {{ index $netinstances "type" | default "mac-vrf" }}
      description: {{ index $netinstances "description" | default "whatever" }}
      vxlan-interface:
        - name: vxlan2.{{ index $netinstances "vni" }}
      interface:
        - name: irb0.{{ index $netinstances "vni" }}        
      protocols:
        bgp-evpn:
          bgp-instance:
            - id: 2
              admin-state: enable
              vxlan-interface: vxlan2.{{ index $netinstances "vni" }}
              evi: {{ index $netinstances "evi" }}
        bgp-vpn:
          bgp-instance:
            - id: 2
              route-target:
                export-rt: target:65123:{{ index $netinstances "evi" }}
                import-rt: target:65123:{{ index $netinstances "evi" }}                           
{{- end }}

updates:
{{ $target := index .Vars .TargetName }}
{{- range $netinstances := index $target "network-instances" }}
  - path: "/network-instance[name=l3evpn]"
    encoding: "json_ietf"
    value:
      interface:
        - name: irb0.{{ index $netinstances "vni" }}                                   
{{- end }}