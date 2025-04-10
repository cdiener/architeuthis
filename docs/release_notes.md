Those are the changes to `architeuthis` starting with version 0.3.0.

## 0.4.0

`architeuthis mapping filter` now allows the `--format` argument.

Switches the default highest rank in lineage annotaion to kingdom (`K__`) as superkingdom
has been removed from the newest NCBI taxonomy.

The annotation leaf is now detected from the back, making it resistant to missing
higher ranks.

`architeuthis mapping filter` logs will now indicate if a custom taxonomy directory is used.

## 0.3.1

Fixes the merging for Kraken2 files.

Add more checks during parsing.

## 0.3.0

Now supports Kraken2 output generated with the `--use-names` flag.

Adds documentation with mkdocs.
