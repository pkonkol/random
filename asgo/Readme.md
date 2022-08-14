*asgo* is primarily supposed to generate DB of ASes and their prefixes, geolocate
them based on whois data and then generate a website with the results.
The database may be later used to implement scanning by filters, but I'm not
planning on doing that now.

Old idea, deprecated:
```md
*asgo* is a program supposed to run masscan on autonomous systems deemed
interesting and then perform more accurate scans on hosts deemed interesting.
Results will be stored in a DB while ASes will be put on an interactive map
based on location retrieved from whois/organistaion/address fields.
```