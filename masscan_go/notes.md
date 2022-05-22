* Ripe API allows to get all AS numbers (and more)
    `https://stat.ripe.net/data/ris-asns/data.json?list_asns=true`
    * only PL ASes
        `https://stat.ripe.net/data/country-asns/data.json?resource=pl&lod=1`
* caida ASRank API https://api.asrank.caida.org/v2/docs eg.
    `https://api.asrank.caida.org/v2/restful/asns/?first=7000&offset=0`
    `https://api.asrank.caida.org/v2/restful/asns/12831`

* BGP checking sites
    https://bgpview.io/asn/12831#info
    https://bgp.he.net/AS12831#_asinfo
    https://ipinfo.io/AS12831#block-summary (allows for checking hsoted domains)
    https://bgp.tools/as/12831#asinfo

# Mongo filters & commands
* `db.as.find({country: "PL"}, {as_number: 1, name: 1, rank:1, number_asns:1, number_addresses:1, _id: 0}).sort({number_addresses: 1})`