The Data for the experiment was generated and downloaded from [Overpass turbo](http://overpass-turbo.eu) which is a web based data mining tool for OpenStreetMap. More information on [https://wiki.openstreetmap.org/wiki/Overpass_turbo](https://wiki.openstreetmap.org/wiki/Overpass_turbo).

Since we are delimiting us to parks, the following query was executed on different maps:

```xml
/*
This has been generated by the overpass-turbo wizard.
The original search was:
“leisure=park”
*/
[out:xml][timeout:25];
// gather results
(
  way["leisure"="park"]({{bbox}});
  relation["leisure"="park"]({{bbox}});
);
// print results
out body;
>;
out skel qt;
```