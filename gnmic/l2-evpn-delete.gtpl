deletes:
{{ $target := index .Vars .TargetName }}
{{- range $netinstances := index $target "network-instances" }}
  - "/network-instance[name={{ index $netinstances "name" }}]"
  - "/tunnel-interface[name=vxlan2]/vxlan-interface[index={{ index $netinstances "vni" }}]"
{{- end }}
