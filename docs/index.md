<img src="architeuthis.webp" width="60%">

[![Go](https://github.com/cdiener/architeuthis/actions/workflows/go.yml/badge.svg)](https://github.com/cdiener/architeuthis/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/cdiener/architeuthis/graph/badge.svg?token=KIeBFhClXY)](https://codecov.io/gh/cdiener/architeuthis)

> *architeuthis* is named after *Architeuthis dux*, the giant squid. It also sounds
> like "archi-do-this", so giving instructions to your pet kraken.

*architeuthis* is a fast standalone command to supplement the Kraken suite of software tools
such like Kraken2, KrakenUniq, and Bracken. I saw myself repeatedly rewriting the same
code in my pipelines when dealing with Kraken output, like merging files or maninpulating
lineage annotations. It also adds some functionality to dive deeper into the individual
k-mer classifications for reads.

## Main functionality

1. Merge or combine outputs from Kraken/Bracken across many samples efficiently
2. Add complete taxonomic lineage annotation to Bracken outputs
3. Analyze mapping across taxa, e.g.
    - How often did reads that matched one taxon also match another?
    - Are there cross-domain matches in my data set?
4. Filter Kraken mappings by several metrics
    - mapping consistency
    - multiple mappings on the final taxonomic rank
    - mapping entropy on the final taxonomic rank

## Usage

```text
Usage:
  architeuthis [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  lineage     Add lineage information to Bracken output.
  mapping     Analyze read and k-mer level mapping.
  merge       Merge various output files related to Kraken.

Flags:
      --db string   path to the Kraken database [optional]
  -h, --help        help for architeuthis
```
