`architeuthis` is also capable of scoring and filtering reads on several metrics that go
a little father than the already available confidence scoring in Kraken2.

## Post-filtering reads

You can use the `mapping filter` to keep only those reads in a Kraken run that have
high quality mappings and drop the rest. This can sometimes avoid common issues such as
cross-domain mapping.

!!! tip "A diverse Kraken DB"
    This will work best when using the most diverse Kraken DB available, meaning a
    database that includes as many of the organisms in the habitat as feasible.

The supported metrics are:

Consistency
: Consistency denotes the fraction of kmer-level taxonomy assignments that are contained
  in the final read classification. So it qunatifies how many of the alternative
  classifications are contained in a single phylogentic path.

Multiplicity
: Multiplicity is the number of unique kmers assignments at the same rank r as the read
  assignment. So a read classified on the species level with a multiplicity of 2 matches
  2 unique species.

Entropy
: Entropy is the Shannon index of the kmer assignments on the same rank as the final read
  assignment. So it measures how surprising or distinct the alternative classifications are.
  One can see it as an abundance weighted multiplicity.


## Usage

To filter reads on one or more metric use

```bash
architeuthis mapping filter \
    --min-consistency 0.9 \
    --max-multiplicity 2 \
    --max-entropy 0.1 \
    --out filtered.k2 \
    my_sample.k2
```

This example is using the default parameters which should lead to fairly high quality
read classifications. The output is a valid Kraken output and a strict subset of the input
file. `mapping filter` supports the `--data-dir` option (see below).

## Scoring reads

It is also possible to only output the metrics for all classified reads in a sample.

## Usage

```bash
architeuthis mapping score my_sample.k2 --out scores.csv
```

## Output

```csv
ample_id,read_id,taxid,name,rank,n_kmers,consistency,confidence,multiplicity,entropy
testdata/negative,165179_NZ_CP102288.1_598818_598628_1_0_0_0_0:0:0_0:0:0_f59,165179,s__Segatella copri,s,153,1,1,1,0
testdata/negative,47678_NZ_CP081920.1_2131436_2131626_0_1_0_0_0:0:0_0:0:0_4749,816,g__Bacteroides,g,145,1,1,1,0
testdata/negative,821_NZ_CP103067.1_1529728_1529923_0_1_0_0_2:0:0_1:0:0_4c09,909656,g__Phocaeicola,g,68,1,1,1,0
testdata/negative,299767_NZ_CP099310.1_741506_741350_1_0_0_0_1:0:0_0:0:0_162,547,g__Enterobacter,g,82,0.975609756097561,0.96,2,0.167944147734173
testdata/negative,328813_NZ_AP019738.1_1486487_1486329_1_0_0_0_0:0:0_2:0:0_2ffc,328813,s__Alistipes onderdonkii,s,116,1,1,1,0
testdata/negative,418240_NZ_CP102267.1_4677848_4677952_0_1_0_0_0:0:0_3:0:0_66a6,1121115,s__Blautia wexlerae,s,185,1,1,1,0
testdata/negative,562_NZ_CP038408.1_5034459_5034717_0_1_0_0_0:0:0_0:0:0_f,543,f__Enterobacteriaceae,f,183,1,1,1,0
testdata/negative,821_NZ_CP103067.1_3126223_3126146_1_0_0_0_0:0:0_2:0:0_a0a1,909656,g__Phocaeicola,g,161,1,1,1,0
testdata/negative,821_NZ_CP043529.1_3737695_3737718_0_1_0_0_0:0:0_1:0:0_247c,821,s__Phocaeicola vulgatus,s,135,1,1,1,0
testdata/negative,46503_NZ_CP085927.1_2378513_2378638_0_1_0_0_1:0:0_1:0:0_8897,46503,s__Parabacteroides merdae,s,137,1,1,1,0
testdata/negative,39486_NZ_CP102279.1_1096816_1096614_1_0_0_0_1:0:0_1:0:0_2649,186803,f__Lachnospiraceae,f,108,1,1,1,0
testdata/negative,820_NZ_CP072255.1_61761_61514_1_0_0_0_0:0:0_1:0:0_1e488,820,s__Bacteroides uniformis,s,176,1,1,1,0
[...]
```

This also reports the Kraken confidence score using the provided taxonomy dump.

### Specifying the NCBI Taxonomy dump

You can use any downloaded [NCBI Taxonomy dump](https://ftp.ncbi.nlm.nih.gov/pub/taxonomy/taxdump.tar.gz)
for lineage annotation by specifying the `--data-dir` option, for instance:

```bash
architeuthis mapping score --data-dir /my/taxdump/ my_sample.k2 --out scores.csv
```

### Specifying the lineage format

You can specify the lineage format using [the taxonkit syntax](https://bioinf.shenwei.me/taxonkit/usage/#reformat).
The defualt lineage format is `{k};{p};{c};{o};{f};{g};{s}` which are the canonical ranks down
to species level. However, you could change this. For instance, to only keep genus and species:

```bash
architeuthis mapping score --format "{g};{s}" my_sample.k2 --out scores.csv