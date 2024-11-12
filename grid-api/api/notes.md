# notes

- more/less and after/before include the equal state
- slice-valued filters takes params separated by `,`
- string-valued filters is searching for a pattern and case-insensitive
- null dependant filters (used/rented) is not supported. use contract=0
- farms can be filtered based on nodes, and nodes on gpus
- if sorting param is true> asc. if false> desc. or nil.

## pkg job

- validate the inputs
- query with the query client
- convert to the output

- do nothing about routes/handlers/mw
- it is not aware about the http
