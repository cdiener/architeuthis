The `mapping` module contains tools to analyze the k-mer-level mapping results from
Kraken output. It also contains commands to [filter Kraken output](filter.md).

## K-mer mapping

The `kmer` subcommand allows to summarize mapping results in detail by resolving on
the k-mer least common ancestor assignments made in individual reads and summarizing
on the final classification. Its major purpose is to see which alternative mappings
exist for any given taxon in the final read-level Kraken classification. It can be
used to answer the following questions for instance:

- How many bacterial reads could also be classified as human?
- Are species well-resolved or do those reads also map to closely related species?

## Usage

To generate a map for all reads in a sample use

```bash
architeuthis mapping kmers my_sample.k2 --out mappings.csv `
```

## Output

This will look somehwat like this:

```csv
sample_id,classification,n_reads,taxid,n_kmers
test,28111,12,9749,5
test,28111,12,171549,7
test,28111,12,0,246
test,28111,12,4498,4
test,28111,12,815,3
test,28111,12,131567,1
test,28111,12,976,1
test,28111,12,28111,1728
test,28111,12,816,90
[...]
```

`sample_id` identified the classified file, `classification` the final taxon ID assigned
by Kraken2, `n_reads` denotes the number of reads assigned that taxon, `taxid` denotes the
taxon ID of the specific k-mers in the reads, and `n_kmers` denotes how many k-mers were assigned
that specific taxon ID. So for instance in the example aboved the final classification
was `28111` (*Bacteroides egghertii*) for 12 reads. The majority of individual k-mers (1728) were
also assigned that taxon, whereas for instace 90 k-mers were assigned to 816 which is the Bacteroides
genus.

## Taxonomic mapping summary

The mapping analyses can also be summarized in a taxonomy-centric manner by collapsing
on individual ranks. This is done with the `summary` subcommand.

## Usage

To generate summarize the mappings use

```bash
architeuthis mapping summary my_sample.k2 --out mapping_summary.csv
```

## Output

This will look like this:

```csv
sample_id,classification,lineage,total_reads,name,rank,kmers,in_lineage
test,815,k__Bacteria;p__Bacteroidota;c__Bacteroidia;o__Bacteroidales;f__Bacteroidaceae;g__;s__,13,p__Bacteroidota,p,2055,1
test,815,k__Bacteria;p__Bacteroidota;c__Bacteroidia;o__Bacteroidales;f__Bacteroidaceae;g__;s__,13,c__Bacteroidia,c,2055,1
test,815,k__Bacteria;p__Bacteroidota;c__Bacteroidia;o__Bacteroidales;f__Bacteroidaceae;g__;s__,13,o__Bacteroidales,o,2048,1
test,815,k__Bacteria;p__Bacteroidota;c__Bacteroidia;o__Bacteroidales;f__Bacteroidaceae;g__;s__,13,f__Bacteroidaceae,f,1794,1
test,815,k__Bacteria;p__Bacteroidota;c__Bacteroidia;o__Bacteroidales;f__Bacteroidaceae;g__;s__,13,k__Eukaryota,k,3,0
test,815,k__Bacteria;p__Bacteroidota;c__Bacteroidia;o__Bacteroidales;f__Bacteroidaceae;g__;s__,13,o__Fabales,o,3,0
test,815,k__Bacteria;p__Bacteroidota;c__Bacteroidia;o__Bacteroidales;f__Bacteroidaceae;g__;s__,13,f__Fabaceae,f,3,0
test,815,k__Bacteria;p__Bacteroidota;c__Bacteroidia;o__Bacteroidales;f__Bacteroidaceae;g__;s__,13,k__Bacteria,k,2096,1
test,815,k__Bacteria;p__Bacteroidota;c__Bacteroidia;o__Bacteroidales;f__Bacteroidaceae;g__;s__,13,p__Streptophyta,p,3,0
test,815,k__Bacteria;p__Bacteroidota;c__Bacteroidia;o__Bacteroidales;f__Bacteroidaceae;g__;s__,13,c__Magnoliopsida,c,3,0
```

`sample_id` identified the classified file, `classification` the final taxon ID assigned
by Kraken2, `lineage` the full lineage of the final classification,
`total_reads` denotes the number of reads assigned that taxon, `name` denotes the
taxonomic name of the specific k-mers in the reads, `rank` the taxonomic rank of the k-mer classificaton,
`n_kmers` denotes how many k-mers were assigned
that specific taxon, and `in_lineage` tells you whether that specific classification is contained within
the final read-leval classification or not (discordant mapping). So for instance in the example above you can
see that within the *Bacteroidaceae* 3 k-mers were mapped to Eukaryotes.

## Options

As for the [lineage](lineage.md) command you can use the `--data-dir` and `--format` options
to control the location of the NCBI taxonomy dumps and to specify the ranks included in the analysis.

For instance, to use a custom taxonomy and only the kingdom level you would use:

```
architeuthis mapping summary --data-dir /my/taxonomy --format "{k}" --out my_summary.csv my_sample.k2
```

!!! warning "Restrictions for the format"
    Note that `architeuthis mapping summary` only supports the `;` separator in the `--format`
    argument.