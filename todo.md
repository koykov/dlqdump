# TODO

* ~~split dump queue to separate parts~~
  * ~~dumping queue~~
  * ~~scheduled restore component~~
* ~~rename Size to Capacity~~
* ~~rename TimeLimit to FlushInterval~~
* ~~combine common encoding methods to one universal encoder~~
* think about Encoder/Decoder signatures
  * json.Marshaller
  * xml.Marshaller
  * gob.Encoder
* ~~implement abstract dumper (like afero)~~
  * ~~interface~~
  * ~~file dump implementation~~
    * ~~move Directory and FileMask params to file dumper~~
    * ~~write dump directly to file~~
