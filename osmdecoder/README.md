# omsdecoder

Inspired by [Lukas Rist](https://github.com/glaslos/go-osm).

osmdecoder decodes .osm files into structs provided by the package. Fields with attribute values are decoded to scalar or string pointers so that we can check if the field holds a zeroValue or not.

Protocol Buffers and FlatBuffers conversations are provided in osmpb and osmfbs package respectively.