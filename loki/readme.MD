# LOKI

loki.go acts as an API server that receives a request and returns it's  json response
Loki is the UI for ODIN and Wheagle. Designed with scalling in mind with a domain driven architecture
TO Add a feature, add it into the domain after defining its data types in definers, add it as a service to call it from the given domain then add url/views in the handlers
If a services does not need to be in the doamin,then *I don't know we'll figure this one out when it happens*

UI - Handles all templates for web apps
  handlers
  templates
